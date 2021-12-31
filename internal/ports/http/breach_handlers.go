package http

import (
	"encoding/json"
	"fmt"
	"gobreach/internal/domains/breach"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type BreachResponse struct {
	Domain       string    `json:"domain"`
	BreachedInfo []string  `json:"breachedInfo"`
	BreachDate   time.Time `json:"breachedDate"`
	BreachSource string    `json:"breachedSource"`
}

type breachHandler struct {
	service BreachServer
	logger  *log.Logger
}


func newBreachRouter(s *BreachServer, l *log.Logger) http.Handler {
	bh := &breachHandler{service: *s, logger: l}
	return addbreachRoutes(bh)
}

func addbreachRoutes(bh *breachHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/", bh.getBreaches)
	return r
}

func (bh *breachHandler) getBreaches(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		et := fmt.Sprintf("%v: 'email' query param required", http.StatusText(http.StatusBadRequest))
		http.Error(w, et, http.StatusBadRequest)
		return
	}

	bS, err := bh.service.GetByEmail(email)
	if err != nil {
		switch code := err.ErrCode; code {
		case breach.BreachValidationErr:
			bh.logger.Printf("error retrieving breaches by email: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			bh.logger.Printf("error retrieving breaches by email: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	rs := []BreachResponse{}
	for _, b := range bS {
		rs = append(rs, BreachResponse{Domain: b.Domain, BreachedInfo: b.BreachedInfo,
			BreachDate: b.BreachDate, BreachSource: b.BreachSource})
	}
	data, errr := json.Marshal(rs)
	if err != nil {
		bh.logger.Printf("error marshalling result to JSON: %v", errr)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
