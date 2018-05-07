package kube

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	DefaultResource = "default"
)

// MapToRuntimeObject maps the resource type string to the actual resource
func MapToRuntimeObject(resourceType string) runtime.Object {
	rType, ok := ResourceMap[resourceType]
	if !ok {
		return ResourceMap[DefaultResource]
	}
	return rType
}

var ResourceMap = map[string]runtime.Object{
	"pods":    &v1.Pod{},
	"default": nil,
}
