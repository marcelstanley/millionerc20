// Copyright (c) Marcel Moura
// SPDX-License-Identifier: MIT (see LICENSE)

package client

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"log"
	"time"

	"github.com/gligneul/eggroll"
	"github.com/gligneul/eggroll/eggtypes"
	"github.com/marcelstanley/millionerc20"
	dapp_image "github.com/marcelstanley/millionerc20/image"
)

const (
	TIMEOUT       = 300 * time.Second
	INVALID_INDEX = -1
	// XXX The image should be stored in a DB or something similar
	DAPP_IMAGE_NAME = "/tmp/.dapp_image.png"
)

var INVALID_CONFIRMATION = millionerc20.Confirmation{false}

var dappImage *image.RGBA

func init() {
	//Check for local file and either load it or create a new DApp image
	var err error
	dappImage, err = dapp_image.Load(DAPP_IMAGE_NAME)
	if err != nil {
		dappImage, err = dapp_image.Load("assets/blank-dapp-image.png")
		if err != nil {
			dappImage = image.NewRGBA(image.Rect(0, 0, millionerc20.MAX_DIMENSION, millionerc20.MAX_DIMENSION))
			log.Println("created new DApp image")
		}
	}
}

func GetDappImage() *image.RGBA {
	return dappImage
}

func AddToDappImage(img *image.RGBA) error {
	draw.Draw(dappImage, img.Bounds(), img, image.Pt(0, 0), draw.Over)

	err := dapp_image.Save(DAPP_IMAGE_NAME, dappImage)
	if err != nil {
		return fmt.Errorf("image not saved: %v", err)
	}

	return nil
}

func Send(input *millionerc20.MetaImage) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	client, signer, err := eggroll.NewDevClient(ctx, millionerc20.Codecs())
	if err != nil {
		return INVALID_INDEX, fmt.Errorf("failed to create client: %v", err)
	}

	inputIndex, err := client.SendInput(ctx, signer, input)
	if err != nil {
		return INVALID_INDEX, fmt.Errorf("failed to send input: %v", err)
	}
	return inputIndex, nil
}

func Check(inputIndex int) (*millionerc20.Confirmation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	client, _, err := eggroll.NewDevClient(ctx, millionerc20.Codecs())
	if err != nil {
		return &INVALID_CONFIRMATION, fmt.Errorf("failed to create client: %v", err)
	}

	result, err := client.WaitFor(ctx, inputIndex)
	if err != nil {
		return &INVALID_CONFIRMATION, fmt.Errorf("failed to wait for input: %v", err)
	}

	confirm, decoded := client.DecodeReturn(result).(*millionerc20.Confirmation)
	if !decoded {
		return &INVALID_CONFIRMATION, fmt.Errorf("inputIndex %v rejected by DApp", inputIndex)
	}
	return confirm, nil
}

// TODO Do we need to pass the confirmation over to the caller?
func SendImageAndWait(pos image.Point, img *image.RGBA) (*millionerc20.Confirmation, error) {
	// Translate image
	img.Rect = img.Rect.Add(pos)
	log.Printf("sending image with bounds %v", img.Rect.Bounds())
	input := &millionerc20.MetaImage{img.Rect}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	client, signer, err := eggroll.NewDevClient(ctx, millionerc20.Codecs())
	if err != nil {
		return &INVALID_CONFIRMATION, fmt.Errorf("failed to create client: %v", err)
	}

	inputIndex, err := client.SendInput(ctx, signer, input)
	if err != nil {
		return &INVALID_CONFIRMATION, fmt.Errorf("failed to send input: %v", err)
	}

	result, err := client.WaitFor(ctx, inputIndex)
	if err != nil {
		return &INVALID_CONFIRMATION, fmt.Errorf("failed to wait for input: %v", err)
	}

	if result.Status == eggtypes.CompletionStatusRejected {
		// TODO Think about either having the contract state contain the error
		// message returned by Contract.Advance or check if eggroll may pass
		// the error along for rejected inputs
		fmt.Printf("Result status: %v", result.Status)
		return &INVALID_CONFIRMATION, fmt.Errorf("image with bounds %v rejected", input.Rect)
	}

	confirm, decoded := client.DecodeReturn(result).(*millionerc20.Confirmation)
	if !decoded {
		return &INVALID_CONFIRMATION, fmt.Errorf("undecoded result: %v", result)
	}

	if err := AddToDappImage(img); err != nil {
		return &INVALID_CONFIRMATION, fmt.Errorf("could not update DApp image")
	}

	return confirm, nil
}
