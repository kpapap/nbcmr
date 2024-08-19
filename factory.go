package nbcmr

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

var (typeStr = component.MustNewType("nbcmr"))
const (defaultInterval = "1m")
const (defaultConfigMapName = "nbcmr-cm")


// NewFactory creates a new receiver factory for the nbcmr receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithLogs(createLogsReceiver, component.StabilityLevelUndefined),
	)
}

// createLogsReceiver creates a new instance of the nbcmr receiver.
func createLogsReceiver(_ context.Context, settings receiver.Settings, cfg component.Config, consumer consumer.Logs) (receiver.Logs, error) {
	// Create the new receiver
	rCfg := cfg.(*Config)
	return newNbcmrReceiver(rCfg, consumer, settings)
}

// createDefaultConfig returns the default configuration for the nbcmr receiver.
func createDefaultConfig() component.Config {
	return &Config{
		Interval: defaultInterval,
		ConfigMapName: defaultConfigMapName,
	}
}


