package tests

import (
	"testing"

	"github.com/makasim/flowstate/testcases"
	"github.com/makasim/flowstatesrv/srvdriver"
)

func TestDelay_Paused(t *testing.T) {
	defer startSrv(t)()

	d := srvdriver.New(`http://127.0.0.1:8080`)

	testcases.Delay_Paused(t, d, d)
}
