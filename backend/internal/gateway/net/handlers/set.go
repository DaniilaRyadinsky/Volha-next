package handlers

import (
	"net/http"
	"strconv"

	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	"github.com/labstack/echo/v4"
)

const notByCategory = "CLR"

const successPag = 0
const maxPaginationLimit = 100

func pag(c echo.Context) (int, int, int, *view.ErrorResponse) {
	s := c.QueryParam("offset")
	if s == "" {
		return 0, 0, http.StatusBadRequest, &view.ErrorResponse{Error: "bad offset"}
	}
	lim := c.QueryParam("limit")
	if lim == "" {
		return 0, 0, http.StatusBadRequest, &view.ErrorResponse{Error: "bad limit"}
	}

	offset, err := strconv.Atoi(s)
	if err != nil {
		return 0, 0, http.StatusBadRequest, &view.ErrorResponse{Error: "offset must be int"}
	}
	limit, err := strconv.Atoi(lim)
	if err != nil {
		return 0, 0, http.StatusBadRequest, &view.ErrorResponse{Error: "limit must be int"}
	}

	if offset < 0 || limit <= 0 {
		return 0, 0, http.StatusBadRequest, &view.ErrorResponse{Error: "bad pagination"}
	}
	if limit > maxPaginationLimit {
		return 0, 0, http.StatusBadRequest, &view.ErrorResponse{Error: "limit is too large"}
	}
	return offset, limit, successPag, nil
}
