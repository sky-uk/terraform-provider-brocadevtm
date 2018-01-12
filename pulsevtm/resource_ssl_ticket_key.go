package pulsevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sky-uk/go-pulse-vtm/api"
	"net/http"
)

func resourceSSLTicketKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSSLTicketKeySet,
		Read:   resourceSSLTicketKeyRead,
		Update: resourceSSLTicketKeySet,
		Delete: resourceSSLTicketKeyDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the SSL Ticket Key",
			},
			"algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The algorithm used to encrypt session tickets",
				Default:      "aes_256_cbc_hmac_sha256",
				ValidateFunc: validation.StringInSlice([]string{"aes_256_cbc_hmac_sha256"}, false),
			},
			"identifier": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A 16-byte key identifier, with each byte encoded as two hexadecimal digits.",
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The session ticket encryption key, with each byte encoded as two hexadecimal digits.",
			},
			"validity_end": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The latest time at which this key may be used to encrypt new session tickets. Given as number of seconds since the epoch",
				ValidateFunc: validation.IntAtLeast(0),
			},
			"validity_start": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The earliest time at which this key may be used to encrypt new session tickets. Given as number of seconds since the epoch",
				ValidateFunc: validation.IntAtLeast(0),
			},
		},
	}
}

func resourceSSLTicketKeySet(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	sslTicketKeyConfiguration := make(map[string]interface{})
	sslTicketKeyPropertiesConfiguration := make(map[string]interface{})
	sslTicketKeyBasicConfiguration := make(map[string]interface{})

	name := d.Get("name").(string)

	sslTicketKeyBasicConfiguration["algorithm"] = d.Get("algorithm").(string)
	sslTicketKeyBasicConfiguration["id"] = d.Get("identifier").(string)
	sslTicketKeyBasicConfiguration["key"] = d.Get("key").(string)
	sslTicketKeyBasicConfiguration["validity_end"] = d.Get("validity_end").(int)
	sslTicketKeyBasicConfiguration["validity_start"] = d.Get("validity_start").(int)

	sslTicketKeyPropertiesConfiguration["basic"] = sslTicketKeyBasicConfiguration
	sslTicketKeyConfiguration["properties"] = sslTicketKeyPropertiesConfiguration

	err := client.Set("ssl/ticket_keys", name, &sslTicketKeyConfiguration, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] PulseVTM error whilst creating SSL Ticket Key %s: %v", name, err)
	}
	d.SetId(name)

	return resourceSSLTicketKeyRead(d, m)
}

func resourceSSLTicketKeyRead(d *schema.ResourceData, m interface{}) error {

	name := d.Id()
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	sslTicketKeyConfiguration := make(map[string]interface{})

	err := client.GetByName("ssl/ticket_keys", name, &sslTicketKeyConfiguration)

	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("[ERROR] PulseVTM error whilst retrieving  SSL Ticket Key %s: %v", name, err)
	}

	sslTicketKeyPropertiesConfiguration := sslTicketKeyConfiguration["properties"].(map[string]interface{})
	sslTicketKeyBasicConfiguration := sslTicketKeyPropertiesConfiguration["basic"].(map[string]interface{})

	for _, key := range []string{
		"algorithm",
		"identifier",
		"key",
		"validity_end",
		"validity_start",
	} {
		err := d.Set(key, sslTicketKeyBasicConfiguration[getTrueName(key)])
		if err != nil {
			return fmt.Errorf("[ERROR] PulseVTM SSL Ticket Key error whilst setting attribute %s: %v", key, err)
		}
	}

	return nil
}

func resourceSSLTicketKeyDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("ssl/ticket_keys", d, m)
}

func getTrueName(attribute string) string {
	if attribute == "identifier" {
		return "id"
	}
	return attribute
}
