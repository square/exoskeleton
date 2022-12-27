package exoskeleton

import (
	"path/filepath"
	"runtime"
)

var fixtures string

func init() {
	_, testfile, _, _ := runtime.Caller(0)
	fixtures = filepath.Join(testfile, "..", "fixtures")
}
