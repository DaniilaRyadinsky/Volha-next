package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func requestContext(c echo.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request().Context(), 3*time.Second)
}

func queryID(c echo.Context) (string, int, view.ErrorResponse) {
	id := c.QueryParam("id")
	if id == "" {
		return "", http.StatusBadRequest, view.ErrorResponse{Error: "missing id"}
	}
	return id, http.StatusOK, view.ErrorResponse{}
}

//INTERNAL IN EVERY
//GetAll ErrUnknownType 	Unimplemented
//Get ErrUnknownType ErrNotFound 	Unimplemented
//GetPhotosByProductColor ErrNotFound 	NotFound
//Create ErrInvalidType ErrUnknownType ErrAlreadyExists ErrForeignKey 	InvalidArgument Unimplemented AlreadyExists FailedPrecondition
//Update ErrInvalidType ErrUnknownType ErrNotFound 		InvalidArgument Unimplemented NotFound
//Delete ErrUnknownType ErrNotFound 	Unimplemented NotFound
//DeleteProductColorPhotos ErrNotFound 		NotFound
//GetDictionaries
//FilterProducts
//Products all the def but without ErrInvalidType ErrUnknownType

func (a *Apis) mapGRPCError(ctx context.Context, op string, err error) (int, view.ErrorResponse) {
	st, ok := status.FromError(err)
	if !ok {
		log.Error(ctx, "non-grpc product service error", "operation", op, "error", err)
		return http.StatusBadGateway, view.ErrorResponse{Error: "check logs on service"}
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return http.StatusBadRequest, view.ErrorResponse{Error: st.Message()}
	case codes.FailedPrecondition:
		return http.StatusBadRequest, view.ErrorResponse{Error: "input error"}
	case codes.AlreadyExists:
		return http.StatusConflict, view.ErrorResponse{Error: "already exists"}
	case codes.NotFound:
		return http.StatusNotFound, view.ErrorResponse{Error: "not found"}
	case codes.Unauthenticated:
		return http.StatusUnauthorized, view.ErrorResponse{Error: "unauthenticated"}
	case codes.PermissionDenied:
		return http.StatusForbidden, view.ErrorResponse{Error: "permission denied"}
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests, view.ErrorResponse{Error: "resource exhausted"}
	case codes.Unavailable:
		return http.StatusServiceUnavailable, view.ErrorResponse{Error: "service unavailable"}
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout, view.ErrorResponse{Error: "service timeout"}
	case codes.Unimplemented, codes.Internal, codes.Unknown:
		log.Error(ctx, "product service error", "operation", op, "grpc_code", st.Code().String(), "error", err)
		return http.StatusBadGateway, view.ErrorResponse{Error: "check logs on service"}
	default:
		log.Error(ctx, "unexpected product service error", "operation", op, "grpc_code", st.Code().String(), "error", err)
		return http.StatusBadGateway, view.ErrorResponse{Error: "check logs on service"}
	}
}

func (a *Apis) cleanDictionariesAfterMutation(ctx context.Context, op string) {
	if err := a.rds.CleanDictionaries(ctx); err != nil {
		log.Error(ctx, "clean dictionaries cache failed", "operation", op, "error", err)
	}
}
