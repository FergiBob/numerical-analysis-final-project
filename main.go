package main

import (
	"fmt"
	"math"
	"numerical-analysis-final-project/internal"
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
	internal.PrintTestCircle(f, r, n)

	// Test 3: Multi varible integration
	fmt.Print("\n\n\n")
	fmt.Println("------------------------------------------------")
	fmt.Println("      Test 2: Area of sphere with radius 1      ")
	fmt.Println("------------------------------------------------")
	internal.PrintTestSphere(1)

	fmt.Print("\n\n\n")
}
