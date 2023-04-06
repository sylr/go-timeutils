package timeutils

import (
	"math/rand"
	"testing"
	"time"
)

func TestNewInterval(t *testing.T) {
	start := time.Now()
	end := start.Add(time.Hour)

	defer func() {
		if err := recover(); err != nil {
			t.Errorf("NewInterval(%s, %s) should not have panicked", start, end)
		}
	}()

	_ = NewInterval(start, end)
}

func TestNewIntervalPanic(t *testing.T) {
	start := time.Now()
	end := start.Add(-time.Hour)

	defer func() {
		if err := recover(); err == nil {
			t.Errorf("NewInterval(%s, %s) should have panicked", start, end)
		}
	}()

	_ = NewInterval(start, end)
}

func TestIntervalInclude(t *testing.T) {
	start := time.Now()
	end := start.Add(time.Hour)

	tests := []struct {
		name     string
		input    time.Time
		expected bool
	}{
		{
			name:     `Input equals start`,
			input:    start,
			expected: true,
		},
		{
			name:     `Input equals end`,
			input:    end,
			expected: false,
		},
		{
			name:     `Input before start`,
			input:    start.Add(-time.Hour),
			expected: false,
		},
		{
			name:     `Input after end`,
			input:    end.Add(time.Hour),
			expected: false,
		},
	}

	i := NewInterval(start, end)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if i.Include(test.input) != test.expected {
				t.Errorf("%s.Include(%s) was expected to be %t", i, test.input.Format(time.RFC3339), test.expected)
			}
		})
	}
}

func TestIntervalEqual(t *testing.T) {
	start := time.Now()
	end := start.Add(time.Hour)

	tests := []struct {
		name     string
		input    Interval
		expected bool
	}{
		{
			name:     `Input equals interval`,
			input:    NewInterval(start, end),
			expected: true,
		},
		{
			name:     `Input's start negatively differs`,
			input:    NewInterval(start.Add(-time.Second), end),
			expected: false,
		},
		{
			name:     `Input's start positevely differs`,
			input:    NewInterval(start.Add(time.Second), end),
			expected: false,
		},
		{
			name:     `Input's end negatively differs`,
			input:    NewInterval(start, end.Add(-time.Second)),
			expected: false,
		},
		{
			name:     `Input's end positevely differs`,
			input:    NewInterval(start, end.Add(time.Second)),
			expected: false,
		},
		{
			name:     `Input is bigger`,
			input:    NewInterval(start.Add(-time.Second), end.Add(time.Second)),
			expected: false,
		},
		{
			name:     `Input is smaller`,
			input:    NewInterval(start.Add(time.Second), end.Add(-time.Second)),
			expected: false,
		},
	}

	i := NewInterval(start, end)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if i.Equal(test.input) != test.expected {
				t.Errorf("%s.Equal(%s) was expected to be %t", i, test.input, test.expected)
			}
		})
	}
}

func TestIntervalEngulf(t *testing.T) {
	start := time.Now()
	end := start.Add(time.Hour)

	tests := []struct {
		name     string
		input    Interval
		expected bool
	}{
		{
			name:     `Input equals interval`,
			input:    NewInterval(start, end),
			expected: true,
		},
		{
			name:     `Input's start negatively differs`,
			input:    NewInterval(start.Add(-time.Second), end),
			expected: false,
		},
		{
			name:     `Input's start positevely differs`,
			input:    NewInterval(start.Add(time.Second), end),
			expected: true,
		},
		{
			name:     `Input's end negatively differs`,
			input:    NewInterval(start, end.Add(-time.Second)),
			expected: true,
		},
		{
			name:     `Input's end positevely differs`,
			input:    NewInterval(start, end.Add(time.Second)),
			expected: false,
		},
		{
			name:     `Input is bigger`,
			input:    NewInterval(start.Add(-time.Second), end.Add(time.Second)),
			expected: false,
		},
		{
			name:     `Input is smaller`,
			input:    NewInterval(start.Add(time.Second), end.Add(-time.Second)),
			expected: true,
		},
	}

	i := NewInterval(start, end)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if i.Engulf(test.input) != test.expected {
				t.Errorf("%s.Engulf(%s) was expected to be %t", i, test.input, test.expected)
			}
		})
	}
}

func TestIntervalOverlap(t *testing.T) {
	start := time.Now().Truncate(time.Second).Truncate(time.Minute).Truncate(time.Hour)
	end := start.Add(time.Hour)

	tests := []struct {
		name     string
		input    Interval
		expected bool
	}{
		{
			name:     `Input equals interval`,
			input:    NewInterval(start, end),
			expected: true,
		},
		{
			name:     `Input's start negatively differs`,
			input:    NewInterval(start.Add(-time.Second), end),
			expected: true,
		},
		{
			name:     `Input's start positevely differs`,
			input:    NewInterval(start.Add(time.Second), end),
			expected: true,
		},
		{
			name:     `Input's end negatively differs`,
			input:    NewInterval(start, end.Add(-time.Second)),
			expected: true,
		},
		{
			name:     `Input's end positevely differs`,
			input:    NewInterval(start, end.Add(time.Second)),
			expected: true,
		},
		{
			name:     `Input is bigger`,
			input:    NewInterval(start.Add(-time.Second), end.Add(time.Second)),
			expected: true,
		},
		{
			name:     `Input is smaller`,
			input:    NewInterval(start.Add(time.Second), end.Add(-time.Second)),
			expected: true,
		},
		{
			name:     `Input is before`,
			input:    NewInterval(start.Add(-time.Hour), start.Add(-time.Second)),
			expected: false,
		},
		{
			name:     `Input is after`,
			input:    NewInterval(end.Add(time.Second), end.Add(time.Hour)),
			expected: false,
		},
	}

	i := NewInterval(start, end)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if i.Overlap(test.input) != test.expected {
				t.Errorf("%s.Overlap(%s) was expected to be %t", i, test.input, test.expected)
			}
		})
	}
}

