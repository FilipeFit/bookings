package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FilipeFit/bookings/internal/handlers"
	"github.com/FilipeFit/bookings/internal/render"

	"github.com/FilipeFit/bookings/internal/config"
	"github.com/alexedwards/scs/v2"
)

const portNumer = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Printf("Starting the application in port %s", portNumer)

	srv := &http.Server{
		Addr:    portNumer,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}
