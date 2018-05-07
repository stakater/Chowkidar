package main

import (
	"log"
	"os"

	config "github.com/stakater/Chowkidar/pkg/config"
	"github.com/stakater/Chowkidar/pkg/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.Println("Starting Chowkidar")
	// create the clientset
	clientset, err := getClient()
	if err != nil {
		log.Fatal(err)
	}

	// get the Controller config file
	config := getControllerConfig()

	// creating the controller
	for _, c := range config.Controllers {
		controller, err := controller.NewController(clientset, c)
		if err != nil {
			log.Printf("Error occured while creating controller. Reason: %s", err.Error())
			continue
		}
		// Now let's start the controller
		stop := make(chan struct{})
		defer close(stop)
		go controller.Run(1, stop)
	}

	// Wait forever
	select {}
}

// gets the client for k8s, if ~/.kube/config exists so get that config else incluster config
func getClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	//If file exists so use that config settings
	if _, err := os.Stat(kubeconfigPath); err == nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else { //Use Incluster Configuration
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

// get the yaml configuration for the controller
func getControllerConfig() config.Config {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if len(configFilePath) == 0 {
		//Default config file is placed in configs/ folder
		configFilePath = "configs/config.yaml"
	}
	configuration, err := config.ReadConfig(configFilePath)
	if err != nil {
		log.Panic(err)
	}
	return configuration
}
