package middle

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var DefaultSwagger = &Swagger{
	Path:     "/swagger/swagger.json",
	Filename: "./docs/swagger.json",
}

type Swagger struct {
	Path     string //json路由 例/swagger/swagger.json
	Filename string //json文件名称
}

func (this *Swagger) UI(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf(swaggerUI, this.Path)))
	w.WriteHeader(200)
}

func (this *Swagger) Json(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(this.Filename)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	defer f.Close()
	w.WriteHeader(200)
	io.Copy(w, f)
}

var (
	swaggerUI = `<!DOCTYPE html>
        <html>
          <head>
            <title>SwaggerUI</title>
            <!-- needed for adaptive design -->
            <meta charset="utf-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1">
            <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
            <style>
              body {
                margin: 0;
                padding: 0;
              }
            </style>
          </head>
          <body>
            <redoc spec-url='%s'></redoc>
            <script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"> </script>
          </body>
        </html>`
)
