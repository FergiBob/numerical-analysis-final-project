package internal

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/bradhe/stopwatch"
)

type UnaryFunction func(float64) float64

func randomfloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// UnaryMonteCarloIntegration approximates the integral of a unaryfunction on the range of a to b by using a number n of random points
// Input: UnaryFunction, min, max, iterations, mode of calculation
// Output: Intergral, Stanadard Error, Seconds to calculate
func UnaryMonteCarloIntegration(f UnaryFunction, a, b float64, n int, mode string) (float64, float64, float64) {
	var sum, sum_squared float64
	watch := stopwatch.Start()
	if mode == "multi-threaded" {
		sum, sum_squared = ParallelMonteCarlo(f, a, b, n)
	} else {
		sum, sum_squared = MonteCarlo(f, a, b, n)
	}

	// calculate integral estimate
	mean := sum / float64(n)
	integral := (b - a) * mean

	// calculate variance and standard deviation required for precision calculation
	// variance formula: E[X^2] - (E[X])^2
	variance := (sum_squared / float64(n)) - (math.Pow(mean, 2))
	standard_deviation := math.Sqrt(variance)

	standard_error := (b - a) * (standard_deviation / math.Sqrt(float64(n)))

	watch.Stop()

	return integral, standard_error, float64(watch.Milliseconds())
}

// MonteCarlo approximates the integral of a unaryfunction on the range of a to b by using a number n of random points
// Input: UnaryFunction, min, max, iterations
// Output: Sum of f(X), Sum of f(X)^2
func MonteCarlo(f UnaryFunction, a, b float64, n int) (float64, float64) {
	var sum float64 = 0.0
	var sum_squared float64 = 0.0

	// sum from 1 to n of f(x_i)
	for range n - 1 {
		var x float64 = randomfloat(a, b)
		y := f(x)
		sum += y
		sum_squared += math.Pow(y, 2)
	}
	return sum, sum_squared
}

// ParallelMonteCarlo approximates the integral of a unaryfunction on the range of a to b by using a number n of random points
// The process is sped up in this instance through the use of go routines
// Input: UnaryFunction, min, max, iterations
// Output: Sum of f(X), Sum of f(X)^2
func ParallelMonteCarlo(f UnaryFunction, a, b float64, n int) (float64, float64) {
	numCores := runtime.NumCPU()
	samplesPerCore := n / numCores

	sumChannel := make(chan float64, numCores)
	sqChannel := make(chan float64, numCores)
	var wg sync.WaitGroup

	// divide the work into a go routine per available core
	for i := range numCores {
		wg.Add(1)
		// Pass 'i' into the goroutine to ensure a unique seed
		go func(id int) {
			defer wg.Done()
			var localSum float64
			var localSumSq float64

			// Unique seed per goroutine
			seed := time.Now().UnixNano() + int64(id)
			source := rand.New(rand.NewSource(seed))

			for range samplesPerCore {
				x := a + source.Float64()*(b-a)
				y := f(x)
				localSum += y
				localSumSq += y * y
			}
			sumChannel <- localSum
			sqChannel <- localSumSq
		}(i)
	}

	// Closer goroutine
	go func() {
		wg.Wait()
		close(sumChannel)
		close(sqChannel)
	}()

	totalSum := 0.0
	totalSumSq := 0.0
	for s := range sumChannel {
		totalSum += s
	}
	for sSq := range sqChannel {
		totalSumSq += sSq
	}

	return totalSum, totalSumSq
}

// MISER Monte Carlo based on recursive stratified sampling. This technique aims to reduce the
// overall integration error by concentrating integration points in the regions of highest variance.
// Input: UnaryFunction, min, max, iterations, depth
// Output: Intergral, Stanadard Error
func UnaryMISERMonteCarlo(f UnaryFunction, a, b float64, n, d int) (float64, float64) {
	// Check depth and samples
	if d <= 0 || n < 100 {
		integral, standard_error, _ := UnaryMonteCarloIntegration(f, a, b, n, "multi-threaded") // use multi-threaded for speed
		return integral, standard_error
	}

	// Estimate variance at different bisection points to find optimal split
	var bestMid float64
	var bestVar float64 = math.Inf(1)
	var bestN1, bestN2 int
	var bestIntegral1, bestIntegral2 float64
	var bestTestN int

	// Try several candidate midpoints
	for candidate := range 5 {
		mid := a + (b-a)*(float64(candidate)+0.5)/5.0

		// Use 10% of samples for testing (min 20)
		testN := max(n/10, 20)
		integral1, se1 := UnaryMISERMonteCarlo(f, a, mid, testN, 0)
		integral2, se2 := UnaryMISERMonteCarlo(f, mid, b, testN, 0)

		// Calculate combined variance
		totalVar := se1*se1 + se2*se2
		if totalVar < bestVar {
			bestVar = totalVar
			bestMid = mid
			bestN1 = int(float64(n) * se1 / (se1 + se2))
			bestN2 = n - bestN1

			// Ensure minimum samples
			bestN1 = max(bestN1, 50)
			bestN2 = max(bestN2, 50)
			if bestN1+bestN2 > n {
				bestN2 = n - bestN1
			}

			// Store test results to incorporate later
			bestIntegral1 = integral1
			bestIntegral2 = integral2
			bestTestN = testN
		}
	}

	// Recursive integration with optimal split
	integral1, se1 := UnaryMISERMonteCarlo(f, a, bestMid, bestN1, d-1)
	integral2, se2 := UnaryMISERMonteCarlo(f, bestMid, b, bestN2, d-1)

	// Combine results, incorporating test samples
	// Weighted average of test and final results
	testWeight := float64(bestTestN) / float64(n)
	finalIntegral := (1-testWeight)*(integral1+integral2) + testWeight*(bestIntegral1+bestIntegral2)

	// Standard error calculation
	standard_error := math.Sqrt(se1*se1 + se2*se2)

	return finalIntegral, standard_error
}

// helper function used to print out each series of monte-carlo tests
func PrintTest(f UnaryFunction, a, b float64, n int) {

	fmt.Printf("\n	Single-Threaded Monte Carlo Integration")
	integral, errorval, time := UnaryMonteCarloIntegration(f, a, b, n, "single-threaded")
	fmt.Printf("\n		Integral %f", integral)
	fmt.Printf("\n		Error value %f", errorval)
	fmt.Printf("\n		Time Elapsed %f milliseconds", time)

	fmt.Printf("\n\n	Multi-Threaded Monte Carlo Integration")
	integral, errorval, time = UnaryMonteCarloIntegration(f, a, b, n, "multi-threaded")
	fmt.Printf("\n		Integral %f", integral)
	fmt.Printf("\n		Error value %f", errorval)
	fmt.Printf("\n		Time Elapsed %f milliseconds", time)

	fmt.Printf("\n\n	MISERMonteCarlo")
	watch := stopwatch.Start()
	integral_MISER, errorval_MISER := UnaryMISERMonteCarlo(f, a, b, n, 3)
	watch.Stop()
	fmt.Printf("\n		Integral %f", integral_MISER)
	fmt.Printf("\n		Error value %f", errorval_MISER)
	fmt.Printf("\n		Time Elapsed %f milliseconds", float64(watch.Milliseconds()))
}
