package main

import "testing"

func TestCreateJob(t *testing.T) {
	jm, err := newjobmanager()
	if err != nil {
		t.Errorf("Unable to create JobManager")
	}
	_, err = jm.CreateJob()
	if err != nil {
		t.Errorf("Unable to create job")
	}
	//TODO wait for job to complete and delete or have job be auto deleting

}

func TestGetConfig(t *testing.T) {
	var _, err = GetConfig()
	if err != nil {
		t.Errorf("unable to create config")
	}
}

func TestCreateSimpleJobthatExits(t *testing.T) {

}

// test for simple job create and exit

// test for job error state

// test for job that never returns and manually terminated
// test for job already exists
