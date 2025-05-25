package handler

import (
	todo "do-app"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Sign Up
// @Tags auth
// @Description createUser
// @ID createUsers
// @Accept json
// @Produce json
// @Param input body todo.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) SignUp(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil {

		newErrorResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Sign In
// @Tags auth
// @Description auth user
// @ID auth-user
// @Accept json
// @Produce json
// @Param input body signInInput true "login and password"
// @Success 200 {integer} integer 1
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {

		newErrorResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
