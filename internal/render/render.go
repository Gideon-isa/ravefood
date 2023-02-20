package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Gideon-isa/ravefoods/internal/config"
	"github.com/Gideon-isa/ravefoods/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{}

// app is a pointer to AppConfig
var app *config.AppConfig

// NewTemplate set the config for the template package from the main function
// by using a pointer variable which points to the app variable in the main function
func NewTemplate(a *config.AppConfig) {
	app = a
}

// AddDefault adds data for all template
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// adding the CSRF_Token to our template
	td.CSRFToken = nosurf.Token(r)

	return td
}

// renderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {

	var tm map[string]*template.Template

	// if UseCache is false i.e In development mode
	// if UseCache is true i.e In Production mode
	if app.UseCache {
		//then tm gets its templates from the template cache
		// get the template cache from the app config
		// where the app.TemplateCache has been updated
		// by the NewTemplate function been called on the main function
		// already created in the main func in line 22
		tm = app.TemplateCache
	} else {
		// then it will read the template afresh
		tm, _ = CreateTemplate()
	}

	// This creates the template cache everytime we make a request: not efficient
	// Instead we use the app config
	// tm, err := CreateTemplate()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	t, ok := tm[tmpl]

	// if the tmpl does not exist in the map
	if !ok {
		log.Fatal("unable to load template from the template cache")
	}

	// writing the template to a buffer
	buf := new(bytes.Buffer)

	//
	td = AddDefaultData(td, r)

	// write to  the buffer
	t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing to browser")
	}

}

// CreateTemplate
func CreateTemplate() (map[string]*template.Template, error) {

	// tmCache is an an initialize *template.Template{}
	tmCache := map[string]*template.Template{}

	//getting all the pages
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return tmCache, err
	}

	//looping over the pages
	for _, page := range pages {
		//what it returns the full path e.g ./templates/home.page.tmpl etc...
		// get the name of the last file name i.e home.page.tmpl, using the Base func in the filepath package
		name := filepath.Base(page)

		//creating a new and empty template using the New function
		//ParseFiles parses the named file i.e (./templates/home.page.tmpl)
		//and associates the resulting template with the "name"
		//ts is a container template or better a Set of templates. So you can add to it as much as you want
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return tmCache, err
		}

		//getting the layouts file path (e.g base.layout.html)
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return tmCache, err
		}

		//checking if the container of the layout.html is empty or not
		if len(matches) > 0 {
			// parsing the layout.html to a go-template with the ts tempalte together into the ts template set
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return tmCache, err
			}
		}

		//taking the template set ts and adding it to the cache map
		tmCache[name] = ts

	}

	return tmCache, nil

}
