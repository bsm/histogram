package histogram_test

import (
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"

	. "github.com/bsm/histogram/v3"
)

func seedHist(vv ...float64) *Histogram {
	h := New(4)
	for _, v := range vv {
		h.Add(v)
	}
	return h
}

var (
	blank = seedHist()
	hist  = seedHist(39, 15, 43, 7, 43, 36, 47, 6, 40, 49, 41)
)

func expectQ(t *testing.T, q, exp float64) {
	t.Helper()

	if got := hist.Quantile(q); math.Abs(exp-got) > 0.1 {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestHistogram_Quantile(t *testing.T) {
	t.Run("0%", func(t *testing.T) { expectQ(t, 0.0, 6.0) })
	t.Run("25%", func(t *testing.T) { expectQ(t, 0.25, 19.6) })
	t.Run("50%", func(t *testing.T) { expectQ(t, 0.5, 39.8) })
	t.Run("75%", func(t *testing.T) { expectQ(t, 0.75, 44.3) })
	t.Run("95%", func(t *testing.T) { expectQ(t, 0.95, 47.2) })
	t.Run("99%", func(t *testing.T) { expectQ(t, 0.99, 48.2) })
	t.Run("100%", func(t *testing.T) { expectQ(t, 1.0, 49.0) })
}

// inspired by https://github.com/aaw/histosketch/commit/d8284aa#diff-11101c92fbb1d58ccf30ca49764bf202R180
// released into the public domain
func TestHistogram_Quantile_accuracy(t *testing.T) {
	N := 20_000
	Q1 := []float64{0.001, 0.01, 0.1, 0.25, 0.35, 0.65, 0.75, 0.9, 0.99, 0.999}
	Q2 := []float64{0.0001, 0.9999}

	for seed := int64(0); seed < 16; seed++ {
		r := rand.New(rand.NewSource(seed))
		h := New(100)           // histogram
		x := make([]float64, N) // exact

		for i := 0; i < N; i++ {
			num := r.NormFloat64() * 1
			h.Add(num)
			x[i] = num
		}
		sort.Float64s(x)

		for _, q := range Q1 {
			tQ := h.Quantile(q)
			xQ := x[int(float64(len(x))*q)]

			// allow ±2%
			if re := math.Abs((tQ - xQ) / xQ); re > 0.02 {
				t.Errorf("s.Quantile(%v) (got %.3f, want %.3f with seed = %v)", q, tQ, xQ, seed)
			}
		}

		for _, q := range Q2 {
			tQ := h.Quantile(q)
			xQ := x[int(float64(len(x))*q)]

			// allow ±10%
			if re := math.Abs((tQ - xQ) / xQ); re > 0.1 {
				t.Errorf("s.Quantile(%v) (got %.3f, want %.3f with seed = %v)", q, tQ, xQ, seed)
			}
		}
	}
}

func TestHistogram_Quantile_nan(t *testing.T) {
	if !math.IsNaN(blank.Quantile(0.5)) {
		t.Errorf("expected IsNaN to be true")
	}
	if !math.IsNaN(hist.Quantile(-0.1)) {
		t.Errorf("expected IsNaN to be true")
	}
	if !math.IsNaN(hist.Quantile(1.1)) {
		t.Errorf("expected IsNaN to be true")
	}
}

func TestHistogram_Sum(t *testing.T) {
	if !math.IsNaN(blank.Sum()) {
		t.Errorf("expected IsNaN to be true")
	}

	if exp, got := 366.0, hist.Sum(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
func TestHistogram_Mean(t *testing.T) {
	if !math.IsNaN(blank.Mean()) {
		t.Errorf("expected IsNaN to be true")
	}

	if exp, got := 33.27272727272727, hist.Mean(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestHistogram_Count(t *testing.T) {
	if exp, got := 0, blank.Count(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}

	if exp, got := 11, hist.Count(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestHistogram_Weight(t *testing.T) {
	if exp, got := 0.0, blank.Weight(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}

	if exp, got := 11.0, hist.Weight(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestHistogram_Min(t *testing.T) {
	if !math.IsNaN(blank.Min()) {
		t.Errorf("expected IsNaN to be true")
	}

	if exp, got := 6.0, hist.Min(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestHistogram_Max(t *testing.T) {
	if !math.IsNaN(blank.Max()) {
		t.Errorf("expected IsNaN to be true")
	}

	if exp, got := 49.0, hist.Max(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestHistogram_Copy(t *testing.T) {
	dupe := hist.Copy(nil)
	if !reflect.DeepEqual(hist, dupe) {
		t.Fatalf("expected %v, got %v", hist, dupe)
	}

	dupe.Add(34)
	if exp, got := 12, dupe.Count(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
	if exp, got := 11, hist.Count(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestHistogram_Merge(t *testing.T) {
	his2 := seedHist(11, 2, 3, 14, 7, 4)
	if exp, got := 41.0, his2.Sum(); exp != got {
		t.Fatalf("expected %v, got %v", exp, got)
	}

	his2.MergeWith(hist)
	if exp, got := 407.0, his2.Sum(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
	if exp, got := 366.0, hist.Sum(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestHistogram_Bin(t *testing.T) {
	if exp, got := 4, hist.NumBins(); exp != got {
		t.Fatalf("expected %v, got %v", exp, got)
	}

	v, w := hist.Bin(0)
	if exp := 6.5; exp != v {
		t.Errorf("expected %v, got %v", exp, v)
	}
	if exp := -2.0; exp != w {
		t.Errorf("expected %v, got %v", exp, w)
	}

	v, w = hist.Bin(1)
	if exp := 15.0; exp != v {
		t.Errorf("expected %v, got %v", exp, v)
	}
	if exp := 1.0; exp != w {
		t.Errorf("expected %v, got %v", exp, w)
	}

	v, w = hist.Bin(2)
	if exp := 39.0; exp != v {
		t.Errorf("expected %v, got %v", exp, v)
	}
	if exp := -4.0; exp != w {
		t.Errorf("expected %v, got %v", exp, w)
	}

	v, w = hist.Bin(3)
	if exp := 45.5; exp != v {
		t.Errorf("expected %v, got %v", exp, v)
	}
	if exp := -4.0; exp != w {
		t.Errorf("expected %v, got %v", exp, w)
	}
}
