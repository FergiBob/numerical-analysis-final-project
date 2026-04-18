package internal

import (
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/bradhe/stopwatch"
)

// MultiUnaryFunction takes a slice of coordinates representing a point in N-dimensional space
type MultiFunction func([]float64) float64

type Bounds struct {
	Min []float64 // e.g., [x_min, y_min, z_min]
	Max []float64 // e.g., [x_max, y_max, z_max]
}

// Helper to calculate the n-dimensional volume
func calculateVolume(b Bounds) float64 {
	vol := 1.0
	for i := 0; i < len(b.Min); i++ {
		vol *= (b.Max[i] - b.Min[i])
	}
	return vol
}

func ParallelMultiMonteCarlo(f MultiFunction, b Bounds, n int) (float64, float64) {
	numCores := runtime.NumCPU()
	samplesPerCore := n / numCores
	dims := len(b.Min)

	sumChannel := make(chan float64, numCores)
	sqChannel := make(chan float64, numCores)
	var wg sync.WaitGroup

	for i := range numCores {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			var localSum, localSumSq float64
			source := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))

			// Reusable slice for coordinates to save on memory allocation
			point := make([]float64, dims)

			for range samplesPerCore {
				// Generate a random point in N-dimensions
				for d := range dims {
					point[d] = b.Min[d] + source.Float64()*(b.Max[d]-b.Min[d])
				}

				y := f(point)
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

func MultiMonteCarloIntegration(f MultiFunction, b Bounds, n int) (float64, float64, float64) {
	var sum, sum_squared float64
	watch := stopwatch.Start()

	// Perform the sampling
	sum, sum_squared = ParallelMultiMonteCarlo(f, b, n)

	// Calculate Hyper-volume
	volume := 1.0
	for i := 0; i < len(b.Min); i++ {
		volume *= (b.Max[i] - b.Min[i])
	}

	// Calculate integral estimate
	mean := sum / float64(n)
	// Scale the average value by the total volume
	integral := volume * mean

	// Calculate variance and standard deviation
	variance := (sum_squared / float64(n)) - (mean * mean)
	standard_deviation := math.Sqrt(variance)

	// Calculate Standard Error
	standard_error := volume * (standard_deviation / math.Sqrt(float64(n)))

	watch.Stop()

	return integral, standard_error, float64(watch.Milliseconds())
}

// MultiMonteCarloIntegration performs multi-dimensional MISER Monte Carlo integration
// using recursive stratified sampling to reduce variance in high-dimensional spaces.
// Input: MultiFunction, Bounds, iterations, depth
// Output: Integral, Standard Error, Variance
func MultiMISERMonteCarloIntegration(f MultiFunction, b Bounds, n, d int) (float64, float64) {
	// Check depth and samples
	if d <= 0 || n < 100 {
		est, err, _ := MultiMonteCarloIntegration(f, b, n)
		return est, err
	}

	// Find dimension with highest variance
	dims := len(b.Min)
	var bestDim int
	var maxVar float64 = -1

	// Test each dimension for highest variance
	for dim := range dims {
		// Create test bounds that split this dimension
		testBounds := Bounds{
			Min: make([]float64, dims),
			Max: make([]float64, dims),
		}
		copy(testBounds.Min, b.Min)
		copy(testBounds.Max, b.Max)

		// Split this dimension in half
		mid := (b.Min[dim] + b.Max[dim]) / 2
		testBounds.Min[dim] = mid
		testBounds.Max[dim] = b.Max[dim]

		// Test variance in this half
		_, se1 := MultiMISERMonteCarloIntegration(f, testBounds, n/10, 0)
		testBounds.Min[dim] = b.Min[dim]
		testBounds.Max[dim] = mid
		_, se2 := MultiMISERMonteCarloIntegration(f, testBounds, n/10, 0)

		totalVar := se1*se1 + se2*se2
		if totalVar > maxVar {
			maxVar = totalVar
			bestDim = dim
		}
	}

	// Split the selected dimension
	mid := (b.Min[bestDim] + b.Max[bestDim]) / 2

	// Create two sub-bounds
	bounds1 := Bounds{
		Min: make([]float64, dims),
		Max: make([]float64, dims),
	}
	bounds2 := Bounds{
		Min: make([]float64, dims),
		Max: make([]float64, dims),
	}
	copy(bounds1.Min, b.Min)
	copy(bounds1.Max, b.Max)
	copy(bounds2.Min, b.Min)
	copy(bounds2.Max, b.Max)

	bounds1.Max[bestDim] = mid
	bounds2.Min[bestDim] = mid

	// Calculate sample sizes based on variance
	testBounds1 := bounds1
	testBounds2 := bounds2
	_, se1 := MultiMISERMonteCarloIntegration(f, testBounds1, n/10, 0)
	_, se2 := MultiMISERMonteCarloIntegration(f, testBounds2, n/10, 0)

	n1 := int(float64(n) * se1 / (se1 + se2))
	n2 := n - n1

	// Ensure minimum samples
	n1 = max(n1, 50)
	n2 = max(n2, 50)
	if n1+n2 > n {
		n2 = n - n1
	}

	// Recursive integration
	integral1, se1 := MultiMISERMonteCarloIntegration(f, bounds1, n1, d-1)
	integral2, se2 := MultiMISERMonteCarloIntegration(f, bounds2, n2, d-1)

	// Combine results
	totalIntegral := integral1 + integral2
	totalVariance := (se1*se1*float64(n1) + se2*se2*float64(n2)) / float64(n)
	totalStandardError := math.Sqrt(totalVariance)

	return totalIntegral, totalStandardError
}
