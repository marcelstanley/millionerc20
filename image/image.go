// Copyright (c) Marcel Moura
// SPDX-License-Identifier: MIT (see LICENSE)

package image

import (
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"
)

func Save(filename string, img image.Image) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}

	log.Printf("saved image to %v\n", f.Name())
	return nil
}

func Load(filename string) (*image.RGBA, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	log.Printf("loaded image from %v\n", f.Name())
	return Decode(f)
}

func Decode(r io.Reader) (*image.RGBA, error) {
	img, err := png.Decode(r)
	if err != nil {
		return nil, err
	}

	b := img.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)

	return rgba, nil
}
