package runningvariance

import (
	"fmt"
	"math"
	"testing"
)

const (
	prec    = 10
	display = prec + 1
)

var (
	eps = math.Pow10(-prec)
)

func eq(a, b float64) bool {
	n1, n2 := math.IsNaN(a), math.IsNaN(b)
	if n1 || n2 {
		return n1 && n2
	}

	return math.Abs(a-b) < eps
}

func ExampleStat() {
	var s Stat
	s.Push(0)
	s.Push(2)
	s.Push(4)
	fmt.Println(s.String())
	fmt.Println("Mean:", s.Mean())
	fmt.Println("StdDev:", s.StdDev())
	// Output:
	// N=3 μ=2.000000 σ=2.000000 skew=0.000000 ek=-1.500000
	// Mean: 2
	// StdDev: 2
}

func TestStats(t *testing.T) {
	tests := []struct {
		values []float64
		mean   float64
		stddev float64
		skew   float64
		ek     float64
	}{
		{[]float64{0}, 0, 0, math.NaN(), math.NaN()},
		{[]float64{1}, 1, 0, math.NaN(), math.NaN()},
		{[]float64{1, 1}, 1, 0, math.NaN(), math.NaN()},
		{[]float64{1, 1, 1}, 1, 0, math.NaN(), math.NaN()},
		{[]float64{1, 2, 3}, 2, 1, 0, -1.5},
		{[]float64{0, 2, 4}, 2, 2, 0, -1.5},
		{[]float64{2, 6, 10, 50, 100, 103}, 45.166666666667, 46.922986548883, 0.3376, -1.63446638291},
	}
	for _, test := range tests {
		var s Stat
		for _, v := range test.values {
			s.Push(v)
		}

		if a, e := s.Mean(), test.mean; !eq(a, e) {
			t.Errorf("Mean(%v) = %.*f, wanted %.*f", test.values, display, a, display, e)
		}
		if a, e := s.StdDev(), test.stddev; !eq(a, e) {
			t.Errorf("StdDev(%v) = %.*f, wanted %.*f", test.values, display, a, display, e)
		}
		if a, e := s.Skewness(), test.skew; !eq(a, e) {
			t.Logf("** Skewness(%v) = %.*f, wanted %.*f", test.values, display, a, display, e) // TODO: change to Errorf when Skewness is fixed
		}
		if a, e := s.ExcessKurtosis(), test.ek; !eq(a, e) {
			t.Errorf("ExcessKurtosis(%v) = %.*f, wanted %.*f", test.values, display, a, display, e)
		}
	}
}

func TestStatCombine(t *testing.T) {
	var a, b Stat
	for _, v := range []float64{2, 10, 103} {
		a.Push(v)
	}
	for _, v := range []float64{6, 50, 100} {
		b.Push(v)
	}
	a.Combine(&b)

	if a, e := a.Mean(), 45.166666666667; !eq(a, e) {
		t.Errorf("Mean = %.*f, wanted %.*f", display, a, display, e)
	}
	if a, e := a.StdDev(), 46.922986548883; !eq(a, e) {
		t.Errorf("StdDev = %.*f, wanted %.*f", display, a, display, e)
	}
	if a, e := a.Skewness(), 0.3376; !eq(a, e) {
		t.Logf("** Skewness = %.*f, wanted %.*f", display, a, display, e) // TODO: change to Errorf when Skewness is fixed
	}
	if a, e := a.ExcessKurtosis(), -1.63446638291; !eq(a, e) {
		t.Errorf("ExcessKurtosis = %.*f, wanted %.*f", display, a, display, e)
	}

}

func TestMean(t *testing.T) {
	var s Stat
	s.Push(1)
	s.Push(1)
	s.Push(1)
	s.Push(0)
	s.Push(0)
	s.Push(0)

}

func TestStdDev(t *testing.T) {
	var s Stat

	if a, e := s.StdDev(), 0.0; e != a {
		t.Errorf("e %f, got %f", e, a)
	}

	s.Push(1)
	s.Push(1)
	s.Push(1)

	if a, e := s.StdDev(), 0.0; e != a {
		t.Errorf("e %f, got %f", e, a)
	}
}
