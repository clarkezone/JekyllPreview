package jobmanager

import (
	"context"
	"fmt"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	kl "temp.com/JekyllBlogPreview/kubelayer"
)

type ResourseStateType int

const (
	Create = 0
	Update
	Delete
)

type jobnotifier func(*batchv1.Job, ResourseStateType)

type Jobmanager struct {
	current_config    *rest.Config
	current_clientset kubernetes.Interface
	ctx               context.Context
	cancel            context.CancelFunc
	jobnotifiers      map[string]jobnotifier
}

func Newjobmanager(incluster bool, namespace string) (*Jobmanager, error) {
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
	created := jm.startWatchers(namespace)
	if created {
		return jm, nil
	}
	return nil, fmt.Errorf("unable to create jobmanager; startwatchers failed")
}

func newjobmanagerwithclient(internal bool, clientset kubernetes.Interface, namespace string) (*Jobmanager, error) {
	jm, err := newjobmanagerinternal(internal)
	if err != nil {
		return nil, err
	}

	jm.current_clientset = clientset

	//TODO only if we want watchers
	created := jm.startWatchers(namespace)
	if created {
		return jm, nil
	}
	return nil, fmt.Errorf("unable to create jobmanaer; startwatchers failed")
}

func newjobmanagerinternal(incluster bool) (*Jobmanager, error) {
	jm := Jobmanager{}

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

func (jm *Jobmanager) startWatchers(namespace string) bool {
	// We will create an informer that writes added pods to a channel.
	//	pods := make(chan *v1.Pod, 1)
	//informers := informers.NewSharedInformerFactory(jm.current_clientset, 0) // when watching in global scope, we need clusterrole / clusterrolebinding not role / rolebinding in the rbac setup
	var info informers.SharedInformerFactory
	if namespace == "" {
		info = informers.NewSharedInformerFactory(jm.current_clientset, 0) // when watching in global scope, we need clusterrole / clusterrolebinding not role / rolebinding in the rbac setup
	} else {
		info = informers.NewSharedInformerFactoryWithOptions(jm.current_clientset, 0, informers.WithNamespace(namespace))
	}
	podInformer := info.Core().V1().Pods().Informer()
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

	jobInformer := info.Batch().V1().Jobs().Informer()

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
	err := jobInformer.SetWatchErrorHandler(func(r *cache.Reflector, err error) {
		// your code goes here
		log.Printf("Bed Shat %v", err.Error())
		jm.cancel()
	})
	if err != nil {
		panic(err)
	}
	info.Start(jm.ctx.Done())

	// Ensuring that the informer goroutine have warmed up and called List before
	// we send any events to it.
	result := cache.WaitForCacheSync(jm.ctx.Done(), podInformer.HasSynced)
	result2 := cache.WaitForCacheSync(jm.ctx.Done(), jobInformer.HasSynced)
	if !result || !result2 {
		log.Printf("Bed Shat")
		return false
	}
	return true
}

func (jm *Jobmanager) CreateJob(name string, namespace string, image string, command []string, args []string, notifier jobnotifier) (*batchv1.Job, error) {
	//TODO: if job exists, delete it
	job, err := kl.CreateJob(jm.current_clientset, name, namespace, image, command, args, true)
	if err != nil {
		return nil, err
	}
	if notifier != nil {
		jm.jobnotifiers[string(job.Name)] = notifier
	}
	return job, nil
}

func (jm *Jobmanager) DeleteJob(name string) error {
	return kl.DeleteJob(jm.current_clientset, name)
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

func (jm *Jobmanager) Close() {
	jm.cancel()
}
