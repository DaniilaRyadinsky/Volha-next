package domain

import productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

type Brand struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func BrandFromRpc(b *productsRPC.Brand) *Brand {
	if b == nil {
		return nil
	}
	return &Brand{
		Id:    b.GetId(),
		Title: b.GetTitle(),
	}
}

func BrandToRpc(b *Brand) *productsRPC.Brand {
	if b == nil {
		return nil
	}
	return &productsRPC.Brand{
		Id:    b.Id,
		Title: b.Title,
	}
}

func BrandsFromRpc(items []*productsRPC.Brand) []*Brand {
	res := make([]*Brand, 0, len(items))
	for _, b := range items {
		if val := BrandFromRpc(b); val != nil {
			res = append(res, val)
		}
	}
	return res
}

func BrandsToRpc(items []*Brand) []*productsRPC.Brand {
	res := make([]*productsRPC.Brand, 0, len(items))
	for _, b := range items {
		if val := BrandToRpc(b); val != nil {
			res = append(res, val)
		}
	}
	return res
}
