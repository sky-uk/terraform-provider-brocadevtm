package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/sky-uk/go-brocade-vtm/api"
	"github.com/sky-uk/go-brocade-vtm/api/model/3.8/traffic_ip_group"
	"github.com/sky-uk/terraform-provider-brocadevtm/brocadevtm/util"
	"net/http"
	"regexp"
)

func resourceTrafficIPGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTrafficIPGroupCreate,
		Read:   resourceTrafficIPGroupRead,
		Update: resourceTrafficIPGroupUpdate,
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

func buildIPMapping(ipMappingBlock *schema.Set) []trafficIpGroups.IPMapping {

	ipMappingObjectList := make([]trafficIpGroups.IPMapping, 0)

	for _, ipMappingItem := range ipMappingBlock.List() {

		ipMapping := ipMappingItem.(map[string]interface{})
		ipMappingObject := trafficIpGroups.IPMapping{}

		if ip, ok := ipMapping["ip"].(string); ok {
			ipMappingObject.IP = ip
		}
		if trafficManager, ok := ipMapping["traffic_manager"].(string); ok {
			ipMappingObject.TrafficManager = trafficManager
		}
		ipMappingObjectList = append(ipMappingObjectList, ipMappingObject)
	}
	return ipMappingObjectList
}

func getTrafficIPGroupAttributeList(mapName string) []string {

	attributes := []string{}

	switch mapName {
	case "basic":
		attributes = []string{"enabled",
			"hash_source_port",
			"ip_assignment_mode",
			"ipaddresses",
			"keeptogether",
			"location",
			"mode",
			"multicast",
			"note",
			"rhi_bgp_metric_base",
			"rhi_bgp_passive_metric_offset",
			"rhi_ospfv2_metric_base",
			"rhi_ospfv2_passive_metric_offset",
			"rhi_protocols",
			"slaves"}
	case "ip_mapping":
		attributes = []string{"ip", "traffic_manager"}
	}
	return attributes
}

func resourceTrafficIPGroupCreate(d *schema.ResourceData, m interface{}) error {

	config := m.(map[string]interface{})
	client := config["jsonClient"].(*api.Client)
	trafficIPGroupConfiguration := make(map[string]interface{})
	trafficIPGroupPropertiesConfiguration := make(map[string]interface{})

	name := d.Get("name").(string)

	trafficIPGroupBasicConfiguration := make(map[string]interface{})
	trafficIPGroupBasicConfiguration = util.AddSimpleGetAttributesToMap(d, trafficIPGroupBasicConfiguration, "", getTrafficIPGroupAttributeList("basic"))
	if v, ok := d.GetOk("ip_mapping"); ok {
		builtList, err := util.BuildListMaps(v.(*schema.Set), getTrafficIPGroupAttributeList("ip_mapping"))
		if err != nil {
			return err
		}
		trafficIPGroupBasicConfiguration["ip_mapping"] = builtList
	}
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

	trafficIPGroupPropertiesConfiguration["basic"] = trafficIPGroupBasicConfiguration
	trafficIPGroupConfiguration["properties"] = trafficIPGroupPropertiesConfiguration

	err := client.Set("traffic_ip_groups", name, trafficIPGroupConfiguration, nil)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst creating %s: %v", name, err)
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

	trafficIPGroupPropertiesConfiguration := trafficIPGroupConfiguration["properties"].(map[string]interface{})
	trafficIPGroupBasicConfiguration := trafficIPGroupPropertiesConfiguration["basic"].(map[string]interface{})
	util.SetSimpleAttributesFromMap(d, trafficIPGroupBasicConfiguration, "", getMonitorMapAttributeList("basic"))
	util.SetSimpleAttributesFromMap(d, trafficIPGroupBasicConfiguration, "", []string{"machines"})

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

	d.SetId(name)
	return resourceTrafficIPGroupRead(d, m)
}

func resourceTrafficIPGroupDelete(d *schema.ResourceData, m interface{}) error {
	return DeleteResource("traffic_ip_groups", d, m)
}
