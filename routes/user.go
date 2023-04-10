package routes

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/BimaAdi/ginGormBoilerplate/core"
	"github.com/BimaAdi/ginGormBoilerplate/models"
	"github.com/BimaAdi/ginGormBoilerplate/repository"
	"github.com/BimaAdi/ginGormBoilerplate/schemas"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func userRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/user")

	users.GET("/", GetAllUserRoute)

	users.GET("/:userId", GetDetailUserRoute)

	users.POST("/", CreateUserRoute)

	users.PUT("/:userId", UpdateUserRoute)

	users.DELETE("/:userId", DeleteUserRoute)
}

// Get All User
//
//	@Summary		Get All User
//	@Description	Get All User
//	@Tags			User
//	@Produce		json
//	@Param			page		query		int	false	"page"
//	@Param			page_size	query		int	false	"page"
//	@Success		200			{object}	schemas.UserPaginateResponse
//	@Failure		400			{object}	schemas.BadRequestResponse
//	@Failure		401			{object}	schemas.UnauthorizedResponse
//	@Failure		500			{object}	schemas.InternalServerErrorResponse
//	@Security		OAuth2Password
//	@Router			/user/ [get]
func GetAllUserRoute(c *gin.Context) {
	// Authorize User
	_, err := core.GetUserFromAuthorizationHeader(models.DBConn, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, schemas.UnauthorizedResponse{
			Message: "Invalid/Expired token",
		})
		return
	}

	// Get Query Parameter
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	search := c.Query("search")
	pageInt, errPage := strconv.Atoi(page)
	pageSizeInt, errPageSize := strconv.Atoi((pageSize))
	if errPage != nil || errPageSize != nil {
		errorResponse := []map[string]string{}
		if errPage != nil {
			x := map[string]string{
				"page": "invalid page, page should integer",
			}
			errorResponse = append(errorResponse, x)
		}

		if errPageSize != nil {
			x := map[string]string{
				"page_size": "invalid page_size, page_size should integer",
			}
			errorResponse = append(errorResponse, x)
		}

		c.JSON(http.StatusUnprocessableEntity, schemas.UnprocessableEntityResponse{
			Message: errorResponse,
		})
		return
	}
	var searchNilable *string = nil
	if search != "" {
		searchNilable = &search
	}

	users, numData, numPage, err := repository.GetPaginatedUser(models.DBConn, pageInt, pageSizeInt, searchNilable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
		return
	}

	arrayDetailUser := []schemas.UserDetailResponse{}
	for _, item := range users {
		arrayDetailUser = append(arrayDetailUser, schemas.UserDetailResponse{
			Id:       item.ID,
			Username: item.Username,
			Email:    item.Email,
			IsActive: item.IsActive,
		})
	}

	c.JSON(http.StatusOK, schemas.UserPaginateResponse{
		Counts:    int(numData),
		PageCount: int(numPage),
		PageSize:  pageSizeInt,
		Page:      pageInt,
		Results:   arrayDetailUser,
	})
}

// Get Detail User
//
//	@Summary		Get Detail User
//	@Description	Get detail user
//	@Tags			User
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	schemas.UserDetailResponse
//	@Failure		400	{object}	schemas.BadRequestResponse
//	@Failure		404	{object}	schemas.NotFoundResponse
//	@Failure		500	{object}	schemas.InternalServerErrorResponse
//	@Security		OAuth2Password
//	@Router			/user/{id} [get]
func GetDetailUserRoute(c *gin.Context) {
	// Authorize User
	_, err := core.GetUserFromAuthorizationHeader(models.DBConn, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, schemas.UnauthorizedResponse{
			Message: "Invalid/Expired token",
		})
		return
	}

	// Get Params
	userId := c.Params.ByName("userId")
	if !core.IsValidUUID(userId) {
		c.JSON(http.StatusNotFound, schemas.NotFoundResponse{
			Message: "user not found",
		})
		return
	}

	user, err := repository.GetUserById(models.DBConn, userId)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, schemas.NotFoundResponse{
				Message: "user not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schemas.UserDetailResponse{
		Id:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		IsActive:    user.IsActive,
		IsSuperuser: user.IsSuperuser,
	})
}

