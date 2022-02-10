package main

import (
	"log"
	"testing"

	batchv1 "k8s.io/api/batch/v1"
)

func RunTestJob(completechannel chan struct{}, deletechannel chan struct{}, t *testing.T, command []string, notifier func(*batchv1.Job, ResourseStateType)) {
	jm, err := newjobmanager()
	defer jm.close()
	if err != nil {
		t.Fatalf("Unable to create JobManager")
	}

	_, err = jm.CreateJob("alpinetest", "alpine", command, nil, notifier)
	if err != nil {
		t.Fatalf("Unable to create job %v", err)
	}
	<-completechannel
	log.Println("Complated; attempting delete")
	err = jm.DeleteJob("alpinetest")
	if err != nil {
		t.Fatalf("Unable to delete job %v", err)
	}
	log.Println(("Deleted."))
	<-deletechannel
}

func TestCreateAndSucceed(t *testing.T) {
	completechannel := make(chan struct{})
	deletechannel := make(chan struct{})
	notifier := (func(job *batchv1.Job, typee ResourseStateType) {
		log.Printf("Got job in outside world %v", typee)

		if completechannel != nil && typee == Update && job.Status.Active == 0 && job.Status.Failed > 0 {
			log.Printf("Error detected")
			close(completechannel)
			completechannel = nil //avoid double close
		}

		if typee == Delete {
			log.Printf("Deleted")
			close(deletechannel)
		}
	})
	command := []string{"error"}
	RunTestJob(completechannel, deletechannel, t, command, notifier)
}

func TestCreateAndFail(t *testing.T) {
	jm, err := newjobmanager()
	defer jm.close()
	if err != nil {
		t.Fatalf("Unable to create JobManager")
	}
	completechannel := make(chan struct{})
	deletechannel := make(chan struct{})
	notifier := (func(job *batchv1.Job, typee ResourseStateType) {
		log.Printf("Got job in outside world %v Active %v Failed %v", typee, job.Status.Active, job.Status.Failed)

		if completechannel != nil && typee == Update && job.Status.Active == 0 && job.Status.Failed > 0 {
			log.Printf("Error detected")
			close(completechannel)
			completechannel = nil //avoid double close
		}

		if typee == Delete {
			log.Printf("Deleted")
			close(deletechannel)
		}
	})
	command := []string{"error"}
	_, err = jm.CreateJob("alpinetest", "alpine", command, nil, notifier)
	if err != nil {
		t.Fatalf("Unable to create job %v", err)
	}
	<-completechannel
	log.Println("Complated; attempting delete")
	err = jm.DeleteJob("alpinetest")
	if err != nil {
		t.Fatalf("Unable to delete job %v", err)
	}
	log.Println(("Deleted."))
	<-deletechannel
	//TODO: [x] add delete function
	//TODO: Move logic into test for succeeded / failed job incl delete.. does it work with mock
	//TODO: ability to inject volumes
	//TODO: support for namespace
	//TODO:         verbose logging
	//TODO:             Conditional log statements
	//TODO:             Environment variable
	//TODO: flag for job to autodelete
	//TODO: test that verifies auto delete
	//TODO: Ensure error if job with same name already exists
}

func TestGetConfig(t *testing.T) {
	var _, err = GetConfig()
	if err != nil {
		t.Errorf("unable to create config")
	}
	//TODO flag for job to autodelete
	//TODO wait for error exit
}

func TestCreateJobExitsError(t *testing.T) {

}

// test for other objects created doesn't fire job completion
// test for simple job create and exit

// test for job error state

// test for job that never returns and manually terminated
// test for job already exists
