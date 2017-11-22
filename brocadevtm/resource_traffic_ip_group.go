package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
	"log"
	"net/http"
	"regexp"
)

func resourceTrafficIPGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTrafficIPGroupSet,
		Read:   resourceTrafficIPGroupRead,
		Update: resourceTrafficIPGroupSet,
		Delete: resourceTrafficIPGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Traffic IP group name",
				Required:    true,
				ForceNew:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the traffic IP group should be enabled",
				Default:     false,
				Optional:    true,
			},
			"hash_source_port": {
				Type:        schema.TypeBool,
				Description: "Whether or not the source port should be taken into account when deciding which traffic manager should handle a request.",
				Default:     false,
				Optional:    true,
			},
			"ip_assignment_mode": {
				Type:        schema.TypeString,
				Description: "Configure how traffic IPs are assigned to traffic managers in single hosted mode",
				Computed:    true,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"alphabetic",
					"balanced",
				}, false),
			},
			"ip_mapping": {
				Type:        schema.TypeSet,
				Description: "Table matching traffic IPs to machines which should host the IPs",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:         schema.TypeString,
							Description:  "Traffic IP address",
							Optional:     true,
							ValidateFunc: util.ValidateIP,
						},
						"traffic_manager": {
							Type:        schema.TypeString,
							Description: "The name of the traffic manager which should host the IP",
							Optional:    true,
						},
					},
				},
			},
			"ipaddresses": {
				Type:        schema.TypeSet,
				Description: "List of IP addresses to raise on the traffic managers",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"keeptogether": {
				Type:        schema.TypeBool,
				Description: "Whether or not all traffic IPs are raised on a single traffic manager",
				Default:     false,
				Optional:    true,
			},
			"location": {
				Type:        schema.TypeInt,
				Description: "The location where the traffic IP group is based",
				Default:     0,
				Optional:    true,
			},
			"machines": {
				Type:        schema.TypeSet,
				Description: "List of traffic managers on which to raise this traffic IP - automatically retrieved from vTM",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"mode": {
				Type:        schema.TypeString,
				Description: "The method used to distribute traffic IPs across machines in the cluster",
				Optional:    true,
				Computed:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"singlehosted",
					"ec2elastic",
					"ec2vpcelastic",
					"ec2vpcprivate",
					"multihosted",
					"rhi",
				}, false),
			},
			"multicast": {
				Type:         schema.TypeString,
				Description:  "Multicast IP address",
				Optional:     true,
				ValidateFunc: validateTrafficIPGroupMulticastIP,
			},
			"note": {
				Type:        schema.TypeString,
				Description: "A note to attach to this traffic IP group",
				Optional:    true,
			},
			"rhi_bgp_metric_base": {
				Type:         schema.TypeInt,
				Description:  "Base BGP routing metric",
				Default:      10,
				Optional:     true,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"rhi_bgp_passive_metric_offset": {
				Type:         schema.TypeInt,
				Description:  "BGP routing metric offset",
				Default:      10,
				Optional:     true,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"rhi_ospfv2_metric_base": {
				Type:         schema.TypeInt,
				Description:  "OSPFv2 routing metric",
				Default:      10,
				Optional:     true,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"rhi_ospfv2_passive_metric_offset": {
				Type:         schema.TypeInt,
				Description:  "OSPFv2 routing metric offset",
				Default:      10,
				Optional:     true,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"rhi_protocols": {
				Type:        schema.TypeString,
				Description: "List of protocols ro be used for RHI",
				Optional:    true,
				Computed:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"ospf",
					"bgp",
				}, false),
			},
			"slaves": {
				Type:        schema.TypeSet,
				Description: "List of passive traffic managers",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func validateTrafficIPGroupMulticastIP(v interface{}, k string) (ws []string, errors []error) {

	multicastIP := v.(string)
	validMulticastIPs := regexp.MustCompile(`^2[2-3][0-9]\.[0-9]+\.[0-9]+\.[0-9]+$`)
	if !validMulticastIPs.MatchString(multicastIP) {
		errors = append(errors, fmt.Errorf("%q must be a valid multicast IP (224.0.0.0 - 239.255.255.255)", k))
	}
	return
}

// getTrafficManagers : retrieves a list of traffic managers from the traffic manager set in the environment variable.
func getTrafficManagers(m interface{}) ([]string, error) {

	var trafficManagers []string
	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()

	trafficManagerList, err := client.GetAllResources("traffic_managers")
	if err != nil {
		return trafficManagers, fmt.Errorf("BrocadeVTM Traffic Managers error whilst retrieving the list of Traffic Managers: %v", err)
	}

	for _, trafficManagerItem := range trafficManagerList {
		trafficManagerName := trafficManagerItem["name"].(string)
		trafficManagers = append(trafficManagers, trafficManagerName)
	}
	return trafficManagers, nil
}

func basicKeys() []string {
	return []string{
		"enabled",
		"hash_source_port",
		"ip_assignment_mode",
		"ipaddresses",
		"ip_mapping",
		"keeptogether",
		"location",
		"machines",
		"mode",
		"multicast",
		"note",
		"rhi_bgp_metric_base",
		"rhi_bgp_passive_metric_offset",
		"rhi_ospfv2_metric_base",
		"rhi_ospfv2_passive_metric_offset",
		"rhi_protocols",
		"slaves",
	}
}

func getSection(d *schema.ResourceData, sectionName string, properties map[string]interface{}, keys []string) error {
	m, err := util.GetAttributesToMap(d, keys)
	if err != nil {
		log.Println("Error getting section ", sectionName, err)
		return err
	}
	properties[sectionName] = m
	return nil
}

func resourceTrafficIPGroupSet(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	trafficIPGroupRequest := make(map[string]interface{})
	trafficIPGroupProperties := make(map[string]interface{})

	name := d.Get("name").(string)

	getSection(d, "basic", trafficIPGroupProperties, basicKeys())

	// If the list of traffic managers (machines isn't provided by the user get it from the traffic manager we're running against.
	if len(trafficIPGroupProperties["basic"].(map[string]interface{})["machines"].([]interface{})) == 0 {
		trafficManagers, err := getTrafficManagers(m)
		if err != nil {
			return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst creating %s: %v", name, err)
		}
		trafficIPGroupProperties["basic"].(map[string]interface{})["machines"] = trafficManagers
	}

	trafficIPGroupRequest["properties"] = trafficIPGroupProperties
	err := client.Set("traffic_ip_groups", name, trafficIPGroupRequest, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst creating %s: %s", name, err)
	}
	d.SetId(name)
	return resourceTrafficIPGroupRead(d, m)
}

func resourceTrafficIPGroupRead(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	client.WorkWithConfigurationResources()
	name := d.Id()

	trafficIPGroupConfiguration := make(map[string]interface{})
	err := client.GetByName("traffic_ip_groups", name, &trafficIPGroupConfiguration)
	if err != nil {
		if client.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst retrieving %s: %v", name, err)
	}
	trafficIPGroupProperties := trafficIPGroupConfiguration["properties"].(map[string]interface{})
	trafficIPGroupBasic := trafficIPGroupProperties["basic"].(map[string]interface{})


	trafficIPGroupPropertiesConfiguration := trafficIPGroupConfiguration["properties"].(map[string]interface{})
	trafficIPGroupBasicConfiguration := trafficIPGroupPropertiesConfiguration["basic"].(map[string]interface{})
	util.SetSimpleAttributesFromMap(d, trafficIPGroupBasicConfiguration, "", getMonitorMapAttributeList("basic"))
	util.SetSimpleAttributesFromMap(d, trafficIPGroupBasicConfiguration, "", []string{"machines"})
	d.Set("ip_mapping", trafficIPGroupBasicConfiguration["ip_mapping"].([]interface{}))


	d.Set("ip_mapping", trafficIPGroupBasicConfiguration["ip_mapping"].([]interface{}))


	ipMappings, err := util.BuildReadListMaps(trafficIPGroupBasicConfiguration["ip_mapping"].(map[string]interface{}), "ip_mapping")
	if err != nil {
		return err
	}
	d.Set("ip_mapping", ipMappings)

	return nil
}

func resourceTrafficIPGroupUpdate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	trafficIPGroupConfiguration := make(map[string]interface{})
	trafficIPGroupPropertiesConfiguration := make(map[string]interface{})
	name := d.Id()

	trafficIPGroupBasicConfiguration := make(map[string]interface{})
	trafficIPGroupBasicConfiguration = util.AddChangedSimpleAttributesToMap(d, trafficIPGroupBasicConfiguration, "", getMonitorMapAttributeList("basic"))
	if d.HasChange("ip_mapping") {
		if v, ok := d.GetOk("ip_mapping"); ok {
			builtList, err := util.BuildListMaps(v.(*schema.Set), getTrafficIPGroupAttributeList("ip_mapping"))
			if err != nil {
				return err
			}
			trafficIPGroupBasicConfiguration["ip_mapping"] = builtList
		}
	}
	if d.HasChange("machines") {
		// Allow the user to override the list of traffic managers. If not specified by the user retrieving them from the traffic manager.
		if v, ok := d.GetOk("machines"); ok {
			trafficIPGroupBasicConfiguration["machines"] = util.BuildStringListFromSet(v.(*schema.Set))
		} else {
			trafficManagers, err := getTrafficManagers(m)
			if err != nil {
				return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst creating %s: %v", name, err)
			}
			trafficIPGroupBasicConfiguration["machines"] = trafficManagers
		}
	}
	trafficIPGroupPropertiesConfiguration["basic"] = trafficIPGroupBasicConfiguration
	trafficIPGroupConfiguration["properties"] = trafficIPGroupPropertiesConfiguration

	err := client.Set("traffic_ip_groups", name, trafficIPGroupConfiguration, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst updating %s: %v", name, err)

	}
	return nil
}

func resourceTrafficIPGroupDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("traffic_ip_groups", d, m)
}
