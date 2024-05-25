package user_repository

import (
	"context"

	user_entity "github.com/kangman53/project-sprint-belibang/entity/user"
)

type UserRepository interface {
	Register(ctx context.Context, req user_entity.User) (string, error)
	Login(ctx context.Context, req user_entity.User) (user_entity.User, error)
}
