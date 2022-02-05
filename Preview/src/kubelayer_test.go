package main

import (
	"errors"
	"testing"

	"k8s.io/client-go/kubernetes/fake"
)

func TestApi(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	PingApi(clientset)
}

func TestCreateJobKubeLayer(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	_, err := CreateJob(clientset)
	if err != nil {
		panic(errors.New("Create job failed."))
	}
}
