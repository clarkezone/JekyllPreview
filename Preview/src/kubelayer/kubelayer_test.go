package kubelayer

import (
	"errors"
	"testing"

	"k8s.io/client-go/kubernetes/fake"
)

func TestApi(t *testing.T) {
	t.Logf("TestApi")
	clientset := fake.NewSimpleClientset()
	PingApi(clientset)
}

func TestCreateJobKubeLayer(t *testing.T) {
	t.Logf("TestCreateJobKubeLayer")
	clientset := fake.NewSimpleClientset()
	_, err := CreateJob(clientset, "testns", "", "", nil, nil, false)
	if err != nil {
		panic(errors.New("Create job failed."))
	}
}
