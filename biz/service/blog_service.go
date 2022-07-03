package service

import "cc_blog/biz/model"

type BlogService interface {
	GetArticles(path string) (*model.Articles, error)
	GetCategories() (*model.Categories, error)
	GetNavbar() (*model.Navbar, error)
	GetSlider() (*model.Slider, error)
	GetArticleDetail(name string) (*model.Article, error)
}
