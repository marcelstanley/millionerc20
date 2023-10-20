// Copyright (c) Marcel Moura
// SPDX-License-Identifier: MIT (see LICENSE)

package cmd

import (
	"image"
	"log"

	"github.com/marcelstanley/millionerc20"
	"github.com/marcelstanley/millionerc20/client"
	cli_image "github.com/marcelstanley/millionerc20/image"
	"github.com/spf13/cobra"
)

var (
	addImageCmd = &cobra.Command{
		Use:   "addImage",
		Short: "Add an image to the Million ERC-20 DApp",
		Long: `Add an image to the Million ERC-20 DApp.

The DApp keeps a 1000 px x 1000 px image which has its pixels for sale.

A client may request an image file to be painted in the DApp image.

Only blank spaces are available for painting.
A request for paint an image  over an already painted image will be rejected.
		`,
		Run: addImage,
	}

	argFilename string
	argX        int
	argY        int
)

func init() {
	addImageCmd.Flags().StringVarP(&argFilename, "filename", "f", "", "File name of the image to be submitted")
	addImageCmd.MarkFlagRequired("filename")
	addImageCmd.Flags().IntVarP(&argX, "x-coordinate", "x", 0, "Pixel coordinate X")
	addImageCmd.Flags().IntVarP(&argY, "y-coordinate", "y", 0, "Pixel coordinate Y")

	rootCmd.AddCommand(addImageCmd)
}

func addImage(cmd *cobra.Command, args []string) {
	img, err := cli_image.Load(argFilename)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded %v\n", argFilename)

	// New rect translated by (argX, argY)
	rect := img.Rect.Add(image.Pt(argX, argY))

	index, err := client.Send(&millionerc20.MetaImage{rect})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("sent image to DApp (inputIndex: %v)\n", index)

	_, err = client.Check(index)
	if err != nil {
		log.Printf("image rejected: %v", err)
		return
	}
	log.Printf("image accepted\n")

	err = cli_image.AddToDappImage(rect, img)
	if err != nil {
		log.Printf("DApp Image could not be updated: %v", err)
		return
	}
}
