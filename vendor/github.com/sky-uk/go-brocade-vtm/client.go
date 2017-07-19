package brocadevtm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/sky-uk/go-brocade-vtm/api"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// NewVTMClient  Creates a new vtmClient object.
func NewVTMClient(
	url string,
	user string,
	password string,
	ignoreSSL bool,
	debug bool,
	headers map[string]string,
) *VTMClient {

	vtmClient := new(VTMClient)
	vtmClient.URL = url
	vtmClient.User = user
	vtmClient.Password = password
	vtmClient.IgnoreSSL = ignoreSSL
	vtmClient.Debug = debug
	vtmClient.Headers = headers

	return vtmClient
}

// VTMClient struct.
type VTMClient struct {
	URL       string
	User      string
	Password  string
	IgnoreSSL bool
	Debug     bool
	Headers   map[string]string
}

func (vtmClient *VTMClient) formatRequestPayload(api api.VTMApi) (io.Reader, error) {

	var requestPayload io.Reader

	var reqBytes []byte
	if api.RequestObject() != nil {
		var err error
		contentType := vtmClient.Headers["Content-Type"]
		if contentType == "application/json" {
			reqBytes, err = json.Marshal(api.RequestObject())
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
		}
		if contentType == "application/xml" {
			reqBytes, err = xml.Marshal(api.RequestObject())
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
		}
		if contentType == "application/octet-stream" {
			reqBytes = api.RequestObject().([]byte)
		}

		requestPayload = bytes.NewReader(reqBytes)
	}

	if vtmClient.Debug {
		log.Println("--------------------------------------------------------------")
		log.Println("Request payload:")
		log.Println(string(reqBytes))
		log.Println("--------------------------------------------------------------")
	}

	return requestPayload, nil
}

// Do - makes the API call.
func (vtmClient *VTMClient) Do(api api.VTMApi) error {

	requestURL := fmt.Sprintf("%s%s", vtmClient.URL, api.Endpoint())

	if vtmClient.Headers == nil {
		vtmClient.Headers = make(map[string]string)
	}

	_, ok := vtmClient.Headers["Content-Type"]
	if !ok {
		vtmClient.Headers["Content-Type"] = "application/json"
	}

	requestPayload, err := vtmClient.formatRequestPayload(api)
	if err != nil {
		return err
	}

	if vtmClient.Debug {
		log.Println("requestURL:", requestURL)
	}
	req, err := http.NewRequest(api.Method(), requestURL, requestPayload)
	if err != nil {
		log.Println("ERROR building the request: ", err)
		return err
	}

	req.SetBasicAuth(vtmClient.User, vtmClient.Password)

	for headerKey, headerValue := range vtmClient.Headers {
		req.Header.Set(headerKey, headerValue)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: vtmClient.IgnoreSSL},
	}
	httpClient := &http.Client{Transport: tr}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Println("ERROR executing request: ", err)
		return err
	}
	defer res.Body.Close()
	return vtmClient.handleResponse(api, res)
}

func (vtmClient *VTMClient) handleResponse(apiObj api.VTMApi, res *http.Response) error {
	apiObj.SetStatusCode(res.StatusCode)
	bodyText, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ERROR reading response: ", err)
		return err
	}

	apiObj.SetRawResponse(bodyText)

	if vtmClient.Debug {
		log.Println(string(bodyText))
	}

	if isJSON(res.Header.Get("Content-Type")) {
		if apiObj.StatusCode() >= http.StatusOK && apiObj.StatusCode() < http.StatusBadRequest {
			if len(bodyText) > 0 {
				JSONerr := json.Unmarshal(bodyText, apiObj.ResponseObject())
				if JSONerr != nil {
					log.Println("ERROR unmarshalling response: ", JSONerr)
					return JSONerr
				}
			}
			return nil
		}

		if len(bodyText) > 0 {
			var errObj api.ReqError
			err := json.Unmarshal(bodyText, &errObj)
			if err != nil {
				log.Printf("Error unmarshalling error response:\n%v", err)
			}
			return errors.New(errObj.Error.ErrorText)
		}
	}
	if isXML(res.Header.Get("Content-Type")) {
		if apiObj.StatusCode() >= http.StatusOK && apiObj.StatusCode() < http.StatusBadRequest {
			if len(bodyText) > 0 {
				xmlErr := xml.Unmarshal(bodyText, apiObj.ResponseObject())
				if xmlErr != nil {
					log.Println("ERROR unmarshalling response: ", xmlErr)
					return xmlErr
				}
			}
			return nil
		}

		if len(bodyText) > 0 {
			var errObj api.ReqError
			err := xml.Unmarshal(bodyText, &errObj)
			if err != nil {
				log.Printf("Error unmarshalling error response:\n%v", err)
			}
			return errors.New(errObj.Error.ErrorText)
		}

	} else {
		data := []byte(bodyText)
		apiObj.SetRawResponse(data)
	}
	return nil
}

func isXML(contentType string) bool {
	return strings.Contains(strings.ToLower(contentType), "/xml")
}

func isJSON(contentType string) bool {
	return strings.Contains(strings.ToLower(contentType), "/json")
}
