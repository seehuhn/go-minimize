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

package minimize

import (
	"math"
	"testing"
)

func Quadratic(x []float64) float64 {
	res := 0.0
	for _, xi := range x {
		res += xi * xi
	}
	return res
}

func Rosenbrock(x []float64) float64 {
	p := 1 - x[0]
	q := x[1] - x[0]*x[0]
	res := p*p + 100*q*q
	return res
}

func Himmelblau(x []float64) float64 {
	p := x[0]*x[0] + x[1] - 11
	q := x[0] + x[1]*x[1] - 7
	return p*p + q*q
}

func Zero(c []float64) float64 {
	return 0
}

func TestMinimize(t *testing.T) {
	targets := []struct {
		name string
		f    func([]float64) float64
	}{
		{"Quadratic", Quadratic},
		{"Rosenbrock", Rosenbrock},
		{"Himmelblau", Himmelblau},
		{"Zero", Zero},
	}
	for _, x := range []float64{-100, -2.1, -2, -1, 0, 1, 2, 50, 200} {
		for _, y := range []float64{-200, -2, -1, 0, 1, 2, 5, 100} {
			x0 := []float64{x, y}

			for _, target := range targets {
				min := Function(target.f, x0, 0.1)
				if target.f(min) >= 1e-6 {
					t.Error(target.name+":", min, "->", target.f(min))
				}
			}
		}
	}
}

func TestQuadratic(t *testing.T) {
	x0 := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	min := Function(Quadratic, x0, 0.1)
	for i, xi := range min {
		if math.Abs(xi) >= 1e-6 {
			t.Errorf("wrong result for x_%d: expected 0, got %f", i, xi)
		}
	}
}

func BenchmarkZero(b *testing.B) {
	x0 := []float64{1, 2, 3, 4, 5, 6}
	for i := 0; i < b.N; i++ {
		_ = Function(Zero, x0, 1.0)
	}
}

func BenchmarkQuadratic(b *testing.B) {
	x0 := []float64{1, 2, 3, 4, 5, 6}
	for i := 0; i < b.N; i++ {
		_ = Function(Quadratic, x0, 1.0)
	}
}
