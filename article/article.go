package article

import (
	"github.com/xordspar0/squirrelbot/pocket"
)

type Article struct {
	Url   string
	Title string
}

func NewArticle(url string) *Article {
	a := &Article{
		Url: url,
	}

	return a
}

func (a *Article) Save(pocketKey, pocketUserToken string) error {
	pocketClient := pocket.NewClient(pocketKey, pocketUserToken)
	title, err := pocketClient.Add(a.Url)

	a.Title = title

	return err
}
