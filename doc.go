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

// Package minimize implements the Nelder–Mead simplex algorithm for
// minimization.
//
// The specific variant of the algorithm is the one described in
// Jeffrey C. Lagarias, James A. Reeds, Margaret H. Wright, and Paul
// E. Wright: Convergence Properties of the Nelder-Mead Simplex Method
// In Low Dimensions.  SIAM J. Optim, Vol. 9 (1998), No. 1,
// pp. 112-147.  https://doi.org/10.1137/S1052623496303470
package minimize
