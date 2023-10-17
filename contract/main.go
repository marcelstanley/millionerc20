// Copyright (c) Marcel Moura
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"fmt"
	"millionerc20"

	"github.com/gligneul/eggroll"
	"github.com/kelindar/bitmap"
)

type MillionErc20Contract struct {
	eggroll.DefaultContract
	Bmap bitmap.Bitmap
}

func (c *MillionErc20Contract) Codecs() []eggroll.Codec {
	return millionerc20.Codecs()
}

func (c *MillionErc20Contract) Advance(env eggroll.Env) (any, error) {
	switch input := env.DecodeInput().(type) {
	case *millionerc20.MetaImage:
		return c.saveImage(env, input)
	default:
		return nil, fmt.Errorf("invalid input: %v", input)
	}
}

// Set user image bits into the DApp bitmap
func (c *MillionErc20Contract) saveImage(env eggroll.Env, img *millionerc20.MetaImage) (*millionerc20.Confirmation, error) {
	env.Logf("received user image with '%v'\n", img)
	rect := img.Rect
	start := rect.Min.X + rect.Min.Y*millionerc20.MAX
	end := rect.Max.X*millionerc20.MAX + rect.Max.Y
	width := rect.Dx()

	i := start
	j := i
	for {
		pixel := uint32(j)
		if c.Bmap.Contains(pixel) {
			return &millionerc20.Confirmation{Value: false}, fmt.Errorf("bit %v already set", pixel)
		}
		c.Bmap.Set(pixel)

		if j >= end {
			// reached last pixel to be set
			break
		} else if j >= i+width {
			// break 2D line by skipping image width
			i = i + millionerc20.MAX
			j = i
		} else {
			j += 1
		}
	}
	return &millionerc20.Confirmation{Value: true}, nil
}

func main() {
	eggroll.Roll(&MillionErc20Contract{})
}
