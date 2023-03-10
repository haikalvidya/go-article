package delivery

import (
	"net/http"

	"github.com/haikalvidya/go-article/internal/delivery/payload"
	"github.com/haikalvidya/go-article/pkg/common"
	"github.com/haikalvidya/go-article/pkg/utils"

	"github.com/labstack/echo/v4"
)

type userDelivery deliveryType

func (d *userDelivery) RegisterUser(c echo.Context) error {
	res := common.Response{}
	req := &payload.RegisterUserRequest{}

	c.Bind(req)

	if err := c.Validate(req); err != nil {
		res.Error = utils.GetErrorValidation(err)
		res.Status = false
		res.Message = "Failed Registration"
		return c.JSON(http.StatusBadRequest, res)
	}

	registRes, err := d.Usecase.User.Register(req)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	res.Message = "Success Registration"
	res.Data = registRes
	res.Status = true

	return c.JSON(http.StatusOK, res)
}

func (d *userDelivery) LoginUser(c echo.Context) error {
	res := common.Response{}
	req := &payload.LoginUserRequest{}

	c.Bind(req)

	if err := c.Validate(req); err != nil {
		res.Error = utils.GetErrorValidation(err)
		res.Status = false
		res.Message = "Failed Login"
		return c.JSON(http.StatusBadRequest, res)
	}

	registRes, err := d.Usecase.User.Login(req)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		if err.Error() == payload.ERROR_USER_NOT_FOUND {
			return c.JSON(http.StatusUnauthorized, res)
		} else {
			return c.JSON(http.StatusPaymentRequired, res)
		}
	}

	res.Message = "Success Login"
	res.Data = registRes
	res.Status = true

	return c.JSON(http.StatusOK, res)
}

func (d *userDelivery) LogoutUser(c echo.Context) error {
	res := common.Response{}

	userId := d.Middleware.JWT.GetUserIdFromJwt(c)

	err := d.Usecase.User.Logout(userId)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	res.Message = "Success Logout"
	return c.JSON(http.StatusOK, res)
}

func (d *userDelivery) DeleteUser(c echo.Context) error {
	res := common.Response{}
	userId := d.Middleware.JWT.GetUserIdFromJwt(c)

	err := d.Usecase.User.DeleteAccount(userId)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	res.Message = "Success Delete User"
	return c.JSON(http.StatusOK, res)
}

func (d *userDelivery) UpdateUser(c echo.Context) error {
	res := common.Response{}
	req := &payload.UpdateUserRequest{}

	c.Bind(req)

	if err := c.Validate(req); err != nil {
		res.Error = utils.GetErrorValidation(err)
		res.Status = false
		res.Message = "Failed Update User"
		return c.JSON(http.StatusBadRequest, res)
	}

	userId := d.Middleware.JWT.GetUserIdFromJwt(c)

	err := d.Usecase.User.UpdateUser(userId, req)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	res.Status = true
	res.Message = "Success Update User"
	return c.JSON(http.StatusOK, res)
}

func (d *userDelivery) GetUser(c echo.Context) error {
	res := common.Response{}
	userId := d.Middleware.JWT.GetUserIdFromJwt(c)

	user, err := d.Usecase.User.GetUser(userId)
	if err != nil {
		res.Status = false
		res.Message = err.Error()
		return c.JSON(http.StatusBadRequest, res)
	}

	res.Message = "Success Get User"
	res.Data = user
	return c.JSON(http.StatusOK, res)
}
