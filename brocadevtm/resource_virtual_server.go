package brocadevtm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
	"regexp"
	"github.com/sky-uk/go-brocade-vtm"
	"github.com/sky-uk/go-brocade-vtm/api/virtualserver"
	"net/http"
	"encoding/json"
)

func resourceVirtualServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualServerCreate,
		Read:   resourceVirtualServerRead,
		Update: resourceVirtualServerUpdate,
		Delete: resourceVirtualServerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type: schema.TypeString,
				Description: "Name of the virtual server",
				Required: true,
			},
			"enabled": {
				Type: schema.TypeBool,
				Description: "Whether the virtual server should be enabled",
				Optional: true,
				Computed: true,
			},
			"listen_on_any": {
				Type: schema.TypeBool,
				Description: "Whether the virtual server should listen on any",
				Optional: true,
				Computed: true,
			},
			"listen_traffic_ips": {
				Type:        schema.TypeList,
				Description: "List of traffic IPs to listen on",
				Optional: true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"pool": {
				Type: schema.TypeString,
				Description: "Name of the pool to use with the virtual server",
				Required: true,
			},
			"port":{
				Type: schema.TypeInt,
				Description: "Port the virtual server should listen on",
				Required: true,
				ValidateFunc: validateVirtualServerUnsignedInteger,
			},
			"protocol": {
				Type: schema.TypeString,
				Description: "Protocol to use with the virtual server",
				Optional: true,
				Computed: true,
				ValidateFunc: validateVirtualServerProtocol,
			},
			"request_rules": {
				Type: schema.TypeList,
				Description: "A list of request rules",
				Optional: true,
			},
			"ssl_decrypt": {
				Type: schema.TypeBool,
				Description: "Whether to enable or disable SSL",
				Optional: true,
				Computed: true,
			},
			"connection_keepalive": {
				Type: schema.TypeBool,
				Description: "Whether to enable keepalive for remote clients",
				Optional: true,
				Computed: true,
			},
			"connection_keepalive_timeout": {
				Type: schema.TypeInt,
				Description: "Keepalive timeout for idle connections",
				Optional: true,
				Computed: true,
				ValidateFunc: validateVirtualServerUnsignedInteger,
			},
			"connection_max_client_buffer": {
				Type: schema.TypeInt,
				Description: "Max memory in bytes for stored client data",
				Optional: true,
				Computed: true,
				ValidateFunc: validateVirtualServerUnsignedInteger,
			},
			"connection_max_server_buffer": {
				Type: schema.TypeInt,
				Description: "Max memory in bytes for stored server data",
				Optional: true,
				Computed: true,
				ValidateFunc: validateVirtualServerUnsignedInteger,
			},
			"connection_max_transaction_duration": {
				Type: schema.TypeInt,
				Description: "Max amount of time a transaction can take",
				Optional: true,
				Computed: true,
				ValidateFunc: validateVirtualServerUnsignedInteger,
			},
			"connection_server_first_banner": {
				Type: schema.TypeString,
				Description: "Banner to send for server first protocols",
				Optional: true,
			},
			"connection_timeout": {
				Type: schema.TypeInt,
				Description: "Time to wait before closing a connection when no additional data has been sent",
				Optional: true,
				Computed: true,
				ValidateFunc: validateVirtualServerUnsignedInteger,
			},
			"ssl_server_cert_default": {
				Type: schema.TypeString,
				Description: "Default SSL certificate",
				Optional: true,
			},
			"ssl_server_cert_host_mapping_host": {
				Type: schema.TypeString,
				Description: "Which host the SSL certificate refers to",
				Optional: true,
			},
			"ssl_server_cert_host_mapping_alt_certificates": {
				Type: schema.TypeList,
				Description: "SSL server certificates for a particular destination IP",
				Optional: true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ssl_server_cert_host_mapping_certificate": {
				Type: schema.TypeString,
				Description: "The SSL server certificate for a particular destination",
				Optional: true,
			},
			"ssl_support_ssl2": {
				Type: schema.TypeString,
				Description: "Whether or not SSLv2 is enabled for this virtual server",
				Optional: true,
				Computed: true,
				ValidateFunc: validateVirtualServerUseSSLSupport,
			},
			"ssl_support_ssl3": {
				Type: schema.TypeString,
				Description: "Whether or not SSLv3 is enabled for this virtual server",
				Optional: true,
				Computed: true,
				ValidateFunc: validateVirtualServerUseSSLSupport,
			},
			"ssl_support_tls1": {
				Type: schema.TypeString,
				Description: "Whether or not TLSv1.0 is enabled for this virtual server",
				Optional: true,
				Computed: true,
				ValidateFunc: validateVirtualServerUseSSLSupport,
			},
			"ssl_support_tls1_1": {
				Type: schema.TypeString,
				Description: "Whether or not TLSv1.1 is enabled for this virtual server",
				Optional: true,
				Computed: true,
				ValidateFunc: validateVirtualServerUseSSLSupport,
			},
			"ssl_support_tls1_2": {
				Type: schema.TypeString,
				Description: "Whether or not TLSv1.2 is enabled for this virtual server",
				Optional: true,
				Computed: true,
				ValidateFunc: validateVirtualServerUseSSLSupport,
			},
		},
	}
}

func validateVirtualServerUnsignedInteger(v interface{}, k string) (ws []string, errors []error) {
	ttl := v.(int)
	if ttl < 0 {
		errors = append(errors, fmt.Errorf("%q can't be negative", k))
	}
	return
}

