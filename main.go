package main

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	gomail "gopkg.in/mail.v2"
)

//go:embed static
var embededFiles embed.FS

//go:embed templates
var templatesFiles embed.FS

type Config struct {
	Addr     string `envconfig:"PORT"`
	Live     bool   `envconfig:"LIVE" default:"true"`
	SMTPPass string `envconfig:"SMTP_USER" default:"b8ea4a5d52519a11b01edad1f7cc5854"`
	SMTPUser string `envconfig:"SMTP_USER" default:"f06c3e95b93b1c291f38435d7b127816"`
}

type Page struct {
	URI       string
	Templates []string
	Data      interface{}
}

type TemplateRegistry struct {
	templates map[string]*template.Template
}

func (t *TemplateRegistry) Add(e *echo.Echo, page *Page) *echo.Route {
	if _, ok := t.templates[page.URI]; ok {
		panic(fmt.Sprintf("template %s alreadyd defined", page.URI))
	}

	if t.templates == nil {
		t.templates = make(map[string]*template.Template)
	}

	fs, _ := fs.Sub(templatesFiles, "templates")
	t.templates[page.URI] = template.Must(template.ParseFS(fs, page.Templates...))

	return e.GET(page.URI, func(c echo.Context) error {
		return c.Render(http.StatusOK, page.URI, page.Data)
	})
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}

	return tmpl.ExecuteTemplate(w, "base.html", data)
}

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

	templates := &TemplateRegistry{}
	templates.Add(e, &Page{URI: "/", Templates: []string{"index.html", "parts/base.html"}})
	templates.Add(e, &Page{URI: "/booking/ok/", Templates: []string{"booking-ok.html", "parts/base.html"}})
	templates.Add(e, &Page{URI: "/booking/error/", Templates: []string{"booking-error.html", "parts/base.html"}})

	e.Renderer = templates
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

func Render(template string, data interface{}) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, template, data)
	}
}
