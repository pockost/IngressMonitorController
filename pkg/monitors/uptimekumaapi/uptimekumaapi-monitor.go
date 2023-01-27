package uptimekumaapi

import (
	"time"
	"encoding/json"
	"errors"
	"fmt"
	Http "net/http"
	"strconv"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	endpointmonitorv1alpha1 "github.com/stakater/IngressMonitorController/v2/api/v1alpha1"
	"github.com/stakater/IngressMonitorController/v2/pkg/config"
	"github.com/stakater/IngressMonitorController/v2/pkg/http"
	"github.com/stakater/IngressMonitorController/v2/pkg/models"
)

var log = logf.Log.WithName("uptimekumaapi-monitor")

type UptimeKumaApiMonitorService struct {
	apiUrl         string
	apiUsername    string
	apiPassword    string
	apiAccessToken string
}

func (service *UptimeKumaApiMonitorService) GenerateHttpClient(route string) (client *http.HttpClient, headers map[string]string) {

	cli := http.CreateHttpClient(service.apiUrl + route)

	// Construct headers
	head := make(map[string]string)
	head["Authorization"] = fmt.Sprintf("Bearer %s", service.apiAccessToken)
	head["Accept"] = "application/json"

	// Test si authentification is valide
	client_test := http.CreateHttpClient(service.apiUrl + "/info/")
	response := client_test.GetUrl(head, nil)
	if response.StatusCode != Http.StatusOK {
		service.RetreiveAccessToken()
		head["Authorization"] = fmt.Sprintf("Bearer %s", service.apiAccessToken)
	} else {
	}

	return cli, head
}

func (service *UptimeKumaApiMonitorService) GetAll() []models.Monitor {

	var monitors []models.Monitor

	route := "/monitors/"

	client, headers := service.GenerateHttpClient(route)

	// Request monitor list
	response := client.GetUrl(headers, nil)

	if response.StatusCode == Http.StatusOK {
		var f UptimeKumaApiMonitorListResponse
		err := json.Unmarshal(response.Bytes, &f)
		if err != nil {
			log.Error(err, "Unable to unmarshal response")
		}

		// Build a list of Monitor
		for _, monitor := range f.Monitors {
			var m models.Monitor
			m.ID = strconv.Itoa(monitor.Id)
			m.Name = monitor.Name
			m.URL = monitor.Url

			config := &endpointmonitorv1alpha1.UptimeKumaApiConfig{
				Interval:            monitor.Interval,
				RetryInterval:       monitor.RetryInterval,
				ResendInterval:      monitor.ResendInterval,
				MaxRetries:          monitor.MaxRetries,
				Method:              monitor.Method,
				IgnoreTLS:           monitor.IgnoreTLS,
				UpsideDown:          monitor.UpsideDown,
				MaxRedirects:        monitor.MaxRedirects,
				AcceptedStatusCodes: monitor.AcceptedStatusCodes,
				SSLExpire:           monitor.SSLExpire,
			}
			m.Config = config
			monitors = append(monitors, m)
		}

	} else {
		log.Error(nil, fmt.Sprintf("Unable to retreive monitors list (%d): %s", response.StatusCode, response.Bytes))
	}

	return monitors
}

func (service *UptimeKumaApiMonitorService) Add(m models.Monitor) {
	route := "/monitors/"

	client, headers := service.GenerateHttpClient(route)

	// Retrieve provider configuration
	providerConfig, _ := m.Config.(*endpointmonitorv1alpha1.UptimeKumaApiConfig)

	// Construct Body
	var uptimeKumaApiMonitor UptimeKumaApiMonitor
	// Default data
	uptimeKumaApiMonitor.Type = "http"
	uptimeKumaApiMonitor.Method = "GET"
	uptimeKumaApiMonitor.Interval = 60
	uptimeKumaApiMonitor.RetryInterval = 60
	uptimeKumaApiMonitor.MaxRetries = 0
	uptimeKumaApiMonitor.ResendInterval = 0
	uptimeKumaApiMonitor.IgnoreTLS = false
	uptimeKumaApiMonitor.UpsideDown = false
	uptimeKumaApiMonitor.MaxRedirects = 10
	uptimeKumaApiMonitor.AcceptedStatusCodes = []string{"200-299"}
	uptimeKumaApiMonitor.SSLExpire = true

	// Overrided data
	if providerConfig != nil {
		if providerConfig.Interval != 0 {
			uptimeKumaApiMonitor.Interval = providerConfig.Interval
		}
		if providerConfig.RetryInterval != 0 {
			uptimeKumaApiMonitor.RetryInterval = providerConfig.RetryInterval
		}
		if providerConfig.ResendInterval != 0 {
			uptimeKumaApiMonitor.ResendInterval = providerConfig.ResendInterval
		}
		if providerConfig.MaxRetries != 0 {
			uptimeKumaApiMonitor.MaxRetries = providerConfig.MaxRetries
		}
		if providerConfig.Method != "" {
			uptimeKumaApiMonitor.Method = providerConfig.Method
		}
		if providerConfig.IgnoreTLS {
			uptimeKumaApiMonitor.IgnoreTLS = providerConfig.IgnoreTLS
		}
		if providerConfig.UpsideDown {
			uptimeKumaApiMonitor.UpsideDown = providerConfig.UpsideDown
		}
		if providerConfig.MaxRedirects != 0 {
			uptimeKumaApiMonitor.MaxRedirects = providerConfig.MaxRedirects
		}
		if providerConfig.AcceptedStatusCodes != nil {
			uptimeKumaApiMonitor.AcceptedStatusCodes = providerConfig.AcceptedStatusCodes
		}
		if !providerConfig.SSLExpire {
			uptimeKumaApiMonitor.SSLExpire = providerConfig.SSLExpire
		}
	}

	// Dynamic data
	uptimeKumaApiMonitor.Name = m.Name
	uptimeKumaApiMonitor.Url = m.URL
	body, err := json.Marshal(uptimeKumaApiMonitor)

	if err != nil {
		log.Error(err, "Unable to marshal json")
	}

	// Add monitor
	response := client.PostUrl(headers, body)

	// Handle error
	if response.StatusCode != Http.StatusOK {
		var f UptimeKumaApiErrorResponse
		err := json.Unmarshal(response.Bytes, &f)
		if err != nil {
			log.Error(err, "Unable to unmarshal JSON error response")
		}
		log.Error(nil, fmt.Sprintf("Unable to create monitor %s %s: %v", m.Name, m.ID, f.Detail))
	}
}

