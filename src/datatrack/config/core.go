// Package config parses the configuration from the environment and stores it
// in the Env variable.
package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// EnvT is the type representing the configuration environment.
type EnvT struct {
	DEBUG          bool
	Databasepath   string
	Insecureclient bool
}

// Env is the parsed environment set when the Data Track was started.
var Env = EnvT{
	DEBUG:          false,
	Databasepath:   "datatrack.db",
	Insecureclient: true,
}

// StaticPath is the relative path to the static HTML folder.
const StaticPath = "static/"

// APIURL is the prefix for the Data Track API.
const APIURL = "/v1"

// CountURL is the API suffix for only returning the count of an API query.
const CountURL = "/count"

// RangeURL is the API suffix pattern for requesting a range of replies.
const RangeURL = "/range/:first/:last"

// ChronologicalURL is the API pattern for requesting that replies are sorted chronologically.
const ChronologicalURL = "/chronological"

// ReverseURL is the API pattern for requesting that replies have their output reversed.
const ReverseURL = "/reverse"

// ImplicitURL is the API pattern for requesting only implicitly disclosed data.
const ImplicitURL = "/implicit"

// ExplicitURL is the API pattern for requesting only explicitly disclosed data.
const ExplicitURL = "/explicit"

// Configure reads the configuration from the environment into the Env variable.
func Configure() error {
	log.SetFlags(log.Flags() | log.Lmicroseconds)
	if err := envconfig.Process("datatrack", &Env); err != nil {
		log.Fatalf("failed to parse environment configuration (%s)", err)
	}
	return nil
}
