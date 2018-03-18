// Copyright 2017 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package restorable

import (
	"github.com/hajimehoshi/ebiten/internal/affine"
	"github.com/hajimehoshi/ebiten/internal/graphics"
)

var (
	quadFloat32Num     = graphics.QuadVertexSizeInBytes() / 4
	theVerticesBackend = &verticesBackend{}
)

type verticesBackend struct {
	backend []float32
	head    int
}

func (v *verticesBackend) get() []float32 {
	const num = 1024
	if v.backend == nil {
		v.backend = make([]float32, quadFloat32Num*num)
	}
	s := v.backend[v.head : v.head+quadFloat32Num]
	v.head += quadFloat32Num
	if v.head+quadFloat32Num > len(v.backend) {
		v.backend = nil
		v.head = 0
	}
	return s
}

func vertices(width, height int, sx0, sy0, sx1, sy1 int, geo *affine.GeoM, color *affine.ColorM) []float32 {
	if sx0 >= sx1 || sy0 >= sy1 {
		return nil
	}
	if sx1 <= 0 || sy1 <= 0 {
		return nil
	}

	vs := theVerticesBackend.get()

	x0, y0 := 0.0, 0.0
	x1, y1 := float64(sx1-sx0), float64(sy1-sy0)

	// it really feels like we should be able to cache this computation
	// but it may not matter.
	w := 1
	h := 1
	for w < width {
		w *= 2
	}
	for h < height {
		h *= 2
	}
	wf := float32(w)
	hf := float32(h)
	u0, v0, u1, v1 := float32(sx0)/wf, float32(sy0)/hf, float32(sx1)/wf, float32(sy1)/hf

	x, y := geo.Apply32(x0, y0)
	// Vertex coordinates
	vs[0] = x
	vs[1] = y

	// Texture coordinates: first 2 values indicates the actual coodinate, and
	// the second indicates diagonally opposite coodinates.
	// The second is needed to calculate source rectangle size in shader programs.
	vs[2] = u0
	vs[3] = v0
	vs[4] = u1
	vs[5] = v1

	// and the same for the other three coordinates
	x, y = geo.Apply32(x1, y0)
	vs[26] = x
	vs[27] = y
	vs[28] = u1
	vs[29] = v0
	vs[30] = u0
	vs[31] = v1

	x, y = geo.Apply32(x0, y1)
	vs[52] = x
	vs[53] = y
	vs[54] = u0
	vs[55] = v1
	vs[56] = u1
	vs[57] = v0

	x, y = geo.Apply32(x1, y1)
	vs[78] = x
	vs[79] = y
	vs[80] = u1
	vs[81] = v1
	vs[82] = u0
	vs[83] = v0

	cBody, cTranslate := color.UnsafeElements()
	for i, v := range cBody {
		vs[6+i] = v
		vs[32+i] = v
		vs[58+i] = v
		vs[84+i] = v
	}
	for i, v := range cTranslate {
		vs[22+i] = v
		vs[48+i] = v
		vs[74+i] = v
		vs[100+i] = v
	}

	return vs
}
