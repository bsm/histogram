package histogram_test

import (
	"math/rand"
	"testing"

	"github.com/bsm/histogram"
)

func BenchmarkHistogram_Add(b *testing.B) {
	h := histogram.New(16)
	r := rand.New(rand.NewSource(0))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Add(r.NormFloat64())
	}
}

func BenchmarkHistogram_Quantile(b *testing.B) {
	h := histogram.New(16)
	r := rand.New(rand.NewSource(0))
	for i := 1; i < 10000; i++ {
		h.Add(r.NormFloat64())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Quantile(0.95)
	}
}
