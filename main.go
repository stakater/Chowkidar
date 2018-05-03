package main

import (
	"fmt"
	"log"
	"os"

	config "github.com/stakater/Chowkidar/pkg/config"
	"github.com/stakater/Chowkidar/pkg/controller"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// create the clientset
	clientset, err := getClient()
	if err != nil {
		log.Fatal(err)
	}
	nodes, err := clientset.Core().Nodes().List(meta_v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Nodes:  ")
	for _, node := range nodes.Items {
		fmt.Println(node.Name)
	}
	pods, err := clientset.Core().Pods("").List(meta_v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Pods:  ")
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}

	config := getControllerConfig()
	controller := controller.NewController(clientset, config.Controllers[0])
	// Now let's start the controller
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)

	// Wait forever
	select {}
}

func getClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	//If file exists so use that
	if _, err := os.Stat(kubeconfigPath); err == nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
func getControllerConfig() config.Config {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if len(configFilePath) == 0 {
		configFilePath = "configs/config.yaml"
	}

	configuration := config.ReadConfig(configFilePath)
	return configuration
}
