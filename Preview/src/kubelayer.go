package main

import (
	"context"
	"fmt"
	"log"

	"k8s.io/client-go/kubernetes"

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

func PingApi(clientset kubernetes.Interface) {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}

// TODO: namespace, name, container image etc
func CreateJob(clientset kubernetes.Interface, name string, image string, command []string, args []string, always bool) (*batchv1.Job, error) {
	jobsClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)

	//TODO hook up pull policy
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: int32Ptr(1),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: apiv1.PodSpec{
					//Volumes: []apiv1.Volume{},
					Containers: []apiv1.Container{
						{
							Name:            name,
							Image:           image,
							ImagePullPolicy: "Always",
							//TODO: command and args optional
							//Command:         command,
							//Args:            args,
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
				},
			},
		},
	}
	if command != nil {
		job.Spec.Template.Spec.Containers[0].Command = command
	}
	if args != nil {
		job.Spec.Template.Spec.Containers[0].Args = args
	}
	result, err := jobsClient.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	log.Printf("Created job %v.\n", result.GetObjectMeta().GetName())
	return job, nil
}

func DeleteJob(clientset kubernetes.Interface, name string) error {
	jobsClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)
	meta := metav1.DeleteOptions{
		TypeMeta:           metav1.TypeMeta{},
		GracePeriodSeconds: new(int64),
		Preconditions:      &metav1.Preconditions{},
	}
	return jobsClient.Delete(context.TODO(), name, meta)
}

func int32Ptr(i int32) *int32 { return &i }
