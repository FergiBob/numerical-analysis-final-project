package main

import (
	"fmt"
	"math"
	"numerical-analysis-final-project/internal"

	"github.com/bradhe/stopwatch"
)

func main() {

	fmt.Printf("\n\n\n") // for readability
	// Set up for tests
	n := 10000000

	// Test 0:
	var f internal.UnaryFunction = func(x float64) float64 { return math.Pow(x, 2) }
	a := 0.0
	b := 2.0

	fmt.Println("------------------------------------------------")
	fmt.Println("                Test 0: f(x) = x^2              ")
	fmt.Println("------------------------------------------------")
	fmt.Println("EXPECTED INTEGRAL: 2.666666666666666")
	internal.PrintTest(f, a, b, n)

	fmt.Print("\n\n\n")
	fmt.Println("------------------------------------------------")
	fmt.Println("      Test 1: Area of circle with radius 1      ")
	fmt.Println("------------------------------------------------")
	fmt.Println("EXPECTED INTEGRAL: 3.1415926")

	r := 1.0
	var g internal.UnaryFunction = func(x float64) float64 { return math.Sqrt(math.Pow(r, 2) - math.Pow(x, 2)) }

	fmt.Printf("\n	Single-Threaded Monte Carlo Integration")
	integral, errorval, time := internal.UnaryMonteCarloIntegration(g, -r, r, n, "single-threaded")
	fmt.Printf("\n		Integral %f", 2*integral)
	fmt.Printf("\n		Error value %f", 2*errorval)
	fmt.Printf("\n		Time Elapsed %f milliseconds", time)

	fmt.Printf("\n\n	Multi-Threaded Monte Carlo Integration")
	integral, errorval, time = internal.UnaryMonteCarloIntegration(g, -r, r, n, "multi-threaded")
	fmt.Printf("\n		Integral %f", 2*integral)
	fmt.Printf("\n		Error value %f", 2*errorval)
	fmt.Printf("\n		Time Elapsed %f milliseconds", time)

	fmt.Printf("\n\n	MISERMonteCarlo")
	watch := stopwatch.Start()
	integral_MISER, errorval_MISER := internal.UnaryMISERMonteCarlo(g, -r, r, n, 3)
	watch.Stop()
	fmt.Printf("\n		Integral %f", 2*integral_MISER)
	fmt.Printf("\n		Error value %f", 2*errorval_MISER)
	fmt.Printf("\n		Time Elapsed %f milliseconds", float64(watch.Milliseconds()))

	// Test 2: Integration of e^x over [0,1]

	// Test 3: neutron transport equation
	fmt.Print("\n\n\n") // for readability
}
