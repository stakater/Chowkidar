package criterion

import (
	"testing"

	"github.com/stakater/Chowkidar/internal/pkg/config"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_matchesCriterion(t *testing.T) {
	type args struct {
		obj       interface{}
		criterion config.Criterion
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "CheckingPodsWithContainers&Resources",
			args: args{
				obj: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod-test",
						Namespace: "asd",
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							v1.Container{
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
				},
				criterion: config.Criterion{Identifiers: []string{"resourceExists"}},
			},
			want: false,
		},
		{
			name: "CheckingPodsWithContainersResourcesRequestsOnly",
			args: args{
				obj: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod-test",
						Namespace: "asd",
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							v1.Container{
								Resources: v1.ResourceRequirements{
									Requests: v1.ResourceList{
										"cpu":    resource.Quantity{},
										"memory": resource.Quantity{},
									},
								},
							},
						},
					},
				},
				criterion: config.Criterion{Identifiers: []string{"resourceExists"}},
			},
			want: true,
		},
		{
			name: "CheckingPodsWithContainersResourcesLimitsOnly",
			args: args{
				obj: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod-test",
						Namespace: "asd",
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							v1.Container{
								Resources: v1.ResourceRequirements{
									Limits: v1.ResourceList{
										"cpu":    resource.Quantity{},
										"memory": resource.Quantity{},
									},
								},
							},
						},
					},
				},
				criterion: config.Criterion{Identifiers: []string{"resourceExists"}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesCriterion(tt.args.obj, tt.args.criterion); got != tt.want {
				t.Errorf("matchesCriterion() = %v, want %v", got, tt.want)
			}
		})
	}
}
