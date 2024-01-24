package telemetry

import (
	"os"
	"strconv"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

var (
	EnvCorrelationIdKey  = "CORRELATION_ID"
	EnvSubscriptionIdKey = "BP_SUBSCRIPTION_ID"
	EnvCallerIdKey       = "CALLER_ID"

	EventLifecyclePhase = "lifecycle.phase"
	EventBuildpackBuild = "buildpack.build"

	InstrumentationKey = "55a59292-dbe1-4610-9be1-5ad2440bbbd7"
)

type TelemetrySender interface {
	Shutdown() error
	Log(metricName string, keyValuePair ...string) func(err error)
}

type Metric struct {
	Name         string
	Properties   map[string]string
	Measurements map[string]float64
}

type aisender struct {
	c appinsights.TelemetryClient
}

func NewAISender(ikey string) TelemetrySender {
	return &aisender{c: appinsights.NewTelemetryClient(ikey)}
}

func (a aisender) send(m Metric) {
	event := appinsights.NewEventTelemetry(m.Name)
	event.Properties = m.Properties
	event.Measurements = m.Measurements
	if event.Properties == nil {
		event.Properties = make(map[string]string)
	}
	event.Properties["correlationId"] = os.Getenv(EnvCorrelationIdKey)
	event.Properties["subscriptionId"] = os.Getenv(EnvSubscriptionIdKey)
	event.Properties["callerId"] = os.Getenv(EnvCallerIdKey)
	a.c.Track(event)
}

func (a aisender) Shutdown() error {
	select {
	case <-a.c.Channel().Close(10 * time.Second):
		// If we got here, then all telemetry was submitted
		// successfully, and we can proceed to exiting.
	case <-time.After(30 * time.Second):
	}
	return nil
}

func (s aisender) Log(name string, dims ...string) func(err error) {
	timer := Timer{}
	timer.Start()

	return func(err error) {
		metric := Metric{
			Name:         name,
			Properties:   make(map[string]string),
			Measurements: make(map[string]float64),
		}
		if err != nil {
			metric.Properties["result"] = "failed"
			metric.Properties["exitCode"] = strconv.Itoa(ExitCode(err))
			metric.Properties["errorMessage"] = err.Error()
		} else {
			metric.Properties["result"] = "success"
		}

		for i := 0; i < len(dims); i = i + 2 {
			metric.Properties[dims[i]] = dims[i+1]
		}

		metric.Measurements["durationInMs"] = float64(timer.Stop())
		s.send(metric)
	}
}
