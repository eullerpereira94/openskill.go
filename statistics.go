package openskill

import (
	"math"

	"github.com/samber/lo"
)

const mean float64 = 0.0
const variance float64 = 1.0

// phiMinor
func pdf(x float64) float64 {
	standardDeviation := math.Sqrt(variance)
	m := standardDeviation * math.Sqrt(2*math.Pi)
	e := math.Exp(-math.Pow(x-mean, 2) / (2 * variance))
	return e / m
}

// phiMajor
func cdf(x float64) float64 {
	standardDeviation := math.Sqrt(variance)
	return 0.5 * math.Erfc(-(x-mean)/(standardDeviation*math.Sqrt(2)))
}

// phiMajorInverse
func ppf(x float64) float64 {
	standardDeviation := math.Sqrt(variance)
	return mean - standardDeviation*math.Sqrt(2)*math.Erfcinv(2*x)
}

func v(x, t float64) float64 {
	xt := x - t
	denom := cdf(xt)
	epsilon := math.Nextafter(1.0, 2.0) - 1.0

	if denom < epsilon {
		return -xt
	}
	return pdf(xt) / denom
}

func vt(x, t float64) float64 {
	xx := math.Abs(x)
	b := cdf((t - xx)) - cdf((-t - xx))

	if b < 1e-5 {
		if x < 0 {
			return -x - t
		}
		return -x + t
	}

	a := pdf((-t - xx)) - pdf((t - xx))
	return lo.Ternary(x < 0, a, -a) / b
}

func w(x, t float64) float64 {
	xt := x - t
	denom := cdf(xt)
	epsilon := math.Nextafter(1.0, 2.0) - 1.0

	if denom < epsilon {
		if x < 0 {
			return 1
		}
		return 0
	}
	return v(x, t) * (v(x, t) + xt)
}

func wt(x, t float64) float64 {
	xx := math.Abs(x)
	b := cdf((t - xx)) - cdf((-t - xx))
	epsilon := math.Nextafter(1.0, 2.0) - 1.0

	if b < epsilon {
		return 1.0
	}

	return ((t-xx)*pdf(t-xx)+(t+xx)*pdf(-t-xx))/b + vt(x, t)*vt(x, t)
}
