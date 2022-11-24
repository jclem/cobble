package test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	out := new(strings.Builder)
	cmd := exec.Command("bin/cobble", "run", "-d", "../examples", "echo-b")
	cmd.Stdout = out
	cmd.Stderr = out
	err := cmd.Run()
	assert.NoError(t, err)

	assert.Equal(t, "Running task: echo-a\nA\nRunning task: echo-b\nB\n", out.String())
}
