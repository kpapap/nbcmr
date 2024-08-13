package nbcmr

import (
	"context"
	"log"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)


type nbcmrReceiver struct {
	config *Config
	logger *zap.Logger
	nextConsumer consumer.Logs
}

func (c *nbcmrReceiver) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (c *nbcmrReceiver) ConsumeLogs(ctx context.Context, ld consumer.Logs) error {
	return nil
}

func (r *nbcmrReceiver) Start(ctx context.Context, host component.Host) error {
// Start the receiver
readconfigmaps()	
return nil
}


// Shutdown shuts down the receiver.
func (r *nbcmrReceiver) Shutdown(ctx context.Context) error {
	// Shutdown the receiver
	// Log a message indicating that the receiver is shutting down.
	log.Println("Shutting down receiver")
	// Return nil to indicate that the receiver shut down successfully.
	return nil
}

