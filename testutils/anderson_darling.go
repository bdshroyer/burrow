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
	pValue, err := ADPValue(aStat * (1.0 + 4.0/N - 25/(N*N)))
	if err != nil {
		return -1.0, err
	}

	fix, err := ADErrFix(float64(len(samples)), pValue)
	if err != nil {
		return -1.0, err
	}

	return pValue + fix, nil
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

// ADErrFix provides a correction term that de-resolves an infinite AD distribution
// (or at least the approximation we use) down to a resolution of N. Takes a sample
// size N and an AD p-value p as inputs.
//
// Returns 0 if N <= 0, or if p is not a value between 0 and 1.
func ADErrFix(N, p float64) (float64, error) {
	if p < 0.0 || p > 1.0 {
		return 0.0, fmt.Errorf("Probability must be between 0 and 1.")
	}

	if N <= 0 {
		return 0.0, fmt.Errorf("N must be greater than 0.")
	}

	N2 := N * N
	N3 := N2 * N
	cN := leftCrossing(N)

	if p < cN {
		return (0.0037/N3 + 0.00078/N2 + 0.00006/N) * g1(p/cN), nil
	}

	if p < 0.8 {
		return (0.04213/N + 0.01365/N2) * g2(p-cN) / (0.8 - cN), nil
	}

	return g3(p) / N, nil
}

/****** ErrFix support functions ******/

// Sourced from Section 3 of Marsaglia & Marsaglia [1]. This gives the first
// (non-zero) intercept of the error function.
func leftCrossing(N float64) float64 {
	return 0.01265 + 0.1757/N
}

func g1(p float64) float64 {
	return math.Sqrt(p) * (1 - p) * (49*p - 102)
}

func g2(p float64) float64 {
	gSeed := 1.91864
	gPolynomial := []float64{8.259, 14.458, 14.6538, 6.54034}
	gIntercept := -0.00022633

	polynomial := gSeed
	for i := 0; i < len(gPolynomial); i++ {
		polynomial = gPolynomial[i] - polynomial*p
	}

	return gIntercept + polynomial*p
}

func g3(p float64) float64 {
	gSeed := 255.7844
	gPolynomial := []float64{1116.360, 1950.646, 1705.091, 745.2337}
	gIntercept := -130.2137

	polynomial := gSeed
	for i := 0; i < len(gPolynomial); i++ {
		polynomial = gPolynomial[i] - polynomial*p
	}

	return gIntercept + polynomial*p
}
