package domain

import productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

type Country struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Friendly string `json:"friendly"`
}

func CountryFromRpc(c *productsRPC.Country) *Country {
	if c == nil {
		return nil
	}
	return &Country{
		Id:       c.GetId(),
		Title:    c.GetTitle(),
		Friendly: c.GetFriendly(),
	}
}

func CountryToRpc(c *Country) *productsRPC.Country {
	if c == nil {
		return nil
	}
	return &productsRPC.Country{
		Id:       c.Id,
		Title:    c.Title,
		Friendly: c.Friendly,
	}
}

func CountriesFromRpc(items []*productsRPC.Country) []*Country {
	res := make([]*Country, 0, len(items))
	for _, c := range items {
		if val := CountryFromRpc(c); val != nil {
			res = append(res, val)
		}
	}
	return res
}

func CountriesToRpc(items []*Country) []*productsRPC.Country {
	res := make([]*productsRPC.Country, 0, len(items))
	for _, c := range items {
		if val := CountryToRpc(c); val != nil {
			res = append(res, val)
		}
	}
	return res
}
