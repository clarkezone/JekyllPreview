package main

import (
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type jobmanager struct {
	current_config    *rest.Config
	current_clientset *kubernetes.Clientset
}

func newjobmanager() (*jobmanager, error) {
	jm := jobmanager{}
	config, err := GetConfig()
	if config == nil {
		return nil, err
	}
	jm.current_config = config

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	jm.current_clientset = clientset

	return &jm, nil
}

func (lcm *jobmanager) CreateJob() (*batchv1.Job, error) {
	return CreateJob(lcm.current_clientset)
}

func GetConfig() (*rest.Config, error) {
	kubepath := "/users/jamesclarke/.kube/config"
	var kubeconfig *string = &kubepath
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	return config, err
}
