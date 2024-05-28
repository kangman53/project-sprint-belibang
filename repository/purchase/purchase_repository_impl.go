package purchase_repository

import (
	"context"
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
	var totalPrice, estimateDeliveryTime int
	var distance float32
	var insertOrderQuery, selectOrderItemQuery, insertOrderItemQuery, itemIds []string

	for i, order := range req.Order {
		if *order.IsStartingPoint && startingPoint == 0 {
			startingPoint += 1
			merchantId = &order.MechantId
		} else if *order.IsStartingPoint && startingPoint == 1 {
			return purchase_entity.Purchase{}, exceptions.BadRequestException("more than 1 starting point")
		}

		insertOrderQuery = append(insertOrderQuery, fmt.Sprintf(`insert_order_%d AS (
			INSERT INTO orders (merchant_id, is_starting_point, purchase_id) 
			VALUES ('%s', %t, (SELECT purchase_id FROM purchase_insert))
			RETURNING id as order_id
		)`, i, order.MechantId, *order.IsStartingPoint))

		selectOrderItemQuery = nil
		for j, orderItem := range order.OrderItem {
			itemIds = append(itemIds, "'"+orderItem.ItemId+"'")
			selectOrderItemQuery = append(selectOrderItemQuery, fmt.Sprintf("SELECT '%s' AS item_id, %d AS quantity", orderItem.ItemId, orderItem.Quantity))

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
		VALUES ('%s', 0, 0, 0, %f, %f) RETURNING id as purchase_id
	),`, req.UserId, req.Latitude, req.Longitude)
	combinedQuery := fmt.Sprintf("%s\n%s,\n\ninsert_order_item AS (INSERT INTO order_items (item_id, quantity, order_id)\n%s) SELECT purchase_id FROM purchase_insert;", purchaseQuery, strings.Join(insertOrderQuery, ", "), strings.Join(insertOrderItemQuery, "\nUNION ALL\n"))
	if err := repository.DBpool.QueryRow(ctx, combinedQuery).Scan(&purchaseId); err != nil {
		return purchase_entity.Purchase{}, err
	}

	// update total_price, distance, and estiamted_delivery_time query
	updateQuery := fmt.Sprintf(`
	UPDATE purchases
	SET total_price = (
		SELECT COALESCE(SUM(items.price * order_items.quantity), 0) AS total_price
		FROM purchases
		JOIN orders ON purchases.id = orders.purchase_id
		JOIN order_items ON orders.id = order_items.order_id
		JOIN items ON order_items.item_id = items.id
		WHERE items.id in  (%s)
		AND purchases.id = '%s'
	),
	distance = (
		SELECT ROUND (6371 * ACOS(
            COS(RADIANS(%f)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS(%f)) + 
            SIN(RADIANS(%f)) * SIN(RADIANS(latitude))
        )::numeric, 2)
		FROM merchants
		WHERE id = '%s'
	),
	estimated_delivery_time = (
		SELECT ((6371 * ACOS(
			COS(RADIANS(%f)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS(%f)) + 
			SIN(RADIANS(%f)) * SIN(RADIANS(latitude))
		)) / 40) * 60
		FROM merchants
		WHERE id = '%s'
	)
	WHERE id = '%s'
	AND (
		SELECT 6371 * ACOS(
            COS(RADIANS(%f)) * COS(RADIANS(latitude)) * COS(RADIANS(longitude) - RADIANS(%f)) + 
            SIN(RADIANS(%f)) * SIN(RADIANS(latitude))
        )
		FROM merchants
		WHERE id = '%s'
	) <= 30
	returning total_price, distance, estimated_delivery_time;
	`, strings.Join(itemIds, ", "), purchaseId, req.Latitude, req.Longitude, req.Latitude, *merchantId, req.Latitude, req.Longitude, req.Latitude, *merchantId, purchaseId, req.Latitude, req.Longitude, req.Latitude, *merchantId)
	if err := repository.DBpool.QueryRow(ctx, updateQuery).Scan(&totalPrice, &distance, &estimateDeliveryTime); err != nil {
		return purchase_entity.Purchase{}, err
	}

	return purchase_entity.Purchase{
		TotalPrice:                     totalPrice,
		Distance:                       distance,
		EstimatedDeliveryTimeInMinutes: estimateDeliveryTime,
		CalculatedEstimateId:           purchaseId,
	}, nil
}
