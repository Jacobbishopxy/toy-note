package controller

import (
	"net/http"
	"strconv"
	"toy-note/api/entity"

	"github.com/gin-gonic/gin"
)

func getPaginationFromQuery(ctx *gin.Context) (entity.Pagination, error) {
	// get pagination's page from query string
	pageQuery := ctx.Query("page")
	page, err := strconv.ParseInt(pageQuery, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return entity.Pagination{}, err
	}

	// get pagination's size from query string
	sizeQuery := ctx.Query("size")
	size, err := strconv.ParseInt(sizeQuery, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return entity.Pagination{}, err
	}

	return entity.Pagination{
		Page: int(page),
		Size: int(size),
	}, nil

}
