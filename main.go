package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

func main() {
	var err error

	inFile, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Panicln(err)
	}

	input, err := os.Open(inFile)
	if err != nil {
		log.Panicln(err)
	}
	defer input.Close()

	var img image.Image
	var outFile string

	switch {
	case strings.HasSuffix(inFile, ".png"):
		outFile = strings.Replace(inFile, ".png", ".ico", 1)
		img, err = png.Decode(input)
	case strings.HasSuffix(inFile, ".jpg"):
		fallthrough
	case strings.HasSuffix(inFile, ".jpeg"):
		outFile = strings.NewReplacer(".jpg", ".ico", ".jpeg", ".ico").Replace(inFile)
		img, err = jpeg.Decode(input)
	default:
		log.Panicln("unsupported file format")
	}

	if err != nil {
		log.Panicln(err)
	}

	// convert image to square and center input image.
	// create square canvas
	_max := int(math.Max(float64(img.Bounds().Dx()), float64(img.Bounds().Dy())))
	src := image.NewRGBA(image.Rect(0, 0, _max, _max))

	padding_t := (_max - int(img.Bounds().Dy())) / 2
	padding_l := (_max - int(img.Bounds().Dx())) / 2

	// add input, source, image into the canvas.
	draw.Draw(
		src,
		image.Rect(padding_l, padding_t, _max-padding_l, _max-padding_t),
		img,
		image.Pt(0, 0),
		draw.Over,
	)

	// resize
	// create destination canvas.
	dst := image.NewRGBA(image.Rect(0, 0, 256, 256))

	// scale
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	// write
	// ICO encode flow: https://github.com/biessek/golang-ico/blob/master/writer.go
	pngBuffer := new(bytes.Buffer)
	pngWriter := bufio.NewWriter(pngBuffer)

	// write prepared icon to writer
	err = png.Encode(pngWriter, dst)
	if err != nil {
		log.Panicln(err)
	}
	err = pngWriter.Flush() // clear
	if err != nil {
		log.Panicln(err)
	}

	// prepare byte buffer
	byteBuffer := new(bytes.Buffer)
	err = binary.Write(
		byteBuffer,
		binary.LittleEndian,
		struct {
			Zero   uint16
			Type   uint16
			Number uint16
		}{
			0, 1, 1,
		})
	if err != nil {
		log.Panicln(err)
	}

	err = binary.Write(
		byteBuffer,
		binary.LittleEndian,
		struct {
			Width   byte
			Height  byte
			Palette byte
			_       byte
			Plane   uint16
			Bits    uint16
			Size    uint32
			Offset  uint32
		}{
			Width:  uint8(dst.Bounds().Dx()),
			Height: uint8(dst.Bounds().Dy()),
			Plane:  1,
			Bits:   32,
			Size:   uint32(len(pngBuffer.Bytes())),
			Offset: 22,
		},
	)
	if err != nil {
		log.Panicln(err)
	}

	// create output file, overwrite if exists
	output, err := os.Create(outFile)
	if err != nil {
		log.Panicln(err)
	}
	defer output.Close()

	// write icon
	_, err = output.Write(byteBuffer.Bytes())
	if err != nil {
		log.Panicln(err)
	}
	_, err = output.Write(pngBuffer.Bytes())
	if err != nil {
		log.Panicln(err)
	}
}
