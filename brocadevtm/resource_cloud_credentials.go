package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
	"net/http"
)

func resourceCloudCredentials() *schema.Resource {
	return &schema.Resource{

		Create: resourceCloudCredentialsCreate,
		Read:   resourceCloudCredentialsRead,
		Update: resourceCloudCredentialsUpdate,
		Delete: resourceCloudCredentialsDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: util.NoZeroValues,
			},
			"api_server": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Cloud server hostname or IP address.",
			},
			"cloud_api_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "The traffic manager creates and destroys nodes via API calls. This setting specifies (in seconds) how long to wait for such calls to complete.",
				Default:      200,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"cred1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The first part of the credentials for the cloud user. Typically this is some variation on the username concept.",
			},
			"cred2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The second part of the credentials for the cloud user. Typically this is some variation on the password concept.",
				Sensitive:   true,
			},
			"cred3": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The third part of the credentials for the cloud user. Typically this is some variation on the authentication token concept.",
				Sensitive:   true,
			},
			"script": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The script to call for communication with the cloud API.",
			},
			"update_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "The traffic manager will periodically check the status of the cloud through an API call. This setting specifies the interval between such updates.",
				Default:      30,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
		},
	}
}

func resourceCloudCredentialsCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	cloudCredentialsConfiguration := make(map[string]interface{})
	cloudCredentialsPropertiesConfiguration := make(map[string]interface{})
	cloudCredentialsBasicConfiguration := make(map[string]interface{})

	name := d.Get("name").(string)

	if v, ok := d.GetOk("api_server"); ok {
		cloudCredentialsBasicConfiguration["api_server"] = v.(string)
	}

	cloudCredentialsBasicConfiguration["cloud_api_timeout"] = uint(d.Get("cloud_api_timeout").(int))

	if v, ok := d.GetOk("cred1"); ok {
		cloudCredentialsBasicConfiguration["cred1"] = v.(string)
	}
	if v, ok := d.GetOk("cred2"); ok {
		cloudCredentialsBasicConfiguration["cred2"] = v.(string)
	}
	if v, ok := d.GetOk("cred3"); ok {
		cloudCredentialsBasicConfiguration["cred3"] = v.(string)
	}
	if v, ok := d.GetOk("script"); ok {
		cloudCredentialsBasicConfiguration["script"] = v.(string)
	}

	cloudCredentialsBasicConfiguration["update_interval"] = uint(d.Get("update_interval").(int))

	cloudCredentialsPropertiesConfiguration["basic"] = cloudCredentialsBasicConfiguration
	cloudCredentialsConfiguration["properties"] = cloudCredentialsPropertiesConfiguration

	err := client.Set("cloud_api_credentials", name, &cloudCredentialsConfiguration, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] BrocadeVTM error whilst creating Cloud API Credentials %s: %v", name, err)
	}
	d.SetId(name)
	return resourceCloudCredentialsRead(d, m)
}

func resourceCloudCredentialsUpdate(d *schema.ResourceData, m interface{}) error {
	name := d.Id()
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	cloudCredentialsConfiguration := make(map[string]interface{})
	cloudCredentialsPropertiesConfiguration := make(map[string]interface{})
	cloudCredentialsBasicConfiguration := make(map[string]interface{})

	if d.HasChange("api_server") {
		cloudCredentialsBasicConfiguration["api_server"] = d.Get("api_server").(string)
	}
	if d.HasChange("cloud_api_timeout") {
		cloudCredentialsBasicConfiguration["cloud_api_timeout"] = uint(d.Get("cloud_api_timeout").(int))
	}
	if d.HasChange("cred1") {
		cloudCredentialsBasicConfiguration["cred1"] = d.Get("cred1").(string)
	}
	if d.HasChange("cred2") {
		cloudCredentialsBasicConfiguration["cred2"] = d.Get("cred2").(string)
	}
	if d.HasChange("cred3") {
		cloudCredentialsBasicConfiguration["cred3"] = d.Get("cred3").(string)
	}
	if d.HasChange("script") {
		cloudCredentialsBasicConfiguration["script"] = d.Get("script").(string)
	}
	if d.HasChange("update_interval") {
		cloudCredentialsBasicConfiguration["update_interval"] = uint(d.Get("update_interval").(int))
	}

	cloudCredentialsPropertiesConfiguration["basic"] = cloudCredentialsBasicConfiguration
	cloudCredentialsConfiguration["properties"] = cloudCredentialsPropertiesConfiguration

	err := client.Set("cloud_api_credentials", name, &cloudCredentialsConfiguration, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] BrocadeVTM error whilst updating Cloud API Credentials %s: %v", name, err)
	}

	return resourceCloudCredentialsRead(d, m)
}

func resourceCloudCredentialsRead(d *schema.ResourceData, m interface{}) error {
	name := d.Id()
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	cloudCredentialsConfiguration := make(map[string]interface{})

	err := client.GetByName("cloud_api_credentials", name, &cloudCredentialsConfiguration)

	if client.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("[ERROR] BrocadeVTM error whilst retrieving Cloud API Credentials %s: %v", name, err)
	}

	cloudCredentialsPropertiesConfiguration := cloudCredentialsConfiguration["properties"].(map[string]interface{})
	cloudCredentialsBasicConfiguration := cloudCredentialsPropertiesConfiguration["basic"].(map[string]interface{})

	for _, key := range []string{
		"api_server",
		"cloud_api_timeout",
		"cred1",
		"cred2",
		"cred3",
		"script",
		"update_interval",
	} {
		err := d.Set(key, cloudCredentialsBasicConfiguration[key])
		if err != nil {
			return fmt.Errorf("[ERROR] BrocadeVTM Cloud Credentials error whilst setting key %s: %v", key, err)
		}
	}

	return nil
}

func resourceCloudCredentialsDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("cloud_api_credentials", d, m)
}
