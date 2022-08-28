package testutils

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
)

// AndersonDarlingTest takes a sample of floating-point numbers and performs
// the Anderson-Darling test for numerical distribution. It returns a p-value
// of the resulting test statistic, with a high p-value (>= 0.95, conventionally)
// indicating that the null hypothesis of a normal distribution can be rejected.
//
// Returns a negative p-value and an error if the sample is empty.
//
// This method of computing the Anderson-Darling test uses the approximation from:
// [1] Marsaglia, John & Marsaglia, George. (2004). Evaluating the Anderson-Darling Distribution. Journal of Statistical Software. 09. 10.18637/jss.v009.i02.
func AndersonDarlingTest(samples []float64) (float64, error) {
	aStat, err := ADStatistic(samples)
	if err != nil {
		return -1.0, err
	}

	N := float64(len(samples))
	return ADPValue(aStat * (1.0 + 4.0/N - 25/(N*N)))
}

// ADStatistic computes the a-statistic of the given sample. Currently only
// supports normal distributions. Returns an error on an empty sample.
func ADStatistic(rawSamples []float64) (float64, error) {
	if len(rawSamples) == 0 {
		return -1, fmt.Errorf("Cannot compute a-statistic on an empty sample")
	}

	mu, sigma := stat.PopMeanStdDev(rawSamples, nil)

	// Normal.CDF() sometimes computes a value of infinity based on a zero variance,
	// and sometimes it computes a NaN. I simplify things with an explicit return.
	if sigma == 0.0 {
		return math.Inf(1), nil
	}

	norm := distuv.Normal{Mu: mu, Sigma: sigma}

	// Computing the a-statistic requires the input to be sorted in ascending order
	samples := make([]float64, len(rawSamples))
	copy(samples, rawSamples)
	sort.Float64s(samples)

	N := float64(len(samples))

	Asquared := -N

	for i := 0; i < int(N); i++ {
		coeff := (2.0*float64(i+1) - 1) / N
		j := int(N) - 1 - i
		Asquared -= coeff * (math.Log(norm.CDF(samples[i])) + math.Log(norm.Survival(samples[j])))
	}

	return math.Sqrt(Asquared), nil
}

// ADPValue computes and returns the p-value of the given a-statistic.
// Returns -1 and error if the z-statistic is 0 or less.
//
// This computation uses the adinf(z) approximation in Marsaglia & Marsaglia [1].
// The approximation is split into two regimes based on whether or not the test
// statistic z >= 2.
func ADPValue(z float64) (float64, error) {
	if z <= 0.0 {
		return -1.0, fmt.Errorf("AD-statistic must be greater than 0")
	}

	if z < 2.0 {
		return adsProbabilityLow(z), nil
	}

	return adsProbabilityHigh(z), nil
}

// adsProbabilityLow Covers the 0 < z < 2 regime of the Marsaglia & Marsaglia
// approximation. It takes an a-statistic value as an input and returns a
// p-value of range [0,1] as an output. The approximation error err < 2e-6 [1].
func adsProbabilityLow(z float64) float64 {
	adSeed := 0.00168691
	adIntercept := 2.00012
	adPolynomial := []float64{0.0116720, 0.0347962, 0.0649821, 0.247105}
	adExponent := -1.2337141

	polynomial := adSeed

	for i := 0; i < len(adPolynomial); i++ {
		polynomial = adPolynomial[i] - polynomial*z
	}

	polynomial = adIntercept + polynomial*z

	return math.Pow(z, -0.5) * math.Exp(adExponent/z) * polynomial
}

// adsProbabilityHigh overs the 2 <= z <= inf regime of the Marsaglia & Marsaglia
// approximation. Takes an a-statistic value as an input and returns a p-value of
// range [0,1] as an output. The approximation error err < 8e-7 [1].
func adsProbabilityHigh(z float64) float64 {
	adSeed := 0.0003146
	adPolynomial := []float64{0.008056, 0.082433, 0.43424, 2.30695}
	adIntercept := 1.0776

	polynomial := adSeed
	for i := 0; i < len(adPolynomial); i++ {
		polynomial = adPolynomial[i] - polynomial*z
	}

	polynomial = adIntercept - polynomial*z

	return math.Exp(-math.Exp(polynomial))
}
