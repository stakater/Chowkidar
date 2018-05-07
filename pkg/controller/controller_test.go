package controller

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stakater/Chowkidar/pkg/config"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	clientSet = getTestClient()
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
func TestControllerForPods(t *testing.T) {

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
	namespace := "dev"
	podName := "testpod-chowkidar"
	pod := createPod(namespace, podName)
	result, err := clientSet.CoreV1().Pods(namespace).Create(pod)
	if err != nil {
		panic(err)
	}
	log.Printf("Created Pod %q.\n", result.GetObjectMeta().GetName())

	time.Sleep(30 * time.Second)

	controller.clientset.CoreV1().Pods(namespace).Delete(podName, &metav1.DeleteOptions{})

}
func createPod(namespace string, podName string) *v1.Pod {
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

// 	// {
// 	// 	name: "ControllerForWrongType",
// 	// 	args: args{
// 	// 		clientset:        getTestClient(),
// 	// 		controllerConfig: configuration.Controllers[0],
// 	// 	},
// 	// 	wantErr: true,
// 	// },
// }
// for _, tt := range tests {
// 	t.Run(tt.name, func(t *testing.T) {
// 		got, err := NewController(tt.args.clientset, tt.args.controllerConfig)
// 		if (err != nil) != tt.wantErr {
// 			t.Errorf("NewController() error = %v, wantErr %v", err, tt.wantErr)
// 			return
// 		}
// 		if !reflect.DeepEqual(got, tt.want) {
// 			t.Errorf("NewController() = %v, want %v", got, tt.want)
// 		}
// 	})
// }
func getTestClient() *kubernetes.Clientset {
	var config *rest.Config
	kubeconfigPath := os.Getenv("HOME") + "/.kube/config"
	if _, err := os.Stat(kubeconfigPath); err == nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}
	client, _ := kubernetes.NewForConfig(config)
	return client

}
