package merchant_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	merchant_entity "github.com/kangman53/project-sprint-belibang/entity/merchant"
)

type merchantRepositoryImpl struct {
	DBpool *pgxpool.Pool
}

func NewMerchantRepository(dbPool *pgxpool.Pool) MerchantRepository {
	return &merchantRepositoryImpl{
		DBpool: dbPool,
	}
}

func (repository *merchantRepositoryImpl) Add(ctx context.Context, merchant merchant_entity.Merchant) (string, error) {
	var merchantId string
	query := "INSERT INTO merchants (name, category, image_url, latitude, longitude) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	if err := repository.DBpool.QueryRow(ctx, query, merchant.Name, merchant.Category, merchant.ImageUrl, merchant.Latitude, merchant.Longitude).Scan(&merchantId); err != nil {
		return "", err
	}

	return merchantId, nil
}

func (repository *merchantRepositoryImpl) SearchNearby(ctx context.Context, searchQuery merchant_entity.SearchNearbyMerchantQuery, latitude float64, longitude float64) (*[]merchant_entity.SearchNearbyMerchantData, error) {
	query := `SELECT 
	json_build_object('merchantId', m.id, 'name', m.name, 'merchantCategory', m.category, 'imageUrl', m.image_url, 'location', json_build_object('lat', m.latitude, 'long', m.longitude), 'createdAt', to_char(m.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')) merchant,
	jsonb_agg(json_build_object('itemId', i.id, 'name', i.name, 'productCategory', i.category, 'price', i.price, 'imageUrl', i.image_url, 'createdAt', to_char(i.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'))) items
	FROM items i
		JOIN merchants m ON i.merchant_id = m.id`

	var whereClause []string
	var searchParams []interface{}

	if searchQuery.MerchantId != "" {
		whereClause = append(whereClause, fmt.Sprintf("i.merchant_id = $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.MerchantId)
	}
	if searchQuery.Name != "" {
		whereClause = append(whereClause, fmt.Sprintf("m.name ~* $%d OR i.name ~* $%d", len(searchParams)+1, len(searchParams)+2))
		searchParams = append(searchParams, searchQuery.Name)
		searchParams = append(searchParams, searchQuery.Name)
	}
	if searchQuery.MerchantCategory != "" {
		whereClause = append(whereClause, fmt.Sprintf("m.category = $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.MerchantCategory)
	}

	if len(whereClause) > 0 {
		query += " WHERE " + strings.Join(whereClause, " AND ")
	}

	// Haversine Formula to calculate distance
	locationFilter := fmt.Sprintf(`(
        6371 * ACOS(
            COS(RADIANS(%f)) * COS(RADIANS(m.latitude)) * COS(RADIANS(m.longitude) - RADIANS(%f)) + 
            SIN(RADIANS(%f)) * SIN(RADIANS(m.latitude))
        )
    )`, latitude, longitude, latitude)

	// <= 1000 -> not bigger that 1000 km
	query += " GROUP BY m.id HAVING " + locationFilter + " <= 1000" + " ORDER BY " + locationFilter

	// handle limit or offset negative
	if searchQuery.Limit < 0 {
		searchQuery.Limit = 5
	}
	if searchQuery.Offset < 0 {
		searchQuery.Offset = 0
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", searchQuery.Limit, searchQuery.Offset)
	rows, err := repository.DBpool.Query(ctx, query, searchParams...)
	if err != nil {
		return &[]merchant_entity.SearchNearbyMerchantData{}, err
	}
	defer rows.Close()

	medicalRecords, err := pgx.CollectRows(rows, pgx.RowToStructByName[merchant_entity.SearchNearbyMerchantData])
	if err != nil {
		return &[]merchant_entity.SearchNearbyMerchantData{}, err
	}

	return &medicalRecords, nil
}
