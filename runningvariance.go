/*
Package runningvariance computes running mean, variance, standard deviation,
skewness and kurtosis using O(1) memory.

Based on code by John D Cook, see:

	- https://www.johndcook.com/blog/skewness_kurtosis/
	- http://www.johndcook.com/blog/standard_deviation/
	- Knuth TAOCP vol 2, 3rd edition, page 232

TODO: port correct skewness computation and possibly something else from
Julia OnlineStats library, https://github.com/joshday/OnlineStats.jl
*/
package runningvariance

import (
	"fmt"
	"math"
)

// Stat assumulates the data required for computing the statistics.
type Stat struct {
	N int64

	M1, M2, M3, M4 float64
}

// String implements Stringer.
func (s *Stat) String() string {
	return fmt.Sprintf("N=%d μ=%f σ=%f skew=%f ek=%f", s.N, s.Mean(), s.StdDev(), s.Skewness(), s.ExcessKurtosis())
}

// Push updates the statistics after adding a new value to the series.
func (s *Stat) Push(x float64) {
	n1 := float64(s.N)
	s.N++
	n := float64(s.N)

	delta := x - s.M1
	delta_n := delta / n
	delta_n2 := delta_n * delta_n
	term1 := delta * delta_n * n1
	s.M1 += delta_n
	s.M4 += term1*delta_n2*(n*n-3*n+3) + 6*delta_n2*s.M2 - 4*delta_n*s.M3
	s.M3 += term1*delta_n*(n-2) - 3*delta_n*s.M2
	s.M2 += term1
}

func (s *Stat) Mean() float64 {
	return s.M1
}

func (s *Stat) Variance() float64 {
	if s.N > 1 {
		return s.M2 / (float64(s.N) - 1.0)
	} else {
		return 0.0
	}
}

func (s *Stat) StdDev() float64 {
	return math.Sqrt(s.Variance())
}

/*
Skewness returns the skewness, a measure of the asymmetry of the probability
distribution.

For a simple distibution with a single peak, positive skewness means the peak is
closer to the left side, and negavive skewness means the peak is closer to the
right side. A zero value means that the tails on both sides of the mean
balance out overall.

WARNING: currently seems to be returning incorrect results, more work needed.
*/
func (s *Stat) Skewness() float64 {
	return math.Sqrt(float64(s.N)) * s.M3 / math.Pow(s.M2, 1.5)
}

/*
ExcessKurtosis returns the kurtosis of the data minus 3 (the “excess kurtosis”),
which gives an idea about how tail-heavy the distibution is.

Positive excess kurtotis means the distribution has a fatter tail than
the normal distribution. Similarly, negative excess kurtotis means a thinner
tail.
*/
func (s *Stat) ExcessKurtosis() float64 {
	return float64(s.N)*s.M4/(s.M2*s.M2) - 3.0
}

func Combined(a, b *Stat) Stat {
	var c Stat

	c.N = a.N + b.N
	an := float64(a.N)
	bn := float64(b.N)
	cn := float64(c.N)

	delta := b.M1 - a.M1
	delta2 := delta * delta
	delta3 := delta * delta2
	delta4 := delta2 * delta2

	c.M1 = (an*a.M1 + bn*b.M1) / cn
	c.M2 = a.M2 + b.M2 + delta2*an*bn/cn

	c.M3 = a.M3 + b.M3 + delta3*an*bn*(an-bn)/(cn*cn)
	c.M3 += 3.0 * delta * (an*b.M2 - bn*a.M2) / cn

	c.M4 = a.M4 + b.M4 + delta4*an*bn*(an*an-an*bn+bn*bn)/
		(cn*cn*cn)
	c.M4 += 6.0*delta2*(an*an*b.M2+bn*bn*a.M2)/(cn*cn) +
		4.0*delta*(an*b.M3-bn*a.M3)/cn

	return c
}

func (s *Stat) Combine(b *Stat) {
	*s = Combined(s, b)
}
