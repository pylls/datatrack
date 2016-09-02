package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goji/param"
	"github.com/marcelfarres/datatrack/database"
	"github.com/marcelfarres/datatrack/model"
	"github.com/zenazn/goji/web"
)

func userHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	if u, err := database.GetUser(); err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
	} else if s, err := json.Marshal(u); err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-type", "application/json")
		fmt.Fprintf(w, "%s", s)
	}
}

func updateUserHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var user model.User
	if err := param.Parse(r.Form, &user); err != nil {
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}
	database.SetUser(user)
}
