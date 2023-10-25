// Copyright (c) Marcel Moura
// SPDX-License-Identifier: MIT (see LICENSE)

package client

import (
	"context"
	"fmt"
	"time"

	"github.com/gligneul/eggroll"
	"github.com/gligneul/eggroll/eggtypes"
	"github.com/marcelstanley/millionerc20"
)

const (
	TIMEOUT       = 300 * time.Second
	INVALID_INDEX = -1
)

var INVALID_CONFIRMATION = millionerc20.Confirmation{false}

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
func SendAndCheck(input *millionerc20.MetaImage) (*millionerc20.Confirmation, error) {
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
	return confirm, nil
}
