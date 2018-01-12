package api

import (
	"errors"
	"log"
	"os"
)

//GetClient - returns an API client
func GetClient() (*Client, error) {

	server, ok := os.LookupEnv("PULSEVTM_SERVER")
	if ok == false || server == "" {
		return nil, errors.New("[ERROR] PULSEVTM_SERVER env var not set")
	}

	username, ok := os.LookupEnv("PULSEVTM_USERNAME")
	if ok == false {
		return nil, errors.New("[ERROR] PULSEVTM_USERNAME env var not set")
	}

	password, ok := os.LookupEnv("PULSEVTM_PASSWORD")
	if ok == false {
		return nil, errors.New("[ERROR] PULSEVTM_PASSWORD env var not set")
	}

	apiVersion, ok := os.LookupEnv("PULSEVTM_API_VERSION")
	if ok == false {
		log.Println("The env var PULSEVTM_API_VERSION is not set, defaulting to 5.1")
		apiVersion = "5.1"
	}

	params := Params{
		APIVersion: apiVersion,
		Server:     server,
		Username:   username,
		Password:   password,
		IgnoreSSL:  true,
		Debug:      true,
	}

	return Connect(params)
}
