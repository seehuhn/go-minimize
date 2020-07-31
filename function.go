package minimize

import "math"

type wrapper struct {
	f func([]float64) float64
	n int

	cache []float64
	tmp   []float64
}

func (w *wrapper) Get(x []float64) float64 {
	if len(w.cache) == 0 {
		w.n = len(x)
		w.cache = make([]float64, (w.n+1)*(w.n+4))
		w.tmp = make([]float64, w.n+1)

		// mark all initial cache entries as invalid
		for i := 0; i < len(w.cache); i++ {
			w.cache[i] = math.NaN()
		}
	}

	n := w.n
	stride := n + 1
	cache := w.cache
	cacheLen := len(cache)
search:
	for base := 0; base+n < cacheLen; base += stride {
		for i := 0; i < n; i++ {
			if cache[base+i] != x[i] {
				continue search
			}
		}

		// move the entry to the front
		copy(w.tmp, cache[base:base+n+1])
		copy(cache[stride:base+stride], cache[0:base])
		copy(cache[0:stride], w.tmp)
		return cache[n]
	}

	res := w.f(x)

	copy(cache[stride:], cache[0:])
	copy(cache[:n], x)
	cache[n] = res
	return res
}

// Function finds an (approximate) local minimum of `f` near `x0`. The
// parameter `ε` gives the size of the initial simplex.
//
// This is a wrapper around `Minimize()`, with caching of returned function
// values to avoid unnecessary calls to `f`.
func Function(f func([]float64) float64, x0 []float64, ε float64) []float64 {
	w := &wrapper{f: f}
	return Minimize(func(x, y []float64) bool {
		return w.Get(x) < w.Get(y)
	}, x0, ε)
}
