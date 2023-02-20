// Package main is the main entry of the application
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Gideon-isa/ravefoods/internal/config"
	"github.com/Gideon-isa/ravefoods/internal/handlers"
	"github.com/Gideon-isa/ravefoods/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":9000"

// app is the variable of type config,AppConfig that holds all the AppConfig values
var app config.AppConfig

// session variable
var session *scs.SessionManager

func main() {

	// change this to true when in production
	app.InProduction = false

	// initializing the session package manager
	session = scs.New()
	session.Lifetime = time.Hour * 20

	// availability of the cookie after the browser is closed
	session.Cookie.Persist = true

	//
	session.Cookie.SameSite = http.SameSiteLaxMode

	// making cookie encrypted
	session.Cookie.Secure = app.InProduction

	// assigning or storing the values of the session to the AppConfig struct
	// making it accessible to any package needed
	app.Session = session

	// tc holds the map of templates parsed by the render.CreateTemplate()
	tc, err := render.CreateTemplate()
	if err != nil {
		log.Fatal(err)
	}

	// TemplateCache field of the app variable is assigned the tc values
	app.TemplateCache = tc

	// In development we build the template each time we make a chane
	// by setting the UseCache to false
	app.UseCache = false

	// This assigns the struct app variable to "app" on the rennder.go in the render package
	// making the contents (template cache) available to the render package
	// it gives the render package access to the template
	// by using the pointer to the app variable for effiency
	render.NewTemplate(&app)

	// this creates the Repository variable
	repo := handlers.NewRepo(&app)
	// this feed to the handlers the app config and db data
	// by giving access to the Repo variable in the handlers package
	handlers.NewHandlers(repo)

	fmt.Printf("Starting the server on port %s...\n", portNumber)

	serv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = serv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
