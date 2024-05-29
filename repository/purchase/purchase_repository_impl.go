package purchase_repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	purchase_entity "github.com/kangman53/project-sprint-belibang/entity/Purchase"
	"github.com/kangman53/project-sprint-belibang/exceptions"
)

type purchaseRepositoryImpl struct {
	DBpool *pgxpool.Pool
}

func NewPurchaseRepository(dbPool *pgxpool.Pool) PurchaseRepository {
	return &purchaseRepositoryImpl{
		DBpool: dbPool,
	}
}

func (repository *purchaseRepositoryImpl) Estimate(ctx context.Context, req purchase_entity.Purchase) (purchase_entity.Purchase, error) {
	startingPoint := 0
	var purchaseId, calculatedEstimateId string
	var merchantId *string
	var distance float32
	var totalPrice, estimateDeliveryTime, insertedItemCount int
	var insertOrderQuery, selectOrderItemQuery, insertOrderItemQuery, itemIds []string
	var values []interface{}

	tx, err := repository.DBpool.Begin(ctx)
	if err != nil {
		return purchase_entity.Purchase{}, err
	}
	defer tx.Rollback(ctx)

	for i, order := range req.Order {
		if *order.IsStartingPoint && startingPoint == 0 {
			startingPoint += 1
			merchantId = &order.MechantId
		} else if *order.IsStartingPoint && startingPoint == 1 {
			return purchase_entity.Purchase{}, exceptions.BadRequestException("more than 1 starting point")
		}

		insertOrderQuery = append(insertOrderQuery, fmt.Sprintf(`insert_order_%d AS (
			INSERT INTO orders (merchant_id, is_starting_point, purchase_id) 
			VALUES ($%d, $%d, (SELECT purchase_id FROM purchase_insert))
			RETURNING id as order_id
		)`, i, len(values)+1, len(values)+2))
		values = append(values, order.MechantId, *order.IsStartingPoint)

		selectOrderItemQuery = nil
		for j, orderItem := range order.OrderItem {
			itemIds = append(itemIds, "'"+orderItem.ItemId+"'")
			selectOrderItemQuery = append(selectOrderItemQuery, fmt.Sprintf("SELECT $%d AS item_id, CAST($%d AS INTEGER) AS quantity", len(values)+1, len(values)+2))
			values = append(values, orderItem.ItemId, orderItem.Quantity)

			if j+1 == len(order.OrderItem) {
				insertOrderItemQuery = append(insertOrderItemQuery, fmt.Sprintf(`SELECT item_id, quantity, (SELECT order_id FROM insert_order_%d)
					FROM (
						%s
					) AS order_item_insert_%d
				WHERE EXISTS (
					SELECT 1
						FROM items
					WHERE items.id = order_item_insert_%d.item_id
				)`, i, strings.Join(selectOrderItemQuery, " UNION ALL "), j, j))
			}
		}
	}

	purchaseQuery := fmt.Sprintf(`WITH 
	purchase_insert AS (
		INSERT INTO purchases (user_id, total_price, estimated_delivery_time, distance, latitude, longitude)
		VALUES ($%d, 0, 0, 0, $%d, $%d) RETURNING id as purchase_id, calculated_estimate_id
	),`, len(values)+1, len(values)+2, len(values)+3)
	values = append(values, req.UserId, req.Latitude, req.Longitude)

	combinedQuery := fmt.Sprintf("%s\n%s, insert_order_items AS (INSERT INTO order_items (item_id, quantity, order_id)\n%s\nRETURNING id) SELECT count(insert_order_items), purchase_insert.purchase_id, purchase_insert.calculated_estimate_id FROM purchase_insert JOIN insert_order_items ON 1=1	GROUP BY purchase_insert.purchase_id, purchase_insert.calculated_estimate_id;", purchaseQuery, strings.Join(insertOrderQuery, ", "), strings.Join(insertOrderItemQuery, "\nUNION ALL\n"))

	if err := tx.QueryRow(ctx, combinedQuery, values...).Scan(&insertedItemCount, &purchaseId, &calculatedEstimateId); err != nil {
		return purchase_entity.Purchase{}, err
	}

	if insertedItemCount != len(itemIds) {
		tx.Rollback(ctx)
		return purchase_entity.Purchase{}, errors.New("itemId not found")
	}

	// update total_price, distance, and estiamted_delivery_time query
	values = nil
	updateQuery := fmt.Sprintf(`
	UPDATE purchases
	SET total_price = (
		SELECT COALESCE(SUM(items.price * order_items.quantity), 0) AS total_price
		FROM purchases
		JOIN orders ON purchases.id = orders.purchase_id
		JOIN order_items ON orders.id = order_items.order_id
		JOIN items ON order_items.item_id = items.id
		WHERE items.id in  (%s)
		AND purchases.id = $%d
	),
	distance = (
		SELECT ROUND (6371 * ACOS(
            COS(RADIANS($%d)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS($%d)) + 
            SIN(RADIANS($%d)) * SIN(RADIANS(latitude))
        )::numeric, 2)
		FROM merchants
		WHERE id = $%d
	),
	estimated_delivery_time = (
		SELECT ((6371 * ACOS(
			COS(RADIANS($%d)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS($%d)) + 
			SIN(RADIANS($%d)) * SIN(RADIANS(latitude))
		)) / 40) * 60
		FROM merchants
		WHERE id = $%d
	)
	WHERE id = $%d
	AND (
		SELECT 6371 * ACOS(
            COS(RADIANS($%d)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS($%d)) + 
            SIN(RADIANS($%d)) * SIN(RADIANS(latitude))
        )
		FROM merchants
		WHERE id = $%d
	) <= 3
	returning total_price, distance, estimated_delivery_time;
	`, strings.Join(itemIds, ", "), len(values)+1, len(values)+2, len(values)+3, len(values)+4, len(values)+5, len(values)+6, len(values)+7, len(values)+8, len(values)+9, len(values)+10, len(values)+11, len(values)+12, len(values)+13, len(values)+14)
	values = append(values, purchaseId, req.Latitude, req.Longitude, req.Latitude, *merchantId, req.Latitude, req.Longitude, req.Latitude, *merchantId, purchaseId, req.Latitude, req.Longitude, req.Latitude, *merchantId)

	if err := tx.QueryRow(ctx, updateQuery, values...).Scan(&totalPrice, &distance, &estimateDeliveryTime); err != nil {
		return purchase_entity.Purchase{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return purchase_entity.Purchase{}, err
	}

	return purchase_entity.Purchase{
		TotalPrice:                     totalPrice,
		Distance:                       distance,
		EstimatedDeliveryTimeInMinutes: estimateDeliveryTime,
		CalculatedEstimateId:           calculatedEstimateId,
	}, nil
}

func (repository *purchaseRepositoryImpl) Order(ctx context.Context, req purchase_entity.Purchase) (purchase_entity.Purchase, error) {
	var purchaseId string
	query := "UPDATE purchases SET status = 'ordered' WHERE calculated_estimate_id = $1 AND user_id = $2 AND status = 'pending' RETURNING id"
	if err := repository.DBpool.QueryRow(ctx, query, req.CalculatedEstimateId, req.UserId).Scan(&purchaseId); err != nil {
		if err == pgx.ErrNoRows {
			return purchase_entity.Purchase{}, errors.New("not found")
		}
		return purchase_entity.Purchase{}, err
	}

	return purchase_entity.Purchase{
		Id: purchaseId,
	}, nil
}

func (repository *purchaseRepositoryImpl) HistoryOrder(ctx context.Context, searchQuery purchase_entity.SearcHistoryOrderQuery, userId string) (*[]purchase_entity.SearchHistoryOrderResponse, error) {
	query := `SELECT jsonb_build_object(
				'orderId', p.id,
				'orders', jsonb_agg(
					jsonb_build_object(
						'merchant', jsonb_build_object(
							'merchantId', m.id,
							'name', m.name,
							'merchantCategory', m.category,
							'imageUrl', m.image_url,
							'location', jsonb_build_object(
								'lat', m.latitude,
								'long', m.longitude
							),
							'createdAt', to_char(m.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
						),
						'items', jsonb_agg(
							jsonb_build_object(
								'itemId', i.id,
								'name', i.name,
								'productCategory', i.category,
								'price', i.price,
								'quantity', oi.quantity,
								'imageUrl', i.image_url,
								'createdAt', to_char(i.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
							)
						)
					)
				)
			) AS result
	FROM purchases p
	JOIN orders o ON p.id = o.purchase_id
	JOIN merchants m ON o.merchant_id = m.id
	JOIN order_items oi ON o.id = oi.order_id
	JOIN items i ON oi.item_id = i.id
	`

	var whereClause = []string{` WHERE p.user_id = $1`}
	var searchParams = []interface{}{userId}

	if searchQuery.MerchantId != "" {
		whereClause = append(whereClause, fmt.Sprintf("m.id = $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.MerchantId)
	}
	if searchQuery.Name != "" {
		whereClause = append(whereClause, fmt.Sprintf("(m.name ~* $%d OR i.name ~* $%d)", len(searchParams)+1, len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.Name)
	}
	if searchQuery.MerchantCategory != "" {
		whereClause = append(whereClause, fmt.Sprintf("m.category = $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.MerchantCategory)
	}

	query += strings.Join(whereClause, " AND ")
	query += " GROUP BY p.id"

	if searchQuery.Limit < 0 {
		searchQuery.Limit = 5
	}
	if searchQuery.Offset < 0 {
		searchQuery.Offset = 0
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", searchQuery.Limit, searchQuery.Offset)
	rows, err := repository.DBpool.Query(ctx, query, searchParams...)
	if err != nil {
		return &[]purchase_entity.SearchHistoryOrderResponse{}, err
	}
	defer rows.Close()

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[purchase_entity.SearchHistoryOrderResponse])
	if err != nil {
		return &[]purchase_entity.SearchHistoryOrderResponse{}, err
	}

	return &orders, nil

}
