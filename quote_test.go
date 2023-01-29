
package quote

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// assert fails the test if the condition is false.
func assert(t *testing.T, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: "+msg+"\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		t.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(t *testing.T, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: unexpected error: %s\n", filepath.Base(file), line, err.Error())
		t.FailNow()
	}