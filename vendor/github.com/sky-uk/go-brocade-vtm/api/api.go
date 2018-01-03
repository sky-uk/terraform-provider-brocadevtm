package api

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/sky-uk/go-rest-api"
)

const defaultAPIVersion = "3.8"
const configPath = "config/active"
const apiPrefix = "/api/tm"

// Params - connection parameters
type Params struct {
	Username   string
	Password   string
	Server     string
	APIVersion string
	IgnoreSSL  bool
	Debug      bool
	Timeout    time.Duration
	Headers    map[string]string
}

// Client - the Brocade vTM Client struct
type Client struct {
	VersionsSupported []string
	restClient        rest.Client
	currentVersion    string
	RootPath          string
	currentServer     string
	params            Params
	StatusCode        int
}

// WorkWithStatus - sets the root path to work with status resources
func (client *Client) WorkWithStatus() {
	client.RootPath = apiPrefix + "/" + client.currentVersion + "/status"
	if client.params.Debug {
		log.Println("[DEBUG] Current Path: ", client.RootPath)
	}
}

// GetStatistics - returns all statistics...
func (client *Client) GetStatistics(node string) (map[string]interface{}, error) {
	client.WorkWithStatus()
	path := client.RootPath + "/" + node + "/statistics"
	all := make(map[string]interface{})
	err := client.TraverseTree(path, all)
	return all, err
}

// GetState - get a node state
func (client *Client) GetState(node string) (map[string]interface{}, error) {
	client.WorkWithStatus()
	path := client.RootPath + "/" + node + "/state"
	state := make(map[string]interface{})
	api := rest.NewBaseAPI(
		http.MethodGet,
		path,
		nil,
		&state,
		new(VTMError),
	)
	err := client.request(api)
	return state, err
}

// GetInformation - returns all information...
func (client *Client) GetInformation(node string) (map[string]interface{}, error) {
	client.WorkWithStatus()
	path := client.RootPath + "/" + node + "/information"
	all := make(map[string]interface{})
	api := rest.NewBaseAPI(
		http.MethodGet,
		path,
		nil,
		&all,
		new(VTMError),
	)
	err := client.request(api)
	return all, err
}

// WorkWithConfigurationResources - set the root path to work with configuration resources
func (client *Client) WorkWithConfigurationResources() {
	client.RootPath = apiPrefix + "/" + client.currentVersion + "/" + configPath
	if client.params.Debug {
		log.Println("[DEBUG] Current Path: ", client.RootPath)
	}
}

// Connect - connect to the Brocade REST API server
// and get the list of supported API versions
// Returns a new client object if everything is fine
func Connect(params Params) (*Client, error) {
	client := new(Client)
	client.currentVersion = params.APIVersion
	client.params = params

	if params.Headers == nil {
		// if client doesn't pass any header, we only set
		// the content type to be the default one...
		headers := make(map[string]string)
		headers["Content-Type"] = "application/json"
		params.Headers = headers
	}

	if strings.HasPrefix(params.Server, "https") == false {
		params.Server = "https://" + params.Server
	}

	client.restClient = rest.Client{
		URL:       params.Server,
		User:      params.Username,
		Password:  params.Password,
		IgnoreSSL: params.IgnoreSSL,
		Debug:     params.Debug,
		Headers:   params.Headers,
		Timeout:   params.Timeout,
	}

	supportedVersionsMap := make(map[string]interface{})

	if client.currentVersion == "" {

		api := rest.NewBaseAPI(
			http.MethodGet,
			apiPrefix,
			nil,
			&supportedVersionsMap,
			new(VTMError),
		)
		err := client.request(api)
		if err != nil || api.StatusCode() != http.StatusOK {
			log.Println("[ERROR] Error while fetching list of available API versions: ", err)
			return nil, err
		}

		versions := make([]string, 0)
		for _, version := range supportedVersionsMap["children"].([]interface{}) {
			if vAsMap, ok := version.(map[string]interface{}); ok {
				versions = append(versions, vAsMap["name"].(string))
			}
		}

		client.VersionsSupported = versions
		sort.Sort(sort.Reverse(sort.StringSlice(client.VersionsSupported)))
		client.currentVersion = client.VersionsSupported[0]
		log.Println("[DEBUG] Working with REST API Version: ", client.currentVersion)

	}

	return client, nil
}

// GetAllResourceTypes - returns the list of all types of configuration resources
func (client *Client) GetAllResourceTypes() ([]map[string]interface{}, error) {

	// work with an environment
	client.WorkWithConfigurationResources()

	path := client.RootPath
	res := make(map[string]interface{}, 0)

	if client.params.Debug {
		log.Println("[DEBUG] Going to get all resource types, using PATH:\n", path)
	}
	api := rest.NewBaseAPI(
		http.MethodGet,
		path,
		nil,
		&res,
		new(VTMError),
	)
	err := client.request(api)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resTypes := make([]map[string]interface{}, 0)
	for _, item := range res["children"].([]interface{}) {
		resTypes = append(resTypes, item.(map[string]interface{}))
	}
	return resTypes, nil
}

