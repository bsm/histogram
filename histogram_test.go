package histogram

import (
	"math"
	"math/rand"
	"sort"
	"testing"

	. "github.com/bsm/ginkgo"
	. "github.com/bsm/ginkgo/extensions/table"
	. "github.com/bsm/gomega"
)

var _ = Describe("Histogram", func() {
	var blank, subject *Histogram

	BeforeEach(func() {
		blank = seedHist()
		subject = seedHist(39, 15, 43, 7, 43, 36, 47, 6, 40, 49, 41)
	})

	DescribeTable("Quantile",
		func(q float64, x float64) {
			Expect(subject.Quantile(q)).To(BeNumerically("~", x, 0.1))
		},

		Entry("0%", 0.0, 6.0),
		Entry("25%", 0.25, 19.6),
		Entry("50%", 0.5, 39.8),
		Entry("75%", 0.75, 44.3),
		Entry("95%", 0.95, 47.2),
		Entry("99%", 0.99, 48.2),
		Entry("100%", 1.0, 49.0),
	)

	// inspired by https://github.com/aaw/histosketch/commit/d8284aa#diff-11101c92fbb1d58ccf30ca49764bf202R180
	// released into the public domain
	It("should accurately predict quantile", func() {
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
				re := math.Abs((tQ - xQ) / xQ)

				Expect(re).To(BeNumerically("<", 0.02), // allow ±2%
					"s.Quantile(%v) (got %.3f, want %.3f with seed = %v)", q, tQ, xQ, seed,
				)
			}

			for _, q := range Q2 {
				tQ := h.Quantile(q)
				xQ := x[int(float64(len(x))*q)]
				re := math.Abs((tQ - xQ) / xQ)

				Expect(re).To(BeNumerically("<", 0.1), // allow ±10%
					"s.Quantile(%v) (got %.3f, want %.3f with seed = %v)", q, tQ, xQ, seed,
				)
			}
		}
	})

	It("should reject bad quantile inputs", func() {
		Expect(math.IsNaN(blank.Quantile(0.5))).To(BeTrue())
		Expect(math.IsNaN(subject.Quantile(-0.1))).To(BeTrue())
		Expect(math.IsNaN(subject.Quantile(1.1))).To(BeTrue())
	})

	It("should calc sum", func() {
		Expect(math.IsNaN(blank.Sum())).To(BeTrue())
		Expect(subject.Sum()).To(Equal(366.0))
	})

	It("should calc mean", func() {
		Expect(math.IsNaN(blank.Mean())).To(BeTrue())
		Expect(subject.Mean()).To(BeNumerically("~", 33.27, 0.01))
	})

	It("should calc count", func() {
		Expect(blank.Count()).To(Equal(0))
		Expect(subject.Count()).To(Equal(11))
	})

	It("should calc weight", func() {
		Expect(blank.Weight()).To(Equal(0.0))
		Expect(subject.Weight()).To(Equal(11.0))
	})

	It("should calc min", func() {
		Expect(math.IsNaN(blank.Min())).To(BeTrue())
		Expect(subject.Min()).To(Equal(float64(6)))
	})

	It("should copy", func() {
		c1 := subject.Copy(nil)
		Expect(c1).To(Equal(subject))

		t2 := seedHist(1, 2, 3, 4)
		c2 := subject.Copy(t2)
		Expect(c2).To(Equal(subject))
		Expect(c2).To(Equal(t2))
	})

	It("should calc max", func() {
		Expect(math.IsNaN(blank.Max())).To(BeTrue())
		Expect(subject.Max()).To(Equal(float64(49)))
	})

	It("should merge", func() {
		h2 := seedHist(11, 2, 3, 14, 7, 4)
		Expect(h2.Sum()).To(Equal(41.0))
		Expect(h2.bins).To(HaveLen(4))

		h2.MergeWith(subject)
		Expect(h2.Sum()).To(Equal(407.0))
		Expect(h2.bins).To(HaveLen(4))
	})

	It("should add with weight", func() {
		Expect(subject.bins).To(HaveLen(4))
		Expect(subject.bins).To(HaveCap(5))
		Expect(subject.bins).To(Equal([]bin{
			{w: -2, v: 6.5},
			{w: 1, v: 15},
			{w: -4, v: 39},
			{w: -4, v: 45.5},
		}))

		subject.AddWeight(6.5, 2.0)
		subject.AddWeight(15, 3.0)
		Expect(subject.bins).To(Equal([]bin{
			{w: -4, v: 6.5},
			{w: 4, v: 15},
			{w: -4, v: 39},
			{w: -4, v: 45.5},
		}))
	})

	It("should return bin count and data", func() {
		Expect(subject.NumBins()).To(Equal(4))

		v, w := subject.Bin(0)
		Expect(v).To(BeNumerically("==", 6.5))
		Expect(w).To(BeNumerically("==", -2))

		v, w = subject.Bin(1)
		Expect(v).To(BeNumerically("==", 15))
		Expect(w).To(BeNumerically("==", 1))

		v, w = subject.Bin(2)
		Expect(v).To(BeNumerically("==", 39))
		Expect(w).To(BeNumerically("==", -4))

		v, w = subject.Bin(3)
		Expect(v).To(BeNumerically("==", 45.5))
		Expect(w).To(BeNumerically("==", -4))
	})

})

// --------------------------------------------------------------------

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "histogram")
}

func seedHist(vv ...float64) *Histogram {
	h := New(4)
	for _, v := range vv {
		h.Add(v)
	}
	return h
}
