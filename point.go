package reducer

import (
	"errors"
	"reflect"
	"sort"
	"sync"
)

const (
	minLengthForChunking = 100000
	// VisvalingamAlg reduce algorithm method
	VisvalingamAlg = "visvalingam"
)

// Point in 2D Axis
type Point interface {
	X() float64
	Y() float64
}

// New Points
func New(data interface{}) ([]Point, error) {
	var points []Point
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)

		for i := 0; i < s.Len(); i++ {
			p, ok := s.Index(i).Interface().(Point)
			if !ok {
				return nil, errors.New("could not convert to point")
			}
			points = append(points, p)
		}
		return points, nil
	default:
		return nil, errors.New("input is not an interface")
	}
}

// Reduce points
func Reduce(data interface{}, minKeep, chunkSize int, algorithm string) []Point {
	points, err := New(data)
	if err != nil {
		return nil
	}
	if len(points) <= chunkSize*2 || len(points) < minLengthForChunking {

		if algorithm == VisvalingamAlg {
			return Visvalingam(points, minKeep)
		}

	}
	chunksMinKeep := minKeep / chunkSize

	jobs := new(sync.WaitGroup)

	chunksPoints := make(chan []Point, 8)

	for i := 0; i < chunkSize; i++ {
		if i+1 == chunkSize {
			chunksMinKeep += minKeep % chunkSize
		}
		jobs.Add(1)
		go func(index int) {
			defer jobs.Done()
			pts := points[index*chunkSize : (index+1)*chunkSize]
			if algorithm == VisvalingamAlg {
				chunksPoints <- Visvalingam(pts, chunksMinKeep)
			}
		}(i)
	}
	go func() {
		jobs.Wait()
		close(chunksPoints)
	}()

	var newPoints []Point
	for c := range chunksPoints {
		newPoints = append(newPoints, c...)
	}
	sort.Slice(newPoints[:], func(i, j int) bool {
		valI := newPoints[i].X()
		valJ := newPoints[j].X()
		return valI > valJ
	})

	return newPoints
}