func (service *UptimeKumaApiMonitorService) Update(m models.Monitor) {

	// Get old monitor in order to construct field to update
	oldMonitor, err := service.GetByName(m.Name)

	if err != nil {
		log.Error(err, "Unable to find given monitor")
	}

	// Construct Body
	var uptimeKumaApiMonitor UptimeKumaApiMonitor
	bodyHasChange := false
	uptimeKumaApiMonitor.Name = oldMonitor.Name
	uptimeKumaApiMonitor.Type = "http"
	uptimeKumaApiMonitor.Method = "GET"
	uptimeKumaApiMonitor.Url = oldMonitor.URL

	if oldMonitor.URL != m.URL {
		uptimeKumaApiMonitor.Url = m.URL
		bodyHasChange = true
	}

	if bodyHasChange {

		route := fmt.Sprintf("/monitors/%s", oldMonitor.ID)

		client, headers := service.GenerateHttpClient(route)

		body, err := json.Marshal(uptimeKumaApiMonitor)

		if err != nil {
			log.Error(err, "Unable to marshal json")
		}

		client.RequestWithHeaders("PATCH", body, headers)
	}

}

func (service *UptimeKumaApiMonitorService) GetByName(name string) (*models.Monitor, error) {
	monitors := service.GetAll()
	for _, monitor := range monitors {
		if monitor.Name == name {
			return &monitor, nil
		}
	}
	return nil, errors.New("Unable to find given monitor")
}

func (service *UptimeKumaApiMonitorService) Remove(m models.Monitor) {

	route := fmt.Sprintf("/monitors/%s", m.ID)

	client, headers := service.GenerateHttpClient(route)

	// Delete monitor
	response := client.DeleteUrl(headers, nil)

	// Handle error
	if response.StatusCode != Http.StatusOK {
		var f UptimeKumaApiErrorResponse
		err := json.Unmarshal(response.Bytes, &f)
		if err != nil {
			log.Error(err, "Unable to unmarshal JSON error response")
		}
		log.Error(nil, fmt.Sprintf("Unable to remove monitor %s %s: %v", m.Name, m.ID, f.Detail))
	}
}

func (service *UptimeKumaApiMonitorService) RetreiveAccessToken() {

	route := "/login/access-token/"
	client := http.CreateHttpClient(service.apiUrl + route)
	body := fmt.Sprintf("grant_type=&username=%s&password=%s&scope=&client_id=&client_secret=", service.apiUsername, service.apiPassword)
	response := client.PostUrlEncodedFormBody(body)

	if response.StatusCode == Http.StatusOK {
		var f UptimeKumaApiAuthenticationResponse
		err := json.Unmarshal(response.Bytes, &f)
		if err != nil {
			log.Error(err, "Unable to unmarshal JSON login response")
		}
		service.apiAccessToken = f.AccessToken
	} else {
		log.Error(nil, "Unable to authenticate")
		// Retry in 5 seconde
		time.Sleep(5 * time.Second)
		service.RetreiveAccessToken()
	}
}

func (service *UptimeKumaApiMonitorService) Setup(p config.Provider) {
	service.apiUsername = p.Username
	service.apiPassword = p.Password
	service.apiUrl = p.ApiURL

	service.RetreiveAccessToken()

	// Authenticate user and save access token

}

func (monitor *UptimeKumaApiMonitorService) Equal(oldMonitor models.Monitor, newMonitor models.Monitor) bool {
	if oldMonitor.Name != newMonitor.Name {
		return false
	}
	if oldMonitor.URL != newMonitor.URL {
		return false
	}
	return true
}
