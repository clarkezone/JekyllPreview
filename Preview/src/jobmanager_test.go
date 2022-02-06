package main

import "testing"

func TestCreateJobExists(t *testing.T) {
	jm, err := newjobmanager()
	defer jm.close()
	if err != nil {
		t.Errorf("Unable to create JobManager")
	}
	command := []string{"error"}
	_, err = jm.CreateJob("alpinetest", "alpine", command, []string{""})
	if err != nil {
		t.Errorf("Unable to create job %v", err)
	}
	//TODO: wait for successful exit
	//TODO:    confirm watcher events fire just for job
	//TODO:         [x] Main does a create job with defer context and sighandler, exit works
	//TODO:         [x] Ensure pod events fire
	//TODO:         [x] Ensure job events fire (start / end)
	//TODO:         [x] watcher started
	//TODO:         [x] inject command and args to inject error
	//TODO:         [x] command optional to enable both success and failure
	//TODO:    add hook or dequer that watches for job event and confirms success on exit
	//TODO:    add delete function
	//TODO:    Ensure that watcher threads exit on sighandler with logging
	//TODO: flag for job to autodelete
	//TODO: test that verifies auto delete
	//TODO: test that verifies deliberate job and or pod failure
	//TODO: ability to inject volumes
	//TODO:         verbose logging
	//TODO:             Conditional log statements
	//TODO:             Environment variable
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
