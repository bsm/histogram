package histogram_test

import (
	"fmt"

	"github.com/bsm/histogram/v2"
)

func ExampleHistogram() {
	// Create a new instance
	h := histogram.New(16)

	// Add values
	for i := 1; i < 100; i++ {
		h.Add(float64(i))
	}

	fmt.Printf("min  : %.1f\n", h.Min())
	fmt.Printf("max  : %.1f\n", h.Max())
	fmt.Printf("sum  : %.0f\n", h.Sum())
	fmt.Printf("mean : %.1f\n", h.Mean())
	fmt.Printf("p50  : %.1f\n", h.Quantile(0.5))
	fmt.Printf("p95  : %.1f\n", h.Quantile(0.95))

	// Output:
	// min  : 1.0
	// max  : 99.0
	// sum  : 4950
	// mean : 50.0
	// p50  : 49.8
	// p95  : 94.6
}
