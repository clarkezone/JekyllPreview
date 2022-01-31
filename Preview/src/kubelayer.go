package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func GetConfig() *rest.Config {
	kubepath := "/users/jamesclarke/.kube/config"
	var kubeconfig *string = &kubepath
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return config
}

func PingApi(config *rest.Config) {
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	if len(pods.Items) == 0 {
		panic(errors.New("no pods found"))
	}
}

// TODO: namespace, name, container image etc
func CreateJob(config *rest.Config) error {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	jobsClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "MarkdownRenderJob",
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: int32Ptr(2),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: apiv1.PodSpec{
					//Volumes: []apiv1.Volume{},
					Containers: []apiv1.Container{
						{
							Name:            "nginx",
							Image:           "nginx",
							ImagePullPolicy: "Always",
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
				},
			},
		},
	}
	result, err := jobsClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	log.Printf("Created job %v.\n", result.GetObjectMeta().GetName())
	return nil
}
