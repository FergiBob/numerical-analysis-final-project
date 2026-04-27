package internal

import (
	"fmt"
	"math"

	"github.com/bradhe/stopwatch"
)

// helper function used to print out each series of monte-carlo tests
func PrintTest(f UnaryFunction, a, b float64, n int) {

	fmt.Printf("\n	Single-Threaded Monte Carlo Integration")
	integral, errorval, time := UnaryMonteCarloIntegration(f, a, b, n, "single-threaded")
	fmt.Printf("\n		Integral %f", integral)
	fmt.Printf("\n		Error value %f", errorval)
	fmt.Printf("\n		Time Elapsed %f ms", time)

	fmt.Printf("\n\n	Multi-Threaded Monte Carlo Integration")
	integral, errorval, time = UnaryMonteCarloIntegration(f, a, b, n, "multi-threaded")
	fmt.Printf("\n		Integral %f", integral)
	fmt.Printf("\n		Error value %f", errorval)
	fmt.Printf("\n		Time Elapsed %f ms", time)

	fmt.Printf("\n\n	MISERMonteCarlo")
	watch := stopwatch.Start()
	integral_MISER, errorval_MISER := UnaryMISERMonteCarlo(f, a, b, n, 3)
	watch.Stop()
	fmt.Printf("\n		Integral %f", integral_MISER)
	fmt.Printf("\n		Error value %f", errorval_MISER)
	fmt.Printf("\n		Time Elapsed %f ms", float64(watch.Milliseconds()))
}

func PrintTestCircle(f UnaryFunction, r float64, n int) {

	var g UnaryFunction = func(x float64) float64 { return math.Sqrt(math.Pow(r, 2) - math.Pow(x, 2)) }

	fmt.Printf("\n	Single-Threaded Monte Carlo Integration")
	integral, errorval, time := UnaryMonteCarloIntegration(g, -r, r, n, "single-threaded")
	fmt.Printf("		Integral %f", 2*integral)
	fmt.Printf("\n		Error value %f", 2*errorval)
	fmt.Printf("\n		Time Elapsed %f ms", time)

	fmt.Printf("\n\n	Multi-Threaded Monte Carlo Integration")
	integral, errorval, time = UnaryMonteCarloIntegration(g, -r, r, n, "multi-threaded")
	fmt.Printf("		Integral %f", 2*integral)
	fmt.Printf("\n		Error value %f", 2*errorval)
	fmt.Printf("\n		Time Elapsed %f ms", time)

	fmt.Printf("\n\n	MISERMonteCarlo")
	watch := stopwatch.Start()
	integral_MISER, errorval_MISER := UnaryMISERMonteCarlo(g, -r, r, n, 3)
	watch.Stop()
	fmt.Printf("\n		Integral %f", 2*integral_MISER)
	fmt.Printf("\n		Error value %f", 2*errorval_MISER)
	fmt.Printf("\n		Time Elapsed %f ms", float64(watch.Milliseconds()))
}

func PrintTestSphere(radius float64) {

	// Define the bounds (a cube that perfectly fits the sphere)
	bounds := Bounds{
		Min: []float64{-1.0, -1.0, -1.0},
		Max: []float64{1.0, 1.0, 1.0},
	}

	// Define the indicator function
	// Equation of a sphere: x² + y² + z² <= r²
	sphereFunc := func(p []float64) float64 {
		sumSq := 0.0
		for _, val := range p {
			sumSq += val * val
		}

		if sumSq <= radius*radius {
			return 1.0 // "Hit" - Inside the sphere
		}
		return 0.0 // "Miss" - Outside the sphere
	}

	// Run the integration
	n := 1000000 // 1 million samples
	fmt.Printf("EXPECTED VOLUME:  %.6f\n", (4.0/3.0)*math.Pi*math.Pow(radius, 3))

	fmt.Println("\n	MultiMonteCarloIntegration Test")
	volume, err, ms := MultiMonteCarloIntegration(sphereFunc, bounds, n)

	fmt.Printf("		Estimated Volume: %.6f", volume)
	fmt.Printf("\n		Standard Error:   %.6f", err)
	fmt.Printf("\n		Time Elapsed:       %.2f ms", ms)

	fmt.Println("\n\n	MultiMISERMonteCarloIntegration Test")
	watch := stopwatch.Start()
	volume, err = MultiMISERMonteCarloIntegration(sphereFunc, bounds, n, 5)
	watch.Stop()
	fmt.Printf("		Estimated Volume: %.6f", volume)
	fmt.Printf("\n		Standard Error:   %.6f", err)
	fmt.Printf("\n		Time Elapsed:     %.2f ms", float64(watch.Milliseconds()))
}

func PrintTest20DObject() {
	// Define 20D bounds for a hypercube from -1 to 1 in each dimension
	dims := 20
	bounds := Bounds{
		Min: make([]float64, dims),
		Max: make([]float64, dims),
	}
	for i := range dims {
		bounds.Min[i] = -1.0
		bounds.Max[i] = 1.0
	}

	hypercubeFunc := func(p []float64) float64 {
		// All points in the hypercube are "inside"
		return 1.0
	}

	n := 10000000 // 10 million samples
	fmt.Printf("   EXPECTED VOLUME: %.6f\n", math.Pow(2, float64(dims)))
	volume, err, ms := MultiMonteCarloIntegration(hypercubeFunc, bounds, n)
	fmt.Printf("   Estimated Volume: %.6f\n", volume)
	fmt.Printf("   Standard Error:   %.6f\n", err)
	fmt.Printf("   Time Elapsed:     %.2f ms\n", ms)

	fmt.Println("\n\n   MultiMISERMonteCarloIntegration Test")
	watch := stopwatch.Start()
	volume, err = MultiMISERMonteCarloIntegration(hypercubeFunc, bounds, n, 5)
	watch.Stop()
	fmt.Printf("   Estimated Volume: %.6f\n", volume)
	fmt.Printf("   Standard Error:   %.6f\n", err)
	fmt.Printf("   Time Elapsed:     %.2f ms\n", float64(watch.Milliseconds()))
}
