package reducers

import (
	"math"
)

// Visvalingam algorithm
func Visvalingam(points []Point, minKeep int) []Point {
	removed := 0
	pointsCount := len(points)

	if pointsCount <= minKeep {
		return points
	}

	// build the initial minheap linked list.
	heap := minHeap(make([]*visItem, 0, pointsCount))

	linkedListStart := &visItem{
		area:       math.Inf(1),
		pointIndex: 0,
	}
	heap.Push(linkedListStart)

	// internal path items
	items := make([]visItem, pointsCount, pointsCount)

	previous := linkedListStart
	for i := 1; i < pointsCount-1; i++ {
		item := &items[i]

		item.area = doubleTriangleArea(points[i-1], points[i], points[i+1])
		item.pointIndex = i
		item.previous = previous

		heap.Push(item)
		previous.next = item
		previous = item
	}

	// final item
	endItem := &items[pointsCount-1]
	endItem.area = math.Inf(1)
	endItem.pointIndex = pointsCount - 1
	endItem.previous = previous

	previous.next = endItem
	heap.Push(endItem)

	// run through the reduction process
	for len(heap) > 0 {
		current := heap.Pop()

		if pointsCount-removed <= minKeep {
			break
		}

		next := current.next
		previous := current.previous

		// remove current element from linked list
		previous.next = current.next
		next.previous = current.previous
		removed++

		// figure out the new areas
		if previous.previous != nil {
			area := doubleTriangleArea(
				points[previous.previous.pointIndex],
				points[previous.pointIndex],
				points[next.pointIndex],
			)

			area = math.Max(area, current.area)
			heap.Update(previous, area)
		}

		if next.next != nil {
			area := doubleTriangleArea(
				points[previous.pointIndex],
				points[next.pointIndex],
				points[next.next.pointIndex],
			)

			area = math.Max(area, current.area)
			heap.Update(next, area)
		}
	}

	item := linkedListStart
	newPoints := make([]Point, 0, len(heap)+2)

	for item != nil {
		newPoints = append(newPoints, points[item.pointIndex])
		item = item.next
	}
	return newPoints
}

// Stuff to create the priority queue, or min heap.
// Rewriting it here, vs using the std lib, resulted in a 10x performance bump!
type minHeap []*visItem

type visItem struct {
	area       float64 // triangle area
	pointIndex int     // index of point in original path

	// to keep a virtual linked list to help rebuild the triangle areas as we remove points.
	next     *visItem
	previous *visItem

	index int // interal index in heap, for removal and update
}

func (h *minHeap) Push(item *visItem) {
	item.index = len(*h)
	*h = append(*h, item)
	h.up(item.index)
}

func (h *minHeap) Pop() *visItem {
	removed := (*h)[0]
	lastItem := (*h)[len(*h)-1]
	(*h) = (*h)[:len(*h)-1]

	if len(*h) > 0 {
		lastItem.index = 0
		(*h)[0] = lastItem
		h.down(0)
	}

	return removed
}

func (h minHeap) Update(item *visItem, area float64) {
	if item.area > area {
		// area got smaller
		item.area = area
		h.up(item.index)
	} else {
		// area got larger
		item.area = area
		h.down(item.index)
	}
}

func (h *minHeap) Remove(item *visItem) {
	i := item.index

	lastItem := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]

	if i != len(*h) {
		lastItem.index = i
		(*h)[i] = lastItem

		if lastItem.area < item.area {
			h.up(i)
		} else {
			h.down(i)
		}
	}
}

func (h minHeap) up(i int) {
	object := h[i]
	for i > 0 {
		up := ((i + 1) >> 1) - 1
		parent := h[up]

		if parent.area <= object.area {
			// parent is smaller so we're done fixing up the heap.
			break
		}

		// swap nodes
		parent.index = i
		h[i] = parent

		object.index = up
		h[up] = object

		i = up
	}
}

func (h minHeap) down(i int) {
	object := h[i]
	for {
		right := (i + 1) << 1
		left := right - 1

		down := i
		child := h[down]

		// swap with smallest child
		if left < len(h) && h[left].area < child.area {
			down = left
			child = h[down]
		}

		if right < len(h) && h[right].area < child.area {
			down = right
			child = h[down]
		}

		// non smaller, so quit
		if down == i {
			break
		}

		// swap the nodes
		child.index = i
		h[child.index] = child

		object.index = down
		h[down] = object

		i = down
	}
}

func doubleTriangleArea(a, b, c Point) float64 {
	return math.Abs((b.Y()-a.Y())*(c.X()-a.X()) - (b.X()-a.X())*(c.Y()-a.Y()))
}
