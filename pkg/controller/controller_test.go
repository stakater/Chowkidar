package controller

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stakater/Chowkidar/pkg/config"
	"github.com/stakater/Chowkidar/pkg/kube"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

var (
	clientSet, _ = kube.GetClient()
)

func TestControllerWithWrongTypeShouldNotCreate(t *testing.T) {
	configFilePath := "../../configs/testConfigs/WrongTypeConfig.yaml"
	configuration, err := config.ReadConfig(configFilePath)
	_, err = NewController(clientSet, configuration.Controllers[0])
	if err != nil {
		log.Printf("Unable to create NewController error = %v", err)
		return
	}
}

// Creating a Controller for Pod with Default Action without Resources so messages printed
func TestControllerForPodWithoutResourcesDefaultAction(t *testing.T) {

	configFilePath := "../../configs/testConfigs/CorrectConfig.yaml"
	configuration, err := config.ReadConfig(configFilePath)
	controller, err := NewController(clientSet, configuration.Controllers[0])
	if err != nil {
		log.Printf("Unable to create NewController error = %v", err)
		return
	}
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)
	time.Sleep(5 * time.Second)
	namespace := "test"
	podName := "testpod-withoutresources-chowkidar"
	pod := podWithoutResources(namespace, podName)
	result, err := clientSet.CoreV1().Pods(namespace).Create(pod)
	if err != nil {
		panic(err)
	}
	log.Printf("Created Pod %q.\n", result.GetObjectMeta().GetName())

	time.Sleep(10 * time.Second)
	log.Printf("Deleting Pod %q.\n", result.GetObjectMeta().GetName())
	controller.clientset.CoreV1().Pods(namespace).Delete(podName, &metav1.DeleteOptions{})
	time.Sleep(10 * time.Second)
}

// Creating a Controller for Pod with Default Action with Resources so no message printed
func TestControllerForPodWithResourcesDefaultAction(t *testing.T) {

	configFilePath := "../../configs/testConfigs/CorrectConfig.yaml"
	configuration, err := config.ReadConfig(configFilePath)
	controller, err := NewController(clientSet, configuration.Controllers[0])
	if err != nil {
		log.Printf("Unable to create NewController error = %v", err)
		return
	}
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)
	time.Sleep(5 * time.Second)
	namespace := "test"
	podName := "testpod-withresources-chowkidar"
	pod := podWithResources(namespace, podName)
	result, err := clientSet.CoreV1().Pods(namespace).Create(pod)
	if err != nil {
		panic(err)
	}
	log.Printf("Created Pod %q.\n", result.GetObjectMeta().GetName())

	time.Sleep(10 * time.Second)

	log.Printf("Deleting Pod %q.\n", result.GetObjectMeta().GetName())
	controller.clientset.CoreV1().Pods(namespace).Delete(podName, &metav1.DeleteOptions{})
	time.Sleep(10 * time.Second)

}

// Creating a Controller for Updating Pod with Default Action without Resources so messages printed
func TestControllerForUpdatePodShouldUpdateAndSendMessage(t *testing.T) {

	configFilePath := "../../configs/testConfigs/CorrectConfig.yaml"
	configuration, err := config.ReadConfig(configFilePath)
	controller, err := NewController(clientSet, configuration.Controllers[0])

	if err != nil {
		log.Printf("Unable to create NewController error = %v", err)
		return
	}
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)
	time.Sleep(5 * time.Second)
	namespace := "test"
	podName := "testpod-withoutresources-chowkidar"
	podClient := clientSet.CoreV1().Pods(namespace)
	pod := podWithoutResources(namespace, podName)
	pod, err = podClient.Create(pod)
	if err != nil {
		panic(err)
	}
	log.Printf("Created Pod %q.\n", pod.GetObjectMeta().GetName())
	time.Sleep(5 * time.Second)
	pod.Spec.Containers[0].Resources = v1.ResourceRequirements{
		Limits: v1.ResourceList{
			"cpu":    resource.Quantity{},
			"memory": resource.Quantity{},
		},
	}
	log.Printf("Updating Pod %q.\n", pod.GetObjectMeta().GetName())
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		pod, err = podClient.Get(podName, metav1.GetOptions{})
		if err != nil {

		}
		_, updateErr := podClient.Update(pod)
		return updateErr
	})
	if retryErr != nil {
		controller.clientset.CoreV1().Pods(namespace).Delete(podName, &metav1.DeleteOptions{})
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	time.Sleep(10 * time.Second)
	log.Printf("Deleting Pod %q.\n", pod.GetObjectMeta().GetName())
	controller.clientset.CoreV1().Pods(namespace).Delete(podName, &metav1.DeleteOptions{})
	time.Sleep(10 * time.Second)
}

func podWithResources(namespace string, podName string) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				v1.Container{
					Image: "hello-world",
					Name:  "testcontainer",
					Resources: v1.ResourceRequirements{
						Limits: v1.ResourceList{
							"cpu":    resource.Quantity{},
							"memory": resource.Quantity{},
						},
						Requests: v1.ResourceList{
							"cpu":    resource.Quantity{},
							"memory": resource.Quantity{},
						},
					},
				},
			},
		},
	}
}
func podWithoutResources(namespace string, podName string) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				v1.Container{
					Image: "hello-world",
					Name:  "testcontainer",
				},
			},
		},
	}
}
