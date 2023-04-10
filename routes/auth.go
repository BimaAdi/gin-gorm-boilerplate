package routes

import (
	"net/http"

	"github.com/BimaAdi/ginGormBoilerplate/core"
	"github.com/BimaAdi/ginGormBoilerplate/models"
	"github.com/BimaAdi/ginGormBoilerplate/repository"
	"github.com/BimaAdi/ginGormBoilerplate/schemas"
	"github.com/gin-gonic/gin"
)

func authRoutes(rq *gin.RouterGroup) {
	auths := rq.Group("/auth")

	auths.POST("/login", authLoginRoute)
	auths.POST("/logout", authLogoutRoute)
}

// Login
//
//	@Summary		Login
//	@Description	login
//	@Tags			Auth
//	@Produce		json
//	@Param			payload	formData	schemas.LoginFormRequest	true	"form data"
//	@Success		200		{object}	schemas.LoginResponse
//	@Failure		400		{object}	schemas.BadRequestResponse
//	@Failure		500		{object}	schemas.InternalServerErrorResponse
//	@Router			/auth/login [post]
func authLoginRoute(c *gin.Context) {
	// Get data from form
	formRequest := schemas.LoginFormRequest{}
	err := c.ShouldBind(&formRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{
			Message: err.Error(),
		})
		return
	}

	// Get User
	user, err := repository.GetUserByUsername(models.DBConn, formRequest.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{
			Message: "invalid credentials",
		})
		return
	}

	// Check Password
	if !core.CheckPasswordHash(formRequest.Password, user.Password) {
		c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{
			Message: "invalid credentials",
		})
		return
	}

	// Generate JWT token
	token, err := core.GenerateJWTTokenFromUser(models.DBConn, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schemas.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	})
}

// Logout
//
//	@Summary		Logout
//	@Description	logout
//	@Tags			Auth
//	@Produce		json
//	@Success		200	{object}	schemas.LogoutResponse
//	@Failure		400	{object}	schemas.UnauthorizedResponse
//	@Failure		500	{object}	schemas.InternalServerErrorResponse
//	@Security		OAuth2Password
//	@Router			/auth/logout [post]
func authLogoutRoute(c *gin.Context) {
	// Authorize User
	user, err := core.GetUserFromAuthorizationHeader(models.DBConn, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, schemas.UnauthorizedResponse{
			Message: "Invalid/Expired token",
		})
		return
	}

	c.JSON(http.StatusOK, schemas.LogoutResponse{
		Email:    user.Email,
		Username: user.Username,
	})
}
