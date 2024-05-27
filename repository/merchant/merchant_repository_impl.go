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

func (repository *merchantRepositoryImpl) Search(ctx context.Context, searchQuery merchant_entity.SearchMerchantQuery) (*[]merchant_entity.SearchMerchantData, error) {
	query := `SELECT id, name, category, image_url, 
	JSON_BUILD_OBJECT('lat', latitude, 'long', longitude) AS location, created_at`

	var whereClause []string
	var searchParams []interface{}

	if searchQuery.Id != "" {
		whereClause = append(whereClause, fmt.Sprintf("id = $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.Id)
	}
	if searchQuery.Name != "" {
		whereClause = append(whereClause, fmt.Sprintf("name ~* $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.Name)
	}
	if searchQuery.Category != "" {
		whereClause = append(whereClause, fmt.Sprintf("category = $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.Category)
	}

	if len(whereClause) > 0 {
		query += " WHERE " + strings.Join(whereClause, " AND ")
	}

	var orderBy string
	if strings.ToLower(searchQuery.CreatedAt) == "asc" {
		orderBy = ` ORDER BY created_at ASC`
	} else {
		orderBy = ` ORDER BY created_at DESC`
	}
	query += orderBy

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
		return &[]merchant_entity.SearchMerchantData{}, err
	}
	defer rows.Close()

	merchants, err := pgx.CollectRows(rows, pgx.RowToStructByName[merchant_entity.SearchMerchantData])
	if err != nil {
		return &[]merchant_entity.SearchMerchantData{}, err
	}

	return &merchants, nil

}
