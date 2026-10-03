package domain

import productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

type Article struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	Img          string `json:"img"`
	Text         string `json:"text"`
	CreationTime int64  `json:"creation_time"`
}

func ArticleFromRpc(b *productsRPC.Article) *Article {
	if b == nil {
		return nil
	}
	return &Article{
		Id:           b.GetId(),
		Title:        b.GetTitle(),
		Img:          b.GetImg(),
		Text:         b.GetText(),
		CreationTime: b.GetCreationTime(),
	}
}

func ArticleToRpc(b *Article) *productsRPC.Article {
	if b == nil {
		return nil
	}
	return &productsRPC.Article{
		Id:           b.Id,
		Title:        b.Title,
		Img:          b.Img,
		Text:         b.Text,
		CreationTime: b.CreationTime,
	}
}

func ArticlesFromRpc(items []*productsRPC.Article) []*Article {
	res := make([]*Article, 0, len(items))
	for _, b := range items {
		if val := ArticleFromRpc(b); val != nil {
			res = append(res, val)
		}
	}
	return res
}

func ArticlesToRpc(items []*Article) []*productsRPC.Article {
	res := make([]*productsRPC.Article, 0, len(items))
	for _, b := range items {
		if val := ArticleToRpc(b); val != nil {
			res = append(res, val)
		}
	}
	return res
}
