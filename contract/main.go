// Copyright (c) Marcel Moura
// SPDX-License-Identifier: MIT (see LICENSE)

package main

import (
	"fmt"
	"image"

	"github.com/gligneul/eggroll"
	"github.com/kelindar/bitmap"
	"github.com/marcelstanley/millionerc20"
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
		receipt, err := c.saveImage(env, input)
		if err != nil {
			return nil, fmt.Errorf("image rejected: %v", err)
		}
		return receipt, nil
	default:
		return nil, fmt.Errorf("invalid input: %v", input)
	}
}

// Set user image bits into the DApp bitmap
func (c *MillionErc20Contract) saveImage(env eggroll.Env, img *millionerc20.MetaImage) (*millionerc20.Confirmation, error) {
	if !isValid(img.Rect) {
		return /*&millionerc20.Confirmation{Value: false}*/ nil, fmt.Errorf("image '%v' is out of bounds", img.Rect)
	}

	rect := img.Rect
	start := rect.Min.X + rect.Min.Y*millionerc20.MAX_DIMENSION
	end := rect.Max.X*millionerc20.MAX_DIMENSION + rect.Max.Y
	width := rect.Dx()

	i := start
	j := i
	for {
		pixel := uint32(j)
		if c.Bmap.Contains(pixel) {
			return /*&millionerc20.Confirmation{Value: false}*/ nil, fmt.Errorf("bit %v already set", pixel)
		}
		c.Bmap.Set(pixel)

		if j >= end {
			// reached last pixel to be set
			break
		} else if j >= i+width {
			// break 2D line by skipping image width
			i = i + millionerc20.MAX_DIMENSION
			j = i
		} else {
			j += 1
		}
	}
	env.Logf("image saved at '%v'\n", img)
	return &millionerc20.Confirmation{Value: true}, nil
}

func isValid(rect image.Rectangle) bool {
	if rect.Min.X < 0 || rect.Min.Y < 0 || rect.Max.X > millionerc20.MAX_DIMENSION || rect.Max.Y > millionerc20.MAX_DIMENSION {
		return false
	}
	return true
}

func main() {
	eggroll.Roll(&MillionErc20Contract{})
}
