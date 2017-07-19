package brocadevtm

import (
	"github.com/sky-uk/go-brocade-vtm"
	"log"
)

// Config is a struct for containing the provider parameters.
type Config struct {
	Debug       bool
	Insecure    bool
	VTMUser     string
	VTMPassword string
	VTMServer   string
	Headers     map[string]string
}

// Client returns a new client for accessing VMWare vSphere.
func (c *Config) Client() (*brocadevtm.VTMClient, error) {
	log.Printf("[INFO] Brocade vTM Client configured for URL: %s", c.VTMServer)
	vtmClient := brocadevtm.NewVTMClient("https://"+c.VTMServer, c.VTMUser, c.VTMPassword, c.Insecure, c.Debug, c.Headers)
	return vtmClient, nil
}
