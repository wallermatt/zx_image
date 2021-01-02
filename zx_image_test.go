package main_test

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wallermatt/z80/zx_image"
)

const TestSnapshotFile = "testData/testSnapshot.sna"

func Test_ReadSnapshot(t *testing.T) {
	s := zx_image.ReadSnapshot(TestSnapshotFile)
	assert.Equal(t, 49179, len(s))
	assert.Equal(t, uint8(0x3f), s[0])
	assert.Equal(t, uint8(0x0), s[len(s)-1])
}

func Test_LoadScrMemory(t *testing.T) {
	s := zx_image.ReadSnapshot(TestSnapshotFile)
	scrMemory := zx_image.LoadScrMemory(s)
	assert.Equal(t, uint8(0x0), scrMemory[0][0])
	assert.Equal(t, uint8(0xff), scrMemory[110][25])
	assert.Equal(t, uint8(0x0), scrMemory[191][31])
}

func Test_LoadScrAttributes(t *testing.T) {
	s := zx_image.ReadSnapshot(TestSnapshotFile)
	scrAttributes := zx_image.LoadScrAttributes(s)
	assert.Equal(t, uint8(0x38), scrAttributes[0][0])
	assert.Equal(t, uint8(0x3b), scrAttributes[10][25])
	assert.Equal(t, uint8(0x38), scrAttributes[23][31])
}

func Test_GetPaperAndInk(t *testing.T) {
	paper, ink := zx_image.GetPaperAndInk(0x00)
	assert.Equal(t, color.RGBA{0, 0, 0, 255}, paper)
	assert.Equal(t, color.RGBA{0, 0, 0, 255}, ink)

	paper, ink = zx_image.GetPaperAndInk(0x50)
	assert.Equal(t, color.RGBA{255, 0, 0, 255}, paper)
	assert.Equal(t, color.RGBA{0, 0, 0, 255}, ink)

	paper, ink = zx_image.GetPaperAndInk(0xff)
	assert.Equal(t, color.RGBA{255, 255, 255, 255}, paper)
	assert.Equal(t, color.RGBA{255, 255, 255, 255}, ink)

	paper, ink = zx_image.GetPaperAndInk(0x10)
	assert.Equal(t, color.RGBA{215, 0, 0, 255}, paper)
	assert.Equal(t, color.RGBA{0, 0, 0, 255}, ink)
}

func Test_GetScrMemoryFromXY(t *testing.T) {
	s := zx_image.ReadSnapshot(TestSnapshotFile)
	scrMemory := zx_image.LoadScrMemory(s)

	value := zx_image.GetScrMemoryFromXY(0, 0, scrMemory)
	assert.Equal(t, uint8(0x0), value)

	value = zx_image.GetScrMemoryFromXY(0, 120, scrMemory)
	assert.Equal(t, uint8(0xff), value)

	value = zx_image.GetScrMemoryFromXY(21, 99, scrMemory)
	assert.Equal(t, uint8(0x1f), value)
}

func Test_BuildImage(t *testing.T) {
	s := zx_image.ReadSnapshot(TestSnapshotFile)
	scrMemory := zx_image.LoadScrMemory(s)
	scrAttributes := zx_image.LoadScrAttributes(s)
	img := zx_image.BuildImage(scrMemory, scrAttributes)
	assert.Equal(t, 196608, len(img.Pix))
	assert.Equal(t, []byte{0xff, 0xd7, 0x0, 0xd7}, img.Pix[122879:122883])
}
