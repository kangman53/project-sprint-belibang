package controller

import (
	user_entity "github.com/kangman53/project-sprint-belibang/entity/user"
	exc "github.com/kangman53/project-sprint-belibang/exceptions"
	user_service "github.com/kangman53/project-sprint-belibang/service/user"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService user_service.UserService
}

func NewUserController(userService user_service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (controller UserController) Register(ctx *fiber.Ctx) error {
	userReq := new(user_entity.UserRegisterRequest)
	if err := ctx.BodyParser(userReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}
	resp, err := controller.UserService.Register(ctx, *userReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(resp)

}

func (controller UserController) Login(ctx *fiber.Ctx) error {
	userReq := new(user_entity.UserLoginRequest)
	if err := ctx.BodyParser(userReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}

	resp, err := controller.UserService.Login(ctx, *userReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
