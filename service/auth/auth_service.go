package auth_service

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	GenerateToken(ctx context.Context, userId string, role string) (string, error)
	AuthorizeRole(role string) fiber.Handler
	GetValidUser(ctx *fiber.Ctx) (string, error)
}
