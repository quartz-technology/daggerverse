package dockercompose

import "path"

type Volume struct {
	origin string
	target string
}

func (v *Volume) Name() string {
	return path.Base(v.origin)
}

func (v *Volume) Origin() string {
	return v.origin
}

func (v *Volume) Target() string {
	return v.target
}