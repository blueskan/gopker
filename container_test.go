package gopker_test

import (
	"testing"
	"github.com/blueskan/gopker"
	"log"
)

func TestContainer_Start(t *testing.T) {
	container, err := gopker.Container("hello-world")
	if err != nil {
		t.Fail()
	}

	_, err = container.Start()
	if err != nil {
		t.Fail()
	}

	log.Printf("Test passed!")
}

func TestContainer_Stop(t *testing.T) {
	container, err := gopker.Container("alpine")
	if err != nil {
		t.Fail()
	}

	_, err = container.Start()
	if err != nil {
		t.Fail()
	}

	err = container.Stop()
	if err != nil {
		t.Fail()
	}

	log.Printf("Test passed!")
}