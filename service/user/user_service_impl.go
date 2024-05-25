package user_service

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	user_entity "github.com/kangman53/project-sprint-belibang/entity/user"
	exc "github.com/kangman53/project-sprint-belibang/exceptions"
	helpers "github.com/kangman53/project-sprint-belibang/helpers"
	userRep "github.com/kangman53/project-sprint-belibang/repository/user"
	authService "github.com/kangman53/project-sprint-belibang/service/auth"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type userServiceImpl struct {
	UserRepository userRep.UserRepository
	AuthService    authService.AuthService
	Validator      *validator.Validate
}

func NewUserService(userRepository userRep.UserRepository, authService authService.AuthService, validator *validator.Validate) UserService {
	return &userServiceImpl{
		UserRepository: userRepository,
		AuthService:    authService,
		Validator:      validator,
	}
}

func (service *userServiceImpl) Register(ctx *fiber.Ctx, req user_entity.UserRegisterRequest) (user_entity.UserResponse, error) {
	// validate by rule we defined in _request_entity.go
	if err := service.Validator.Struct(req); err != nil {
		return user_entity.UserResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	hashPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return user_entity.UserResponse{}, err
	}
	role := strings.Split(ctx.OriginalURL(), "/")[1]
	user := user_entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashPassword,
		Role:     role,
	}

	userContext := ctx.UserContext()
	userId, err := service.UserRepository.Register(userContext, user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return user_entity.UserResponse{}, exc.ConflictException("User with this username/email already registered")
		}
		return user_entity.UserResponse{}, err
	}

	token, err := service.AuthService.GenerateToken(userContext, userId, role)
	if err != nil {
		return user_entity.UserResponse{}, err
	}

	return user_entity.UserResponse{
		Token: token,
	}, nil
}

func (service *userServiceImpl) Login(ctx *fiber.Ctx, req user_entity.UserLoginRequest) (user_entity.UserResponse, error) {
	// validate by rule we defined in _request_entity.go
	if err := service.Validator.Struct(req); err != nil {
		return user_entity.UserResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	role := strings.Split(ctx.OriginalURL(), "/")[1]
	user := user_entity.User{
		Username: req.Username,
		Password: req.Password,
		Role:     role,
	}

	fmt.Println(role)

	userContext := ctx.UserContext()
	userLogin, err := service.UserRepository.Login(userContext, user)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return user_entity.UserResponse{}, exc.NotFoundException("User is not found")
		}
		return user_entity.UserResponse{}, err
	}

	if _, err = helpers.ComparePassword(userLogin.Password, req.Password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return user_entity.UserResponse{}, exc.BadRequestException("Invalid password")
		}

		return user_entity.UserResponse{}, err
	}

	token, err := service.AuthService.GenerateToken(userContext, userLogin.Id, role)
	if err != nil {
		return user_entity.UserResponse{}, err
	}

	return user_entity.UserResponse{
		Token: token,
	}, nil
}
