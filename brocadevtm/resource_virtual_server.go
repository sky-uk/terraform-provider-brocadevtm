package brocadevtm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/sky-uk/go-brocade-vtm"
	"github.com/sky-uk/go-brocade-vtm/api/virtualserver"
	"net/http"
	"regexp"
)

func resourceVirtualServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualServerCreate,
		Read:   resourceVirtualServerRead,
		Update: resourceVirtualServerUpdate,
		Delete: resourceVirtualServerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the virtual server",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether the virtual server should be enabled",
				Optional:    true,
				Computed:    true,
			},
			"listen_on_any": {
				Type:        schema.TypeBool,
				Description: "Whether the virtual server should listen on any",
				Optional:    true,
				Computed:    true,
			},
			"listen_traffic_ips": {
				Type:        schema.TypeList,
				Description: "List of traffic IPs to listen on",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"pool": {
				Type:        schema.TypeString,
				Description: "Name of the pool to use with the virtual server",
				Required:    true,
			},
			"port": {
				Type:         schema.TypeInt,
				Description:  "Port the virtual server should listen on",
				Required:     true,
				ValidateFunc: validateVirtualServerUnsignedInteger,
			},
			"protocol": {
				Type:         schema.TypeString,
				Description:  "Protocol to use with the virtual server",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateVirtualServerProtocol,
			},
			"request_rules": {
				Type:        schema.TypeList,
				Description: "A list of request rules",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ssl_decrypt": {
				Type:        schema.TypeBool,
				Description: "Whether to enable or disable SSL",
				Optional:    true,
				Computed:    true,
			},
			"connection_keepalive": {
				Type:        schema.TypeBool,
				Description: "Whether to enable keepalive for remote clients",
				Optional:    true,
				Computed:    true,
			},
			"connection_keepalive_timeout": {
				Type:         schema.TypeInt,
				Description:  "Keepalive timeout for idle connections",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateVirtualServerUnsignedInteger,
			},
			"connection_max_client_buffer": {
				Type:         schema.TypeInt,
				Description:  "Max memory in bytes for stored client data",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateVirtualServerUnsignedInteger,
			},
			"connection_max_server_buffer": {
				Type:         schema.TypeInt,
				Description:  "Max memory in bytes for stored server data",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateVirtualServerUnsignedInteger,
			},
			"connection_max_transaction_duration": {
				Type:         schema.TypeInt,
				Description:  "Max amount of time a transaction can take",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateVirtualServerUnsignedInteger,
			},
			"connection_server_first_banner": {
				Type:        schema.TypeString,
				Description: "Banner to send for server first protocols",
				Optional:    true,
			},
			"connection_timeout": {
				Type:         schema.TypeInt,
				Description:  "Time to wait before closing a connection when no additional data has been sent",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateVirtualServerUnsignedInteger,
			},
			"ssl_server_cert_default": {
				Type:        schema.TypeString,
				Description: "Default SSL certificate",
				Optional:    true,
			},
			"ssl_support_ssl2": {
				Type:         schema.TypeString,
				Description:  "Whether or not SSLv2 is enabled for this virtual server",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateVirtualServerUseSSLSupport,
			},
			"ssl_support_ssl3": {
				Type:         schema.TypeString,
				Description:  "Whether or not SSLv3 is enabled for this virtual server",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateVirtualServerUseSSLSupport,
			},
			"ssl_support_tls1": {
				Type:         schema.TypeString,
				Description:  "Whether or not TLSv1.0 is enabled for this virtual server",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateVirtualServerUseSSLSupport,
			},
			"ssl_support_tls1_1": {
				Type:         schema.TypeString,
				Description:  "Whether or not TLSv1.1 is enabled for this virtual server",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateVirtualServerUseSSLSupport,
			},
			"ssl_support_tls1_2": {
				Type:         schema.TypeString,
				Description:  "Whether or not TLSv1.2 is enabled for this virtual server",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateVirtualServerUseSSLSupport,
			},
			"ssl_server_cert_host_mapping": {
				Type:        schema.TypeList,
				Description: "Host specific SSL server certificate mappings",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ssl_server_cert_host": {
							Type:        schema.TypeString,
							Description: "Which host the SSL certificate refers to",
							Optional:    true,
						},
						"ssl_server_alt_certs": {
							Type:        schema.TypeList,
							Description: "SSL server certificates for a particular destination IP",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ssl_server_cert": {
							Type:        schema.TypeString,
							Description: "The SSL server certificate for a particular destination",
							Optional:    true,
						},
					},
				},
			},
			"ocsp_enable": {
				Type:        schema.TypeBool,
				Description: "Whether or not the traffic manager should use OCSP to check the revocation status of client certificates",
				Optional:    true,
				Computed:    true,
			},
			"ocsp_issuers": {
				Type:        schema.TypeList,
				Description: "A table of certificate issuer specific OCSP settings",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"issuer": {
							Type:        schema.TypeString,
							Description: "The name of an issuer",
							Optional:    true,
							Computed:    true,
						},
						"aia": {
							Type:        schema.TypeBool,
							Description: "Whether the traffic manager should use AIA information",
							Optional:    true,
							Computed:    true,
						},
						"nonce": {
							Type:         schema.TypeString,
							Description:  "How to use the OCSP nonce extension, which protects against OCSP replay attacks",
							Optional:     true,
							ValidateFunc: validateVirtualServerOCSPNonce,
						},
						"required": {
							Type:         schema.TypeString,
							Description:  "Whether we should do an OCSP check for this issuer",
							Optional:     true,
							ValidateFunc: validateVirtualServerOCSPRequired,
						},
						"responder_cert": {
							Type:        schema.TypeString,
							Description: "The expected responder certificate",
							Optional:    true,
						},
						"signer": {
							Type:        schema.TypeString,
							Description: "The certificate with which to sign the request",
							Optional:    true,
						},
						"url": {
							Type:        schema.TypeString,
							Description: "Which OCSP responders this virtual server should use to verify client certificates",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func validateVirtualServerOCSPRequired(v interface{}, k string) (ws []string, errors []error) {
	ocspRequired := v.(string)
	ocspRequiredOptions := regexp.MustCompile(`^(none|optional|strict)$`)
	if !ocspRequiredOptions.MatchString(ocspRequired) {
		errors = append(errors, fmt.Errorf("%q must be one of none, optional, strict", k))
	}
	return
}

func validateVirtualServerOCSPNonce(v interface{}, k string) (ws []string, errors []error) {
	nonce := v.(string)
	nonceOptions := regexp.MustCompile(`^(off|on|strict)$`)
	if !nonceOptions.MatchString(nonce) {
		errors = append(errors, fmt.Errorf("%q must be one of off, on or strict", k))
	}
	return
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

func buildStringList(strings interface{}) []string {
	stringList := make([]string, len(strings.([]interface{})))
	for idx, stringValue := range strings.([]interface{}) {
		stringList[idx] = stringValue.(string)
	}
	return stringList
}

func buildSSLCertMapping(sslCertMapping []interface{}) []virtualserver.CertItem {
	certItemList := make([]virtualserver.CertItem, len(sslCertMapping))
	var certItem virtualserver.CertItem

	for idx, value := range sslCertMapping {
		mappingItem := value.(map[string]interface{})
		certItem.Host = mappingItem["ssl_server_cert_host"].(string)
		certItem.Certificate = mappingItem["ssl_server_cert"].(string)
		certItem.AltCertificates = buildStringList(mappingItem["ssl_server_alt_certs"])
		certItemList[idx] = certItem
	}
	return certItemList
}

func buildSSLOCSPIssuers(ocspIssuers []interface{}) []virtualserver.OCSPIssuer {
	ocspIssuerList := make([]virtualserver.OCSPIssuer, len(ocspIssuers))
	var ocspIssuer virtualserver.OCSPIssuer

	for idx, value := range ocspIssuers {
		issuerItem := value.(map[string]interface{})
		ocspIssuer.Issuer = issuerItem["issuer"].(string)
		ocspIssuerAIA := issuerItem["aia"].(bool)
		ocspIssuer.AIA = &ocspIssuerAIA
		ocspIssuer.Nonce = issuerItem["nonce"].(string)
		ocspIssuer.Required = issuerItem["required"].(string)
		ocspIssuer.ResponderCert = issuerItem["responder_cert"].(string)
		ocspIssuer.Signer = issuerItem["signer"].(string)
		ocspIssuer.URL = issuerItem["url"].(string)
		ocspIssuerList[idx] = ocspIssuer
	}
	return ocspIssuerList
}

func resourceVirtualServerCreate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*brocadevtm.VTMClient)
	var virtualServerName string
	var virtualServer virtualserver.VirtualServer

	if v, ok := d.GetOk("name"); ok && v != "" {
		virtualServerName = v.(string)
	}
	if v, ok := d.GetOk("enabled"); ok {
		virtualServerEnabled := v.(bool)
		virtualServer.Properties.Basic.Enabled = &virtualServerEnabled
	}
	if v, _ := d.GetOk("listen_on_any"); v != "" {
		virtualServerListenAny := v.(bool)
		virtualServer.Properties.Basic.ListenOnAny = &virtualServerListenAny
	}
	if v, ok := d.GetOk("listen_traffic_ips"); v != "" && ok {
		virtualServer.Properties.Basic.ListenOnTrafficIps = buildStringList(v)
	}
	if v, ok := d.GetOk("pool"); ok && v != "" {
		virtualServer.Properties.Basic.Pool = v.(string)
	}
	if v, ok := d.GetOk("port"); ok && v != "" {
		virtualServerPort := v.(int)
		virtualServer.Properties.Basic.Port = uint(virtualServerPort)
	}
	if v, ok := d.GetOk("protocol"); ok && v != "" {
		virtualServer.Properties.Basic.Protocol = v.(string)
	}
	if v, ok := d.GetOk("request_rules"); ok && v != "" {
		virtualServer.Properties.Basic.RequestRules = buildStringList(v)
	}
	if v, ok := d.GetOk("ssl_decrypt"); ok {
		virtualServerSSLDeCrypt := v.(bool)
		virtualServer.Properties.Basic.SslDecrypt = &virtualServerSSLDeCrypt
	}
	if v, _ := d.GetOk("connection_keepalive"); v != "" {
		virtalServerConnectionKeepalive := v.(bool)
		virtualServer.Properties.Connection.Keepalive = &virtalServerConnectionKeepalive
	}
	if v, ok := d.GetOk("connection_keepalive_timeout"); ok {
		virtualServerConnectionKeepaliveTimeout := v.(int)
		virtualServer.Properties.Connection.KeepaliveTimeout = uint(virtualServerConnectionKeepaliveTimeout)
	}
	if v, ok := d.GetOk("connection_max_client_buffer"); ok {
		virtualServerConnectionMaxClientBuffer := v.(int)
		virtualServer.Properties.Connection.MaxClientBuffer = uint(virtualServerConnectionMaxClientBuffer)
	}
	if v, ok := d.GetOk("connection_max_server_buffer"); ok {
		virtualServerConnectionMaxServerBuffer := v.(int)
		virtualServer.Properties.Connection.MaxServerBuffer = uint(virtualServerConnectionMaxServerBuffer)
	}
	if v, ok := d.GetOk("connection_max_transaction_duration"); ok {
		virtualServerConnectionMaxTransActionDuration := v.(int)
		virtualServer.Properties.Connection.MaxTransactionDuration = uint(virtualServerConnectionMaxTransActionDuration)
	}
	if v, ok := d.GetOk("connection_server_first_banner"); ok && v != "" {
		virtualServer.Properties.Connection.ServerFirstBanner = v.(string)
	}
	if v, ok := d.GetOk("connection_timeout"); ok {
		virtualServerConnectionTimeout := v.(int)
		virtualServer.Properties.Connection.Timeout = uint(virtualServerConnectionTimeout)
	}
	if v, ok := d.GetOk("ssl_server_cert_default"); ok && v != "" {
		virtualServer.Properties.Ssl.ServerCertDefault = v.(string)
	}
	if v, ok := d.GetOk("ssl_support_ssl2"); ok && v != "" {
		virtualServer.Properties.Ssl.SslSupportSsl2 = v.(string)
	}
	if v, ok := d.GetOk("ssl_support_ssl3"); ok && v != "" {
		virtualServer.Properties.Ssl.SslSupportSsl3 = v.(string)
	}
	if v, ok := d.GetOk("ssl_support_tls1"); ok && v != "" {
		virtualServer.Properties.Ssl.SslSupportTLS1 = v.(string)
	}
	if v, ok := d.GetOk("ssl_support_tls1_1"); ok && v != "" {
		virtualServer.Properties.Ssl.SslSupportTLS1_1 = v.(string)
	}
	if v, ok := d.GetOk("ssl_support_tls1_2"); ok && v != "" {
		virtualServer.Properties.Ssl.SslSupportTLS1_2 = v.(string)
	}
	if v, ok := d.GetOk("ssl_server_cert_host_mapping"); ok && v != "" {
		virtualServer.Properties.Ssl.ServerCertHostMap = buildSSLCertMapping(v.([]interface{}))
	}
	if v, ok := d.GetOk("ocsp_enable"); ok {
		virtualServerOCSPEnable := v.(bool)
		virtualServer.Properties.Ssl.OCSPEnable = &virtualServerOCSPEnable
	}
	if v, ok := d.GetOk("ocsp_issuers"); ok {
		virtualServer.Properties.Ssl.OCSPIssuers = buildSSLOCSPIssuers(v.([]interface{}))
	}

	createAPI := virtualserver.NewCreate(virtualServerName, virtualServer)
	err := vtmClient.Do(createAPI)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Brocade vTM Virtual Server Create failed for %s with error: %+v", virtualServerName, err))
	}
	if createAPI.StatusCode() != http.StatusCreated {
		return fmt.Errorf(fmt.Sprintf("Brocade vTM Virtual Server Create failed for %s with http status code != 201 - error: %+v", virtualServerName, createAPI.GetResponse()))
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
		return fmt.Errorf(fmt.Sprintf("Brocade vTM Virtual Server Read failed for %s with error: %+v", virtualServerName, err))
	}
	if getSingleAPI.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	virtualServer = *getSingleAPI.GetResponse()
	d.SetId(virtualServerName)
	d.Set("enabled", *virtualServer.Properties.Basic.Enabled)
	d.Set("listen_on_any", *virtualServer.Properties.Basic.ListenOnAny)
	d.Set("listen_traffic_ips", virtualServer.Properties.Basic.ListenOnTrafficIps)
	d.Set("pool", virtualServer.Properties.Basic.Pool)
	d.Set("port", virtualServer.Properties.Basic.Port)
	d.Set("protocol", virtualServer.Properties.Basic.Protocol)
	d.Set("request_rules", virtualServer.Properties.Basic.RequestRules)
	d.Set("ssl_decrypt", *virtualServer.Properties.Basic.SslDecrypt)
	d.Set("connection_keepalive", *virtualServer.Properties.Connection.Keepalive)
	d.Set("connection_keepalive_timeout", virtualServer.Properties.Connection.KeepaliveTimeout)
	d.Set("connection_max_client_buffer", virtualServer.Properties.Connection.MaxClientBuffer)
	d.Set("connection_max_server_buffer", virtualServer.Properties.Connection.MaxServerBuffer)
	d.Set("connection_max_transaction_duration", virtualServer.Properties.Connection.MaxTransactionDuration)
	d.Set("connection_server_first_banner", virtualServer.Properties.Connection.ServerFirstBanner)
	d.Set("connection_timeout", virtualServer.Properties.Connection.Timeout)
	d.Set("ssl_server_cert_default", virtualServer.Properties.Ssl.ServerCertDefault)
	d.Set("ssl_support_ssl2", virtualServer.Properties.Ssl.SslSupportSsl2)
	d.Set("ssl_support_ssl3", virtualServer.Properties.Ssl.SslSupportSsl3)
	d.Set("ssl_support_tls1", virtualServer.Properties.Ssl.SslSupportTLS1)
	d.Set("ssl_support_tls1_1", virtualServer.Properties.Ssl.SslSupportTLS1_1)
	d.Set("ssl_support_tls1_2", virtualServer.Properties.Ssl.SslSupportTLS1_2)
	d.Set("ssl_server_cert_host_mapping", virtualServer.Properties.Ssl.ServerCertHostMap)
	d.Set("ocsp_enable", virtualServer.Properties.Ssl.OCSPEnable)
	d.Set("ocsp_issuers", virtualServer.Properties.Ssl.OCSPIssuers)

	return nil
}

