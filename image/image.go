// Copyright (c) Marcel Moura
// SPDX-License-Identifier: MIT (see LICENSE)

package image

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/marcelstanley/millionerc20"
)

const DAPP_IMAGE_NAME = ".dapp_image.png"

var dappImage *image.RGBA

func init() {
	//Check for local file and either load it or create a new one as the DApp image
	var err error
	dappImage, err = Load(DAPP_IMAGE_NAME)
	if err != nil {
		dappImage = image.NewRGBA(image.Rect(0, 0, millionerc20.MAX_DIMENSION, millionerc20.MAX_DIMENSION))
		log.Printf("created %v\n", DAPP_IMAGE_NAME)
	}
}

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

	return nil
}

func Load(filename string) (*image.RGBA, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	b := img.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)

	return rgba, nil
}

func AddToDappImage(rect image.Rectangle, img image.Image) error {
	draw.Draw(dappImage, rect, img, image.Pt(0, 0), draw.Over)

	err := Save(DAPP_IMAGE_NAME, dappImage)
	if err != nil {
		return fmt.Errorf("image not saved: %v", err)
	}

	return nil
}
