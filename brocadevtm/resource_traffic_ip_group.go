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
			"hash_source_port": {
				Type:        schema.TypeBool,
				Description: "Whether or not the source port should be taken into account when deciding which traffic manager should handle a request.",
				Default:     false,
				Optional:    true,
			},
			"ip_assignment_mode": {
				Type:         schema.TypeString,
				Description:  "Configure how traffic IPs are assigned to traffic managers in single hosted mode",
				Computed:     true,
				Optional:     true,
				ValidateFunc: validateIPAssignmentMode,
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
				Type:         schema.TypeInt,
				Description:  "The location where the traffic IP group is based",
				Default:      0,
				Optional:     true,
				ValidateFunc: util.ValidateUnsignedInteger,
			},
			"machines": {
				Type:        schema.TypeSet,
				Description: "List of traffic managers on which to raise this traffic IP - automatically retrieved from vTM",
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"mode": {
				Type:         schema.TypeString,
				Description:  "The method used to distribute traffic IPs across machines in the cluster",
				Optional:     true,
				Computed:     true,
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
				Type:         schema.TypeString,
				Description:  "List of protocols ro be used for RHI",
				Optional:     true,
				Computed:     true,
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

func validateIPAssignmentMode(v interface{}, k string) (ws []string, errors []error) {

	assignmentMode := v.(string)
	assignmentModeOptions := regexp.MustCompile(`^(alphabetic|balanced)$`)
	if !assignmentModeOptions.MatchString(assignmentMode) {
		errors = append(errors, fmt.Errorf("%q must be one of alphabetic or balanced", k))
	}
	return
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

	enabled := d.Get("enabled").(bool)
	trafficIPGroup.Properties.Basic.Enabled = &enabled

	hashSourcePort := d.Get("hash_source_port").(bool)
	trafficIPGroup.Properties.Basic.HashSourcePort = &hashSourcePort

	if v, ok := d.GetOk("ip_assignment_mode"); ok && v != "" {
		trafficIPGroup.Properties.Basic.IPAssignmentMode = v.(string)
	}
	if v, ok := d.GetOk("ip_mapping"); ok {
		trafficIPGroup.Properties.Basic.IPMapping = buildIPMapping(v.(*schema.Set))
	}
	if v, ok := d.GetOk("ipaddresses"); ok {
		trafficIPGroup.Properties.Basic.IPAddresses = util.BuildStringListFromSet(v.(*schema.Set))
	}

	keepTogether := d.Get("keeptogether").(bool)
	trafficIPGroup.Properties.Basic.KeepTogether = &keepTogether

	location := d.Get("location").(int)
	trafficIPGroup.Properties.Basic.Location = &location

	// Allow the user to override the list of traffic managers. If not specified by the user retrieving them from the traffic manager.
	if v, ok := d.GetOk("machines"); ok {
		trafficIPGroup.Properties.Basic.Machines = util.BuildStringListFromSet(v.(*schema.Set))
	} else {
		trafficManagers, err := getTrafficManagers(m)
		if err != nil {
			return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst creating %s: %v", name, err)
		}
		trafficIPGroup.Properties.Basic.Machines = trafficManagers
	}
	if v, ok := d.GetOk("mode"); ok {
		trafficIPGroup.Properties.Basic.Mode = v.(string)
	}
	if v, ok := d.GetOk("multicast"); ok && v != "" {
		trafficIPGroup.Properties.Basic.Multicast = v.(string)
	}
	if v, ok := d.GetOk("note"); ok && v != "" {
		trafficIPGroup.Properties.Basic.Note = v.(string)
	}

	rhiBgpMetricBase := uint(d.Get("rhi_bgp_metric_base").(int))
	trafficIPGroup.Properties.Basic.RhiBgpMetricBase = &rhiBgpMetricBase

	rhiBgpPassiveMetricOffset := uint(d.Get("rhi_bgp_passive_metric_offset").(int))
	trafficIPGroup.Properties.Basic.RhiBgpPassiveMetricOffset = &rhiBgpPassiveMetricOffset

	rhiOspfv2MetricBase := uint(d.Get("rhi_ospfv2_metric_base").(int))
	trafficIPGroup.Properties.Basic.RhiOspfv2MetricBase = &rhiOspfv2MetricBase

	rhiOspfv2PassiveMetricOffset := uint(d.Get("rhi_ospfv2_passive_metric_offset").(int))
	trafficIPGroup.Properties.Basic.RhiOspfv2PassiveMetricOffset = &rhiOspfv2PassiveMetricOffset

	if v, ok := d.GetOk("rhi_protocols"); ok && v != "" {
		trafficIPGroup.Properties.Basic.RhiProtocols = v.(string)
	}
	if v, ok := d.GetOk("slaves"); ok {
		trafficIPGroup.Properties.Basic.Slaves = util.BuildStringListFromSet(v.(*schema.Set))
	}

	createAPI := trafficIpGroups.NewCreate(name, trafficIPGroup)
	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst creating %s: %v", name, createAPI.ErrorObject())
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
		return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst retrieving %s: %v", name, getAPI.ErrorObject())
	}

	trafficIPGroup = *getAPI.ResponseObject().(*trafficIpGroups.TrafficIPGroup)
	d.Set("name", name)
	d.Set("enabled", *trafficIPGroup.Properties.Basic.Enabled)
	d.Set("hash_source_port", *trafficIPGroup.Properties.Basic.HashSourcePort)
	d.Set("ip_assignment_mode", trafficIPGroup.Properties.Basic.IPAssignmentMode)
	d.Set("ip_mapping", trafficIPGroup.Properties.Basic.IPMapping)
	d.Set("ipaddresses", trafficIPGroup.Properties.Basic.IPAddresses)
	d.Set("keeptogether", *trafficIPGroup.Properties.Basic.KeepTogether)
	d.Set("location", *trafficIPGroup.Properties.Basic.Location)
	d.Set("machines", trafficIPGroup.Properties.Basic.Machines)
	d.Set("mode", trafficIPGroup.Properties.Basic.Mode)
	d.Set("multicast", trafficIPGroup.Properties.Basic.Multicast)
	d.Set("note", trafficIPGroup.Properties.Basic.Note)
	d.Set("rhi_bgp_metric_base", *trafficIPGroup.Properties.Basic.RhiBgpMetricBase)
	d.Set("rhi_bgp_passive_metric_offset", *trafficIPGroup.Properties.Basic.RhiBgpPassiveMetricOffset)
	d.Set("rhi_ospfv2_metric_base", *trafficIPGroup.Properties.Basic.RhiOspfv2MetricBase)
	d.Set("rhi_ospfv2_passive_metric_offset", *trafficIPGroup.Properties.Basic.RhiOspfv2PassiveMetricOffset)
	d.Set("rhi_protocols", trafficIPGroup.Properties.Basic.RhiProtocols)
	d.Set("slaves", trafficIPGroup.Properties.Basic.Slaves)

	return nil
}

func resourceTrafficIPGroupUpdate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*rest.Client)
	name := d.Id()
	var trafficIPGroup trafficIpGroups.TrafficIPGroup
	hasChanges := false

	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		trafficIPGroup.Properties.Basic.Enabled = &enabled
		hasChanges = true
	}
	if d.HasChange("hash_source_port") {
		hashSourcePort := d.Get("hash_source_port").(bool)
		trafficIPGroup.Properties.Basic.HashSourcePort = &hashSourcePort
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
		if v, ok := d.GetOk("ipaddresses"); ok {
			trafficIPGroup.Properties.Basic.IPAddresses = util.BuildStringListFromSet(v.(*schema.Set))
		}
		hasChanges = true
	}

	if d.HasChange("keeptogether") {
		keepTogether := d.Get("keeptogether").(bool)
		trafficIPGroup.Properties.Basic.KeepTogether = &keepTogether
		hasChanges = true
	}
	if d.HasChange("location") {
		if v, ok := d.GetOk("location"); ok {
			location := v.(int)
			trafficIPGroup.Properties.Basic.Location = &location
		}
		hasChanges = true
	}
	if d.HasChange("machines") {
		// Allow the user to override the list of traffic managers. If not specified by the user retrieve them from the traffic manager.
		if v, ok := d.GetOk("machines"); ok {
			trafficIPGroup.Properties.Basic.Machines = util.BuildStringListFromSet(v.(*schema.Set))
		} else {
			trafficManagers, err := getTrafficManagers(m)
			if err != nil {
				return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst updating %s: %v", name, err)
			}
			trafficIPGroup.Properties.Basic.Machines = trafficManagers
		}
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
	if d.HasChange("rhi_bgp_metric_base") {
		rhiBgpPassiveMetricOffset := uint(d.Get("rhi_bgp_metric_base").(int))
		trafficIPGroup.Properties.Basic.RhiBgpMetricBase = &rhiBgpPassiveMetricOffset
		hasChanges = true
	}
	if d.HasChange("rhi_bgp_passive_metric_offset") {
		rhiBgpPassiveMetricOffset := uint(d.Get("rhi_bgp_passive_metric_offset").(int))
		trafficIPGroup.Properties.Basic.RhiBgpPassiveMetricOffset = &rhiBgpPassiveMetricOffset
		hasChanges = true
	}
	if d.HasChange("rhi_ospfv2_metric_base") {
		rhiOspfv2MetricBase := uint(d.Get("rhi_ospfv2_metric_base").(int))
		trafficIPGroup.Properties.Basic.RhiOspfv2MetricBase = &rhiOspfv2MetricBase
		hasChanges = true
	}
	if d.HasChange("rhi_ospfv2_passive_metric_offset") {
		rhiOspfv2PassiveMetricOffset := uint(d.Get("rhi_ospfv2_passive_metric_offset").(int))
		trafficIPGroup.Properties.Basic.RhiOspfv2PassiveMetricOffset = &rhiOspfv2PassiveMetricOffset
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
		updateAPI := trafficIpGroups.NewUpdate(name, trafficIPGroup)
		err := vtmClient.Do(updateAPI)
		if err != nil {
			return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst updating %s: %v", name, updateAPI.ErrorObject())
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
		return fmt.Errorf("BrocadeVTM Traffic IP Group error whilst deleting %s: %v", name, deleteAPI.ErrorObject())
	}
	d.SetId("")
	return nil
}
