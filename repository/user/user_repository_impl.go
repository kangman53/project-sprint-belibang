package user_repository

import (
	"context"

	user_entity "github.com/kangman53/project-sprint-belibang/entity/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepositoryImpl struct {
	DBpool *pgxpool.Pool
}

func NewUserRepository(dbPool *pgxpool.Pool) UserRepository {
	return &userRepositoryImpl{
		DBpool: dbPool,
	}
}

func (repository *userRepositoryImpl) Register(ctx context.Context, user user_entity.User) (string, error) {
	var userId string
	query := "INSERT INTO users (username, email, role, password) VALUES ($1, $2, $3, $4) RETURNING id"
	if err := repository.DBpool.QueryRow(ctx, query, user.Username, user.Email, user.Role, user.Password).Scan(&userId); err != nil {
		return "", err
	}

	return userId, nil
}

func (repository *userRepositoryImpl) Login(ctx context.Context, user user_entity.User) (user_entity.User, error) {
	query := "SELECT id, password FROM users WHERE username = $1 AND role = $2 LIMIT 1"
	row := repository.DBpool.QueryRow(ctx, query, user.Username, user.Role)

	var loggedInUser user_entity.User
	err := row.Scan(&loggedInUser.Id, &loggedInUser.Password)
	if err != nil {
		return user_entity.User{}, err
	}

	return loggedInUser, nil
}
