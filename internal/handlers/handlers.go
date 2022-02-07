package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/javalen/bookings/internal/config"
	"github.com/javalen/bookings/internal/forms"
	"github.com/javalen/bookings/internal/models"
	"github.com/javalen/bookings/internal/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.html", &models.TemplateData{}, r)
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)

	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	log.Printf("THe remote IP is ", remoteIP)
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	}, r)
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "generals.page.html", &models.TemplateData{}, r)
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "majors.page.html", &models.TemplateData{}, r)
}

func (m *Repository) Search(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.html", &models.TemplateData{}, r)
}

func (m *Repository) PostSearch(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("Start day is: %s and EndDate is %s", start, end)))
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplate(w, "make-reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	}, r)
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	reservation := models.Reservation {
		FirstName: r.Form.Get("first_name"),
		LastName: r.Form.Get("last_name"),
		Phone: r.Form.Get("phone"),
		Email: r.Form.Get("email"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)

	if !form.Valid(){
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, "make-reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		}, r)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/summary", http.StatusSeeOther)

}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.html", &models.TemplateData{}, r)
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) AvailabilityPost(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      false,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) Summary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("Cannot get model from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from Session")
		http.Redirect(w,r,"/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, "summary.page.html", &models.TemplateData{
		Data: data,
	}, r)
}