package compose

import (
	"bytes"
	"github.com/goforgery/forgery2"
	"net/http"
)

// The Writer used in place of stackr.Writer for buffering.
type BufferedResponseWriter struct {
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
	// fmt.Println(string(this.Buffer.Bytes()))
	return len, err
}

// WriteHeader buffers the HTTP status code.
func (this *BufferedResponseWriter) WriteHeader(code int) {
	this.Status = code
}

/*
   c := compose.Map{
       "header": func(req, res, next) {
           res.Send("Header string")
       },
       "body": func(req, res, next) {
           res.Render("page.html", "Body string")
       },
       "footer": func(req, res, next) {
           res.End("Footer string")
       },
       "tail": func(req, res, next) {
           res.Write("Tail string")
       },
       "close": func(req, res, next) {
           res.WriteBytes([]byte("Close string"))
       },
   }

   data := c.Execute(req, res, next)
*/
type Map map[string]func(*f.Request, *f.Response, func())

/*
   The worker.
*/
func (this Map) Execute(req *f.Request, res *f.Response, next func()) map[string]string {

	renders := map[string]string{}

	// Grab the res.Writer so we can put it back later.
	w := res.Writer

	// Loop over the items in the map.
	for id, fn := range this {
		// Execute function.
		func(mapId string, mapFn func(*f.Request, *f.Response, func())) {
			// Create a buffer.
			buffer := &BufferedResponseWriter{}
			// Replace res.Writer with BufferedResponseWriter so all the output can be captured.
			res.Writer = buffer
			// Call the function.
			mapFn(req, res, next)
			// Add the buffered data to the renders map.
			if buffer.Buffer != nil {
				renders[mapId] = buffer.Buffer.String()
			}
		}(id, fn)
	}

	// Put the res.Writer back.
	res.Writer = w

	return renders
}
