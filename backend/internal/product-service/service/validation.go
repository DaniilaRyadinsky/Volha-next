package service

import (
	"errors"
	"strings"

	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/rs/xid"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

func validateID(id string) error {
	if id == "" {
		return errors.New("id is empty")
	}
	if _, err := xid.FromString(id); err != nil {
		return errors.New("id not in xid")
	}
	return nil
}

func validateBrand(b *domain.Brand) error {
	if b == nil {
		return errors.New("object is nil")
	}
	if err := validateID(b.Id); err != nil {
		return err
	}
	if b.Title == "" {
		return errors.New("brand title is empty")
	}
	return nil
}

func validateCategory(c *domain.Category) error {
	if c == nil {
		return errors.New("object is nil")
	}
	if err := validateID(c.Id); err != nil {
		return err
	}
	if c.Title == "" {
		return errors.New("category title is empty")
	}
	if c.Uri == "" {
		return errors.New("category uri is empty")
	}
	return nil
}

func validateCountry(c *domain.Country) error {
	if c == nil {
		return errors.New("object is nil")
	}
	if err := validateID(c.Id); err != nil {
		return err
	}
	if c.Title == "" {
		return errors.New("country title is empty")
	}
	return nil
}

func validateMaterial(m *domain.Material) error {
	if m == nil {
		return errors.New("object is nil")
	}
	if err := validateID(m.Id); err != nil {
		return err
	}
	if m.Title == "" {
		return errors.New("material title is empty")
	}
	return nil
}

func validateColor(clr *domain.Color) error {
	if clr == nil {
		return errors.New("object is nil")
	}
	if err := validateID(clr.Id); err != nil {
		return err
	}
	if clr.Title == "" {
		return errors.New("color title is empty")
	}
	if clr.Hex == "" {
		return errors.New("color hex is empty")
	}
	return nil
}

func validateSlide(sl *domain.Slide) error {
	if sl == nil {
		return errors.New("object is nil")
	}
	if err := validateID(sl.Id); err != nil {
		return err
	}
	if sl.Link == "" {
		return errors.New("slide link is empty")
	}
	if sl.Img == "" {
		return errors.New("slide img is empty")
	}
	if sl.Img762 == "" {
		return errors.New("slide img762 is empty")
	}
	return nil
}

func validateArticle(a *domain.Article) error {
	return validateArticleFields(a)
}

func validateArticleUpdate(a *domain.Article) error {
	return validateArticleFields(a)
}

func validateArticleFields(a *domain.Article) error {
	if a == nil {
		return errors.New("object is nil")
	}
	if err := validateID(a.Id); err != nil {
		return err
	}
	if a.Title == "" {
		return errors.New("article title is empty")
	}
	if a.Img == "" {
		return errors.New("article img is empty")
	}

	return nil
}

func validatePCPIDs(productID, colorID string) error {
	if err := validateID(productID); err != nil {
		return err
	}
	if err := validateID(colorID); err != nil {
		return err
	}
	return nil
}

func validatePCP(pcp *domain.ProductColorPhotos) error {
	if pcp == nil {
		return errors.New("object is nil")
	}
	if err := validatePCPIDs(pcp.ProductId, pcp.ColorId); err != nil {
		return err
	}
	if len(pcp.Photos) == 0 {
		return errors.New("photos are empty")
	}
	return nil
}

func validateProductPayload(p *domain.ProductId) error {
	if p == nil {
		return errors.New("object is nil")
	}
	if err := validateID(p.Id); err != nil {
		return err
	}
	if p.Title == "" {
		return errors.New("product title is empty")
	}
	if p.Article == "" {
		return errors.New("product article is empty")
	}
	if !isEightDigitArticle(p.Article) {
		return errors.New("product article must be exactly 8 digits")
	}
	if err := validateID(p.Brand); err != nil {
		return errors.New("brand id is invalid")
	}
	if err := validateID(p.Category); err != nil {
		return errors.New("category id is invalid")
	}
	if err := validateID(p.Country); err != nil {
		return errors.New("country id is invalid")
	}
	if p.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if p.Width <= 0 || p.Height <= 0 || p.Depth <= 0 {
		return errors.New("dimensions must be greater than zero")
	}
	if len(p.Photos) == 0 {
		return errors.New("photos are empty")
	}
	if err := validateIDs("material id", p.Materials); err != nil {
		return err
	}
	if err := validateIDs("color id", p.Colors); err != nil {
		return err
	}
	if err := validateIDs("seem id", p.Seems); err != nil {
		return err
	}
	for _, s := range p.Seems {
		if s == p.Id {
			return errors.New("product cannot be similar to itself")
		}
	}

	return nil
}

func validateProductRange(offset, limit int) error {
	if offset < 0 || limit <= 0 {
		return errors.New("pagination boundaries are invalid")
	}
	if limit > MaxLimit {
		return errors.New("pagination limit is too large")
	}
	return nil
}

func validateFilter(filter *domain.ProductFilter) error {
	if filter == nil {
		return errors.New("filter is nil")
	}
	if filter.Limit <= 0 || filter.Offset < 0 {
		return errors.New("pagination boundaries are invalid")
	}
	if filter.Limit > MaxLimit {
		return errors.New("pagination limit is too large")
	}
	if err := validateMinMax(filter.MinWidth, filter.MaxWidth, "width"); err != nil {
		return err
	}
	if err := validateMinMax(filter.MinHeight, filter.MaxHeight, "height"); err != nil {
		return err
	}
	if err := validateMinMax(filter.MinDepth, filter.MaxDepth, "depth"); err != nil {
		return err
	}
	if err := validateMinMax(filter.MinPrice, filter.MaxPrice, "price"); err != nil {
		return err
	}
	if filter.SortBy != "" && !validSortBy(filter.SortBy) {
		return errors.New("sort_by is invalid")
	}
	if filter.SortOrder != "" && !validSortOrder(filter.SortOrder) {
		return errors.New("sort_order is invalid")
	}
	if err := validateIDs("brand id", filter.Brand); err != nil {
		return err
	}
	if err := validateIDs("category id", filter.Category); err != nil {
		return err
	}
	if err := validateIDs("country id", filter.Country); err != nil {
		return err
	}
	if err := validateIDs("material id", filter.Materials); err != nil {
		return err
	}
	if err := validateIDs("color id", filter.Colors); err != nil {
		return err
	}
	return nil
}

func isEightDigitArticle(article string) bool {
	if len(article) != 8 {
		return false
	}
	for _, r := range article {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func validateIDs(name string, ids []string) error {
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if err := validateID(id); err != nil {
			return errors.New(name + " is invalid")
		}
		if _, ok := seen[id]; ok {
			return errors.New(name + " is duplicated")
		}
		seen[id] = struct{}{}
	}
	return nil
}

func validateMinMax(min, max int32, name string) error {
	if min > 0 && max > 0 && min > max {
		return errors.New(name + " min is greater than max")
	}
	return nil
}

func validSortBy(sortBy string) bool {
	switch sortBy {
	case "price", "width", "height", "depth", "title":
		return true
	default:
		return false
	}
}

func validSortOrder(sortOrder string) bool {
	switch strings.ToLower(sortOrder) {
	case "asc", "desc":
		return true
	default:
		return false
	}
}
