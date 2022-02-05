package main

import (
	"context"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type jobmanager struct {
	current_config    *rest.Config
	current_clientset *kubernetes.Clientset
	ctx               context.Context
	cancel            context.CancelFunc
}

func newjobmanager() (*jobmanager, error) {
	jm := jobmanager{}

	ctx, cancel := context.WithCancel(context.Background())
	jm.ctx = ctx
	jm.cancel = cancel

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

	//TODO only if we want watchers
	jm.startWatchers()
	return &jm, nil
}

func (jm *jobmanager) startWatchers() {
	//watcherStarted := make(chan struct{})
	// A catch-all watch reactor that allows us to inject the watcherStarted channel.
	//	jm.current_clientset.PrependWatchReactor("*", func(action clienttesting.Action) (handled bool, ret watch.Interface, err error) {
	//		gvr := action.GetResource()
	//		ns := action.GetNamespace()
	//		watch, err := jm.current_clientset.Tracker().Watch(gvr, ns)
	//		if err != nil {
	//			return false, nil, err
	//		}
	//		close(watcherStarted)
	//		return true, watch, nil
	//	})

	// We will create an informer that writes added pods to a channel.
	pods := make(chan *v1.Pod, 1)
	informers := informers.NewSharedInformerFactory(jm.current_clientset, 0)
	podInformer := informers.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			log.Printf("pod added: %s/%s", pod.Namespace, pod.Name)
			pods <- pod
		},
	})

	// Make sure informers are running.
	informers.Start(jm.ctx.Done())

	// Ensuring that the informer goroutine have warmed up and called List before
	// we send any events to it.
	cache.WaitForCacheSync(jm.ctx.Done(), podInformer.HasSynced)

	//<-watcherStarted
}

func (jm *jobmanager) CreateJob(name string, image string) (*batchv1.Job, error) {
	return CreateJob(jm.current_clientset, name, image, true)
}

func GetConfig() (*rest.Config, error) {
	kubepath := "/users/jamesclarke/.kube/config"
	var kubeconfig *string = &kubepath
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	return config, err
}

func (jm *jobmanager) close() {
	jm.cancel()
}
