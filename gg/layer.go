// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gg

import (
	"image/color"

	"github.com/aclements/go-gg/table"
)

// LayerLines is like LayerPaths, but sorts the data by the "x"
// property.
func LayerLines() Plotter {
	return func(p *Plot) {
		b := p.mustGetBinding("x")
		if b.isConstant {
			p.Add(LayerPaths())
		} else {
			p.Save().SortBy(b.col).Add(LayerPaths()).Restore()
		}
	}
}

// LayerPaths layers a path connecting successive data points in each
// group. By default the path is stroked, but if the "fill" property
// is set, it produces solid polygons.
//
// The "x" and "y" properties define the location of each point on the
// path. They are connected with straight lines.
//
// The "color" property defines the color of each line segment in the
// path. If two subsequent points have different color value, the
// color of the first is used. "Color" defaults to black.
//
// The "fill" property defines the fill color of the polygon with
// vertices at each data point. Only the fill value of the first point
// in each group is used. "Fill" defaults to transparent.
//
// XXX Perhaps the theme should provide default values for things like
// "color". That would suggest we need to resolve defaults like that
// at render time. Possibly a special scale that gets values from the
// theme could be used to resolve them.
//
// XXX strokeOpacity, fillOpacity, strokeWidth, what other properties
// do SVG strokes have?
//
// XXX Should the set of known styling bindings be fixed, and all
// possible rendering targets have to know what to do with them, or
// should the rendering target be able to have different styling
// bindings they understand (presumably with some reasonable base
// set)? If the renderer can determine the known bindings, we would
// probably just capture the environment here (and make it so a
// captured environment does not change) and hand that to the renderer
// later.
func LayerPaths() Plotter {
	return func(p *Plot) {
		xb, yb := p.mustGetBinding("x"), p.mustGetBinding("y")
		colorb := p.getBinding("color")
		fillb := p.getBinding("fill")

		if colorb == nil {
			// XXX Yuck
			colorb = &binding{isConstant: true, constant: color.Black, scales: map[table.GroupID]Scaler{table.RootGroupID: NewIdentityScale()}}
		}
		if fillb == nil {
			fillb = &binding{isConstant: true, constant: color.Transparent, scales: map[table.GroupID]Scaler{table.RootGroupID: NewIdentityScale()}}
		}

		data := p.Data()
		for _, gid := range data.Groups() {
			switch data.Table(gid).Len() {
			case 0:
				continue

			case 1:
				// TODO: Depending on the stroke cap,
				// this *could* be well-defined.
				Warning.Print("cannot layer path through 1 point; ignoring")
				continue
			}

			// TODO: Check that scales map to the right types.
			//checkScaleRange("LayerPaths", x, ScaleTypeReal)
			//checkScaleRange("LayerPaths", y, ScaleTypeReal)
			//checkScaleRange("LayerPaths", colorb, ScaleTypeColor)
			//checkScaleRange("LayerPaths", fill, ScaleTypeColor)

			// TODO: When I register a mark, perhaps I
			// have to include what group it belongs to so
			// rendering knows which facet to put it in.
			mark := &markPath{
				// TODO: Should these names match the
				// visual property names? color vs
				// stroke.
				p.use("x", xb, gid),
				p.use("y", yb, gid),
				p.use("stroke", colorb, gid),
				p.use("fill", fillb, gid),
			}
			p.marks = append(p.marks, mark)
		}
	}
}
