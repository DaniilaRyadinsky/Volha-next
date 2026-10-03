package service

import (
	"context"
	"time"

	"github.com/autumnterror/volha-backend/internal/product-service/domain"
)

func (s *ProductsService) GetAllBrands(ctx context.Context) ([]*domain.Brand, error) {
	repo, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck("service.GetAllBrands", err)
	}
	return repo.GetAllBrands(ctx)
}

func (s *ProductsService) GetBrand(ctx context.Context, id string) (*domain.Brand, error) {
	const op = "service.GetBrand"
	if err := validateID(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	repo, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	return repo.GetBrand(ctx, id)
}

func (s *ProductsService) CreateBrand(ctx context.Context, b *domain.Brand) error {
	const op = "service.CreateBrand"
	if err := validateBrand(b); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.CreateBrand(ctx, b)
	})
}

func (s *ProductsService) UpdateBrand(ctx context.Context, b *domain.Brand) error {
	const op = "service.UpdateBrand"
	if err := validateBrand(b); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.UpdateBrand(ctx, b)
	})
}

func (s *ProductsService) DeleteBrand(ctx context.Context, id string) error {
	const op = "service.DeleteBrand"
	if err := validateID(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.DeleteBrand(ctx, id)
	})
}

func (s *ProductsService) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	repo, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck("service.GetAllCategories", err)
	}
	return repo.GetAllCategories(ctx)
}

func (s *ProductsService) GetCategory(ctx context.Context, id string) (*domain.Category, error) {
	const op = "service.GetCategory"
	if err := validateID(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	repo, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	return repo.GetCategory(ctx, id)
}

func (s *ProductsService) CreateCategory(ctx context.Context, c *domain.Category) error {
	const op = "service.CreateCategory"
	if err := validateCategory(c); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.CreateCategory(ctx, c)
	})
}

func (s *ProductsService) UpdateCategory(ctx context.Context, c *domain.Category) error {
	const op = "service.UpdateCategory"
	if err := validateCategory(c); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.UpdateCategory(ctx, c)
	})
}

func (s *ProductsService) DeleteCategory(ctx context.Context, id string) error {
	const op = "service.DeleteCategory"
	if err := validateID(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.DeleteCategory(ctx, id)
	})
}

func (s *ProductsService) GetAllCountries(ctx context.Context) ([]*domain.Country, error) {
	repo, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck("service.GetAllCountries", err)
	}
	return repo.GetAllCountries(ctx)
}

func (s *ProductsService) GetCountry(ctx context.Context, id string) (*domain.Country, error) {
	const op = "service.GetCountry"
	if err := validateID(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	repo, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	return repo.GetCountry(ctx, id)
}

func (s *ProductsService) CreateCountry(ctx context.Context, c *domain.Country) error {
	const op = "service.CreateCountry"
	if err := validateCountry(c); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.CreateCountry(ctx, c)
	})
}

func (s *ProductsService) UpdateCountry(ctx context.Context, c *domain.Country) error {
	const op = "service.UpdateCountry"
	if err := validateCountry(c); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.UpdateCountry(ctx, c)
	})
}

func (s *ProductsService) DeleteCountry(ctx context.Context, id string) error {
	const op = "service.DeleteCountry"
	if err := validateID(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.DeleteCountry(ctx, id)
	})
}

func (s *ProductsService) GetAllMaterials(ctx context.Context) ([]*domain.Material, error) {
	repo, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck("service.GetAllMaterials", err)
	}
	return repo.GetAllMaterials(ctx)
}

func (s *ProductsService) GetMaterial(ctx context.Context, id string) (*domain.Material, error) {
	const op = "service.GetMaterial"
	if err := validateID(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	repo, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	return repo.GetMaterial(ctx, id)
}

func (s *ProductsService) CreateMaterial(ctx context.Context, m *domain.Material) error {
	const op = "service.CreateMaterial"
	if err := validateMaterial(m); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.CreateMaterial(ctx, m)
	})
}

func (s *ProductsService) UpdateMaterial(ctx context.Context, m *domain.Material) error {
	const op = "service.UpdateMaterial"
	if err := validateMaterial(m); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.UpdateMaterial(ctx, m)
	})
}

func (s *ProductsService) DeleteMaterial(ctx context.Context, id string) error {
	const op = "service.DeleteMaterial"
	if err := validateID(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.DeleteMaterial(ctx, id)
	})
}

func (s *ProductsService) GetAllColors(ctx context.Context) ([]*domain.Color, error) {
	repo, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck("service.GetAllColors", err)
	}
	return repo.GetAllColors(ctx)
}

func (s *ProductsService) GetColor(ctx context.Context, id string) (*domain.Color, error) {
	const op = "service.GetColor"
	if err := validateID(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	repo, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	return repo.GetColor(ctx, id)
}

func (s *ProductsService) CreateColor(ctx context.Context, c *domain.Color) error {
	const op = "service.CreateColor"
	if err := validateColor(c); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.CreateColor(ctx, c)
	})
}

func (s *ProductsService) UpdateColor(ctx context.Context, c *domain.Color) error {
	const op = "service.UpdateColor"
	if err := validateColor(c); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.UpdateColor(ctx, c)
	})
}

func (s *ProductsService) DeleteColor(ctx context.Context, id string) error {
	const op = "service.DeleteColor"
	if err := validateID(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.DeleteColor(ctx, id)
	})
}

func (s *ProductsService) GetAllSlides(ctx context.Context) ([]*domain.Slide, error) {
	repo, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck("service.GetAllSlides", err)
	}
	return repo.GetAllSlides(ctx)
}

func (s *ProductsService) GetSlide(ctx context.Context, id string) (*domain.Slide, error) {
	const op = "service.GetSlide"
	if err := validateID(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	repo, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	return repo.GetSlide(ctx, id)
}

func (s *ProductsService) CreateSlide(ctx context.Context, sl *domain.Slide) error {
	const op = "service.CreateSlide"
	if err := validateSlide(sl); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.CreateSlide(ctx, sl)
	})
}

func (s *ProductsService) UpdateSlide(ctx context.Context, sl *domain.Slide) error {
	const op = "service.UpdateSlide"
	if err := validateSlide(sl); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.UpdateSlide(ctx, sl)
	})
}

func (s *ProductsService) DeleteSlide(ctx context.Context, id string) error {
	const op = "service.DeleteSlide"
	if err := validateID(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.DeleteSlide(ctx, id)
	})
}

func (s *ProductsService) GetAllArticles(ctx context.Context) ([]*domain.Article, error) {
	repo, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck("service.GetAllArticles", err)
	}
	return repo.GetAllArticles(ctx)
}

func (s *ProductsService) GetArticle(ctx context.Context, id string) (*domain.Article, error) {
	const op = "service.GetArticle"
	if err := validateID(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	repo, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	return repo.GetArticle(ctx, id)
}

func (s *ProductsService) CreateArticle(ctx context.Context, a *domain.Article) error {
	const op = "service.CreateArticle"
	if a != nil {
		cp := *a
		cp.CreationTime = time.Now().UTC().Unix()
		a = &cp
	}
	if err := validateArticle(a); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.CreateArticle(ctx, a)
	})
}

func (s *ProductsService) UpdateArticle(ctx context.Context, a *domain.Article) error {
	const op = "service.UpdateArticle"
	if err := validateArticleUpdate(a); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.UpdateArticle(ctx, a)
	})
}

func (s *ProductsService) DeleteArticle(ctx context.Context, id string) error {
	const op = "service.DeleteArticle"
	if err := validateID(id); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.catalogRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.DeleteArticle(ctx, id)
	})
}
