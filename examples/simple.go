package main

import (
	"github.com/goforgery/compose"
	"github.com/goforgery/forgery2"
	"github.com/goforgery/mustache"
)

func main() {
	app := f.CreateApp()
	app.Engine(".html", mustache.Create())
	app.Get("/", func(req *f.Request, res *f.Response, next func()) {
		c := compose.Map{
			"header": func(req *f.Request, res *f.Response, next func()) {
				res.Send("Header string")
			},
			"body": func(req *f.Request, res *f.Response, next func()) {
				res.Render("body.html", "Title")
			},
			"footer": func(req *f.Request, res *f.Response, next func()) {
				res.End("Footer string")
			},
			"tail": func(req *f.Request, res *f.Response, next func()) {
				res.Write("Tail string")
			},
			"close": func(req *f.Request, res *f.Response, next func()) {
				res.WriteBytes([]byte("Close string"))
			},
		}
		data := c.Execute(req, res, next)
		res.Render("index.html", data)
	})
	app.Listen(3000)
}
