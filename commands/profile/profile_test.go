package profile

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfileCmd_no_Args(t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	assert.Nil(t, NewProfile().Execute())

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	err := w.Close()
	if err != nil {
		panic(err)
	}
	os.Stdout = old // restoring the real stdout
	out := <-outC

	assert.Contains(t, out, "Create, list, switch airflow profiles.")
}
