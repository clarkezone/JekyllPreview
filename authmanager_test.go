package main

import (
	"testing"
)

func TestAuthmanagerTemplateParse(t *testing.T) {
	lrm := NewAuthManager()
	if lrm == nil {
		t.Fatalf("Failed to create authmanager")
	}
}
