package imgtoico

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

type Converter struct {
	From    string
	Resized *image.RGBA

	ext string
}

func NewConverter() *Converter {
	return &Converter{
		From: "",
	}
}

func (c *Converter) Validate() (err error) {
	if c.From == "" {
		log.Panic("from file is required")
	}

	c.From, err = filepath.Abs(c.From)
	if err != nil {
		log.Panic(err)
	}

	// check if file exists
	if _, err := os.Stat(c.From); os.IsNotExist(err) {
		log.Panic(err)
	}

	// check if file is an image
	c.ext = filepath.Ext(c.From)

	// check if file is a supported image format
	switch c.ext {
	// supported file formats
	case ".png":
	case ".jpg":
	case ".jpeg":

	default: // unsupported file format
		log.Panic("unsupported file format")
	}

	return err
}

func createCanvas(img image.Image) (canvas *image.RGBA, maxDimLen int) {
	// convert image to square and center input image.
	// create square canvas
	maxDimLen = int(math.Max(float64(img.Bounds().Dx()), float64(img.Bounds().Dy())))
	canvas = image.NewRGBA(image.Rect(0, 0, maxDimLen, maxDimLen))

	return canvas, maxDimLen
}

func drawToCanvas(canvas *image.RGBA, original image.Image, maxDimLen int) {
	paddingT := (maxDimLen - original.Bounds().Dy()) / 2
	paddingL := (maxDimLen - original.Bounds().Dx()) / 2

	// add input, source, image into the canvas.
	draw.Draw(
		canvas,
		image.Rect(paddingL, paddingT, maxDimLen-paddingL, maxDimLen-paddingT),
		original,
		image.Pt(0, 0),
		draw.Over,
	)
}

func resizeCanvas(canvas *image.RGBA) (resized *image.RGBA) {
	resized = image.NewRGBA(image.Rect(0, 0, 256, 256))

	draw.CatmullRom.Scale(resized, resized.Bounds(), canvas, canvas.Bounds(), draw.Over, nil)

	return resized
}

// ICO encode flow: https://github.com/biessek/golang-ico/blob/master/writer.go
func writeToIcoFile(resized *image.RGBA, outFile string) (err error) {
	pngBuffer := new(bytes.Buffer)
	pngWriter := bufio.NewWriter(pngBuffer)

	// write prepared icon to writer
	err = png.Encode(pngWriter, resized)
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
			Width:  uint8(resized.Bounds().Dx()),
			Height: uint8(resized.Bounds().Dy()),
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

	return err
}

func (c *Converter) Convert() {
	c.Validate()

	fromFile, err := os.Open(c.From)
	if err != nil {
		log.Panic(err)
	}
	defer fromFile.Close()

	var img image.Image

	switch c.ext {
	case ".png":
		img, err = png.Decode(fromFile)
	case ".jpg":
		fallthrough
	case ".jpeg":
		img, err = jpeg.Decode(fromFile)
	}

	if err != nil {
		log.Panic(err)
	}

	canvas, maxDimLen := createCanvas(img)
	drawToCanvas(canvas, img, maxDimLen)

	c.Resized = resizeCanvas(canvas)
}

func (c *Converter) Write() {
	var outFile string

	switch c.ext {
	case ".png":
		outFile = strings.Replace(c.From, ".png", ".ico", 1)
	case ".jpg":
		fallthrough
	case ".jpeg":
		outFile = strings.NewReplacer(".jpg", ".ico", ".jpeg", ".ico").Replace(c.From)
	}

	c.WriteTo(outFile)
}

func (c *Converter) WriteTo(outFile string) {
	err := writeToIcoFile(c.Resized, outFile)
	if err != nil {
		log.Panic(err)
	}
}
