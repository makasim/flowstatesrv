package tests

import (
	"testing"

	"github.com/makasim/flowstate/testcases"
	"github.com/makasim/flowstatesrv/srvdriver"
)

func TestQueue(t *testing.T) {
	defer startSrv(t)()

	d := srvdriver.New(`http://127.0.0.1:8080`)

	testcases.Queue(t, d, d)
}
