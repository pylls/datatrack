package testdata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"

	"github.com/pylls/datatrack/database"
	"github.com/pylls/datatrack/insert"
	"github.com/pylls/datatrack/model"
	"github.com/zenazn/goji/web"
)

func getHandler(file string) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		err := database.Setup()
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		}
		location := file
		if len(file) == 0 {
			location = c.URLParams["f"]
		}
		location = path.Join("testdata/", filepath.Base(filepath.Clean(location)))

		data, err := ioutil.ReadFile(location)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		}

		var ins model.Insert
		err = json.Unmarshal(data, &ins)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		}
		err = insert.DoInsert(&ins)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