// FormatNestedStruct - recursively navigate the tree for any child which value is not a string.
func FormatNestedStruct(str map[string]interface{}, level int) string {
	var outStr string
	outStr += "\n"
	for key, value := range str {
		outStr += strings.Repeat("\t", level) + key
		if vAsStr, ok := value.(string); ok {
			outStr += " : " + vAsStr
		} else {
			level++
			outStr += FormatNestedStruct(value.(map[string]interface{}), level)
		}
		outStr += "\n"
	}

	return outStr
}

// FormatErrorText - formats the error message including attributes error strings
func FormatErrorText(tmErr *VTMError) string {
	retStr := tmErr.ErrorID + ": " + tmErr.ErrorText + "\n"
	retStr += FormatNestedStruct(tmErr.ErrorInfo, 0)
	return retStr
}

func (client *Client) request(api *rest.BaseAPI) error {
	err := client.restClient.Do(api)
	client.StatusCode = client.restClient.StatusCode
	tmErr := api.ErrorObject().(*VTMError)
	if tmErr.ErrorText != "" {
		errStr := FormatErrorText(tmErr)
		err = fmt.Errorf(errStr)
	}
	return err
}

// TraverseTree - retrieves a resource and eventually keep doing it
// for each nested resource
// Fill up the passed slice of resources, returns the first error it
// eventually bumps into
func (client *Client) TraverseTree(url string, resources map[string]interface{}) error {
	res := make(map[string]interface{})

	if url == "" {
		return fmt.Errorf("[ERROR] Invalid path")
	}

	if client.params.Debug {
		log.Println("[DEBUG] Going to get PATH: ", url)
	}
	api := rest.NewBaseAPI(
		http.MethodGet,
		url,
		nil,
		&res,
		new(VTMError),
	)
	err := client.request(api)
	if err != nil {
		log.Println(err)
		return err
	}
	if children, exists := res["children"]; exists {
		for _, item := range children.([]interface{}) {
			if itemAsMap, ok := item.(map[string]interface{}); ok {
				err = client.TraverseTree(itemAsMap["href"].(string), resources)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("[ERROR] Strange...I expected a slice of maps")
			}
		}
	} else {
		resources[url] = res
	}

	return nil
}

// GetAllResources - returns all resources of the specified type
func (client *Client) GetAllResources(resType string) ([]map[string]interface{}, error) {
	path := client.RootPath + "/" + resType
	res := make(map[string]interface{})

	if client.params.Debug {
		log.Println("[DEBUG] Going to get all resources, using PATH: ", path)
	}
	api := rest.NewBaseAPI(
		http.MethodGet,
		path,
		nil,
		&res,
		new(VTMError),
	)
	err := client.request(api)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resources := make([]map[string]interface{}, 0)
	if list, exists := res["children"].([]interface{}); exists {
		for _, item := range list {
			resources = append(resources, item.(map[string]interface{}))
		}
	}
	return resources, nil
}

// GetByName - gets a resource profile given its type and name
func (client *Client) GetByName(resType, resName string, out interface{}) error {
	path := client.RootPath + "/" + resType + "/" + resName
	api := rest.NewBaseAPI(
		http.MethodGet,
		path,
		nil,
		out,
		new(VTMError),
	)
	return client.request(api)
}

// GetByURL - gets a resource profile given its type and URL
func (client *Client) GetByURL(resURL string, out interface{}) error {
	api := rest.NewBaseAPI(
		http.MethodGet,
		resURL,
		nil,
		out,
		new(VTMError),
	)
	return client.request(api)
}

// Set - Sets a resource
// This works only in Configuration environment (statistics/information resources can't be set)
// A new resources gets created if not existent or an existent resource gets updated
// The restClient.StatusCode is set properly to http.StatusCreated or http.StatusOK accordingly
// Returns the created/updated object or an error
func (client *Client) Set(resType, name string, profile interface{}, out interface{}) error {

	// you can only set configuration resources...
	client.WorkWithConfigurationResources()

	path := client.RootPath + "/" + resType + "/" + name
	if out == nil {
		res := make(map[string]interface{})
		out = &res
	}
	api := rest.NewBaseAPI(
		http.MethodPut,
		path,
		profile,
		out,
		new(VTMError),
	)
	return client.request(api)
}

// Delete - deletes a resource
func (client *Client) Delete(resType, name string) error {

	// you can only delete configuration resources...
	client.WorkWithConfigurationResources()
	path := client.RootPath + "/" + resType + "/" + name
	api := rest.NewBaseAPI(http.MethodDelete, path, nil, nil, new(VTMError))
	err := client.request(api)
	if err != nil {
		log.Println(err)
		return err
	}
	if api.StatusCode() != http.StatusNoContent {
		return fmt.Errorf(" [ERROR] Error deleting resource %s, status: %d", name, api.StatusCode())
	}
	return nil
}
