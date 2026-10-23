package api

import (
	"context"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CreateArticle(ctx context.Context, req *productsRPC.Article) (*emptypb.Empty, error) {
	const op = "grpc.CreateArticle"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.CreateArticle(ctx, domain.ArticleFromRpc(req))
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) UpdateArticle(ctx context.Context, req *productsRPC.Article) (*emptypb.Empty, error) {
	const op = "grpc.UpdateArticle"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.UpdateArticle(ctx, domain.ArticleFromRpc(req))
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) DeleteArticle(ctx context.Context, req *productsRPC.Id) (*emptypb.Empty, error) {
	const op = "grpc.DeleteArticle"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.DeleteArticle(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) GetAllArticles(ctx context.Context, req *productsRPC.Pagination) (*productsRPC.ArticleList, error) {
	const op = "grpc.GetAllArticles"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		items, total, err := s.s.GetAllArticles(ctx, int(req.GetOffset()), int(req.GetLimit()))
		return articlePage{items, total}, err
	})
	if err != nil {
		return nil, err
	}

	page := res.(articlePage)
	return &productsRPC.ArticleList{
		Items: domain.ArticlesToRpc(page.items), Total: int32(page.total), Offset: req.GetOffset(), Limit: req.GetLimit(),
	}, nil
}

func (s *ServerAPI) GetArticle(ctx context.Context, req *productsRPC.Id) (*productsRPC.Article, error) {
	const op = "grpc.GetArticle"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.GetArticle(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}

	return domain.ArticleToRpc(res.(*domain.Article)), nil
}

type articlePage struct {
	items []*domain.Article
	total int
}