// Create User
//
//	@Summary		Create User
//	@Description	Create User
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		schemas.UserCreateRequest	true	"Create User"
//	@Success		200		{object}	schemas.UserCreateResponse
//	@Failure		400		{object}	schemas.BadRequestResponse
//	@Failure		500		{object}	schemas.InternalServerErrorResponse
//	@Security		OAuth2Password
//	@Router			/user/ [post]
func CreateUserRoute(c *gin.Context) {
	// Authorize User
	_, err := core.GetUserFromAuthorizationHeader(models.DBConn, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, schemas.UnauthorizedResponse{
			Message: "Invalid/Expired token",
		})
		return
	}

	var newUser schemas.UserCreateRequest
	err = c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{
			Message: err.Error(),
		})
		return
	}

	now := time.Now()
	createdUser, err := repository.CreateUser(
		models.DBConn,
		newUser.Username,
		newUser.Email,
		newUser.Password,
		newUser.IsActive,
		newUser.IsSuperuser,
		now,
		&now,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, schemas.UserCreateResponse{
		Id:          createdUser.ID,
		Username:    createdUser.Username,
		Email:       createdUser.Email,
		IsActive:    createdUser.IsActive,
		IsSuperuser: createdUser.IsSuperuser,
	})
}

// Update User
//
//	@Summary		Update User
//	@Description	Update User
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"User ID"
//	@Param			user	body		schemas.UserUpdateRequest	true	"Update User"
//	@Success		200		{object}	schemas.UserUpdateResponse
//	@Failure		400		{object}	schemas.BadRequestResponse
//	@Failure		404		{object}	schemas.NotFoundResponse
//	@Failure		500		{object}	schemas.InternalServerErrorResponse
//	@Security		OAuth2Password
//	@Router			/user/{id} [put]
func UpdateUserRoute(c *gin.Context) {
	// Authorize User
	_, err := core.GetUserFromAuthorizationHeader(models.DBConn, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, schemas.UnauthorizedResponse{
			Message: "Invalid/Expired token",
		})
		return
	}

	// get input user
	userId := c.Params.ByName("userId")
	if !core.IsValidUUID(userId) {
		c.JSON(http.StatusNotFound, schemas.NotFoundResponse{
			Message: "user not found",
		})
		return
	}
	jsonRequest := schemas.UserUpdateRequest{}
	err = c.BindJSON(&jsonRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, schemas.BadRequestResponse{
			Message: err.Error(),
		})
		return
	}

	// get existing user
	user, err := repository.GetUserById(models.DBConn, userId)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, schemas.NotFoundResponse{
				Message: "user not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// update user
	updatedUser, err := repository.UpdateUser(
		models.DBConn,
		user,
		jsonRequest.Email,
		jsonRequest.Username,
		jsonRequest.Password,
		jsonRequest.IsActive,
		jsonRequest.IsSuperuser,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, schemas.UserUpdateResponse{
		Id:          updatedUser.ID,
		Username:    updatedUser.Username,
		Email:       updatedUser.Email,
		IsActive:    updatedUser.IsActive,
		IsSuperuser: updatedUser.IsSuperuser,
	})
}

// Delete User
//
//	@Summary		Delete User
//	@Description	Delete user
//	@Tags			User
//	@Param			id	path	string	true	"User ID"
//	@Success		204
//	@Failure		404	{object}	schemas.NotFoundResponse
//	@Failure		500	{object}	schemas.InternalServerErrorResponse
//	@Security		OAuth2Password
//	@Router			/user/{id} [delete]
func DeleteUserRoute(c *gin.Context) {
	// Authorize User
	_, err := core.GetUserFromAuthorizationHeader(models.DBConn, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, schemas.UnauthorizedResponse{
			Message: "Invalid/Expired token",
		})
		return
	}

	// get input user
	userId := c.Params.ByName("userId")
	if !core.IsValidUUID(userId) {
		c.JSON(http.StatusNotFound, schemas.NotFoundResponse{
			Message: "user not found",
		})
		return
	}

	// get existing user
	user, err := repository.GetUserById(models.DBConn, userId)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, schemas.NotFoundResponse{
				Message: "user not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
		return
	}

	_, err = repository.DeleteUser(models.DBConn, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.InternalServerErrorResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
