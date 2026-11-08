package generate

import (
	"context"
	"testing"
)

// Recording must be safe without an Argvio client, whatever the inputs.
func TestRecordInput(t *testing.T) {
	for _, fromFile := range []bool{true, false} {
		for _, configFile := range []bool{true, false} {
			recordInput(context.Background(), fromFile, configFile)
		}
	}
}
