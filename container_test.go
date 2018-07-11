package gopker_test

import (
	"testing"
	"github.com/blueskan/gopker"
	"log"
)

func TestContainer_Start(t *testing.T) {
	var container gopker.Container
	var err error

	container, err = gopker.NewContainer("hello-world")
	if err != nil {
		t.Fail()
	}

	_, err = container.Start()
	if err != nil {
		t.Fail()
	}

	log.Printf("Container start test passed!")
}

func TestContainer_Stop(t *testing.T) {
	container, err := gopker.NewContainer("alpine")
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

	log.Printf("Container stop test passed!")
}

func TestContainer_Kill(t *testing.T) {
	container, err := gopker.NewContainer("alpine")
	if err != nil {
		t.Fail()
	}

	_, err = container.Start()
	if err != nil {
		t.Fail()
	}

	err = container.Kill()
	if err != nil {
		t.Fail()
	}

	log.Printf("Container kill test passed!")
}