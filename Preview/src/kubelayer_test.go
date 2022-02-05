package main

import (
	"errors"
	"testing"

	"k8s.io/client-go/kubernetes"
)

func TestGetConfig(t *testing.T) {
	GetConfig()
}

func TestApi(t *testing.T) {
	var config, err = GetConfig()
	if err != nil {
		t.Errorf("unable to create config")
	}
	PingApi(config)

	_, err = kubernetes.NewForConfig(config)
	if err != nil {
		t.Errorf("Unable to create clientset")
	}

}

func TestCreateJobKubeLayer(t *testing.T) {
	var config, _ = GetConfig()
	var clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		t.Errorf("unable to create clientset")
	}
	_, err = CreateJob(clientset)
	if err != nil {
		panic(errors.New("Create job failed."))
	}
}

//TODO this test should use mocks
