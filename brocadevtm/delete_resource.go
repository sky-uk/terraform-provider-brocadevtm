package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

// DeleteResource - Deletes a Brocade vTM Configuration Resource
func DeleteResource(resourceType string, d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	name := d.Id()

	err := client.Delete(resourceType, name)

	if client.StatusCode == http.StatusNoContent || client.StatusCode == http.StatusNotFound {
		return nil
	}
	return fmt.Errorf("BrocadeVTM %s error whilst deleting %s: %v", resourceType, d.Id(), err)
}
