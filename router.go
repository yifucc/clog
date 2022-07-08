package main

import (
	"cc_blog/biz/handler"
	"cc_blog/config"
	"github.com/kataras/iris/v12"
)

// customizeRegister registers customize routers.
func customizedRegister(app *iris.Application) {
	html := iris.HTML("./web/views", ".html")
	html.Delims("{%", "%}")
	html.Layout("layout.html")
	app.RegisterView(html)
	app.HandleDir("/static", iris.Dir("./web/static"))
	app.HandleDir("/img", iris.Dir(config.Conf.ImgDir))
	app.Post("/api/webhook", handler.Webhook)
	root := app.Party("/", handler.Site, handler.Navbar, handler.Categories, handler.Slider)
	root.Get("/", handler.ArticleList)
	docs := root.Party("/docs")
	docs.Get("/{path}", handler.ArticleList)
	docs.Get("/{path}/{name}", handler.ArticleDetail)
}
