package domain

import productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

type Category struct {
	Id          string `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Count       int32  `json:"count"`
	Img         string `json:"img"`
}

func CategoryFromRpc(c *productsRPC.Category) *Category {
	if c == nil {
		return nil
	}
	return &Category{
		Id:          c.GetId(),
		Slug:        c.GetSlug(),
		Title:       c.GetTitle(),
		Description: c.GetDescription(),
		Count:       c.GetCount(),
		Img:         c.GetImg(),
	}
}

func CategoryToRpc(c *Category) *productsRPC.Category {
	if c == nil {
		return nil
	}
	return &productsRPC.Category{
		Id:          c.Id,
		Slug:        c.Slug,
		Title:       c.Title,
		Description: c.Description,
		Count:       c.Count,
		Img:         c.Img,
	}
}

func CategoriesFromRpc(items []*productsRPC.Category) []*Category {
	res := make([]*Category, 0, len(items))
	for _, c := range items {
		if val := CategoryFromRpc(c); val != nil {
			res = append(res, val)
		}
	}
	return res
}

func CategoriesToRpc(items []*Category) []*productsRPC.Category {
	res := make([]*productsRPC.Category, 0, len(items))
	for _, c := range items {
		if val := CategoryToRpc(c); val != nil {
			res = append(res, val)
		}
	}
	return res
}
