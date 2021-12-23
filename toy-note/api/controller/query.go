package controller

import (
	"errors"
	"strconv"
	"time"
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

func getTimeSearchFromQuery(ctx *gin.Context) (entity.TimeSearch, error) {
	// get start time from query string
	startQuery, v := ctx.GetQuery("start")
	if !v {
		return entity.TimeSearch{}, errors.New("start query is required")
	}
	start, err := time.Parse(time.RFC3339, startQuery)
	if err != nil {
		return entity.TimeSearch{}, err
	}

	// get end time from query string
	endQuery, v := ctx.GetQuery("end")
	if !v {
		return entity.TimeSearch{}, errors.New("end query is required")
	}
	end, err := time.Parse(time.RFC3339, endQuery)
	if err != nil {
		return entity.TimeSearch{}, err
	}

	// get time type from query string
	timeTypeQuery := ctx.Query("type")
	tt, err := strconv.ParseInt(timeTypeQuery, 10, 64)
	if err != nil {
		return entity.TimeSearch{}, err
	}

	return entity.TimeSearch{
		Start: start,
		End:   end,
		Type:  entity.TimeType(tt),
	}, nil
}
