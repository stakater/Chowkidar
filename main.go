package main

import (
	"log"
	"os"

	config "github.com/stakater/Chowkidar/pkg/config"
	"github.com/stakater/Chowkidar/pkg/controller"
	"github.com/stakater/Chowkidar/pkg/kube"
)

func main() {
	log.Println("Starting Chowkidar")
	// create the clientset
	clientset, err := kube.GetClient()
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
