package uptimekumaapi

import (
	"testing"

	endpointmonitorv1alpha1 "github.com/stakater/IngressMonitorController/v2/api/v1alpha1"
	"github.com/stakater/IngressMonitorController/v2/pkg/config"
	"github.com/stakater/IngressMonitorController/v2/pkg/models"
	"github.com/stakater/IngressMonitorController/v2/pkg/util"
)

// Load config for UptimeKumaApi service
func setupService(t *testing.T) UptimeKumaApiMonitorService {
	config := config.GetControllerConfigTest()

	service := UptimeKumaApiMonitorService{}

	provider := util.GetProviderWithName(config, "UptimeKumaApi")

	if provider == nil {
		t.Fatalf("The UptimeKumaApi provider is not registred")
	}
	service.Setup(*provider)

	return service
}

// Remove all monitor
func emptyMonitorList(service UptimeKumaApiMonitorService) {
	monitors := service.GetAll()
	for _, monitor := range monitors {
		service.Remove(monitor)
	}
}

func createMonitor(service UptimeKumaApiMonitorService, name string, url string) {
	monitor := models.Monitor{
		Name: name,
		URL:  url,
	}

	service.Add(monitor)
}

func TestDeleteMonitor(t *testing.T) {
	service := setupService(t)
	emptyMonitorList(service)
	createMonitor(service, "test", "test.fr")

	monitor, err := service.GetByName("test")
	if err != nil {
		t.Fatal(err, "Unable to get service")
	}
	service.Remove(*monitor)
}

func TestEqualMonitor(t *testing.T) {

	service := setupService(t)
	emptyMonitorList(service)
	createMonitor(service, "test", "test.fr")

	monitor1, err := service.GetByName("test")
	if err != nil {
		t.Fatal(err, "Unable to get service")
	}

	monitor2, err := service.GetByName("test")
	if err != nil {
		t.Fatal(err, "Unable to get service")
	}

	if !service.Equal(*monitor1, *monitor2) {
		t.Fatal("Test of Equal monitor should return true")
	}

	monitor2.URL = "https://dummy.com"

	if service.Equal(*monitor1, *monitor2) {
		t.Fatal("Test of Non Equal monitor should return false")
	}

	// We should not test the ID as it's an auto generated value
	var monitor3 models.Monitor
	monitor3.Name = monitor1.Name
	monitor3.URL = monitor1.URL

	if !service.Equal(*monitor1, monitor3) {
		t.Fatal("Same monitor with no ID should be considered as same monitor")
	}

	// Test if Name has change
	monitor3.Name = "dummy"

	if service.Equal(*monitor1, monitor3) {
		t.Fatal("Name has change, should return false")
	}

}

func TestSetupMonitor(t *testing.T) {

	service := setupService(t)

	if service.apiUrl == "" {
		t.Fatalf("API Url is not loaded")
	}
	if service.apiPassword == "" {
		t.Fatalf("Password is not loaded")
	}
	if service.apiUsername == "" {
		t.Fatalf("Username is not loaded")
	}
	if service.apiAccessToken == "" {
		t.Fatalf("Access token is not loaded")
	}
}

func TestGetAllMonitorWithEmptyList(t *testing.T) {

	service := setupService(t)

	emptyMonitorList(service)

	monitors := service.GetAll()

	if len(monitors) != 0 {
		t.Error("Unable to load empty service list")
	}
}

func TestGetAllMonitorWithElement(t *testing.T) {

	service := setupService(t)

	emptyMonitorList(service)
	createMonitor(service, "google", "https://www.google.com")

	monitors := service.GetAll()

	if len(monitors) != 1 {
		t.Error("Unable to retreive a list of 1 monitor")
	}

}

func TestAddMonitor(t *testing.T) {
	service := setupService(t)
	emptyMonitorList(service)

	monitor := models.Monitor{
		Name: "google.com",
		URL:  "https://www.google.com",
	}

	service.Add(monitor)

	monitors := service.GetAll()
	if len(monitors) == 0 {
		t.Fatal("No service created")
	}
	if len(monitors) > 1 {
		t.Error("More than one service created")
	}

	if monitors[0].Name != "google.com" {
		t.Error("Created service has not correct name")
	}

	// Test default values
	providerConfig, _ := monitors[0].Config.(*endpointmonitorv1alpha1.UptimeKumaApiConfig)

	if providerConfig == nil {
		t.Error("No provider config returned")
	}

	if providerConfig.Interval != 60 {
		t.Error("Incorrect default value (Interval)")
	}
	if providerConfig.RetryInterval != 60 {
		t.Error("Incorrect default value (RetryInterval)")
	}

	if providerConfig.ResendInterval!= 0 {
		t.Error("Incorrect default value (ResendInterval)")
	}
	if providerConfig.MaxRetries != 0 {
		t.Error("Incorrect default value (MaxRetries)")
	}
	if providerConfig.Method != "GET" {
		t.Error("Incorrect default value (Method)")
	}
	if providerConfig.IgnoreTLS != false {
		t.Error("Incorrect default value (IgnoreTLS)")
	}
	if providerConfig.UpsideDown != false {
		t.Error("Incorrect default value (UpsideDown)")
	}
	if providerConfig.MaxRedirects != 10 {
		t.Error("Incorrect default value (MaxRedirects)")
	}
	acceptedStatusCodesOk := false
	for _, v := range providerConfig.AcceptedStatusCodes {
		if v == "200-299" {
			acceptedStatusCodesOk = true
		}
	}
	if !acceptedStatusCodesOk {
		t.Error("Incorrect default value (AcceptedStatusCodes)")
	}
	if providerConfig.SSLExpire != true {
		t.Error("Incorrect default value (SSLExpire)")
	}
}

