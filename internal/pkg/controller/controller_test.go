package controller

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stakater/Chowkidar/internal/pkg/config"
	"github.com/stakater/Chowkidar/pkg/kube"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

var (
	clientSet, _     = kube.GetClient()
	configFilePath   = "../../../configs/testConfigs/CorrectConfig.yaml"
	configuration, _ = config.ReadConfig(configFilePath)

	podNamePrefix = "testpod-chowkidar"
	letters       = []rune("abcdefghijklmnopqrstuvwxyz")
)

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func TestControllerWithWrongTypeShouldNotCreate(t *testing.T) {
	_, err := NewController(clientSet, configuration.Controllers[0])
	if err != nil {
		log.Printf("Unable to create NewController error = %v", err)
		return
	}
}

// Creating a Controller for Pod with Default Action without Resources so messages printed
func TestControllerForPodWithoutResourcesDefaultAction(t *testing.T) {
	controller, err := NewController(clientSet, configuration.Controllers[0])
	if err != nil {
		log.Printf("Unable to create NewController error = %v", err)
		return
	}
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)
	time.Sleep(10 * time.Second)
	namespace := "test"
	podName := podNamePrefix + "-withoutresources-" + randSeq(5)
	pod := podWithoutResources(namespace, podName)
	result, err := clientSet.CoreV1().Pods(namespace).Create(pod)
	if err != nil {
		panic(err)
	}
	log.Printf("Created Pod %q.\n", result.GetObjectMeta().GetName())

	time.Sleep(10 * time.Second)
	log.Printf("Deleting Pod %q.\n", result.GetObjectMeta().GetName())
	controller.clientset.CoreV1().Pods(namespace).Delete(podName, &metav1.DeleteOptions{})
	time.Sleep(15 * time.Second)
}

// Creating a Controller for Pod with Default Action with Resources so no message printed
func TestControllerForPodWithResourcesDefaultAction(t *testing.T) {
	controller, err := NewController(clientSet, configuration.Controllers[0])
	if err != nil {
		log.Printf("Unable to create NewController error = %v", err)
		return
	}
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)
	time.Sleep(10 * time.Second)
	namespace := "test"
	podName := podNamePrefix + "-withresources-" + randSeq(5)
	pod := podWithResources(namespace, podName)
	result, err := clientSet.CoreV1().Pods(namespace).Create(pod)
	if err != nil {
		panic(err)
	}
	log.Printf("Created Pod %q.\n", result.GetObjectMeta().GetName())
	time.Sleep(15 * time.Second)

	log.Printf("Deleting Pod %q.\n", result.GetObjectMeta().GetName())
	controller.clientset.CoreV1().Pods(namespace).Delete(podName, &metav1.DeleteOptions{})
	time.Sleep(15 * time.Second)

}

// Creating a Controller for Updating Pod with Default Action without Resources so messages printed
func TestControllerForUpdatePodShouldUpdateDefaultAction(t *testing.T) {
	controller, err := NewController(clientSet, configuration.Controllers[0])
	if err != nil {
		log.Printf("Unable to create NewController error = %v", err)
		return
	}
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)
	time.Sleep(10 * time.Second)
	namespace := "test"
	podName := podNamePrefix + "-withoutresources-update-" + randSeq(5)
	podClient := clientSet.CoreV1().Pods(namespace)
	pod := podWithoutResources(namespace, podName)
	pod, err = podClient.Create(pod)
	if err != nil {
		panic(err)
	}
	log.Printf("Created Pod %q.\n", pod.GetObjectMeta().GetName())
	time.Sleep(10 * time.Second)
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
	time.Sleep(15 * time.Second)
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
