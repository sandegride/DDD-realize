package kernel

import (
	"delivery-service/internal/pkg/errs"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Location struct {
	x int
	y int

	valid bool
}

const (
	minX int = 1
	minY int = 1
	maxX int = 10
	maxY int = 10
)

func NewLocation(x int, y int) (Location, error) {
	if x < minX || x > maxX {
		return Location{}, errs.NewValueIsOutOfRangeError("x", x, minX, maxX)
	}
	if y < minY || y > maxY {
		return Location{}, errs.NewValueIsOutOfRangeError("y", y, minY, maxY)
	}

	return Location{
		x: x,
		y: y,

		valid: true,
	}, nil
}

func NewRandomLocation() Location {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	x := r.Intn(maxX-minX+1) + minX
	y := r.Intn(maxY-minY+1) + minY

	location, err := NewLocation(x, y)
	if err != nil {
		panic(fmt.Sprintf("invalid random location: x=%d, y=%d, err=%v", x, y, err))
	}
	return location
}

func (a Location) X() int {
	return a.x
}

func (a Location) Y() int {
	return a.y
}

func (a Location) IsValid() bool {
	return a.valid
}

func (a Location) Equal(other Location) bool {
	return a.x == other.x && a.y == other.y
}

func (l Location) DistanceTo(target Location) (int, error) {
	if !target.IsValid() {
		return 0, errs.NewValueIsRequiredError("location")
	}
	return int(math.Abs(float64(l.x-target.x)) + math.Abs(float64(l.y-target.y))), nil
}

func MinLocation() Location {
	location, err := NewLocation(minX, minY)
	if err != nil {
		panic("invalid min location configuration")
	}
	return location
}

func MaxLocation() Location {
	location, err := NewLocation(maxX, maxY)
	if err != nil {
		panic("invalid max location configuration")
	}
	return location
}
