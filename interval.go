package timeutils

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// Interval defines a time period that is constained by two time boundaries, a
// start time that is part of the interval and an end time which is excluded
// from the interval.
//
//	  |----------i----------[
//	start                  end
//
// i.Include(start) == true
// i.Include(end) == false
type Interval struct {
	Start time.Time
	End   time.Time
}

func NewInterval(start, end time.Time) Interval {
	if start.After(end) {
		panic("interval's start can't be after its end.")
	}

	return Interval{Start: start, End: end}
}

const (
	stringFormat = "Interval{start: %s, end: %s, duration: %s}"
)

func (i Interval) String() string {
	return fmt.Sprintf(stringFormat, i.Start.Format(time.RFC3339), i.End.Format(time.RFC3339), i.Duration())
}

func (i Interval) Duration() time.Duration {
	return i.End.Sub(i.Start)
}

// Include tests if input time is within the interval. Note that if input is
// equal to the end of the interval, then false is returned.
//
// interval:      |------------i------------[
// input:              |
func (i Interval) Include(input time.Time) bool {
	return (i.Start.Before(input) || i.Start.Equal(input)) && i.End.After(input)
}

// Equal tests that the input interval has the time time boundaries as Interval.
//
// interval:      |------------i------------[
// input:         |----------input----------[
func (i Interval) Equal(input Interval) bool {
	return i.Start.Equal(input.Start) && i.End.Equal(input.End)
}

// Engulf tests that the input interval is within Interval. Returns true also if
// both intervals are equal.
//
// interval:      |------------i------------[
// input:              |---input---[
func (i Interval) Engulf(input Interval) bool {
	return (i.Start.Before(input.Start) || i.Start.Equal(input.Start)) &&
		(i.End.After(input.End) || i.End.Equal(input.End))
}

// Overlap tests if input overlaps with Interval. Sharing opposite time
// boundaries is not enough to overlap.
//
// interval:      |------------i------------[
// input:                          |----input---[
// input:    |----input----[
// input:             |---input---[
// input:      |-------------input-------------[
func (i Interval) Overlap(input Interval) bool {
	return !((input.End.Before(i.Start) || input.End.Equal(i.Start)) ||
		(input.Start.After(i.End) || input.Start.Equal(i.End)))
}

// Contiguous tests if input is contiguous to Interval.
//
// interval:           |----------i----------[
// input:                                    |--input--[
// input:    |--input--[
func (i Interval) Contiguous(input Interval) bool {
	return i.End.Equal(input.Start) || i.Start.Equal(input.End)
}

// Sub substracts input to Interval.
//
// interval:      |------------i------------[
// input:                    |------input------[
// output:        |----i'----[
// input:                 |--input--[
// output:        |--i'---[         |--i"---[
func (i Interval) Sub(input Interval) Intervals {
	if input.Equal(i) || input.Engulf(i) {
		return Intervals{}
	}

	if (input.End.Before(i.Start) || input.End.Equal(i.Start)) ||
		(input.Start.After(i.End) || input.Start.Equal(i.End)) {
		return Intervals{
			{Start: i.Start, End: i.End},
		}
	}

	if i.Engulf(input) {
		if i.Start.Equal(input.Start) {
			return Intervals{
				{Start: input.End, End: i.End},
			}
		} else if i.End.Equal(input.End) {
			return Intervals{
				{Start: i.Start, End: input.Start},
			}
		} else {
			return Intervals{
				{Start: i.Start, End: input.Start},
				{Start: input.End, End: i.End},
			}
		}
	}

	if input.Start.Before(i.Start) {
		return Intervals{
			{Start: input.End, End: i.End},
		}
	} else {
		return Intervals{
			{Start: i.Start, End: input.Start},
		}
	}
}

type Intervals []Interval

func (is Intervals) String() string {
	strs := make([]string, 0, len(is))
	for _, s := range is {
		strs = append(strs, s.String())
	}

	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
}

func (is Intervals) Equal(input Intervals) bool {
	if len(is) != len(input) {
		return false
	}

	less := func(i, j Interval) int {
		if i.Start.Before(j.Start) {
			return -1
		} else if i.Start.After(j.Start) {
			return 1
		} else {
			if i.End.Before(j.End) {
				return -1
			} else if i.End.After(j.End) {
				return 1
			} else {
				return 0
			}
		}
	}

	slices.SortFunc(is, less)
	slices.SortFunc(input, less)

	for i := range input {
		if !is[i].Equal(input[i]) {
			return false
		}
	}

	return true
}

func (is Intervals) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}