func TestIntervalSub(t *testing.T) {
	start := time.Now()
	end := start.Add(time.Hour)

	i := NewInterval(start, end)

	tests := []struct {
		name     string
		input    Interval
		expected Intervals
	}{
		{
			name:     `Input equals interval`,
			input:    NewInterval(start, end),
			expected: Intervals{},
		},
		{
			name:     `Input before interval`,
			input:    NewInterval(start.Add(-time.Hour), start.Add(-time.Second)),
			expected: Intervals{i},
		},
		{
			name:     `Input after interval`,
			input:    NewInterval(end.Add(time.Second), end.Add(time.Hour)),
			expected: Intervals{i},
		},
		{
			name:     `Input's start negatively differs`,
			input:    NewInterval(start.Add(-time.Second), end),
			expected: Intervals{},
		},
		{
			name:     `Input's start positively differs`,
			input:    NewInterval(start.Add(time.Second), end),
			expected: Intervals{NewInterval(start, start.Add(time.Second))},
		},
		{
			name:     `Input's end negatively differs`,
			input:    NewInterval(start, end.Add(-time.Second)),
			expected: Intervals{NewInterval(end.Add(-time.Second), end)},
		},
		{
			name:     `Input's start positively differs`,
			input:    NewInterval(start, end.Add(time.Second)),
			expected: Intervals{},
		},
		{
			name:  `Input is smaller`,
			input: NewInterval(start.Add(time.Second), end.Add(-time.Second)),
			expected: Intervals{
				NewInterval(start, start.Add(time.Second)),
				NewInterval(end.Add(-time.Second), end),
			},
		},
		{
			name:     `Input is bigger`,
			input:    NewInterval(start.Add(-time.Second), end.Add(time.Second)),
			expected: Intervals{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := i.Sub(test.input)
			if !actual.Equal(test.expected) {
				t.Errorf("%s.Sub(%s) was expected to be %s, got %s", i, test.input, test.expected, actual)
			}
		})
	}
}

func TestIntervalsEqual(t *testing.T) {
	start := time.Now()
	end := start.Add(time.Hour)

	is := Intervals{
		NewInterval(start, end),
		NewInterval(start, end.Add(-time.Second)),
		NewInterval(start, end.Add(+time.Second)),
		NewInterval(start.Add(-time.Second), end),
		NewInterval(start.Add(+time.Second), end),
		NewInterval(start.Add(-time.Second), end.Add(-time.Second)),
		NewInterval(start.Add(+time.Second), end.Add(+time.Second)),
		NewInterval(start.Add(-time.Second), end.Add(+time.Second)),
		NewInterval(start.Add(+time.Second), end.Add(-time.Second)),
	}

	tests := []struct {
		name     string
		input    Intervals
		expected bool
	}{
		{
			name: `Input equals interval`,
			input: Intervals{
				NewInterval(start, end),
				NewInterval(start, end.Add(-time.Second)),
				NewInterval(start, end.Add(+time.Second)),
				NewInterval(start.Add(-time.Second), end),
				NewInterval(start.Add(+time.Second), end),
				NewInterval(start.Add(-time.Second), end.Add(-time.Second)),
				NewInterval(start.Add(+time.Second), end.Add(+time.Second)),
				NewInterval(start.Add(-time.Second), end.Add(+time.Second)),
				NewInterval(start.Add(+time.Second), end.Add(-time.Second)),
			},
			expected: true,
		},
		{
			name: `Randomly shuffled`,
			input: Intervals{
				NewInterval(start.Add(+time.Second), end.Add(+time.Second)),
				NewInterval(start, end.Add(+time.Second)),
				NewInterval(start, end),
				NewInterval(start.Add(-time.Second), end.Add(+time.Second)),
				NewInterval(start.Add(-time.Second), end),
				NewInterval(start, end.Add(-time.Second)),
				NewInterval(start.Add(+time.Second), end),
				NewInterval(start.Add(-time.Second), end.Add(-time.Second)),
				NewInterval(start.Add(+time.Second), end.Add(-time.Second)),
			},
			expected: true,
		},
		{
			name: `Randomly shuffled`,
			input: Intervals{
				NewInterval(start, end),
				NewInterval(start, end.Add(-time.Minute)),
				NewInterval(start, end.Add(+time.Minute)),
				NewInterval(start.Add(-time.Minute), end),
				NewInterval(start.Add(+time.Minute), end),
				NewInterval(start.Add(-time.Minute), end.Add(-time.Minute)),
				NewInterval(start.Add(+time.Minute), end.Add(+time.Minute)),
				NewInterval(start.Add(-time.Minute), end.Add(+time.Minute)),
				NewInterval(start.Add(+time.Minute), end.Add(-time.Minute)),
			},
			expected: false,
		},
		{
			name: `Not same size`,
			input: Intervals{
				NewInterval(start, end),
				NewInterval(start, end.Add(-time.Minute)),
				NewInterval(start, end.Add(+time.Minute)),
			},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for i := 0; i < 10; i++ {
				rand.Shuffle(len(test.input), is.Swap)
				if is.Equal(test.input) != test.expected {
					t.Errorf("%s.Equal(%s) was expected to be %t", is, test.input, test.expected)
				}
			}
		})
	}
}
