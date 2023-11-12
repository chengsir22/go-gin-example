package article_service

import "github.com/luciferCN22/go-gin-example/models"

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

	PageNum  int
	PageSize int
}

func (a *Article) ExistByID() bool {
	return models.ExistArticleByID(a.ID)
}
