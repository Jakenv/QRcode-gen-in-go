package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/skip2/go-qrcode"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, date interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, date)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("view/*.html")),
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})
	e.POST("/generate", func(c echo.Context) error {
		text := c.FormValue("text")
		code, err := qrcode.Encode(text, qrcode.Medium, 256)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to generate image")
		}
		qrCodeBase64 := base64.StdEncoding.EncodeToString(code)
		imgTag := fmt.Sprintf(`<img src="data:image/png;base64,%s" alt="QR Code">`, qrCodeBase64)
		return c.HTML(http.StatusOK, imgTag)
	})

	e.Logger.Fatal(e.Start(":5001"))
}
