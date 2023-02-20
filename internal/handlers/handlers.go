package handlers

import (
	"net/http"

	"github.com/Gideon-isa/ravefoods/internal/config"
	"github.com/Gideon-isa/ravefoods/internal/models"
	"github.com/Gideon-isa/ravefoods/internal/render"
)

type Repository struct {
	App *config.AppConfig
}

// Repo the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository with the app config
// pointing to the Repository struct
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// HomePage is the home page
// link the Repo to all methods by making the handlers function a receiver
func (m *Repository) HomePage(w http.ResponseWriter, r *http.Request) {
	remoteAdd := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteAdd)
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// SignUp is the sign up page
func (m *Repository) SignUp(w http.ResponseWriter, r *http.Request) {
	theIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap := make(map[string]string)
	stringMap["theRemote"] = theIP

	// m.App.Session.Put(r.Context(), "first-name", r.PostForm.Get("firstName"))
	// m.App.Session.Put(r.Context(), "last-name", r.PostForm.Get("lastName"))
	// m.App.Session.Put(r.Context(), "password", r.PostForm.Get("password"))
	render.RenderTemplate(w, r, "signup.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// NigerianFoods is the nigerian food page
func (m *Repository) NigerianFoods(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Gideon"
	render.RenderTemplate(w, r, "nigerianfoods.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) PostSignUp(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	stringMap := make(map[string]string)

	stringMap["fn"] = r.PostForm.Get("firstName")
	stringMap["ln"] = r.PostForm.Get("lastName")
	stringMap["pw"] = r.PostForm.Get("password")

	render.RenderTemplate(w, r, "postsignup.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