func resourceVirtualServerUpdate(d *schema.ResourceData, m interface{}) error {

	vtmClient := m.(*brocadevtm.VTMClient)
	var virtualServerName string
	var virtualServer virtualserver.VirtualServer
	hasChanges := false

	if v, ok := d.GetOk("name"); ok && v != "" {
		virtualServerName = v.(string)
	}
	if d.HasChange("enabled") {
		virtualServerEnabled := d.Get("enabled").(bool)
		virtualServer.Properties.Basic.Enabled = &virtualServerEnabled
		hasChanges = true
	}

	if d.HasChange("listen_on_any") {
		virtualServerListenAny := d.Get("listen_on_any").(bool)
		virtualServer.Properties.Basic.ListenOnAny = &virtualServerListenAny
		hasChanges = true
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
	if d.HasChange("listen_traffic_ips") {
		if v, ok := d.GetOk("listen_traffic_ips"); ok && v != "" {
			virtualServer.Properties.Basic.ListenOnTrafficIps = buildStringList(v)
		}
		hasChanges = true
	}
	if d.HasChange("protocol") {
		if v, ok := d.GetOk("protocol"); ok && v != "" {
			virtualServer.Properties.Basic.Protocol = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("request_rules") {
		if v, ok := d.GetOk("request_rules"); ok && v != "" {
			virtualServer.Properties.Basic.RequestRules = buildStringList(v)
		}
		hasChanges = true
	}
	if d.HasChange("ssl_decrypt") {
		virtualServerSSLDeCrypt := d.Get("ssl_decrypt").(bool)
		virtualServer.Properties.Basic.SslDecrypt = &virtualServerSSLDeCrypt
		hasChanges = true
	}
	if d.HasChange("connection_keepalive") {
		virtalServerConnectionKeepalive := d.Get("connection_keepalive").(bool)
		virtualServer.Properties.Connection.Keepalive = &virtalServerConnectionKeepalive
		hasChanges = true
	}
	if d.HasChange("connection_keepalive_timeout") {
		if v, ok := d.GetOk("connection_keepalive_timeout"); ok {
			virtualServerConnectionKeepaliveTimeout := v.(int)
			virtualServer.Properties.Connection.KeepaliveTimeout = uint(virtualServerConnectionKeepaliveTimeout)
		}
		hasChanges = true
	}
	if d.HasChange("connection_max_client_buffer") {
		if v, ok := d.GetOk("connection_max_client_buffer"); ok {
			virtualServerConnectionMaxClientBuffer := v.(int)
			virtualServer.Properties.Connection.MaxClientBuffer = uint(virtualServerConnectionMaxClientBuffer)
		}
		hasChanges = true
	}
	if d.HasChange("connection_max_server_buffer") {
		if v, ok := d.GetOk("connection_max_server_buffer"); ok {
			virtualServerConnectionMaxServerBuffer := v.(int)
			virtualServer.Properties.Connection.MaxServerBuffer = uint(virtualServerConnectionMaxServerBuffer)
		}
		hasChanges = true
	}
	if d.HasChange("connection_max_transaction_duration") {
		if v, ok := d.GetOk("connection_max_transaction_duration"); ok {
			virtualServerConnectionMaxTransActionDuration := v.(int)
			virtualServer.Properties.Connection.MaxTransactionDuration = uint(virtualServerConnectionMaxTransActionDuration)
		}
		hasChanges = true
	}
	if d.HasChange("connection_server_first_banner") {
		if v, ok := d.GetOk("connection_server_first_banner"); ok && v != "" {
			virtualServer.Properties.Connection.ServerFirstBanner = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("connection_timeout") {
		if v, ok := d.GetOk("connection_timeout"); ok {
			virtualServerConnectionTimeout := v.(int)
			virtualServer.Properties.Connection.Timeout = uint(virtualServerConnectionTimeout)
		}
		hasChanges = true
	}
	if d.HasChange("ssl_server_cert_default") {
		if v, ok := d.GetOk("ssl_server_cert_default"); ok && v != "" {
			virtualServer.Properties.Ssl.ServerCertDefault = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("ssl_support_ssl2") {
		if v, ok := d.GetOk("ssl_support_ssl2"); ok && v != "" {
			virtualServer.Properties.Ssl.SslSupportSsl2 = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("ssl_support_ssl3") {
		if v, ok := d.GetOk("ssl_support_ssl3"); ok && v != "" {
			virtualServer.Properties.Ssl.SslSupportSsl3 = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("ssl_support_tls1") {
		if v, ok := d.GetOk("ssl_support_tls1"); ok && v != "" {
			virtualServer.Properties.Ssl.SslSupportTLS1 = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("ssl_support_tls1_1") {
		if v, ok := d.GetOk("ssl_support_tls1_1"); ok && v != "" {
			virtualServer.Properties.Ssl.SslSupportTLS1_1 = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("ssl_support_tls1_2") {
		if v, ok := d.GetOk("ssl_support_tls1_2"); ok && v != "" {
			virtualServer.Properties.Ssl.SslSupportTLS1_2 = v.(string)
		}
		hasChanges = true
	}
	if d.HasChange("ssl_server_cert_host_mapping") {
		if v, ok := d.GetOk("ssl_server_cert_host_mapping"); ok && v != "" {
			virtualServer.Properties.Ssl.ServerCertHostMap = buildSSLCertMapping(v.([]interface{}))
		}
		hasChanges = true
	}
	if d.HasChange("ocsp_enable") {
		virtualServerOCSPEnable := d.Get("ocsp_enable").(bool)
		virtualServer.Properties.Ssl.OCSPEnable = &virtualServerOCSPEnable
		hasChanges = true
	}
	if d.HasChange("ocsp_issuers") {
		if v, ok := d.GetOk("ocsp_issuers"); ok {
			virtualServer.Properties.Ssl.OCSPIssuers = buildSSLOCSPIssuers(v.([]interface{}))
		}
		hasChanges = true
	}

	if hasChanges {
		updateAPI := virtualserver.NewUpdate(virtualServerName, virtualServer)
		err := vtmClient.Do(updateAPI)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("Brocade vTM Virtual Server update failed for %s", virtualServerName))
		}
		responseCode := updateAPI.StatusCode()
		if responseCode != http.StatusOK {
			return fmt.Errorf(fmt.Sprintf("Brocade vTM Virtual Server update failed for %s with invalid response code %d - response: %+v", virtualServerName, responseCode, updateAPI.GetResponse()))
		}
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
		return fmt.Errorf(fmt.Sprintf("Brocade vTM Virtual Server delete failed for %s - error: %+v", virtualServerName, err))
	}
	if getVirtualServer.StatusCode() == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	deleteAPI := virtualserver.NewDelete(virtualServerName)
	err = vtmClient.Do(deleteAPI)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Brocade vTM Virtual Server delete failed for %s - error: %+v", virtualServerName, err))
	}
	responseCode := deleteAPI.StatusCode()
	if responseCode != http.StatusNoContent {
		return fmt.Errorf(fmt.Sprintf("Brocade vTM Virtual Server delete returned an invalid http response code %d for %s", responseCode, virtualServerName))
	}

	d.SetId("")
	return nil
}
