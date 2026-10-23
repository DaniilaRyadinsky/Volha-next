package domain

import productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

type Slide struct {
	Id          string `json:"id"`
	Link        string `json:"link"`
	Img         string `json:"img"`
	Img762      string `json:"img762"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Badge       string `json:"subtitle"`
}

func SlideFromRpc(s *productsRPC.Slide) *Slide {
	if s == nil {
		return nil
	}
	return &Slide{
		Id:          s.GetId(),
		Link:        s.GetLink(),
		Img:         s.GetImg(),
		Img762:      s.GetImg762(),
		Title:       s.GetTitle(),
		Description: s.GetDescription(),
		Badge:       s.GetBadge(),
	}
}

func SlideToRpc(s *Slide) *productsRPC.Slide {
	if s == nil {
		return nil
	}
	return &productsRPC.Slide{
		Id:          s.Id,
		Link:        s.Link,
		Img:         s.Img,
		Img762:      s.Img762,
		Title:       s.Title,
		Description: s.Description,
		Badge:       s.Badge,
	}
}

func SlidesFromRpc(items []*productsRPC.Slide) []*Slide {
	res := make([]*Slide, 0, len(items))
	for _, s := range items {
		if val := SlideFromRpc(s); val != nil {
			res = append(res, val)
		}
	}
	return res
}

func SlidesToRpc(items []*Slide) []*productsRPC.Slide {
	res := make([]*productsRPC.Slide, 0, len(items))
	for _, s := range items {
		if val := SlideToRpc(s); val != nil {
			res = append(res, val)
		}
	}
	return res
}
