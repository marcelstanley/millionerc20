// Copyright (c) Marcel Moura
// SPDX-License-Identifier: MIT (see LICENSE)

package millionerc20

import (
	"image"

	"github.com/gligneul/eggroll"
)

const MAX = 1000

type MetaImage struct {
	// Rect defines the placement of the user image
	Rect image.Rectangle
}

type Confirmation struct {
	Value bool
}

func Codecs() []eggroll.Codec {
	return []eggroll.Codec{
		eggroll.NewJSONCodec[MetaImage](),
		eggroll.NewJSONCodec[Confirmation](),
	}
}
