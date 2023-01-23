package uptimekumaapi

import (
	"fmt"
	"reflect"
	"testing"

	endpointmonitorv1alpha1 "github.com/stakater/IngressMonitorController/v2/api/v1alpha1"
	"github.com/stakater/IngressMonitorController/v2/pkg/models"
)

func TestUptimeKumaApiMonitorToBaseMonitorMapper(t *testing.T) {
	uptimeKumaApiMonitorObject := UptimeKumaApiMonitor{
		Id:                  1,
		Type:                "http",
		Name:                "google.com",
		Interval:            42,
		RetryInterval:       43,
		ResendInterval:      44,
		MaxRetries:          45,
		Url:                 "http://dummy.com",
		Method:              "GET",
		IgnoreTLS:           true,
		UpsideDown:          true,
		MaxRedirects:        46,
		AcceptedStatusCodes: []string{"100", "142"},
		SSLExpire:           true,
	}

	monitorObject := UptimeKumaApiMonitorMonitorToBaseMonitorMapper(uptimeKumaApiMonitorObject)

	providerConfig, _ := monitorObject.Config.(*endpointmonitorv1alpha1.UptimeKumaApiConfig)

	if monitorObject.ID != "1" {
		t.Error("Mapper did not map the values correctly (ID)")
	}
	if monitorObject.Name != "google.com" {
		t.Error("Mapper did not map the values correctly (Name)")
	}
	if monitorObject.URL != "http://dummy.com" {
		t.Error("Mapper did not map the values correctly (URL)")
	}
	if providerConfig.Interval != 42 {
		t.Error("Mapper did not map the values correctly (Interval)")
	}
	if providerConfig.RetryInterval != 43 {
		t.Error("Mapper did not map the values correctly (RetryInterval)")
	}
	if providerConfig.ResendInterval != 44 {
		t.Error("Mapper did not map the values correctly (ResendInterval)")
	}
	if providerConfig.MaxRetries != 45 {
		t.Error("Mapper did not map the values correctly (MaxRetries)")
	}
	if providerConfig.Method != "GET" {
		t.Error("Mapper did not map the values correctly (Method)")
	}
	if providerConfig.IgnoreTLS != true {
		t.Error("Mapper did not map the values correctly (IgnoreTLS)")
	}
	if providerConfig.UpsideDown != true {
		t.Error("Mapper did not map the values correctly (UpsideDown)")
	}
	if providerConfig.MaxRedirects != 46 {
		t.Error("Mapper did not map the values correctly (MaxRedirects)")
	}
	acceptedStatusCodesOk := false
	for _, v := range providerConfig.AcceptedStatusCodes {
		if v == "100" {
			acceptedStatusCodesOk = true
		}
	}
	if !acceptedStatusCodesOk {
		t.Error("Mapper did not map the values correctly (AcceptedStatusCodes)")
	}
	acceptedStatusCodesOk = false
	for _, v := range providerConfig.AcceptedStatusCodes {
		if v == "142" {
			acceptedStatusCodesOk = true
		}
	}
	if !acceptedStatusCodesOk {
		t.Error("Mapper did not map the values correctly (AcceptedStatusCodes)")
	}
	if providerConfig.SSLExpire != true {
		t.Error("Mapper did not map the values correctly (SSLExpire)")
	}
}

func TestUptimeKumaApiMonitorsToBaseMonitorsMapper(t *testing.T) {
	m1 := UptimeKumaApiMonitor{
		Id:             1,
		Type:           "http",
		Name:           "myMonitor1",
		Interval:       1,
		RetryInterval:  1,
		ResendInterval: 1,
		MaxRetries:     1,
		Url:            "http://monitor1.tld",
		Method:         "GET",
	}
	m2 := UptimeKumaApiMonitor{
		Id:             2,
		Type:           "http",
		Name:           "myMonitor2",
		Interval:       2,
		RetryInterval:  2,
		ResendInterval: 2,
		MaxRetries:     2,
		Url:            "http://monitor2.tld",
		Method:         "GET",
	}

	// Create correct monitor for testing
	config1 := &endpointmonitorv1alpha1.UptimeKumaApiConfig{
		Type:           "http",
		Interval:       1,
		RetryInterval:  1,
		ResendInterval: 1,
		MaxRetries:     1,
		Method:         "GET",
	}
	config2 := &endpointmonitorv1alpha1.UptimeKumaApiConfig{
		Type:           "http",
		Interval:       2,
		RetryInterval:  2,
		ResendInterval: 2,
		MaxRetries:     2,
		Method:         "GET",
	}
	correctMonitors := []models.Monitor{
		{
			Name:   "myMonitor1",
			ID:     "1",
			URL:    "http://monitor1.tld",
			Config: config1,
		},
		{
			Name:   "myMonitor2",
			ID:     "2",
			URL:    "http://monitor2.tld",
			Config: config2,
		},
	}

	var kMonitors []UptimeKumaApiMonitor
	kMonitors = append(kMonitors, m1)
	kMonitors = append(kMonitors, m2)

	monitors := UptimeKumaApiMonitorMonitorsToBaseMonitorsMapper(kMonitors)

	for index := 0; index < len(monitors); index++ {
		if !reflect.DeepEqual(correctMonitors[index], monitors[index]) {
			t.Error(fmt.Sprintf("Monitor object should be the same (%#v != %#v)", correctMonitors[index].Config, monitors[index].Config))
		}
	}
}
