// seehuhn.de/go/minimize - the simplex algorithm of Nelder and Mead
// Copyright (C) 2019  Jochen Voss <voss@seehuhn.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package minimize // import "seehuhn.de/go/minimize"

import (
	"sort"
)

// The following parameters refer to the description of the method in
// Jeffrey C. Lagarias, James A. Reeds, Margaret H. Wright, and Paul
// E. Wright: Convergence Properties of the Nelder-Mead Simplex Method
// In Low Dimensions.  SIAM J. Optim, Vol. 9 (1998), No. 1,
// pp. 112-147.  https://doi.org/10.1137/S1052623496303470
const (
	ρ = 1   // reflection parameter
	χ = 2   // expansion parameter
	γ = 0.5 // contraction parameter
	σ = 0.5 // shrinkage parameter
)

type state struct {
	LessFn func(x, y []float64) bool
	N      int
	X      []float64
}

func (s *state) Point(i int) []float64 {
	n := s.N
	return s.X[i*n : (i+1)*n]
}

func (s *state) Len() int {
	return s.N + 1
}

func (s *state) Less(i, j int) bool {
	xi := s.Point(i)
	xj := s.Point(j)
	return s.LessFn(xi, xj)
}

func (s *state) Swap(i, j int) {
	n := s.N
	xi := s.Point(i)
	xj := s.Point(j)
	tmp := s.Point(n + 1)
	copy(tmp, xi)
	copy(xi, xj)
	copy(xj, tmp)
}

func (s *state) Init(x []float64, ε float64) {
	for k := 0; k <= s.N; k++ {
		point := s.Point(k)
		copy(point, x)
		if k < s.N {
			point[k] += ε
		}
	}
	sort.Sort(s)
}

// Insert point `src`, when we already know that the new position will
// be one of p_i, ..., p_j.
func (s *state) Insert(src, i, j int) {
	// Invariant: Less(src, i-1) == false, Less(src, j) == true
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		// i ≤ h < j
		if !s.Less(src, h) {
			i = h + 1 // preserves Less(src, i-1) == false
		} else {
			j = h // preserves Less(src, j) == true
		}
	}
	// now i == j is the smallest index where s.Less(src, i) is true

	n := s.N
	copy(s.X[(i+1)*n:(n+1)*n], s.X[i*n:n*n])
	copy(s.Point(i), s.Point(src))
}

// Compute the average of p_0, ..., p_{n-1} and store in p_{n+1}
func (s *state) Centroid() {
	n := s.N
	cent := s.Point(n + 1)
	for i := 0; i < n; i++ {
		cent[i] = 0
	}
	for k := 0; k < n; k++ {
		point := s.Point(k)
		for i := 0; i < n; i++ {
			cent[i] += point[i]
		}
	}
	for i := 0; i < n; i++ {
		cent[i] /= float64(s.N)
	}
}

// let p_a = (1-λ)*p_b + λ*p_c
func (s *state) Shift(a, b, c int, λ float64) {
	n := s.N
	pa := s.Point(a)
	pb := s.Point(b)
	pc := s.Point(c)
	for i := 0; i < n; i++ {
		pa[i] = (1-λ)*pb[i] + λ*pc[i]
	}
}

// Move all points closer to p_0
func (s *state) Shrink() {
	n := s.N
	best := s.Point(0)
	for k := 1; k <= n; k++ {
		pt := s.Point(k)
		for i := 0; i < n; i++ {
			pt[i] = (1-σ)*best[i] + σ*pt[i]
		}
	}
}

// Minimize finds an (approximate) local minimum near `x0`.  The parameter `ε`
// gives the size of the initial simplex.
func Minimize(less func(x, y []float64) bool, x0 []float64, ε float64) []float64 {
	n := len(x0)

	// Allocate an array for the n+1 vertices of the simplex, together
	// with three scratch vertices.
	s := &state{
		LessFn: less,
		N:      n,
		X:      make([]float64, (n+4)*n),
	}
	s.Init(x0, ε)

	shrinkCount := 0
	for step := 0; step < 100000; step++ {
		s.Centroid() // stored in p_{n+1}

		s.Shift(n+2, n+1, n, -ρ) // reflect
		winner := s.Less(n+2, 0)
		if !winner && s.Less(n+2, n-1) {
			s.Insert(n+2, 1, n-1)
			continue
		}

		if winner {
			s.Shift(n+3, n+1, n+2, χ) // expand
			if s.Less(n+3, n+2) {
				s.Insert(n+3, 0, 0)
				shrinkCount = 0
			} else {
				s.Insert(n+2, 0, 0)
			}
			continue
		}

		if s.Less(n+2, n) {
			s.Shift(n+2, n+1, n+2, γ) // outside contraction
		} else {
			s.Shift(n+2, n+1, n, γ) // inside contraction
		}
		if s.Less(n+2, n) {
			s.Insert(n+2, 0, n)
			continue
		}

		s.Shrink()
		sort.Sort(s)
		shrinkCount++

		if shrinkCount > 100 {
			break
		}
	}

	res := make([]float64, n)
	copy(res, s.Point(0))
	return res
}
