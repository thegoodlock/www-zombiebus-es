package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	gomail "gopkg.in/mail.v2"
)

type Config struct {
	Addr     string `envconfig:"PORT"`
	Live     bool   `envconfig:"LIVE" default:"true"`
	SMTPPass string `envconfig:"SMTP_USER" default:"b8ea4a5d52519a11b01edad1f7cc5854"`
	SMTPUser string `envconfig:"SMTP_USER" default:"f06c3e95b93b1c291f38435d7b127816"`
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
	if cfg.Live {
		e.Pre(middleware.HTTPSNonWWWRedirect())
	}

	e.GET("/*", echo.WrapHandler(http.FileServer(getFileSystem(cfg.Live))))
	e.POST("/contact", func(c echo.Context) error {
		if err := sendEmail(c, &cfg); err != nil {
			return err
		}

		return c.String(http.StatusOK, "OK")
	})

	e.Logger.Fatal(e.Start(":" + cfg.Addr))
}

func sendEmail(c echo.Context, cfg *Config) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	comment := c.FormValue("comment")

	m := gomail.NewMessage()
	m.SetHeader("From", "Zombie Bus Escape Experience <reservas@zombiebus.es>")
	m.SetHeader("To", "info@zombiebus.es")
	m.SetHeader("Reply-To", fmt.Sprintf("%s <%s>", name, email))
	m.SetHeader("Subject", "Nuevo comentario desde la web")
	m.SetBody("text/plain", fmt.Sprintf("Nuevo comentario desde la web.\n\nNombre: %s\nEmail: %s\nMensaje: \n%s\n", name, email, comment))

	d := gomail.NewDialer("in-v3.mailjet.com", 587, cfg.SMTPUser, cfg.SMTPPass)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
