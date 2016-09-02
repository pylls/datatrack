package category

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pylls/datatrack/database"
	"github.com/zenazn/goji/web"
)

func categoryHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	cat, err := database.GetCategory(c.URLParams["categoryId"])
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
	} else if s, err := json.Marshal(cat); err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-type", "application/json")
		fmt.Fprintf(w, "%s", s)
	}
}
