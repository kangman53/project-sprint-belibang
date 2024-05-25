package auth_service

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	exc "github.com/kangman53/project-sprint-belibang/exceptions"
)

type authServiceImpl struct {
}

func NewAuthService() AuthService {
	return &authServiceImpl{}
}

func (service *authServiceImpl) GenerateToken(ctx context.Context, userId string, role string) (string, error) {
	// 8 hours
	var expDuration = time.Now().Add(time.Hour * 8).Unix()
	jwtconf := jwt.MapClaims{
		"user_id": userId,
		"exp":     expDuration,
		"role":    role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtconf)
	signToken, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signToken, nil
}

func (service *authServiceImpl) AuthorizeRole(role string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if userRole := ctx.Locals("userRole"); userRole != role {
			return exc.ForbiddenException("Access Forbidden")
		}
		return ctx.Next()
	}
}
