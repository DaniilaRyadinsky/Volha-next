package domain

import productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

type Article struct {
	Id           string   `json:"id"`
	Slug         string   `json:"slug"`
	Title        string   `json:"title"`
	Img          []string `json:"img"`
	Text         string   `json:"text"`
	CreationTime int64    `json:"creationTime"`
}

func ArticleFromRpc(a *productsRPC.Article) *Article {
	if a == nil {
		return nil
	}
	return &Article{Id: a.GetId(), Slug: a.GetSlug(), Title: a.GetTitle(), Img: append([]string{}, a.GetImg()...), Text: a.GetText(), CreationTime: a.GetCreationTime()}
}
func ArticleToRpc(a *Article) *productsRPC.Article {
	if a == nil {
		return nil
	}
	return &productsRPC.Article{Id: a.Id, Slug: a.Slug, Title: a.Title, Img: append([]string{}, a.Img...), Text: a.Text, CreationTime: a.CreationTime}
}
func ArticlesFromRpc(items []*productsRPC.Article) []*Article {
	result := make([]*Article, 0, len(items))
	for _, a := range items {
		if a != nil {
			result = append(result, ArticleFromRpc(a))
		}
	}
	return result
}
func ArticlesToRpc(items []*Article) []*productsRPC.Article {
	result := make([]*productsRPC.Article, 0, len(items))
	for _, a := range items {
		if a != nil {
			result = append(result, ArticleToRpc(a))
		}
	}
	return result
}