func TestAddMonitorWithCustomConfig(t *testing.T) {
	service := setupService(t)
	emptyMonitorList(service)

	config := &endpointmonitorv1alpha1.UptimeKumaApiConfig {
		Type: "http",
		Interval: 20,
		RetryInterval: 21,
		ResendInterval: 22,
		MaxRetries: 23,
		Method: "HEAD",
		IgnoreTLS: true,
		UpsideDown: true,
		MaxRedirects: 24,
		AcceptedStatusCodes: []string{"200", "201"},
		SSLExpire: false,
	}
	monitor := models.Monitor{
		Name: "google.com",
		URL:  "https://www.google.com",
		Config: config,
	}

	service.Add(monitor)

	monitors := service.GetAll()
	if len(monitors) == 0 {
		t.Fatal("No service created")
	}
	if len(monitors) > 1 {
		t.Error("More than one service created")
	}

	if monitors[0].Name != "google.com" {
		t.Error("Created service has not correct name")
	}

	providerConfig, _ := monitors[0].Config.(*endpointmonitorv1alpha1.UptimeKumaApiConfig)

	if providerConfig == nil {
		t.Error("No provider config returned")
	}

	if providerConfig.Interval != 20 {
		t.Error("Unable to store custom provider configuration (Interval)")
	}
	if providerConfig.RetryInterval != 21 {
		t.Error("Unable to store custom provider configuration (RetryInterval)")
	}

	if providerConfig.ResendInterval!= 22 {
		t.Error("Unable to store custom provider configuration (ResendInterval)")
	}
	if providerConfig.MaxRetries != 23 {
		t.Error("Unable to store custom provider configuration (MaxRetries)")
	}
	if providerConfig.Method != "HEAD" {
		t.Error("Unable to store custom provider configuration (Method)")
	}
	if providerConfig.IgnoreTLS != true {
		t.Error("Unable to store custom provider configuration (IgnoreTLS)")
	}
	if providerConfig.UpsideDown != true {
		t.Error("Unable to store custom provider configuration (UpsideDown)")
	}
	if providerConfig.MaxRedirects != 24 {
		t.Error("Unable to store custom provider configuration (IgnoreTLS)")
	}
	acceptedStatusCodesOk := false
	for _, v := range providerConfig.AcceptedStatusCodes {
		if v == "200" {
			acceptedStatusCodesOk = true
		}
	}
	if !acceptedStatusCodesOk {
		t.Error("Unable to store custom provider configuration (AcceptedStatusCodes)")
	}
	acceptedStatusCodesOk = false
	for _, v := range providerConfig.AcceptedStatusCodes {
		if v == "201" {
			acceptedStatusCodesOk = true
		}
	}
	if !acceptedStatusCodesOk {
		t.Error("Unable to store custom provider configuration (AcceptedStatusCodes)")
	}
	if providerConfig.SSLExpire != false {
		t.Error("Unable to store custom provider configuration (SSLExpire)")
	}
}

func TestUpdateMonitorNoChange(t *testing.T) {
	service := setupService(t)
	emptyMonitorList(service)
	createMonitor(service, "google.com", "https://www.google.com")

	monitor, err := service.GetByName("google.com")

	if err != nil {
		t.Fatal("Unable to create monitor")
	}

	service.Update(*monitor)

	monitor, err = service.GetByName("google.com")

	if err != nil {
		t.Fatal("Unable to retreive monitor after update")
	}

	if monitor.Name != "google.com" {
		t.Error("Monitor name is not correct after update")
	}
	if monitor.URL != "https://www.google.com" {
		t.Error("Monitor URL is not correct after update")
	}
}

func TestUpdateMonitor(t *testing.T) {
	service := setupService(t)
	emptyMonitorList(service)
	createMonitor(service, "google.com", "https://www.google.com")

	monitor, err := service.GetByName("google.com")

	if err != nil {
		t.Fatal("Unable to create monitor")
	}

	monitor.URL = "http://dummy.com"

	service.Update(*monitor)

	monitor, err = service.GetByName("google.com")

	if err != nil {
		t.Fatal("Unable to retreive monitor after update")
	}

	if monitor.Name != "google.com" {
		t.Error("Monitor name is not correct after update")
	}
	if monitor.URL != "http://dummy.com" {
		t.Error("Monitor URL is not correct after update")
	}
}
