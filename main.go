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
	Live     bool   `envconfig:"LIVE"`
	SMTPPass string `envconfig:"SMTP_USER" default:"reservas@zombiebus.es"`
	SMTPUser string `envconfig:"SMTP_USER" default:"fQbgMFf4Dk6SsS2"`
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

	fmt.Println(cfg.Live)
	e := echo.New()
	if cfg.Live {
		e.Pre(middleware.HTTPSNonWWWRedirect())
	}

	e.GET("/*", echo.WrapHandler(http.FileServer(getFileSystem(cfg.Live))))
	e.GET("/hello", func(c echo.Context) error {
		sendEmail(&cfg)
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":" + cfg.Addr))
}

func sendEmail(cfg *Config) {
	fmt.Println(cfg.SMTPPass, cfg.SMTPUser)
	m := gomail.NewMessage()
	m.SetHeader("From", "from@gmail.com")
	m.SetHeader("To", "to@example.com")
	m.SetHeader("Subject", "Gomail test subject")
	m.SetBody("text/plain", "This is Gomail test body")

	d := gomail.NewDialer("smtp.gmail.com", 587, cfg.SMTPUser, cfg.SMTPUser)

	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return
}
