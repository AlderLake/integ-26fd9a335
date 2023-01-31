
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
}

// equals fails the test if exp is not equal to act.
func equals(t *testing.T, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\texp: %#v\n\tgot: %#v\n", filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}

// TODO - everything

func TestNewQuoteFromCSV(t *testing.T) {
	symbol := "aapl"
	csv := `datetime,open,high,low,close,volume
2014-07-14 00:00,95.86,96.89,95.65,88.40,42810000.00
2014-07-15 00:00,96.80,96.85,95.03,87.36,45477900.00
2014-07-16 00:00,96.97,97.10,94.74,86.87,53396300.00
2014-07-17 00:00,95.03,95.28,92.57,85.32,57298000.00
2014-07-18 00:00,93.62,94.74,93.02,86.55,49988000.00
2014-07-21 00:00,94.99,95.00,93.72,86.10,39079000.00
2014-07-22 00:00,94.68,94.89,94.12,86.81,55197000.00
2014-07-23 00:00,95.42,97.88,95.17,89.08,92918000.00
2014-07-24 00:00,97.04,97.32,96.42,88.93,45729000.00
2014-07-25 00:00,96.85,97.84,96.64,89.52,43469000.00
2014-07-28 00:00,97.82,99.24,97.55,90.75,55318000.00
2014-07-29 00:00,99.33,99.44,98.25,90.17,43143000.00
2014-07-30 00:00,98.44,98.70,97.67,89.96,33010000.00
2014-07-31 00:00,97.16,97.45,95.33,87.62,56843000.00`
	q, _ := NewQuoteFromCSV(symbol, csv)
	//fmt.Println(q)
	if len(q.Close) != 14 {
		t.Error("Invalid length")
	}
	if q.Close[len(q.Close)-1] != 87.62 {
		t.Error("Invalid last value")
	}
}

func TestNewQuotesFromCSV(t *testing.T) {
	csv := `symbol,datetime,open,high,low,close,volume
spy,2018-07-12 00:00,278.28,279.43,277.60,273.95,60124700.00