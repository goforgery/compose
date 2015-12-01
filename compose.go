package compose

import (
	"bytes"
	"github.com/goforgery/forgery2"
	"net/http"
)

// The Writer used in place of stackr.Writer for buffering.
type BufferedResponseWriter struct {
	Id      string
	Headers http.Header
	Buffer  *bytes.Buffer
	Status  int
}

// Header returns the header map that would have been sent by WriteHeader.
func (this *BufferedResponseWriter) Header() http.Header {
	if this.Headers == nil {
		this.Headers = http.Header{}
	}
	return this.Headers
}

// Write writes the data to the buffer.
func (this *BufferedResponseWriter) Write(b []byte) (int, error) {
	if this.Buffer == nil {
		this.Buffer = &bytes.Buffer{}
	}
	len, err := this.Buffer.Write(b)
	return len, err
}

// WriteHeader buffers the HTTP status code.
func (this *BufferedResponseWriter) WriteHeader(code int) {
	this.Status = code
}

// c := compose.Map{
//     "header": func(req *f.Request, res *f.Response, next func()) {
//         res.Send("Header string")
//     },
//     "body": func(req *f.Request, res *f.Response, next func()) {
//         res.Render("body.html", "Title")
//     },
//     "footer": func(req *f.Request, res *f.Response, next func()) {
//         res.End("Footer string")
//     },
//     "tail": func(req *f.Request, res *f.Response, next func()) {
//         res.Write("Tail string")
//     },
//     "close": func(req *f.Request, res *f.Response, next func()) {
//         res.WriteBytes([]byte("Close string"))
//     },
// }
// data := c.Execute(req, res, next)
type Map map[string]func(*f.Request, *f.Response, func())

// Execute all functions in the given map.
func (this Map) Execute(req *f.Request, res *f.Response, next func()) map[string]string {
	// Grab the res.Writer so we can put it back later.
	w := res.Writer
	// Make chan.
	c := make(chan *BufferedResponseWriter)
	// Loop over the items in the map and dispatch each one.
	for id, fn := range this {
		go this.dispatch(req, res.Clone(), next, id, fn, c)
	}
	// Put the res.Writer back.
	res.Writer = w
	// Create the return map.
	renders := map[string]string{}
	for i := 0; i < len(this); i++ {
		buf := <-c
		for k, v := range buf.Headers {
			if k != "Content-Length" {
				res.Set(k, v[0])
			}
		}
		if buf.Buffer != nil {
			renders[buf.Id] = buf.Buffer.String()
		}
	}
	return renders
}

// Dispatch the call to handle the given function.
func (this Map) dispatch(req *f.Request, res *f.Response, next func(), id string, fn func(*f.Request, *f.Response, func()), c chan *BufferedResponseWriter) {
	// Create a buffer.
	buffer := &BufferedResponseWriter{Id: id}
	// Replace res.Writer with BufferedResponseWriter so all the output can be captured.
	res.Writer = buffer
	// Call the function.
	fn(req, res, next)
	// Return the buffer.
	c <- buffer
}
