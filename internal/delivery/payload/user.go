package payload

type RegisterUserRequest struct {
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=4,max=100"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=4,max=100,eqfield=Password"`
}

type UserWithTokenResponse struct {
	UserInfo *UserInfo `json:"user_info"`
	Token    string    `json:"token"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4,max=100"`
}

type UpdateUserRequest struct {
	Name                 *string `json:"name"`
	Email                *string `json:"email" validate:"omitempty,email"`
	Password             *string `json:"password" validate:"omitempty,min=4,max=100"`
	PasswordConfirmation *string `json:"password_confirmation" validate:"omitempty,min=4,max=100,eqfield=Password"`
}

type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

const (
	ERROR_USER_NOT_FOUND     = "user not found"
	ERROR_USER_EXIST         = "user already exist"
	ERROR_USER_INVALID       = "invalid user"
	ERROR_WRONG_PASSWORD     = "wrong password"
	ERROR_TOKEN_INVALID      = "invalid token"
	ERROR_PASSWORD_NOT_MATCH = "password not match"
	ERROR_USER_NOT_LOGGED_IN = "user not logged in"
	ERROR_AUTHOR_NOT_FOUND   = "author not found"
)
