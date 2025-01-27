package app

import (
	"github.com/kangman53/project-sprint-belibang/controller"
	"github.com/kangman53/project-sprint-belibang/helpers"

	user_repository "github.com/kangman53/project-sprint-belibang/repository/user"
	auth_service "github.com/kangman53/project-sprint-belibang/service/auth"
	file_service "github.com/kangman53/project-sprint-belibang/service/file"
	user_service "github.com/kangman53/project-sprint-belibang/service/user"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterBluePrint(app *fiber.App, dbPool *pgxpool.Pool) {
	validator := validator.New()
	// register custom validator
	helpers.RegisterCustomValidator(validator)

	authService := auth_service.NewAuthService()

	userRepository := user_repository.NewUserRepository(dbPool)
	userService := user_service.NewUserService(userRepository, authService, validator)
	userController := controller.NewUserController(userService)

	fileService := file_service.NewFileService()
	fileController := controller.NewFileController(fileService)

	// Users API
	adminApi := app.Group("/admin")
	adminApi.Post("/register", userController.Register)
	adminApi.Post("/login", userController.Login)
	userApi := app.Group("/user")
	userApi.Post("/login", userController.Login)
	userApi.Post("/register", userController.Register)

	// JWT middleware
	// app.Use(helpers.CheckTokenHeader)
	app.Use(helpers.GetTokenHandler())
	app.Post("/image", authService.AuthorizeRole("admin"), fileController.Upload)
}
