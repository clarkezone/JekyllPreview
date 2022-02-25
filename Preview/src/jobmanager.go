package main

import (
	"context"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type ResourseStateType int

const (
	ResourceState ResourseStateType = 0
	Create                          = 1
	Update                          = 2
	Delete                          = 3
)

type jobnotifier func(*batchv1.Job, ResourseStateType)

type jobmanager struct {
	current_config    *rest.Config
	current_clientset kubernetes.Interface
	ctx               context.Context
	cancel            context.CancelFunc
	jobnotifiers      map[string]jobnotifier
}

func newjobmanager(incluster bool) (*jobmanager, error) {
	jm, err := newjobmanagerinternal(incluster)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(jm.current_config)
	if err != nil {
		return nil, err
	}
	jm.current_clientset = clientset

	//TODO only if we want watchers
	jm.startWatchers()
	return jm, nil
}

func newjobmanagerwithclient(internal bool, clientset kubernetes.Interface) (*jobmanager, error) {

	jm, err := newjobmanagerinternal(internal)
	if err != nil {
		return nil, err
	}

	jm.current_clientset = clientset

	//TODO only if we want watchers
	jm.startWatchers()
	return jm, nil
}

func newjobmanagerinternal(incluster bool) (*jobmanager, error) {
	jm := jobmanager{}

	ctx, cancel := context.WithCancel(context.Background())
	jm.ctx = ctx
	jm.cancel = cancel
	jm.jobnotifiers = make(map[string]jobnotifier)

	config, err := GetConfig(incluster)
	if config == nil {
		return nil, err
	}
	jm.current_config = config
	return &jm, nil
}

func (jm *jobmanager) startWatchers() {
	// We will create an informer that writes added pods to a channel.
	//	pods := make(chan *v1.Pod, 1)
	informers := informers.NewSharedInformerFactory(jm.current_clientset, 0)
	podInformer := informers.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			log.Printf("pod added: %s/%s", pod.Namespace, pod.Name)
			//	pods <- pod
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			log.Printf("pod deleted: %s/%s", pod.Namespace, pod.Name)
		},
	})

	jobInformer := informers.Batch().V1().Jobs().Informer()

	jobInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			job := obj.(*batchv1.Job)
			log.Printf("Job added: %s/%s uid:%v", job.Namespace, job.Name, job.UID)
			if val, ok := jm.jobnotifiers[job.Name]; ok {
				val(job, Create)
			}
		},
		DeleteFunc: func(obj interface{}) {
			job := obj.(*batchv1.Job)
			log.Printf("Job deleted: %s/%s uid:%v", job.Namespace, job.Name, job.UID)
			if val, ok := jm.jobnotifiers[job.Name]; ok {
				val(job, Delete)
				delete(jm.jobnotifiers, job.Name)
			}
		},
		UpdateFunc: func(oldobj interface{}, newobj interface{}) {
			oldjob := oldobj.(*batchv1.Job)
			newjob := newobj.(*batchv1.Job)
			log.Printf("Job updated: %s/%s status:%v uid:%v", oldjob.Namespace, oldjob.Name, newjob.Status, newjob.UID)

			if val, ok := jm.jobnotifiers[newjob.Name]; ok {
				val(newjob, Update)
			}
		},
	})
	// Make sure informers are running.
	informers.Start(jm.ctx.Done())

	// Ensuring that the informer goroutine have warmed up and called List before
	// we send any events to it.
	cache.WaitForCacheSync(jm.ctx.Done(), podInformer.HasSynced)
	cache.WaitForCacheSync(jm.ctx.Done(), jobInformer.HasSynced)
}

func (jm *jobmanager) CreateJob(name string, image string, command []string, args []string, notifier jobnotifier) (*batchv1.Job, error) {
	job, err := CreateJob(jm.current_clientset, name, image, command, args, true)
	if err != nil {
		return nil, err
	}
	if notifier != nil {
		jm.jobnotifiers[string(job.Name)] = notifier
	}
	return job, nil
}

func (jm *jobmanager) DeleteJob(name string) error {
	return DeleteJob(jm.current_clientset, name)
}

func GetConfig(incluster bool) (*rest.Config, error) {
	var config *rest.Config
	var err error
	if incluster {
		config, err = rest.InClusterConfig()
	} else {
		kubepath := "/users/jamesclarke/.kube/config"
		var kubeconfig *string = &kubepath
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	}
	return config, err
}

func (jm *jobmanager) close() {
	jm.cancel()
}
