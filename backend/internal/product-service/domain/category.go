package domain

import productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

type Category struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Uri   string `json:"uri"`
	Img   string `json:"img"`
}

func CategoryFromRpc(c *productsRPC.Category) *Category {
	if c == nil {
		return nil
	}
	return &Category{
		Id:    c.GetId(),
		Title: c.GetTitle(),
		Uri:   c.GetUri(),
		Img:   c.GetImg(),
	}
}

func CategoryToRpc(c *Category) *productsRPC.Category {
	if c == nil {
		return nil
	}
	return &productsRPC.Category{
		Id:    c.Id,
		Title: c.Title,
		Uri:   c.Uri,
		Img:   c.Img,
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
