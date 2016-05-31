package facebook

import (
	"bytes"
	"datatrack/remote/facebook"
	"fmt"
	"net/http"
	"strings"

	"github.com/albrow/forms"
	"github.com/zenazn/goji/web"
)

func facebookHandler(c web.C, w http.ResponseWriter, r *http.Request) {
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
		err = facebook.ParseDataZip(bytes.NewReader(data))
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "missing file or unsuported format", http.StatusBadRequest)
}
