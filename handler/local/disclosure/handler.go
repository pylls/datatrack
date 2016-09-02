package disclosure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/marcelfarres/datatrack/database"
	"github.com/zenazn/goji/web"
)

func disclosureHandler(m mode, op ...operation) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		countOutput := false
		var out []string
		var err error

		// get data based on mode
		switch m {
		case disclosure:
			out, err = database.GetDisclosureIDs()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case disclosureChrono:
			out, err = database.GetDisclosureIDsChrono()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case organization:
			out, err = database.GetDisclosureIDsToOrg(c.URLParams["organizationId"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case organizationChrono:
			out, err = database.GetDisclosureIDsToOrgChrono(c.URLParams["organizationId"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case attribute:
			d, err := database.GetDisclosed(c.URLParams["disclosureId"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			out = d.Attribute
		case implicit:
			out, err = database.GetImplicitDisclosureIDs(c.URLParams["disclosureId"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case implicitChrono:
			out, err = database.GetImplicitDisclosureIDsChrono(c.URLParams["disclosureId"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case downstream:
			out, err = database.GetDownstreamDisclosureIDs(c.URLParams["disclosureId"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case downstreamChrono:
			out, err = database.GetDownstreamDisclosureIDsChrono(c.URLParams["disclosureId"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// go over each operation
		for i := 0; i < len(op); i++ {
			switch op[i] {
			case subset:
				first, err := strconv.Atoi(c.URLParams["first"])
				if err != nil {
					panic(err)
				}
				last, err := strconv.Atoi(c.URLParams["last"])
				if err != nil {
					panic(err)
				}
				if first < 0 || len(out) <= last || last > first {
					http.Error(w, "invalid range", http.StatusBadRequest)
					return
				}
				out = out[first:last]
			case reverse:
				for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
					out[i], out[j] = out[j], out[i]
				}
			case count:
				countOutput = true
			}
		}

		// make pretty output
		if countOutput {
			fmt.Fprintf(w, "%d", len(out))
			return
		}
		s, err := json.Marshal(out)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/json")
		fmt.Fprintf(w, "%s", s)
	}
}

func detailsHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	d, err := database.GetDisclosure(c.URLParams["disclosureId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	j, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-type", "application/json")
	fmt.Fprintf(w, "%s", j)
}
