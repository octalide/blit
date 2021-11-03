package blit

type Orienter struct {
	Parent *Orienter
	Vec
	R float32
}

// Pos gets the position vector relative to the parent
func (o *Orienter) Pos() Vec {
	if o.Parent != nil {
		return o.Parent.Pos().Add(o.Vec)
	}

	return o.Vec
}

// Rot gets the rotation relative to the parent
func (o *Orienter) Rot() float32 {
	if o.Parent != nil {
		return o.Parent.Rot() + o.R
	}

	return o.R
}

// Mat calculates the translation matrix relative to the parent
func (o *Orienter) Mat() Mat {
	// NOTE: both Pos() and Rot() are relative to the parent already
	m := Ident()
	m = m.Pos(o.Pos())
	m = m.Rot(o.Rot())

	return m
}
