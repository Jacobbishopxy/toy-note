package controller

import (
	"errors"
	"strconv"
	"toy-note/api/entity"

	"github.com/gin-gonic/gin"
)

func getPaginationFromQuery(ctx *gin.Context) (entity.Pagination, error) {
	// get pagination's page from query string
	pageQuery, v := ctx.GetQuery("page")
	if !v {
		return entity.Pagination{}, errors.New("page query is required")
	}
	page, err := strconv.ParseInt(pageQuery, 10, 64)
	if err != nil {
		return entity.Pagination{}, err
	}

	// get pagination's size from query string
	sizeQuery, v := ctx.GetQuery("size")
	if !v {
		return entity.Pagination{}, errors.New("size query is required")
	}
	size, err := strconv.ParseInt(sizeQuery, 10, 64)
	if err != nil {
		return entity.Pagination{}, err
	}

	return entity.Pagination{
		Page: int(page),
		Size: int(size),
	}, nil

}
