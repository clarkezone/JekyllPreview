package main

import (
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type localclustermanager struct {
	current_config    *rest.Config
	current_clientset *kubernetes.Clientset
}

func newlcm() (*localclustermanager, error) {
	lcm := localclustermanager{}
	config, err := GetConfig()
	if lcm.current_config == nil {
		return nil, err
	}
	lcm.current_config = config

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	lcm.current_clientset = clientset

	return &lcm, nil
}

func (lcm *localclustermanager) CreateJob() (*batchv1.Job, error) {
	return CreateJob(lcm.current_clientset)
}
