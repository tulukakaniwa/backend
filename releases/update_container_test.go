package releases

import (
	"testing"
)

func TestUpdateContainer(t *testing.T) {
	if !updateDockerImage() {
		t.Fail()
	}
}
