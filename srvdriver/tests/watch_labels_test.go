package tests

import (
	"testing"

	"github.com/makasim/flowstate/testcases"
	"github.com/makasim/flowstatesrv/srvdriver"
)

func TestWatchLabels(t *testing.T) {
	defer startSrv(t)()

	d := srvdriver.New(`h2c://127.0.0.1:8080`)

	testcases.WatchLabels(t, d, d)
}