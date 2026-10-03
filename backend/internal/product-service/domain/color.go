package domain

import productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

type Color struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Hex   string `json:"hex"`
}

func ColorFromRpc(c *productsRPC.Color) *Color {
	if c == nil {
		return nil
	}
	return &Color{
		Id:    c.GetId(),
		Title: c.GetTitle(),
		Hex:   c.GetHex(),
	}
}

func ColorToRpc(c *Color) *productsRPC.Color {
	if c == nil {
		return nil
	}
	return &productsRPC.Color{
		Id:    c.Id,
		Title: c.Title,
		Hex:   c.Hex,
	}
}

func ColorsFromRpc(items []*productsRPC.Color) []*Color {
	res := make([]*Color, 0, len(items))
	for _, c := range items {
		if val := ColorFromRpc(c); val != nil {
			res = append(res, val)
		}
	}
	return res
}

func ColorsToRpc(items []*Color) []*productsRPC.Color {
	res := make([]*productsRPC.Color, 0, len(items))
	for _, c := range items {
		if val := ColorToRpc(c); val != nil {
			res = append(res, val)
		}
	}
	return res
}
