package user_entity

type UserRegisterRequest struct {
	Username string `json:"username" validate:"required,min=5,max=30"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required,min=5,max=30"`
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required,min=5,max=30"`
	Password string `json:"password" validate:"required,min=5,max=30"`
}
