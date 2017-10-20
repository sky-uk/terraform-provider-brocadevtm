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

	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("BrocadeVTM error whilst deleting %s %s: %v", resourceType, name, err)
	}
	return nil
}
