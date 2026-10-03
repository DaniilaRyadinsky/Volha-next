package domain

import productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

type ProductColorPhotos struct {
	ProductId string   `json:"product_id"`
	ColorId   string   `json:"color_id"`
	Photos    []string `json:"photos"`
}

func ProductColorPhotosFromRpc(pcp *productsRPC.ProductColorPhotos) *ProductColorPhotos {
	if pcp == nil {
		return nil
	}
	return &ProductColorPhotos{
		ProductId: pcp.GetProductId(),
		ColorId:   pcp.GetColorId(),
		Photos:    append([]string{}, pcp.GetPhotos()...),
	}
}

func ProductColorPhotosToRpc(pcp *ProductColorPhotos) *productsRPC.ProductColorPhotos {
	if pcp == nil {
		return nil
	}
	return &productsRPC.ProductColorPhotos{
		ProductId: pcp.ProductId,
		ColorId:   pcp.ColorId,
		Photos:    append([]string{}, pcp.Photos...),
	}
}

func ProductColorPhotosListToRpc(items []ProductColorPhotos) []*productsRPC.ProductColorPhotos {
	res := make([]*productsRPC.ProductColorPhotos, 0, len(items))
	for i := range items {
		p := items[i]
		res = append(res, ProductColorPhotosToRpc(&p))
	}
	return res
}
