package main

import "testing"

func TestCreateJobExists(t *testing.T) {
	jm, err := newjobmanager()
	defer jm.close()
	if err != nil {
		t.Errorf("Unable to create JobManager")
	}
	_, err = jm.CreateJob("alpinetest", "alpine")
	if err != nil {
		t.Errorf("Unable to create job %v", err)
	}
	//TODO: wait for successful exit
	//TODO:    confirm watcher events fire just for job
	//TODO:         [x] Main does a create job with defer context and sighandler, exit works
	//TODO:         Ensure pod events fire
	//TODO:         Ensure job events fire (start / end)
	//TODO:         verbose logging
	//TODO:    add hook that watches for job event and confirms success on exit
	//TODO:    Ensure that watcher threads exit on sighandler with logging
	//TODO: flag for job to autodelete
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
