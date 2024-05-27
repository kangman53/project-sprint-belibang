package item_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	item_entity "github.com/kangman53/project-sprint-belibang/entity/item"
	exc "github.com/kangman53/project-sprint-belibang/exceptions"
)

type itemRepositoryImpl struct {
	DBpool *pgxpool.Pool
}

func NewItemRepository(dbPool *pgxpool.Pool) ItemRepository {
	return &itemRepositoryImpl{
		DBpool: dbPool,
	}
}

func (repository *itemRepositoryImpl) Add(ctx context.Context, item item_entity.Item, merchantId string) (string, error) {
	var itemId string
	query := "INSERT INTO items (name, category, image_url, price, merchant_id) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	if err := repository.DBpool.QueryRow(ctx, query, item.Name, item.Category, item.ImageUrl, item.Price, merchantId).Scan(&itemId); err != nil {
		return "", err
	}

	return itemId, nil
}

func (repository *itemRepositoryImpl) Search(ctx context.Context, searchQuery item_entity.SearchItemQuery, merchantId string) (*[]item_entity.SearchItemData, error) {
	// check merchant id exist
	existMerchantQuery := `SELECT id FROM merchants WHERE id = $1`
	if err := repository.DBpool.QueryRow(ctx, existMerchantQuery, merchantId).Scan(new(string)); err != nil {
		if err == pgx.ErrNoRows {
			return &[]item_entity.SearchItemData{}, exc.NotFoundException("Merchant ID does not exist")
		}
		return &[]item_entity.SearchItemData{}, err
	}

	query := `SELECT id, name, category, price, image_url, 
	to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') created_at
	FROM items`

	var whereClause = []string{` WHERE merchant_id = $1`}
	var searchParams = []interface{}{merchantId}

	if searchQuery.ItemId != "" {
		whereClause = append(whereClause, fmt.Sprintf("id = $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.ItemId)
	}
	if searchQuery.Name != "" {
		whereClause = append(whereClause, fmt.Sprintf("name ~* $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.Name)
	}
	if searchQuery.Category != "" {
		whereClause = append(whereClause, fmt.Sprintf("category = $%d", len(searchParams)+1))
		searchParams = append(searchParams, searchQuery.Category)
	}

	query += strings.Join(whereClause, " AND ")

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
		return &[]item_entity.SearchItemData{}, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[item_entity.SearchItemData])
	if err != nil {
		return &[]item_entity.SearchItemData{}, err
	}

	return &items, nil
}
