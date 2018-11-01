package histogram

import (
	"math"
	"math/rand"
	"sort"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Histogram", func() {
	blnk := seed()
	std := seed(39, 15, 43, 7, 43, 36, 47, 6, 40, 49, 41)

	DescribeTable("Quantile",
		func(h *Histogram, q float64, x float64) {
			Expect(h.Quantile(q)).To(BeNumerically("~", x, 0.1))
		},

		Entry("0%", std, 0.0, 6.0),
		Entry("25%", std, 0.25, 19.6),
		Entry("50%", std, 0.5, 39.8),
		Entry("75%", std, 0.75, 44.3),
		Entry("95%", std, 0.95, 47.2),
		Entry("99%", std, 0.99, 48.2),
		Entry("100%", std, 1.0, 49.0),
	)

	// inspired by https://github.com/aaw/histosketch/commit/d8284aa#diff-11101c92fbb1d58ccf30ca49764bf202R180
	// released into the public domain
	It("should accurately predict quantile", func() {
		N := 10000
		Q := []float64{0.0001, 0.001, 0.01, 0.1, 0.25, 0.35, 0.65, 0.75, 0.9, 0.99, 0.999, 0.9999}

		for seed := 0; seed < 10; seed++ {
			r := rand.New(rand.NewSource(int64(seed)))
			s := New(16)            // sketch
			x := make([]float64, N) // exact

			for i := 0; i < N; i++ {
				num := r.NormFloat64()
				s.Add(num)
				x[i] = num
			}
			sort.Float64s(x)

			for _, q := range Q {
				sQ := s.Quantile(q)
				xQ := x[int(float64(len(x))*q)]
				re := math.Abs((sQ - xQ) / xQ)

				Expect(re).To(BeNumerically("<", 0.09),
					"s.Quantile(%v) (got %.3f, want %.3f with seed = %v)", q, sQ, xQ, seed,
				)
			}
		}
	})

	It("should reject bad quantile inputs", func() {
		Expect(math.IsNaN(blnk.Quantile(0.5))).To(BeTrue())
		Expect(math.IsNaN(std.Quantile(-0.1))).To(BeTrue())
		Expect(math.IsNaN(std.Quantile(1.1))).To(BeTrue())
	})

	It("should calc sum", func() {
		Expect(math.IsNaN(blnk.Sum())).To(BeTrue())
		Expect(std.Sum()).To(Equal(366.0))
	})

	It("should calc mean", func() {
		Expect(math.IsNaN(blnk.Mean())).To(BeTrue())
		Expect(std.Mean()).To(BeNumerically("~", 33.27, 0.01))
	})

	It("should calc weight", func() {
		Expect(blnk.Weight()).To(Equal(0.0))
		Expect(std.Weight()).To(Equal(11.0))
	})

	It("should calc min", func() {
		Expect(math.IsNaN(blnk.Min())).To(BeTrue())
		Expect(std.Min()).To(Equal(float64(6)))
	})

	It("should copy", func() {
		c1 := std.Copy(nil)
		Expect(c1).To(Equal(std))

		t2 := seed(1, 2, 3, 4)
		c2 := std.Copy(t2)
		Expect(c2).To(Equal(std))
		Expect(c2).To(Equal(t2))
	})

	It("should calc max", func() {
		Expect(math.IsNaN(blnk.Max())).To(BeTrue())
		Expect(std.Max()).To(Equal(float64(49)))
	})

	It("should merge", func() {
		h2 := seed(11, 2, 3, 14, 7, 4)
		Expect(h2.Sum()).To(Equal(41.0))
		Expect(h2.bins).To(HaveLen(4))

		h2.MergeWith(std)
		Expect(h2.Sum()).To(Equal(407.0))
		Expect(h2.bins).To(HaveLen(4))
	})

	It("should add", func() {
		Expect(std.bins).To(HaveLen(4))
		Expect(std.bins).To(HaveCap(5))
		Expect(std.bins).To(Equal([]bin{
			{w: -2, v: 6.5},
			{w: 1, v: 15},
			{w: -4, v: 39},
			{w: -4, v: 45.5},
		}))
	})

})

// --------------------------------------------------------------------

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "histogram")
}

func seed(vv ...float64) *Histogram {
	h := New(4)
	for _, v := range vv {
		h.Add(v)
	}
	return h
}
