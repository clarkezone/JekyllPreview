package main

import (
	"testing"
)

func TestGetConfig(t *testing.T) {
	GetConfig()
}

func TestApi(t *testing.T) {
	var config = GetConfig()
	PingApi(config)

}
