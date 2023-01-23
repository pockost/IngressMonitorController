package uptimekumaapi

import (
	"strconv"

	endpointmonitorv1alpha1 "github.com/stakater/IngressMonitorController/v2/api/v1alpha1"
	"github.com/stakater/IngressMonitorController/v2/pkg/models"
)

/**
 * Convert a UptimeKumaApiMonitor to base Monitor
 **/
func UptimeKumaApiMonitorMonitorToBaseMonitorMapper(monitor UptimeKumaApiMonitor) * models.Monitor {

	var m models.Monitor

	m.Name = monitor.Name
	m.URL = monitor.Url
	m.ID = strconv.Itoa(monitor.Id)

	var providerConfig endpointmonitorv1alpha1.UptimeKumaApiConfig
	providerConfig.Type = monitor.Type
	providerConfig.Interval = monitor.Interval
	providerConfig.RetryInterval = monitor.RetryInterval
	providerConfig.ResendInterval = monitor.ResendInterval
	providerConfig.MaxRetries = monitor.MaxRetries
	providerConfig.Method = monitor.Method
	providerConfig.IgnoreTLS = monitor.IgnoreTLS
	providerConfig.UpsideDown = monitor.UpsideDown
	providerConfig.MaxRedirects = monitor.MaxRedirects
	providerConfig.AcceptedStatusCodes = monitor.AcceptedStatusCodes
	providerConfig.SSLExpire = monitor.SSLExpire

	m.Config = &providerConfig

	return &m
}

func UptimeKumaApiMonitorMonitorsToBaseMonitorsMapper(uptimeKumaApiMonitors []UptimeKumaApiMonitor) []models.Monitor {
	var monitors []models.Monitor

	for index := 0; index < len(uptimeKumaApiMonitors); index++ {
		monitors = append(monitors, *UptimeKumaApiMonitorMonitorToBaseMonitorMapper(uptimeKumaApiMonitors[index]))
	}

	return monitors
}
