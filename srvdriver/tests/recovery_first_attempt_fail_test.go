package tests

import (
	"testing"

	"github.com/makasim/flowstate/testcases"
	"github.com/makasim/flowstatesrv/srvdriver"
)

func TestRecoveryFirstAttemptFail(t *testing.T) {
	defer startSrv(t)()

	d := srvdriver.New(`http://127.0.0.1:8080`)

	testcases.RecoveryFirstAttemptFail(t, d, d)
}
