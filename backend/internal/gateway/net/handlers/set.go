package handlers

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	"github.com/labstack/echo/v4"
)

const notByCategory = "CLR"
const successPag = 0
const maxPaginationLimit = 100

func pageOffset(page, limit int) (int, int, int, error) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 20
	}
	if page < 1 || limit < 1 || limit > maxPaginationLimit || page-1 > math.MaxInt32/limit {
		return 0, 0, 0, errors.New("invalid page or limit")
	}
	return (page - 1) * limit, limit, page, nil
}
func pag(c echo.Context) (int, int, int, *view.ErrorResponse) {
	page, limit := 0, 0
	var err error
	if value := c.QueryParam("page"); value != "" {
		page, err = strconv.Atoi(value)
		if err != nil {
			return 0, 0, http.StatusBadRequest, &view.ErrorResponse{Error: "page must be int"}
		}
	}
	if value := c.QueryParam("limit"); value != "" {
		limit, err = strconv.Atoi(value)
		if err != nil {
			return 0, 0, http.StatusBadRequest, &view.ErrorResponse{Error: "limit must be int"}
		}
	}
	offset, limit, _, err := pageOffset(page, limit)
	if err != nil {
		return 0, 0, http.StatusBadRequest, &view.ErrorResponse{Error: err.Error()}
	}
	// Accept existing offset callers while page is the public pagination contract.
	if value := c.QueryParam("offset"); value != "" && c.QueryParam("page") == "" {
		offset, err = strconv.Atoi(value)
		if err != nil || offset < 0 || offset > math.MaxInt32 {
			return 0, 0, http.StatusBadRequest, &view.ErrorResponse{Error: "bad offset"}
		}
	}
	return offset, limit, successPag, nil
}
