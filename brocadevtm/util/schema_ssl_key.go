package util

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"net/http"
)

//SchemaSSLKey : Returns an SSL Key Schema
func SchemaSSLKey() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: NoZeroValues,
		},

		"note": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},

		"private": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},

		"public": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},

		"request": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func SSLKeyCreate(d *schema.ResourceData, meta interface{}, keyType string) error {

	name := d.Get("name").(string)

	sslKeyPropertiesConfig := make(map[string]interface{})
	sslKeyBasicConfig := make(map[string]interface{})
	sslKeyConfig := make(map[string]interface{})

	if v, ok := d.GetOk("note"); ok {
		sslKeyBasicConfig["note"] = v.(string)
	}
	if v, ok := d.GetOk("private"); ok {
		sslKeyBasicConfig["private"] = v.(string)
	}
	if v, ok := d.GetOk("public"); ok {
		sslKeyBasicConfig["public"] = v.(string)
	}
	if v, ok := d.GetOk("request"); ok {
		sslKeyBasicConfig["request"] = v.(string)
	}
	sslKeyPropertiesConfig["basic"] = sslKeyBasicConfig
	sslKeyConfig["properties"] = sslKeyPropertiesConfig

	config := meta.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	err := client.Set(keyType, name, sslKeyConfig, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM %s error whilst creating %s: %v", keyType, name, err)
	}
	d.SetId(name)

	return nil
}

func SSLKeyRead(d *schema.ResourceData, meta interface{}, keyType string) error {
	config := meta.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	sslClientKeyConfig := make(map[string]interface{})
	err := client.GetByName(keyType, d.Id(), &sslClientKeyConfig)
	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BrocadeVTM %s error whilst retrieving %s: %v", keyType, d.Id(), err)
	}

	sslClientKeyPropertiesConfig := sslClientKeyConfig["properties"].(map[string]interface{})
	sslClientKeyBasicConfig := sslClientKeyPropertiesConfig["basic"].(map[string]interface{})

	d.Set("note", sslClientKeyBasicConfig["note"])
	d.Set("public", sslClientKeyBasicConfig["public"])
	d.Set("request", sslClientKeyBasicConfig["request"])

	return nil
}

func SSLKeyUpdate(d *schema.ResourceData, meta interface{}, keyType string) error {

	sslKeyPropertiesConfig := make(map[string]interface{})
	sslKeyBasicConfig := make(map[string]interface{})
	sslKeyConfig := make(map[string]interface{})

	if d.HasChange("note") {
		sslKeyBasicConfig["note"] = d.Get("note").(string)
	}
	if d.HasChange("private") {
		sslKeyBasicConfig["private"] = d.Get("private").(string)
	}
	if d.HasChange("public") {
		sslKeyBasicConfig["public"] = d.Get("public").(string)
	}
	if d.HasChange("request") {
		sslKeyBasicConfig["request"] = d.Get("request").(string)
	}

	sslKeyPropertiesConfig["basic"] = sslKeyBasicConfig
	sslKeyConfig["properties"] = sslKeyPropertiesConfig

	config := meta.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	err := client.Set(keyType, d.Id(), sslKeyConfig, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM %s error whilst updating %s: %v", keyType, d.Id(), err)
	}
	return nil
}

func SSLKeyDelete(d *schema.ResourceData, meta interface{}, keyType string) error {
	config := meta.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)

	err := client.Delete(keyType, d.Id())
	if client.StatusCode == http.StatusNoContent || client.StatusCode == http.StatusNotFound {
		return nil
	}
	return fmt.Errorf("BrocadeVTM %s error whilst deleting %s: %v", keyType, d.Id(), err)
}
