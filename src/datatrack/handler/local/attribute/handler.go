package attribute

import (
	"datatrack/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zenazn/goji/web"
)

func attributeHandler(m mode, op ...operation) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		countOutput := false
		var out []string
		var err error

		// get data based on mode
		switch m {
		case attribute:
			out, err = database.GetAttributeIDs()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case explicitAttribute:
			out, err = database.GetExplicitlyDisclosedAttributeIDs()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case implicitAttribute:
			out, err = database.GetImplicitlyDisclosedAttributeIDs()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case organization:
			out, err = database.GetAttributeIDsToOrg(c.URLParams["organizationId"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case explicitOrganization:
			out, err = database.GetExplicitlyDisclosedAttributeIDsToOrg(c.URLParams["organizationId"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		case implicitOrganization:
			out, err = database.GetImplicitlyDisclosedAttributeIDsToOrg(c.URLParams["organizationId"])
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

func attributeFieldHandler(m mode, f field, op ...operation) func(web.C, http.ResponseWriter, *http.Request) {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		countOutput := false
		var out []string
		var err error

		// get data based on mode
		switch m {
		case attribute:
			switch f {
			case thetype:
				out, err = database.GetAttributeTypeIDs()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			case thevalue:
				out, err = database.GetAttributeTypeValues(c.URLParams["typeId"])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

		case explicitAttribute:
			switch f {
			case thetype:
				out, err = database.GetExplicitlyDisclosedAttributeTypeIDs()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			case thevalue:
				out, err = database.GetExplicitlyDisclosedAttributeTypeValues(c.URLParams["typeId"])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

		case implicitAttribute:
			switch f {
			case thetype:
				out, err = database.GetImplicitlyDisclosedAttributeTypeIDs()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			case thevalue:
				out, err = database.GetImplicitlyDisclosedAttributeTypeValues(c.URLParams["typeId"])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

		case organization:
			switch f {
			case thetype:
				out, err = database.GetAttributeTypeIDsToOrg(c.URLParams["organizationId"])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			case thevalue:
				out, err = database.GetAttributeTypeValuesToOrg(c.URLParams["organizationId"], c.URLParams["typeId"])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

		case explicitOrganization:
			switch f {
			case thetype:
				out, err = database.GetExplicitlyDisclosedAttributeTypeIDsToOrg(c.URLParams["organizationId"])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			case thevalue:
				out, err = database.GetExplicitlyDisclosedAttributeTypeValuesToOrg(c.URLParams["organizationId"], c.URLParams["typeId"])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

		case implicitOrganization:
			switch f {
			case thetype:
				out, err = database.GetImplicitlyDisclosedAttributeTypeIDsToOrg(c.URLParams["organizationId"])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			case thevalue:
				out, err = database.GetImplicitlyDisclosedAttributeTypeValuesToOrg(c.URLParams["organizationId"], c.URLParams["typeId"])
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
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
	d, err := database.GetAttribute(c.URLParams["attributeId"])
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
