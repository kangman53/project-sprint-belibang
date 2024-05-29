package app

import (
	"github.com/kangman53/project-sprint-belibang/controller"
	"github.com/kangman53/project-sprint-belibang/helpers"

	item_repository "github.com/kangman53/project-sprint-belibang/repository/item"
	merchant_repository "github.com/kangman53/project-sprint-belibang/repository/merchant"
	purchase_repository "github.com/kangman53/project-sprint-belibang/repository/purchase"
	user_repository "github.com/kangman53/project-sprint-belibang/repository/user"
	auth_service "github.com/kangman53/project-sprint-belibang/service/auth"
	file_service "github.com/kangman53/project-sprint-belibang/service/file"
	item_service "github.com/kangman53/project-sprint-belibang/service/item"
	merchant_service "github.com/kangman53/project-sprint-belibang/service/merchant"
	purchase_service "github.com/kangman53/project-sprint-belibang/service/purchase"
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

	merchantRepository := merchant_repository.NewMerchantRepository(dbPool)
	merchantService := merchant_service.NewMerchantService(merchantRepository, validator)
	merchantController := controller.NewMerchantController(merchantService)

	itemRepository := item_repository.NewItemRepository(dbPool)
	itemService := item_service.NewItemService(itemRepository, validator)
	itemController := controller.NewItemController(itemService)

	purchaseRepository := purchase_repository.NewPurchaseRepository(dbPool)
	purchaseService := purchase_service.NewPurchaseService(purchaseRepository, authService, validator)
	purchaseController := controller.NewPurchaseController(purchaseService)

	// Users API
	adminApi := app.Group("/admin")
	adminApi.Post("/register", userController.Register)
	adminApi.Post("/login", userController.Login)
	userApi := app.Group("/users")
	userApi.Post("/login", userController.Login)
	userApi.Post("/register", userController.Register)

	// JWT middleware
	// app.Use(helpers.CheckTokenHeader)
	app.Use(helpers.GetTokenHandler())
	app.Post("/image", authService.AuthorizeRole("admin"), fileController.Upload)

	// Merchants API
	merchantApi := adminApi.Group("/merchants")
	merchantApi.Post("/", authService.AuthorizeRole("admin"), merchantController.Add)

	// Items API
	itemsApi := merchantApi.Group("/:merchantId/items")
	itemsApi.Post("/", authService.AuthorizeRole("admin"), itemController.Add)

	// Puchase API
	app.Get("/merchants/nearby/:coordinate", authService.AuthorizeRole("users"), merchantController.SearchNearby)
	userApi.Post("/estimate", authService.AuthorizeRole("users"), purchaseController.Estimate)
	userApi.Post("/orders", authService.AuthorizeRole("users"), purchaseController.Order)
	userApi.Get("/orders", authService.AuthorizeRole("users"), purchaseController.HistoryOrder)
}
