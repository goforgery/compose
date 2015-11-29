package compose

import (
	"bytes"
	"github.com/goforgery/forgery2"
	. "github.com/ricallinson/simplebdd"
	"testing"
)

func TestCompose(t *testing.T) {

	var app *f.Application
	var req *f.Request
	var res *f.Response
	var buf *bytes.Buffer
	var next func()

	BeforeEach(func() {
		app, req, res, buf = f.CreateAppMock()
		next = func() {}
	})

	Describe("Execute()", func() {

		It("should return a map with all the functions executed", func() {

			c := Map{
				"header": func(req *f.Request, res *f.Response, next func()) {
					res.Send("Header")
				},
				"empty": func(req *f.Request, res *f.Response, next func()) {
					res.Locals["append"] = "Bar"
					res.End("")
				},
				"body": func(req *f.Request, res *f.Response, next func()) {
					res.Locals["title"] = "Foo"
					res.End("Body")
				},
				"footer": func(req *f.Request, res *f.Response, next func()) {
					res.Locals["append"] = "Bar"
					res.End("Footer")
				},
				"tail": func(req *f.Request, res *f.Response, next func()) {
					res.Write("Tail")
				},
				"close": func(req *f.Request, res *f.Response, next func()) {
					res.WriteBytes([]byte("Close"))
				},
			}

			data := c.Execute(req, res, next)

			AssertEqual(string(data["header"]), "Header")
			AssertEqual(string(data["empty"]), "")
			AssertEqual(string(data["body"]), "Body")
			AssertEqual(res.Locals["title"], "Foo")
			AssertEqual(string(data["footer"]), "Footer")
			AssertEqual(string(data["tail"]), "Tail")
			AssertEqual(string(data["close"]), "Close")
			AssertEqual(res.Locals["append"], "Bar")
		})
	})

	Report(t)
}
