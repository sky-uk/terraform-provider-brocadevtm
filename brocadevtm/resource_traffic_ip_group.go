package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm/api/traffic_ip_group"
	"github.com/sky-uk/go-brocade-vtm/api/traffic_ip_group_manager"
	"github.com/sky-uk/go-rest-api"
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
			"hashsourceport": {
				Type:        schema.TypeBool,
				Description: "Whether or not the source port should be taken into account when deciding which traffic manager should handle a request.",
				Default:     false,
				Optional:    true,
			},
			"ip_assignment_mode": {
				Type:        schema.TypeString,
				Description: "Configure how traffic IPs are assigned to traffic managers in single hosted mode",
				Default:     "balanced",
				Optional:    true,
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
				Type:        schema.TypeList,
				Description: "List of IP addresses to raise on the traffic managers",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"keeptogether": {
				Type:        schema.TypeBool,
				Description: "Whether or not all traffic IPs are raised on a single traffic manager",
				Optional:    true,
				Default:     false,
			},
			"location": {
				Type:        schema.TypeInt,
				Description: "The location where the traffic IP group is based",
				Optional:    true,
			},
			"machines": {
				Type:        schema.TypeList,
				Description: "List of traffic managers on which to raise this traffic IP - automatically retrieved from vTM",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"mode": {
				Type:         schema.TypeString,
				Description:  "The method used to distribute traffic IPs across machines in the cluster",
				Optional:     true,
				Default:      "singlehosted",
				ValidateFunc: validateTrafficIPGroupMode,
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
				Type:        schema.TypeInt,
				Description: "Base BGP routing metric",
				Default:     10,
				Optional:    true,
			},
			"rhi_bgp_passive_metric_offset": {
				Type:        schema.TypeInt,
				Description: "BGP routing metric offset",
				Default:     10,
				Optional:    true,
			},
			"rhi_ospfv2_metric_base": {
				Type:        schema.TypeInt,
				Description: "OSPFv2 routing metric",
				Default:     10,
				Optional:    true,
			},
			"rhi_ospfv2_passive_metric_offset": {
				Type:        schema.TypeInt,
				Description: "OSPFv2 routing metric offset",
				Default:     10,
				Optional:    true,
			},
			"rhi_protocols": {
				Type:         schema.TypeString,
				Description:  "List of protocols ro be used for RHI",
				Default:      "ospf",
				Optional:     true,
				ValidateFunc: validateRHIProtocols,
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

func validateRHIProtocols(v interface{}, k string) (ws []string, errors []error) {
	protocol := v.(string)
	protocolOption := regexp.MustCompile(`^(ospf|bgp)$`)
	if !protocolOption.MatchString(protocol) {
		errors = append(errors, fmt.Errorf("%q must be one of ospf or bgp", k))
	}
	return
}

func validateTrafficIPGroupMode(v interface{}, k string) (ws []string, errors []error) {
	mode := v.(string)
	modeOptions := regexp.MustCompile(`^(singlehosted|ec2elastic|ec2vpcelastic|ec2vpcprivate|multihosted|rhi)$`)
	if !modeOptions.MatchString(mode) {
		errors = append(errors, fmt.Errorf("%q must be one of singlehosted, ec2elastic, ec2vpcelastic, ec2vpcprivate, multihosted or rhi", k))
	}
	return
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
	vtmClient := m.(*rest.Client)
	getTrafficManagersAPI := trafficIpGroupManager.NewGetAll()
	var trafficManagers []string

	err := vtmClient.Do(getTrafficManagersAPI)
	if err != nil {
		return trafficManagers, fmt.Errorf("BrocadeVTM Traffic Managers error whilst retrieving the list of Traffic Managers: %v", err)
	}

	response := getTrafficManagersAPI.ResponseObject().(*trafficIpGroupManager.TrafficManagerChildren)
	for _, trafficManager := range response.Children {
		trafficManagers = append(trafficManagers, trafficManager.Name)
	}
	return trafficManagers, nil
}

func buildIPAddresses(ipAddresses interface{}) []string {
	ipAddressList := make([]string, len(ipAddresses.([]interface{})))
	for idx, ipAddress := range ipAddresses.([]interface{}) {
		ipAddressList[idx] = ipAddress.(string)
	}
	return ipAddressList
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

func resourceTrafficIPGroupCreate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	var trafficIPGroup trafficIpGroups.TrafficIPGroup
	var name string

	if v, ok := d.GetOk("name"); ok && v != "" {
		name = v.(string)
	}
	trafficIPGroup.Properties.Basic.Enabled = d.Get("enabled").(bool)
	trafficIPGroup.Properties.Basic.HashSourcePort = d.Get("hashsourceport").(bool)
	if v, ok := d.GetOk("ip_assignment_mode"); ok && v != "" {
		trafficIPGroup.Properties.Basic.IPAssignmentMode = v.(string)
	}
	if v, ok := d.GetOk("ip_mapping"); ok {
		trafficIPGroup.Properties.Basic.IPMapping = buildIPMapping(v.(*schema.Set))
	}
	if v, ok := d.GetOk("ipaddresses"); ok && v != "" {
		trafficIPGroup.Properties.Basic.IPAddresses = buildIPAddresses(v)
	}
	trafficIPGroup.Properties.Basic.KeepTogether = d.Get("keeptogether").(bool)
	if v, ok := d.GetOk("location"); ok {
		location := v.(int)
		trafficIPGroup.Properties.Basic.Location = &location
	}
	// Retrieve the list of Brocade vTM traffic managers and assign it to Machines
	trafficManagers, err := getTrafficManagers(m)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst creating %s: %v", name, err)
	}
	trafficIPGroup.Properties.Basic.Machines = trafficManagers
	if v, ok := d.GetOk("mode"); ok {
		trafficIPGroup.Properties.Basic.Mode = v.(string)
	}
	if v, ok := d.GetOk("multicast"); ok && v != "" {
		trafficIPGroup.Properties.Basic.Multicast = v.(string)
	}
	if v, ok := d.GetOk("note"); ok && v != "" {
		trafficIPGroup.Properties.Basic.Note = v.(string)
	}
	trafficIPGroup.Properties.Basic.RhiBgpMetricBase = uint(d.Get("rhi_bgp_passive_metric_offset").(int))
	trafficIPGroup.Properties.Basic.RhiBgpPassiveMetricOffset = uint(d.Get("rhi_bgp_passive_metric_offset").(int))
	trafficIPGroup.Properties.Basic.RhiOspfv2MetricBase = uint(d.Get("rhi_ospfv2_metric_base").(int))
	trafficIPGroup.Properties.Basic.RhiOspfv2PassiveMetricOffset = uint(d.Get("rhi_ospfv2_passive_metric_offset").(int))
	if v, ok := d.GetOk("rhi_protocols"); ok && v != "" {
		trafficIPGroup.Properties.Basic.RhiProtocols = v.(string)
	}
	if v, ok := d.GetOk("slaves"); ok {
		trafficIPGroup.Properties.Basic.Slaves = util.BuildStringListFromSet(v.(*schema.Set))
	}

	createAPI := trafficIpGroups.NewCreate(name, trafficIPGroup)
	err = vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst creating %s: %v", name, err)
	}
	d.SetId(name)
	return resourceTrafficIPGroupRead(d, m)
}

func resourceTrafficIPGroupRead(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	var trafficIPGroup trafficIpGroups.TrafficIPGroup
	name := d.Id()

	getAPI := trafficIpGroups.NewGet(name)
	err := vtmClient.Do(getAPI)
	if err != nil {
		if getAPI.StatusCode() == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst retrieving %s: %v", name, err)
	}

	trafficIPGroup = *getAPI.ResponseObject().(*trafficIpGroups.TrafficIPGroup)
	d.Set("name", name)
	d.Set("enabled", trafficIPGroup.Properties.Basic.Enabled)
	d.Set("hashsourceport", trafficIPGroup.Properties.Basic.HashSourcePort)
	d.Set("ip_assignment_mode", trafficIPGroup.Properties.Basic.IPAssignmentMode)
	d.Set("ip_mapping", trafficIPGroup.Properties.Basic.IPMapping)
	d.Set("ipaddresses", trafficIPGroup.Properties.Basic.IPAddresses)
	d.Set("keeptogether", trafficIPGroup.Properties.Basic.KeepTogether)
	d.Set("location", *trafficIPGroup.Properties.Basic.Location)
	d.Set("machines", trafficIPGroup.Properties.Basic.Machines)
	d.Set("mode", trafficIPGroup.Properties.Basic.Mode)
	d.Set("multicast", trafficIPGroup.Properties.Basic.Multicast)
	d.Set("note", trafficIPGroup.Properties.Basic.Note)
	d.Set("rhi_bgp_metric_base", trafficIPGroup.Properties.Basic.RhiBgpMetricBase)
	d.Set("rhi_bgp_passive_metric_offset", trafficIPGroup.Properties.Basic.RhiBgpPassiveMetricOffset)
	d.Set("rhi_ospfv2_metric_base", trafficIPGroup.Properties.Basic.RhiOspfv2MetricBase)
	d.Set("rhi_ospfv2_passive_metric_offset", trafficIPGroup.Properties.Basic.RhiOspfv2PassiveMetricOffset)
	d.Set("rhi_protocols", trafficIPGroup.Properties.Basic.RhiProtocols)
	d.Set("slaves", trafficIPGroup.Properties.Basic.Slaves)

	return nil
}

func resourceTrafficIPGroupUpdate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	name := d.Id()
	var trafficIPGroup trafficIpGroups.TrafficIPGroup
	hasChanges := false

	trafficIPGroup.Properties.Basic.Enabled = d.Get("enabled").(bool)
	if d.HasChange("enabled") {
		hasChanges = true
	}
	trafficIPGroup.Properties.Basic.HashSourcePort = d.Get("hashsourceport").(bool)
	if d.HasChange("hashsourceport") {
		hasChanges = true
	}
	if d.HasChange("ip_assignment_mode") {
		if v, ok := d.GetOk("ip_assignment_mode"); ok && v != "" {
			trafficIPGroup.Properties.Basic.IPAssignmentMode = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("ip_mapping") {
		if v, ok := d.GetOk("ip_mapping"); ok {
			trafficIPGroup.Properties.Basic.IPMapping = buildIPMapping(v.(*schema.Set))
		}
		hasChanges = true
	}
	if d.HasChange("ipaddresses") {
		if v, ok := d.GetOk("ipaddresses"); ok && v != "" {
			trafficIPGroup.Properties.Basic.IPAddresses = buildIPAddresses(v)
		}
		hasChanges = true
	}
	trafficIPGroup.Properties.Basic.KeepTogether = d.Get("keeptogether").(bool)
	if d.HasChange("keeptogether") {
		hasChanges = true
	}
	if d.HasChange("location") {
		if v, ok := d.GetOk("location"); ok {
			location := v.(int)
			trafficIPGroup.Properties.Basic.Location = &location
		}
		hasChanges = true
	}
	if d.HasChange("mode") {
		if v, ok := d.GetOk("mode"); ok && v != "" {
			trafficIPGroup.Properties.Basic.Mode = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("multicast") {
		if v, ok := d.GetOk("multicast"); ok && v != "" {
			trafficIPGroup.Properties.Basic.Multicast = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("note") {
		if v, ok := d.GetOk("note"); ok && v != "" {
			trafficIPGroup.Properties.Basic.Note = v.(string)
		}
		hasChanges = true
	}
	trafficIPGroup.Properties.Basic.RhiBgpMetricBase = uint(d.Get("rhi_bgp_passive_metric_offset").(int))
	if d.HasChange("rhi_bgp_passive_metric_offset") {
		hasChanges = true
	}
	trafficIPGroup.Properties.Basic.RhiBgpPassiveMetricOffset = uint(d.Get("rhi_bgp_passive_metric_offset").(int))
	if d.HasChange("rhi_bgp_passive_metric_offset") {
		hasChanges = true
	}
	trafficIPGroup.Properties.Basic.RhiOspfv2MetricBase = uint(d.Get("rhi_ospfv2_metric_base").(int))
	if d.HasChange("rhi_ospfv2_metric_base") {
		hasChanges = true
	}
	trafficIPGroup.Properties.Basic.RhiOspfv2PassiveMetricOffset = uint(d.Get("rhi_ospfv2_passive_metric_offset").(int))
	if d.HasChange("rhi_ospfv2_passive_metric_offset") {
		hasChanges = true
	}
	if d.HasChange("rhi_protocols") {
		if v, ok := d.GetOk("rhi_protocols"); ok && v != "" {
			trafficIPGroup.Properties.Basic.RhiProtocols = v.(string)
		}
	}
	if d.HasChange("slaves") {
		if v, ok := d.GetOk("slaves"); ok {
			trafficIPGroup.Properties.Basic.Slaves = util.BuildStringListFromSet(v.(*schema.Set))
		}
	}

	if hasChanges {
		// On all updates refresh the list of traffic managers
		trafficManagers, err := getTrafficManagers(m)
		if err != nil {
			return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst updating %s: %v", name, err)
		}
		trafficIPGroup.Properties.Basic.Machines = trafficManagers

		updateAPI := trafficIpGroups.NewUpdate(name, trafficIPGroup)
		err = vtmClient.Do(updateAPI)
		if err != nil {
			return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst updating %s: %v", name, err)
		}
		d.SetId(name)
	}
	return resourceTrafficIPGroupRead(d, m)
}

func resourceTrafficIPGroupDelete(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	name := d.Id()

	deleteAPI := trafficIpGroups.NewDelete(name)
	err := vtmClient.Do(deleteAPI)
	if err != nil && deleteAPI.StatusCode() != http.StatusNotFound {
		return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst deleting %s: %v", name, err)
	}
	d.SetId("")
	return nil
}
