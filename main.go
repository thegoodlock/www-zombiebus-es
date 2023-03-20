package main

import (
	"embed"
	"io/fs"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
)

type Config struct {
	Addr string `envconfig:"PORT"`
	Live bool   `envconfig:"LIVE"`
}

//go:embed static
var embededFiles embed.FS

func getFileSystem(useOS bool) http.FileSystem {
	if useOS {
		return http.FS(os.DirFS("static"))
	}

	fsys, err := fs.Sub(embededFiles, "static")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

func main() {
	var cfg Config
	envconfig.MustProcess("web", &cfg)

	e := echo.New()
	e.GET("/*", echo.WrapHandler(http.FileServer(getFileSystem(cfg.Live))))
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":" + cfg.Addr))
}
