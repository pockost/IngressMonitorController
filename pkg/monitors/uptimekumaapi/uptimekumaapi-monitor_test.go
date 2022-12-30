package uptimekumaapi

import (
	"testing"

	"github.com/stakater/IngressMonitorController/pkg/config"
	"github.com/stakater/IngressMonitorController/pkg/models"
	"github.com/stakater/IngressMonitorController/pkg/util"
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

	if !service.Equal(*monitor1, monitor3) {
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
