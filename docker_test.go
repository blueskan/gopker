package gopker_test

import (
	"github.com/blueskan/gopker"
	"log"
	"testing"
)

func TestDocker_Containers(t *testing.T) {
	helloWorldContainer, err := gopker.Container("hello-world")
	if err != nil {
		t.Fail()
	}
	helloWorldContainer.Start()

	containers, err := gopker.Containers()
	if err != nil {
		t.Fail()
	}

	found := false

	for _, container := range containers {
		if container.Image == "hello-world" {
			found = true
			break
		}
	}

	if !found {
		t.Fail()
	}

	log.Printf("Test passed!")
}
