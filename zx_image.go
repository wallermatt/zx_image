package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type ScrRow [32]byte
type ScrMemory [192]ScrRow
type ScrAttributes [24]ScrRow

type SpecColour struct {
	normal color.RGBA
	bright color.RGBA
}
type SpecColoursType [8]SpecColour

var SpecColours = SpecColoursType{
	SpecColour{color.RGBA{0, 0, 0, 255}, color.RGBA{0, 0, 0, 255}},
	SpecColour{color.RGBA{0, 0, 215, 255}, color.RGBA{0, 0, 255, 255}},
	SpecColour{color.RGBA{215, 0, 0, 255}, color.RGBA{255, 0, 0, 255}},
	SpecColour{color.RGBA{215, 0, 215, 255}, color.RGBA{255, 0, 255, 255}},
	SpecColour{color.RGBA{0, 215, 0, 255}, color.RGBA{0, 255, 0, 255}},
	SpecColour{color.RGBA{0, 215, 215, 255}, color.RGBA{0, 255, 255, 255}},
	SpecColour{color.RGBA{215, 215, 0, 255}, color.RGBA{255, 255, 0, 255}},
	SpecColour{color.RGBA{215, 215, 215, 255}, color.RGBA{255, 255, 255, 255}},
}

const (
	TestSnapshotFile   = "testData/testSnapshot.sna"
	scrMemoryStart     = 27
	scrAttributesStart = 6171
	width              = 256
	height             = 192
)

// ReadSnapshot reads a ZX Spectrum snapshot file into a byte slice.
func ReadSnapshot(f string) []byte {
	s, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}
	return s
}

// LoadScrMemory loads screen memory from snapshot into ScrMemory array of arrays.
func LoadScrMemory(s []byte) ScrMemory {
	snapshotIndex := scrMemoryStart
	scrMemory := ScrMemory{}
	for y, row := range scrMemory {
		for x, _ := range row {
			scrMemory[y][x] = s[snapshotIndex]
			snapshotIndex += 1
		}
	}
	return scrMemory
}

// LoadScrAttributes loads screen attributes from snapshot into ScrAttributes array of arrays.
func LoadScrAttributes(s []byte) ScrAttributes {
	snapshotIndex := scrAttributesStart
	scrAttributes := ScrAttributes{}
	for y, row := range scrAttributes {
		for x, _ := range row {
			scrAttributes[y][x] = s[snapshotIndex]
			snapshotIndex += 1
		}
	}
	return scrAttributes
}

// GetPaperAndInk determines paper and ink colours from byte (attribute memory value).
func GetPaperAndInk(value byte) (paper color.RGBA, ink color.RGBA) {
	bright := (value << 1) / 128
	paperIndex := (value >> 3) % 8
	inkIndex := value % 8
	if bright == 0 {
		paper = SpecColours[paperIndex].normal
		ink = SpecColours[inkIndex].normal
	} else {
		paper = SpecColours[paperIndex].bright
		ink = SpecColours[inkIndex].bright
	}
	return paper, ink
}

// GetScrMemoryFromXY returns screen memory location for given coordinates.
func GetScrMemoryFromXY(x int, y int, scrMemory ScrMemory) byte {
	memX := x / 8
	block := y / 64
	blockOffset := y % 64
	row := blockOffset % 8
	rowOffset := blockOffset / 8
	memY := block*64 + row*8 + rowOffset
	return scrMemory[memY][memX]
}

// GetXPixelFromByte determines for a given x coordinate whether the relevant memory location specifies a pixel
// should be displayed.
func GetXPixelFromByte(x int, memory byte) bool {
	offsetX := x % 8
	pixelPosition := int(math.Pow(2, float64(7-offsetX)))
	return int(memory)&pixelPosition != 0
}

// Builds image of ZX Spectrum's display from specified screen and attribute memory.
// It builds the image pixel by pixel, determining what the background (paper) and
// foreground (ink) colours are for that pixel, then it determines if a pixel is
// printed or not.
func BuildImage(scrMemory ScrMemory, scrAttributes ScrAttributes) *image.RGBA {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			gridX := x / 8
			gridY := y / 8
			paper, ink := GetPaperAndInk(scrAttributes[gridY][gridX])
			memory := GetScrMemoryFromXY(x, y, scrMemory)
			if GetXPixelFromByte(x, memory) {
				img.Set(x, y, ink)
			} else {
				img.Set(x, y, paper)
			}
		}
	}
	return img
}

// SaveImage saves screen image as a .png file.
func SaveImage(img *image.RGBA, imageFilename string) error {
	f, err := os.Create(imageFilename)
	if err != nil {
		return err
	}

	png.Encode(f, img)
	return nil
}

// CreateImageFromSnapshot controls the whole conversion process.
func CreateImageFromSnapshot(snapshotFile string, imageFilename string) {
	s := ReadSnapshot(snapshotFile)
	scrMemory := LoadScrMemory(s)
	scrAttributes := LoadScrAttributes(s)

	img := BuildImage(scrMemory, scrAttributes)
	err := SaveImage(img, imageFilename)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Image %s created from snapshot %s \n", imageFilename, snapshotFile)
}

// Snapshot file defaults to the test file if no arguments, otherwise it's the first argument.
// The image name defaults to the name of the snapshot file (but ends in .png not .sna) unless
// second argument is supplied.
func main() {
	snapshotFile := TestSnapshotFile
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		snapshotFile = argsWithoutProg[0]
	}
	splitF := strings.Split(snapshotFile, "/")
	imageFilename := splitF[len(splitF)-1]
	imageFilename = strings.Split(imageFilename, ".")[0] + ".png"
	if len(argsWithoutProg) > 1 {
		imageFilename = argsWithoutProg[1]
	}
	CreateImageFromSnapshot(snapshotFile, imageFilename)
}