func validateVirtualServerProtocol(v interface{}, k string) (ws []string, errors []error) {
	protocol := v.(string)
	protocolOptions := regexp.MustCompile(`^(client_first|dns|dns_tcp|ftp|http|https|imaps|imapv2|imapv3|imapv4|ldap|ldaps|pop3|pop3s|rtsp|server_first|siptcp|sipudp|smtp|ssl|stream|telnet|udp|udpstreaming)$`)
	if !protocolOptions.MatchString(protocol) {
		errors = append(errors, fmt.Errorf("%q must be one of client_first, dns, dns_tcp, ftp, http, https, imaps, imapv2, imapv3, imapv4, ldap, ldaps, pop3, pop3s, rtsp, server_first, siptcp, sipudp, smtp, ssl, stream, telnet, udp or udpstreaming", k))
	}
	return
}

func validateVirtualServerUseSSLSupport(v interface{}, k string) (ws []string, errors []error) {
	sslUseSSLSupport := v.(string)
	sslUseSSLSupportOptions := regexp.MustCompile(`^(use_default|disabled|enabled)$`)
	if !sslUseSSLSupportOptions.MatchString(sslUseSSLSupport) {
		errors = append(errors, fmt.Errorf("%q must be one of use_default, disabled or enabled", k))
	}
	return
}

func resourceVirtualServerCreate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*brocadevtm.VTMClient)
	var virtualServerName string
	var virtualServer virtualserver.VirtualServer

	if v, ok := d.GetOk("name"); ok && v != "" {
		virtualServerName = v.(string)
	}
	if v, ok := d.GetOk("pool"); ok && v != "" {
		virtualServer.Properties.Basic.Pool = v.(string)
	}
	if v, ok := d.GetOk("port"); ok && v != "" {
		virtualServerPort := v.(int)
		virtualServer.Properties.Basic.Port = uint(virtualServerPort)
	}

	createAPI := virtualserver.NewCreate(virtualServerName, virtualServer)
	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf("Brocade vTM Virtual Server Create failed for %s with error: %+v", virtualServerName, err)
	}
	if createAPI.StatusCode() != http.StatusCreated {
		return fmt.Errorf("Brocade vTM Virtual Server Create failed for %s with http status code != 201", virtualServerName)
	}
	d.SetId(virtualServerName)
	return resourceVirtualServerRead(d, m)
}

func resourceVirtualServerRead(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*brocadevtm.VTMClient)
	var virtualServerName string
	var virtualServer virtualserver.VirtualServer

	if v, ok := d.GetOk("name"); ok && v != "" {
		virtualServerName = v.(string)
	}

	getSingleAPI := virtualserver.NewGetSingle(virtualServerName)
	err := vtmClient.Do(getSingleAPI)
	if err != nil {
		return fmt.Errorf("Brocade vTM Virtual Server Read failed for %s with error: %+v", virtualServerName, err)
	}
	if getSingleAPI.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	virtualServer = *getSingleAPI.GetResponse()
	d.SetId(virtualServerName)
	d.Set("pool", virtualServer.Properties.Basic.Pool)
	d.Set("port", virtualServer.Properties.Basic.Port)
	return nil
}

func resourceVirtualServerUpdate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*brocadevtm.VTMClient)
	var virtualServerName string
	var virtualServer, responseVirtualServer virtualserver.VirtualServer
	hasChanges := false

	if v, ok := d.GetOk("name"); ok && v != "" {
		virtualServerName = v.(string)
	}
	if d.HasChange("pool") {
		if v, ok := d.GetOk("pool"); ok && v != "" {
			virtualServer.Properties.Basic.Pool = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("port") {
		if v, ok := d.GetOk("port"); ok && v != "" {
			virtualServerPort := v.(int)
			virtualServer.Properties.Basic.Port = uint(virtualServerPort)
		}
		hasChanges = true
	}

	if hasChanges {
		updateAPI := virtualserver.NewUpdate(virtualServerName, virtualServer)
		err := vtmClient.Do(updateAPI)
		if err != nil {
			return fmt.Errorf("Brocade vTM Virtual Server update failed for %s", virtualServerName)
		}
		responseCode := updateAPI.StatusCode()
		if responseCode != http.StatusOK {
			return fmt.Errorf("Brocade vTM Virtual Server update failed for %s with invalid response code %d", virtualServerName, responseCode)
		}

		response := updateAPI.GetResponse()
		jsonErr := json.Unmarshal(response, &responseVirtualServer)
		if jsonErr != nil {
			return fmt.Errorf("Brocade vTM Virtual Server update faild for %s while unmarshalling JSON response - raw response: %s", virtualServerName, response)
		}

		d.SetId(virtualServerName)
		d.Set("pool", responseVirtualServer.Properties.Basic.Pool)
		d.Set("port", responseVirtualServer.Properties.Basic.Port)
	}
	return resourceVirtualServerRead(d, m)
}

func resourceVirtualServerDelete(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*brocadevtm.VTMClient)
	var virtualServerName string

	if v, ok := d.GetOk("name"); ok && v != "" {
		virtualServerName = v.(string)
	}
	getVirtualServer := virtualserver.NewGetSingle(virtualServerName)
	err := vtmClient.Do(getVirtualServer)
	if err != nil {
		return fmt.Errorf("Brocade vTM Virtual Server delete failed for %s - error: %+v", virtualServerName, err)
	}
	if getVirtualServer.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	deleteAPI := virtualserver.NewDelete(virtualServerName)
	err = vtmClient.Do(deleteAPI)
	if err != nil {
		return fmt.Errorf("Brocade vTM Virtual Server delete failed for %s - error: %+v", virtualServerName, err)
	}
	responseCode := deleteAPI.StatusCode()
	if responseCode != http.StatusNoContent {
		return fmt.Errorf("Brocade vTM Virtual Server delete returned an invalid http response code %d for %s", responseCode, virtualServerName)
	}

	d.SetId("")
	return nil
}