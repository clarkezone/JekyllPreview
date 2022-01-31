package main

import (
	"errors"
	"testing"
)

func TestGetConfig(t *testing.T) {
	GetConfig()
}

func TestApi(t *testing.T) {
	var config = GetConfig()
	PingApi(config)
}

func TestCreateJob(t *testing.T) {
	var config = GetConfig()
	err := CreateJob(config)
	if err != nil {
		panic(errors.New("Create job failed."))
	}
}
