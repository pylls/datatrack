package googletakeout

import (
	"bytes"
	"datatrack/database"
	"datatrack/model"
	"datatrack/remote/google"
	"fmt"
	"net/http"
	"strings"

	"github.com/albrow/forms"
	"github.com/zenazn/goji/web"
)

func takeoutHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	formData, err := forms.Parse(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}
	file := formData.GetFile("file")
	if file == nil {
		http.Error(w, "missing file", http.StatusInternalServerError)
		return
	}

	data, err := formData.GetFileBytes("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	if strings.HasSuffix(file.Filename, ".zip") {
		err = google.ParseTakeoutZip(bytes.NewReader(data))
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}
		_, err = database.GetUser()
		if err != nil {
			err = database.SetUser(model.User{
				Name:    "Alice",
				Picture: "defaultuser.png",
			})
			if err != nil {
				http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			}
		}
		return
	} else if strings.HasSuffix(file.Filename, ".tgz") {
		err = google.ParseTakeoutGzip(bytes.NewReader(data))
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}
		_, err = database.GetUser()
		if err != nil {
			err = database.SetUser(model.User{
				Name:    "Alice",
				Picture: "defaultuser.png",
			})
			if err != nil {
				http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			}
		}
	}

	http.Error(w, "missing file or unsuported format", http.StatusBadRequest)
}
