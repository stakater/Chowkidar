package main

import (
	"fmt"
	"log"
	"os"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// create the clientset
	clientset, err := getClient(true) // GetClient(UseExternal) true = Getting value from .kube/config
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
}
func getClient(useExternal bool) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
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
