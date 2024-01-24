package telemetry

import (
	"github.com/buildpacks/lifecycle/buildpack"
	"github.com/buildpacks/lifecycle/cmd"
	"os/exec"
	"time"
)

func ExitCode(err error) int {
	if bpe, ok := err.(*buildpack.Error); ok && bpe.Cause() != nil {
		return ExitCode(bpe.Cause())
	}

	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode()
	}
	if ef, ok := err.(*cmd.ErrorFail); ok {
		return ef.Code
	}
	return cmd.CodeForFailed
}

type Timer struct {
	start int64
}

func (t *Timer) Start() {
	t.start = time.Now().UnixMilli()
}

func (t *Timer) Stop() int64 {
	return time.Now().UnixMilli() - t.start
}
