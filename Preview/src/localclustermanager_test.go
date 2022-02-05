package main

import "testing"

func TestCreateJob(t *testing.T) {
	lcm, err := newlcm()
	if err != nil {
		t.Errorf("Unable to create LCM")
	}
	_, err = lcm.CreateJob()
	if err != nil {
		t.Errorf("Unable to create job")
	}
	//TODO wait for job to complete or have job be auto deleting

}
