package nbcmr

import (
	"context"
	"os"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Define the structure of the YAML data
// ConfigMapMapping represents the mapping of a ConfigMap to a namespace.
// It contains the name of the ConfigMap and the namespace it belongs to.
type ConfigMapMapping struct {
	// Name is the name of the ConfigMap.
	Name string `yaml:"name"`
	// Namespace is the namespace in which the ConfigMap is located.
	Namespace string `yaml:"namespace"`
}

type nbcmrReceiver struct {
	cfg      			*Config
	nextConsumer	consumer.Logs
	settings 			receiver.Settings
	shutdownWG  	sync.WaitGroup
}

func newNbcmrReceiver(cfg *Config, nextConsumer consumer.Logs, settings receiver.Settings) (*nbcmrReceiver, error) {
	r := &nbcmrReceiver{
		cfg:        	cfg,
		nextConsumer:	nextConsumer,
		settings:			settings,
	}
	return r, nil
}

func (r *nbcmrReceiver) Start(ctx context.Context, host component.Host) error {
	// Start the receiver
	// Load Kubernetes cluster configuration
	r.settings.Logger.Info("Loading Kubernetes cluster configuration")
	clusterconfig, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		r.settings.Logger.Sugar().Errorf("Error building kubeconfig: %s", clusterconfig)
	}
	r.settings.Logger.Info("Kubernetes cluster configuration loaded")

	// Create Kubernetes client
	r.settings.Logger.Info("Creating Kubernetes client")
	clientset, err := kubernetes.NewForConfig(clusterconfig)
	if err != nil {
		r.settings.Logger.Sugar().Errorf("Error creating Kubernetes client: %s", clientset)
	}
	r.settings.Logger.Info("Kubernetes client created")

	// Read the YAML data from the environment variable
	r.settings.Logger.Info("Reading YAML data from environment variable")
	yamlData := os.Getenv("CONFIGMAP_LIST")
	if yamlData == "" {
		r.settings.Logger.Error("CONFIGMAP_LIST environment variable is not set")
	}
	r.settings.Logger.Info("YAML data read from environment variable")

	// Decode the YAML data into a struct
	r.settings.Logger.Info("Decoding YAML data")
	var configYAML []ConfigMapMapping
	err = yaml.Unmarshal([]byte(yamlData), &configYAML)
	if err != nil {
		r.settings.Logger.Sugar().Errorf("Error decoding YAML data: %s", err)
	}
	r.settings.Logger.Info("YAML data decoded")

	// Create a map of ConfigMap names and their corresponding namespaces
	r.settings.Logger.Info("Creating map of ConfigMap names and namespaces")
	configMapMap := make(map[string]string)

	// Populate the configMapMap from the YAML data
	for _, mapping := range configYAML {
		configMapMap[mapping.Name] = mapping.Namespace
	}
	r.settings.Logger.Info("Map of ConfigMap names and namespaces created")

	// Create a ticker to repeat the code
	r.settings.Logger.Info("Creating ticker")
	// Get the repeat time from environment variable
	repeatTimeStr := os.Getenv("INTERVAL")
	if repeatTimeStr == "" {
		r.settings.Logger.Error("INTERVAL environment variable is not set")
	}
	repeatTime, err := time.ParseDuration(repeatTimeStr)
	if err != nil {
		r.settings.Logger.Sugar().Errorf("Error parsing INTERVAL environment variable: %s", repeatTimeStr)
	}
	ticker := time.NewTicker(repeatTime)
	defer ticker.Stop()
	r.settings.Logger.Info("Ticker created")

	// Run the code in a loop with a ticker
	r.settings.Logger.Info("Starting loop")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			r.settings.Logger.Info("Listing selected ConfigMaps:")
			for name, namespace := range configMapMap {
				r.settings.Logger.Sugar().Infof("Getting ConfigMap %s in namespace %s", name, namespace)
				configmap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
				if err != nil {
					r.settings.Logger.Sugar().Errorf("Error getting ConfigMap %s in namespace %s: %s", name, namespace, err)
					continue
				}
				if configmap != nil {
					r.settings.Logger.Sugar().Infof("Namespace: %s, Name: %s, Data: %v", configmap.Namespace, configmap.Name, configmap.Data)
				}
			}
		}
	}
}

// Shutdown shuts down the receiver.
func (r *nbcmrReceiver) Shutdown(ctx context.Context) error {
	var err error
	r.shutdownWG.Wait()
	// Log a message indicating that the receiver is shutting down.
	r.settings.Logger.Info("Shutting down receiver")
	// Return err to indicate that the receiver shut down successfully.
	return err
}

