package domain

import productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

type Material struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func MaterialFromRpc(m *productsRPC.Material) *Material {
	if m == nil {
		return nil
	}
	return &Material{
		Id:    m.GetId(),
		Title: m.GetTitle(),
	}
}

func MaterialToRpc(m *Material) *productsRPC.Material {
	if m == nil {
		return nil
	}
	return &productsRPC.Material{
		Id:    m.Id,
		Title: m.Title,
	}
}

func MaterialsFromRpc(items []*productsRPC.Material) []*Material {
	res := make([]*Material, 0, len(items))
	for _, m := range items {
		if val := MaterialFromRpc(m); val != nil {
			res = append(res, val)
		}
	}
	return res
}

func MaterialsToRpc(items []*Material) []*productsRPC.Material {
	res := make([]*productsRPC.Material, 0, len(items))
	for _, m := range items {
		if val := MaterialToRpc(m); val != nil {
			res = append(res, val)
		}
	}
	return res
}
