package burrow

// A SampleDistribution is a function that returns a float64, representing a number drawn randomly according to a probability distribution.
type SampleDistribution func() float64

// Sample generates `n` samples drawn from the probability distribution `distro`, which it feeds into the channel returned to the caller. In effect, this is a Python-style generator for distribution samples.
func (distro SampleDistribution) Sample(n uint) chan float64 {
	output := make(chan float64)

	go func() {
		defer close(output)

		for i := uint(0); i < n; i++ {
			output <- distro()
		}
	}()

	return output
}
