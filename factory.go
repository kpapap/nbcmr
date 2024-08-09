package nbcmr

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"
)


var (
	typeStr = component.MustNewType("nbcmr")
)

const (
	defaultInterval = 1 * time.Minute
)

// Factory is the receiver factory for the nbcmr receiver.
type Factory struct{}

// CreateBuild creates a new instance of the nbcmr receiver.
//
// This function takes the base configuration for the receiver, the component host, and a context.
// It returns a receiver.Logs interface that can be used by the OpenTelemetry Collector to
// create and start the receiver.
//
// Parameters:
// - ctx: the context.Context object for the receiver.
// - cfg: the component.Config object for the receiver.
// - host: the component.Host object for the receiver.
//
// Returns:
// - receiver.Logs: the receiver.Logs interface for the nbcmr receiver.
// - error: an error if the configuration is invalid.
func (f Factory) CreateBuild(ctx context.Context, cfg component.Config, host component.Host) (receiver.Logs, error) {
	// Check if the configuration is of type *Config
	nbcmrCfg, ok := cfg.(*Config)
	if !ok {
		// If the configuration is not of type *Config, return an error
		return nil, fmt.Errorf("invalid configuration for nbcmr receiver")
	}

	// Create a new instance of the nbcmr receiver with the configuration and host
	return createLogsReceiver(ctx, nbcmrCfg, host), nil
}


// createLogsReceiver creates a new instance of the nbcmr receiver.
//
// This function takes the base configuration for the receiver, the component host, and a context.
// It returns a receiver.Logs interface that can be used by the OpenTelemetry Collector to
// create and start the receiver.
//
// The context parameter is not used in this function and can be ignored.
//
// The component host parameter is not used in this function and can be ignored.
func createLogsReceiver(_ context.Context, baseCfg component.Config, _ component.Host) receiver.Logs {
	// Cast the base configuration to the specific configuration type used by the nbcmr receiver.
	nbcmrCfg := baseCfg.(*Config)

	// Create a new logger that will write to stdout.
	logger := log.New(os.Stdout, "", 0)

	// Create a new instance of the nbcmr receiver with the configuration and logger.
	logRcvr := &nbcmrReceiver{
		interval:    nbcmrCfg.Interval, // Set the interval to the value from the configuration.
		logger:      logger,            // Set the logger to the stdout logger.
		config:       nbcmrCfg,         // Set the configuration to the configuration passed in.
	}

	// Return the nbcmr receiver.
	return logRcvr
}


// NewFactory creates a new receiver factory for the nbcmr receiver.
//
// This function returns a receiver.Factory that can be used to create new instances of the nbcmr receiver.
// The receiver.Factory is used by the OpenTelemetry Collector to create new instances of the receiver at runtime.
//
// The NewFactory function takes no parameters and returns a receiver.Factory.
func NewFactory() receiver.Factory {
	// Create a new receiver.Factory using the typeStr and createDefaultConfig functions.
	// The receiver.Factory will be used to create new instances of the nbcmr receiver.
	return receiver.NewFactory(
		typeStr,      // The type of the receiver. This is used to identify the receiver in the collector configuration.
		createDefaultConfig, // A function that returns the default configuration for the receiver. This function is used when creating a new receiver instance.
	)
}


// createDefaultConfig returns the default configuration for the nbcmr receiver.
// This function is used when creating a new factory to provide a default configuration
// for the receiver.
func createDefaultConfig() component.Config {
	// Return a new instance of the Config struct with the default interval value set.
	// The Interval field is a string representation of the defaultInterval value.
	return &Config{
		// The Interval field is a string representing the duration between each
		// interval for the receiver to collect and send data to the collector.
		// The default value is 1 minute.
		Interval: defaultInterval.String(), // Set the default interval to 1 minute
	}
}
