// Copyright (c) Marcel Moura
// SPDX-License-Identifier: MIT (see LICENSE)

package millionerc20

import (
	"context"
	"fmt"
	"image"
	"testing"
	"time"

	"github.com/gligneul/eggroll"
	"github.com/gligneul/eggroll/eggeth"
	"github.com/gligneul/eggroll/eggtest"
)

const TEST_TIMEOUT = 300 * time.Second

func TestAllImagesAreAccepted(t *testing.T) {
	tester := eggtest.NewIntegrationTester(t)
	defer tester.Close()

	ctx, cancel := context.WithTimeout(context.Background(), TEST_TIMEOUT)
	defer cancel()

	client, signer, err := eggroll.NewDevClient(ctx, Codecs())
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	inputs := []any{
		&MetaImage{image.Rect(0, 0, 1, 1)},
		&MetaImage{image.Rect(2, 2, 3, 3)},
	}

	for _, input := range inputs {
		if err := hasImageBeenAccepted(t, ctx, client, signer, input); err != nil {
			t.Fatalf("image should have been accepted: %v", err)
		}
	}
}

func TestImageIsNotAccepted(t *testing.T) {
	tester := eggtest.NewIntegrationTester(t)
	defer tester.Close()

	ctx, cancel := context.WithTimeout(context.Background(), TEST_TIMEOUT)
	defer cancel()

	client, signer, err := eggroll.NewDevClient(ctx, Codecs())
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	inputs := []any{
		&MetaImage{image.Rect(0, 0, 1, 1)},
		&MetaImage{image.Rect(1, 1, 2, 2)},
	}

	if err := hasImageBeenAccepted(t, ctx, client, signer, inputs[0]); err != nil {
		t.Fatalf("image should have been accepted: %v", err)
	}

	// Second image must be rejected, as it intersects the first one
	if err := hasImageBeenAccepted(t, ctx, client, signer, inputs[0]); err == nil {
		t.Fatalf("image should have been rejected: %v", err)
	}
}

func hasImageBeenAccepted(
	t *testing.T, ctx context.Context, client *eggroll.Client,
	signer eggeth.Signer, input any) error {

	lastInputIndex, err := client.SendInput(ctx, signer, input)
	if err != nil {
		t.Fatalf("failed to send input: %v", err)
	}

	result, err := client.WaitFor(ctx, lastInputIndex)
	if err != nil {
		t.Fatalf("failed to wait for input: %v", err)
	}

	confirm, ok := client.DecodeReturn(result).(*Confirmation)
	if !ok {
		return fmt.Errorf("expected confirmation, but got %v", client.DecodeReturn(result))
	}

	if confirm.Value != true {
		return fmt.Errorf("invalid value %v; expected %v", confirm.Value, true)
	}
	return nil
}
