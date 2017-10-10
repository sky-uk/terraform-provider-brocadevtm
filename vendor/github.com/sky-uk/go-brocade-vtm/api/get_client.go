package api

import (
	"errors"
	"os"
)

//GetClient - returns an API client
func GetClient() (*Client, error) {

	server, ok := os.LookupEnv("BROCADEVTM_SERVER")
	if ok == false || server == "" {
		return nil, errors.New("BROCADEVTM_SERVER env var not set")
	}

	username, ok := os.LookupEnv("BROCADEVTM_USERNAME")
	if ok == false {
		return nil, errors.New("BROCADEVTM_USERNAME env var not set")
	}

	password, ok := os.LookupEnv("BROCADEVTM_PASSWORD")
	if ok == false {
		return nil, errors.New("BROCADEVTM_PASSWORD env var not set")
	}

	params := Params{
		APIVersion: "3.8",
		Server:     server,
		Username:   username,
		Password:   password,
		IgnoreSSL:  true,
		Debug:      true,
	}

	return Connect(params)
}
