package purchase_repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	var purchaseId string
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
		VALUES ($%d, 0, 0, 0, $%d, $%d) RETURNING id as purchase_id
	),`, len(values)+1, len(values)+2, len(values)+3)
	values = append(values, req.UserId, req.Latitude, req.Longitude)

	combinedQuery := fmt.Sprintf("%s\n%s, insert_order_items AS (INSERT INTO order_items (item_id, quantity, order_id)\n%s\nRETURNING id) SELECT count(insert_order_items), purchase_insert.purchase_id FROM purchase_insert JOIN insert_order_items ON 1=1	GROUP BY purchase_insert.purchase_id;", purchaseQuery, strings.Join(insertOrderQuery, ", "), strings.Join(insertOrderItemQuery, "\nUNION ALL\n"))

	if err := tx.QueryRow(ctx, combinedQuery, values...).Scan(&insertedItemCount, &purchaseId); err != nil {
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
		CalculatedEstimateId:           purchaseId,
	}, nil
}
