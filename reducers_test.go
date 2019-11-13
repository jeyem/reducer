package reducers

import (
	"fmt"
	"testing"
	"time"
)

type testPoint struct {
	name      string
	value     float64
	timestamp int64
}

func (p testPoint) X() float64 {
	return float64(p.timestamp)
}

func (p testPoint) Y() float64 {
	return p.value
}

var testData = []testPoint{
	{"a", 1.1, time.Now().Add(-time.Second * 10).Unix()},
	{"b", 1.1, time.Now().Add(-time.Second * 9).Unix()},
	{"c", 1.31, time.Now().Add(-time.Second * 8).Unix()},
	{"d", 1.21, time.Now().Add(-time.Second * 7).Unix()},
	{"e", 1.31, time.Now().Add(-time.Second * 6).Unix()},
	{"f", 2.134, time.Now().Add(-time.Second * 5).Unix()},
	{"g", 2.12, time.Now().Add(-time.Second * 4).Unix()},
	{"h", 2.08, time.Now().Add(-time.Second * 3).Unix()},
	{"i", 1.96, time.Now().Add(-time.Second * 2).Unix()},
	{"j", 0.008, time.Now().Add(-time.Second * 1).Unix()},
	{"k", 4.34, time.Now().Add(-time.Second).Unix()},
	{"l", 1.198, time.Now().Unix()},
}

func TestVisvalingam(t *testing.T) {
	points, err := New(testData)
	if err != nil {
		t.Error(err)
		return
	}
	pts := Reduce(points, 5, 3, VisvalingamAlg)
	for _, p := range pts {
		fmt.Println(p.(testPoint))
	}

}
