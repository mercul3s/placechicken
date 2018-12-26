package imaging

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name      string
		w, h      int
		c         color.Color
		dstBounds image.Rectangle
		dstPix    []uint8
	}{
		{
			"New 1x1 transparent",
			1, 1,
			color.Transparent,
			image.Rect(0, 0, 1, 1),
			[]uint8{0x00, 0x00, 0x00, 0x00},
		},
		{
			"New 1x2 red",
			1, 2,
			color.RGBA{255, 0, 0, 255},
			image.Rect(0, 0, 1, 2),
			[]uint8{0xff, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0xff},
		},
		{
			"New 2x1 white",
			2, 1,
			color.White,
			image.Rect(0, 0, 2, 1),
			[]uint8{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		{
			"New 3x3 with alpha",
			3, 3,
			color.NRGBA{0x01, 0x23, 0x45, 0x67},
			image.Rect(0, 0, 3, 3),
			[]uint8{
				0x01, 0x23, 0x45, 0x67, 0x01, 0x23, 0x45, 0x67, 0x01, 0x23, 0x45, 0x67,
				0x01, 0x23, 0x45, 0x67, 0x01, 0x23, 0x45, 0x67, 0x01, 0x23, 0x45, 0x67,
				0x01, 0x23, 0x45, 0x67, 0x01, 0x23, 0x45, 0x67, 0x01, 0x23, 0x45, 0x67,
			},
		},
		{
			"New 0x0 white",
			0, 0,
			color.White,
			image.Rect(0, 0, 0, 0),
			nil,
		},
		{
			"New 800x600 custom",
			800, 600,
			color.NRGBA{1, 2, 3, 4},
			image.Rect(0, 0, 800, 600),
			bytes.Repeat([]byte{1, 2, 3, 4}, 800*600),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := New(tc.w, tc.h, tc.c)
			want := image.NewNRGBA(tc.dstBounds)
			want.Pix = tc.dstPix
			if !compareNRGBA(got, want, 0) {
				t.Fatalf("got result %#v want %#v", got, want)
			}
		})
	}
}

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		New(1024, 1024, color.White)
	}
}

func TestClone(t *testing.T) {
	testCases := []struct {
		name string
		src  image.Image
		want *image.NRGBA
	}{
		{
			"Clone NRGBA",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 0, 1),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff},
			},
		},
		{
			"Clone NRGBA64",
			&image.NRGBA64{
				Rect:   image.Rect(-1, -1, 0, 1),
				Stride: 1 * 8,
				Pix: []uint8{
					0x00, 0x00, 0x11, 0x11, 0x22, 0x22, 0x33, 0x33,
					0xcc, 0xcc, 0xdd, 0xdd, 0xee, 0xee, 0xff, 0xff,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff},
			},
		},
		{
			"Clone RGBA",
			&image.RGBA{
				Rect:   image.Rect(-1, -1, 0, 2),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 3),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x55, 0xaa, 0x33, 0xcc, 0xdd, 0xee, 0xff},
			},
		},
		{
			"Clone RGBA64",
			&image.RGBA64{
				Rect:   image.Rect(-1, -1, 0, 2),
				Stride: 1 * 8,
				Pix: []uint8{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x11, 0x11, 0x22, 0x22, 0x33, 0x33,
					0xcc, 0xcc, 0xdd, 0xdd, 0xee, 0xee, 0xff, 0xff,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 3),
				Stride: 1 * 4,
				Pix:    []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x55, 0xaa, 0x33, 0xcc, 0xdd, 0xee, 0xff},
			},
		},
		{
			"Clone Gray",
			&image.Gray{
				Rect:   image.Rect(-1, -1, 0, 1),
				Stride: 1 * 1,
				Pix:    []uint8{0x11, 0xee},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1 * 4,
				Pix:    []uint8{0x11, 0x11, 0x11, 0xff, 0xee, 0xee, 0xee, 0xff},
			},
		},
		{
			"Clone Gray16",
			&image.Gray16{
				Rect:   image.Rect(-1, -1, 0, 1),
				Stride: 1 * 2,
				Pix:    []uint8{0x11, 0x11, 0xee, 0xee},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1 * 4,
				Pix:    []uint8{0x11, 0x11, 0x11, 0xff, 0xee, 0xee, 0xee, 0xff},
			},
		},
		{
			"Clone Alpha",
			&image.Alpha{
				Rect:   image.Rect(-1, -1, 0, 1),
				Stride: 1 * 1,
				Pix:    []uint8{0x11, 0xee},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 2),
				Stride: 1 * 4,
				Pix:    []uint8{0xff, 0xff, 0xff, 0x11, 0xff, 0xff, 0xff, 0xee},
			},
		},
		{
			"Clone YCbCr",
			&image.YCbCr{
				Rect:           image.Rect(-1, -1, 5, 0),
				SubsampleRatio: image.YCbCrSubsampleRatio444,
				YStride:        6,
				CStride:        6,
				Y:              []uint8{0x00, 0xff, 0x7f, 0x26, 0x4b, 0x0e},
				Cb:             []uint8{0x80, 0x80, 0x80, 0x6b, 0x56, 0xc0},
				Cr:             []uint8{0x80, 0x80, 0x80, 0xc0, 0x4b, 0x76},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 6, 1),
				Stride: 6 * 4,
				Pix: []uint8{
					0x00, 0x00, 0x00, 0xff,
					0xff, 0xff, 0xff, 0xff,
					0x7f, 0x7f, 0x7f, 0xff,
					0x7f, 0x00, 0x00, 0xff,
					0x00, 0x7f, 0x00, 0xff,
					0x00, 0x00, 0x7f, 0xff,
				},
			},
		},
		{
			"Clone YCbCr 444",
			&image.YCbCr{
				Y:              []uint8{0x4c, 0x69, 0x1d, 0xb1, 0x96, 0xe2, 0x26, 0x34, 0xe, 0x59, 0x4b, 0x71, 0x0, 0x4c, 0x99, 0xff},
				Cb:             []uint8{0x55, 0xd4, 0xff, 0x8e, 0x2c, 0x01, 0x6b, 0xaa, 0xc0, 0x95, 0x56, 0x40, 0x80, 0x80, 0x80, 0x80},
				Cr:             []uint8{0xff, 0xeb, 0x6b, 0x36, 0x15, 0x95, 0xc0, 0xb5, 0x76, 0x41, 0x4b, 0x8c, 0x80, 0x80, 0x80, 0x80},
				YStride:        4,
				CStride:        4,
				SubsampleRatio: image.YCbCrSubsampleRatio444,
				Rect:           image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
			&image.NRGBA{
				Pix:    []uint8{0xff, 0x0, 0x0, 0xff, 0xff, 0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xff, 0x49, 0xe1, 0xca, 0xff, 0x0, 0xff, 0x0, 0xff, 0xff, 0xff, 0x0, 0xff, 0x7f, 0x0, 0x0, 0xff, 0x7f, 0x0, 0x7f, 0xff, 0x0, 0x0, 0x7f, 0xff, 0x0, 0x7f, 0x7f, 0xff, 0x0, 0x7f, 0x0, 0xff, 0x82, 0x7f, 0x0, 0xff, 0x0, 0x0, 0x0, 0xff, 0x4c, 0x4c, 0x4c, 0xff, 0x99, 0x99, 0x99, 0xff, 0xff, 0xff, 0xff, 0xff},
				Stride: 16,
				Rect:   image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
		},
		{
			"Clone YCbCr 440",
			&image.YCbCr{
				Y:              []uint8{0x4c, 0x69, 0x1d, 0xb1, 0x96, 0xe2, 0x26, 0x34, 0xe, 0x59, 0x4b, 0x71, 0x0, 0x4c, 0x99, 0xff},
				Cb:             []uint8{0x2c, 0x01, 0x6b, 0xaa, 0x80, 0x80, 0x80, 0x80},
				Cr:             []uint8{0x15, 0x95, 0xc0, 0xb5, 0x80, 0x80, 0x80, 0x80},
				YStride:        4,
				CStride:        4,
				SubsampleRatio: image.YCbCrSubsampleRatio440,
				Rect:           image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
			&image.NRGBA{
				Pix:    []uint8{0x0, 0xb5, 0x0, 0xff, 0x86, 0x86, 0x0, 0xff, 0x77, 0x0, 0x0, 0xff, 0xfb, 0x7d, 0xfb, 0xff, 0x0, 0xff, 0x1, 0xff, 0xff, 0xff, 0x1, 0xff, 0x80, 0x0, 0x1, 0xff, 0x7e, 0x0, 0x7e, 0xff, 0xe, 0xe, 0xe, 0xff, 0x59, 0x59, 0x59, 0xff, 0x4b, 0x4b, 0x4b, 0xff, 0x71, 0x71, 0x71, 0xff, 0x0, 0x0, 0x0, 0xff, 0x4c, 0x4c, 0x4c, 0xff, 0x99, 0x99, 0x99, 0xff, 0xff, 0xff, 0xff, 0xff},
				Stride: 16,
				Rect:   image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
		},
		{
			"Clone YCbCr 422",
			&image.YCbCr{
				Y:              []uint8{0x4c, 0x69, 0x1d, 0xb1, 0x96, 0xe2, 0x26, 0x34, 0xe, 0x59, 0x4b, 0x71, 0x0, 0x4c, 0x99, 0xff},
				Cb:             []uint8{0xd4, 0x8e, 0x01, 0xaa, 0x95, 0x40, 0x80, 0x80},
				Cr:             []uint8{0xeb, 0x36, 0x95, 0xb5, 0x41, 0x8c, 0x80, 0x80},
				YStride:        4,
				CStride:        2,
				SubsampleRatio: image.YCbCrSubsampleRatio422,
				Rect:           image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
			&image.NRGBA{
				Pix:    []uint8{0xe2, 0x0, 0xe1, 0xff, 0xff, 0x0, 0xfe, 0xff, 0x0, 0x4d, 0x36, 0xff, 0x49, 0xe1, 0xca, 0xff, 0xb3, 0xb3, 0x0, 0xff, 0xff, 0xff, 0x1, 0xff, 0x70, 0x0, 0x70, 0xff, 0x7e, 0x0, 0x7e, 0xff, 0x0, 0x34, 0x33, 0xff, 0x1, 0x7f, 0x7e, 0xff, 0x5c, 0x58, 0x0, 0xff, 0x82, 0x7e, 0x0, 0xff, 0x0, 0x0, 0x0, 0xff, 0x4c, 0x4c, 0x4c, 0xff, 0x99, 0x99, 0x99, 0xff, 0xff, 0xff, 0xff, 0xff},
				Stride: 16,
				Rect:   image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
		},
		{
			"Clone YCbCr 420",
			&image.YCbCr{
				Y:       []uint8{0x4c, 0x69, 0x1d, 0xb1, 0x96, 0xe2, 0x26, 0x34, 0xe, 0x59, 0x4b, 0x71, 0x0, 0x4c, 0x99, 0xff},
				Cb:      []uint8{0x01, 0xaa, 0x80, 0x80},
				Cr:      []uint8{0x95, 0xb5, 0x80, 0x80},
				YStride: 4, CStride: 2,
				SubsampleRatio: image.YCbCrSubsampleRatio420,
				Rect:           image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
			&image.NRGBA{
				Pix:    []uint8{0x69, 0x69, 0x0, 0xff, 0x86, 0x86, 0x0, 0xff, 0x67, 0x0, 0x67, 0xff, 0xfb, 0x7d, 0xfb, 0xff, 0xb3, 0xb3, 0x0, 0xff, 0xff, 0xff, 0x1, 0xff, 0x70, 0x0, 0x70, 0xff, 0x7e, 0x0, 0x7e, 0xff, 0xe, 0xe, 0xe, 0xff, 0x59, 0x59, 0x59, 0xff, 0x4b, 0x4b, 0x4b, 0xff, 0x71, 0x71, 0x71, 0xff, 0x0, 0x0, 0x0, 0xff, 0x4c, 0x4c, 0x4c, 0xff, 0x99, 0x99, 0x99, 0xff, 0xff, 0xff, 0xff, 0xff},
				Stride: 16,
				Rect:   image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 4, Y: 4}},
			},
		},
		{
			"Clone Paletted",
			&image.Paletted{
				Rect:   image.Rect(-1, -1, 5, 0),
				Stride: 6 * 1,
				Palette: color.Palette{
					color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
					color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
					color.NRGBA{R: 0x7f, G: 0x7f, B: 0x7f, A: 0xff},
					color.NRGBA{R: 0x7f, G: 0x00, B: 0x00, A: 0xff},
					color.NRGBA{R: 0x00, G: 0x7f, B: 0x00, A: 0xff},
					color.NRGBA{R: 0x00, G: 0x00, B: 0x7f, A: 0xff},
				},
				Pix: []uint8{0x0, 0x1, 0x2, 0x3, 0x4, 0x5},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 6, 1),
				Stride: 6 * 4,
				Pix: []uint8{
					0x00, 0x00, 0x00, 0xff,
					0xff, 0xff, 0xff, 0xff,
					0x7f, 0x7f, 0x7f, 0xff,
					0x7f, 0x00, 0x00, 0xff,
					0x00, 0x7f, 0x00, 0xff,
					0x00, 0x00, 0x7f, 0xff,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Clone(tc.src)
			delta := 0
			if _, ok := tc.src.(*image.YCbCr); ok {
				delta = 1
			}
			if !compareNRGBA(got, tc.want, delta) {
				t.Fatalf("got result %#v want %#v", got, tc.want)
			}
		})
	}
}

func TestCrop(t *testing.T) {
	testCases := []struct {
		name string
		src  image.Image
		r    image.Rectangle
		want *image.NRGBA
	}{
		{
			"Crop 2x3 2x1",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
			image.Rect(-1, 0, 1, 1),
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 1),
				Stride: 2 * 4,
				Pix: []uint8{
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Crop(tc.src, tc.r)
			if !compareNRGBA(got, tc.want, 0) {
				t.Fatalf("got result %#v want %#v", got, tc.want)
			}
		})
	}
}

func BenchmarkCrop(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Crop(testdataBranchesJPG, image.Rect(100, 100, 300, 300))
	}
}

func TestCropCenter(t *testing.T) {
	testCases := []struct {
		name string
		src  image.Image
		w, h int
		want *image.NRGBA
	}{
		{
			"CropCenter 2x3 2x1",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
			2, 1,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 1),
				Stride: 2 * 4,
				Pix: []uint8{
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
				},
			},
		},
		{
			"CropCenter 2x3 0x1",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
			0, 1,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 0, 0),
				Stride: 0,
				Pix:    []uint8{},
			},
		},
		{
			"CropCenter 2x3 5x5",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
			5, 5,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 3),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := CropCenter(tc.src, tc.w, tc.h)
			if !compareNRGBA(got, tc.want, 0) {
				t.Fatalf("got result %#v want %#v", got, tc.want)
			}
		})
	}
}

func TestCropAnchor(t *testing.T) {
	testCases := []struct {
		name   string
		src    image.Image
		w, h   int
		anchor Anchor
		want   *image.NRGBA
	}{
		{
			"CropAnchor 4x4 2x2 TopLeft",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 3, 3),
				Stride: 4 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
			2, 2,
			TopLeft,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
				},
			},
		},
		{
			"CropAnchor 4x4 2x2 Top",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 3, 3),
				Stride: 4 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
			2, 2,
			Top,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b,
					0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b,
				},
			},
		},
		{
			"CropAnchor 4x4 2x2 TopRight",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 3, 3),
				Stride: 4 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
			2, 2,
			TopRight,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
				},
			},
		},
		{
			"CropAnchor 4x4 2x2 Left",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 3, 3),
				Stride: 4 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
			2, 2,
			Left,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
				},
			},
		},
		{
			"CropAnchor 4x4 2x2 Center",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 3, 3),
				Stride: 4 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
			2, 2,
			Center,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b,
					0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b,
				},
			},
		},
		{
			"CropAnchor 4x4 2x2 Right",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 3, 3),
				Stride: 4 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
			2, 2,
			Right,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
				},
			},
		},
		{
			"CropAnchor 4x4 2x2 BottomLeft",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 3, 3),
				Stride: 4 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
			2, 2,
			BottomLeft,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37,
				},
			},
		},
		{
			"CropAnchor 4x4 2x2 Bottom",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 3, 3),
				Stride: 4 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
			2, 2,
			Bottom,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b,
					0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b,
				},
			},
		},
		{
			"CropAnchor 4x4 2x2 BottomRight",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 3, 3),
				Stride: 4 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
			2, 2,
			BottomRight,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
		},
		{
			"CropAnchor 4x4 0x0 BottomRight",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 3, 3),
				Stride: 4 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
			0, 0,
			BottomRight,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 0, 0),
				Stride: 0,
				Pix:    []uint8{},
			},
		},
		{
			"CropAnchor 4x4 100x100 BottomRight",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 3, 3),
				Stride: 4 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
			100, 100,
			BottomRight,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 4, 4),
				Stride: 4 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
		},
		{
			"CropAnchor 4x4 1x100 BottomRight",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 3, 3),
				Stride: 4 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
			1, 100,
			BottomRight,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 1, 4),
				Stride: 1 * 4,
				Pix: []uint8{
					0x0c, 0x0d, 0x0e, 0x0f,
					0x1c, 0x1d, 0x1e, 0x1f,
					0x2c, 0x2d, 0x2e, 0x2f,
					0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
		},
		{
			"CropAnchor 4x4 0x100 BottomRight",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 3, 3),
				Stride: 4 * 4,
				Pix: []uint8{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				},
			},
			0, 100,
			BottomRight,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 0, 0),
				Stride: 0,
				Pix:    []uint8{},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := CropAnchor(tc.src, tc.w, tc.h, tc.anchor)
			if !compareNRGBA(got, tc.want, 0) {
				t.Fatalf("got result %#v want %#v", got, tc.want)
			}
		})
	}
}

func TestPaste(t *testing.T) {
	testCases := []struct {
		name string
		src1 image.Image
		src2 image.Image
		p    image.Point
		want *image.NRGBA
	}{
		{
			"Paste 2x3 2x1",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(1, 1, 3, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				},
			},
			image.Pt(-1, 0),
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 3),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
		},
		{
			"Paste 3x4 4x3 bottom right intersection",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 2, 3),
				Stride: 3 * 4,
				Pix: []uint8{
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b,
					0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(1, 1, 5, 4),
				Stride: 4 * 4,
				Pix: []uint8{
					0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf,
					0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf,
					0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xce, 0xcf,
				},
			},
			image.Pt(0, 1),
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 3, 4),
				Stride: 3 * 4,
				Pix: []uint8{
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b,
					0x30, 0x31, 0x32, 0x33, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7,
					0x40, 0x41, 0x42, 0x43, 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7,
				},
			},
		},
		{
			"Paste 3x4 4x3 top left intersection",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 2, 3),
				Stride: 3 * 4,
				Pix: []uint8{
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b,
					0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(1, 1, 5, 4),
				Stride: 4 * 4,
				Pix: []uint8{
					0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf,
					0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf,
					0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xce, 0xcf,
				},
			},
			image.Pt(-3, -2),
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 3, 4),
				Stride: 3 * 4,
				Pix: []uint8{
					0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf, 0x18, 0x19, 0x1a, 0x1b,
					0xc8, 0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xce, 0xcf, 0x28, 0x29, 0x2a, 0x2b,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b,
					0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b,
				},
			},
		},
		{
			"Paste 3x4 4x3 no intersection",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 2, 3),
				Stride: 3 * 4,
				Pix: []uint8{
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b,
					0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(1, 1, 5, 4),
				Stride: 4 * 4,
				Pix: []uint8{
					0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf,
					0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf,
					0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xce, 0xcf,
				},
			},
			image.Pt(-20, 20),
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 3, 4),
				Stride: 3 * 4,
				Pix: []uint8{
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b,
					0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b,
					0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b,
					0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Paste(tc.src1, tc.src2, tc.p)
			if !compareNRGBA(got, tc.want, 0) {
				t.Fatalf("got result %#v want %#v", got, tc.want)
			}
		})
	}
}

func BenchmarkPaste(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Paste(testdataBranchesJPG, testdataFlowersSmallPNG, image.Pt(100, 100))
	}
}

func TestPasteCenter(t *testing.T) {
	testCases := []struct {
		name string
		src1 image.Image
		src2 image.Image
		want *image.NRGBA
	}{
		{
			"PasteCenter 2x3 2x1",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0xff, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(1, 1, 3, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 3),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := PasteCenter(tc.src1, tc.src2)
			if !compareNRGBA(got, tc.want, 0) {
				t.Fatalf("got result %#v want %#v", got, tc.want)
			}
		})
	}
}

func TestOverlay(t *testing.T) {
	testCases := []struct {
		name string
		src1 image.Image
		src2 image.Image
		p    image.Point
		a    float64
		want *image.NRGBA
	}{
		{
			"Overlay 2x3 2x1 1.0",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0x60, 0x00, 0x90, 0xff, 0xff, 0x00, 0x99, 0x7f,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(1, 1, 3, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x20, 0x40, 0x80, 0x7f, 0xaa, 0xbb, 0xcc, 0xff,
				},
			},
			image.Pt(-1, 0),
			1.0,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 3),
				Stride: 2 * 4,
				Pix: []uint8{
					0x00, 0x11, 0x22, 0x33, 0xcc, 0xdd, 0xee, 0xff,
					0x40, 0x1f, 0x88, 0xff, 0xaa, 0xbb, 0xcc, 0xff,
					0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0xff,
				},
			},
		},
		{
			"Overlay 2x2 2x2 0.5",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 1),
				Stride: 2 * 4,
				Pix: []uint8{
					0xff, 0x00, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff,
					0x00, 0x00, 0xff, 0xff, 0x20, 0x20, 0x20, 0x00,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 1),
				Stride: 2 * 4,
				Pix: []uint8{
					0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
					0xff, 0xff, 0x00, 0xff, 0x20, 0x20, 0x20, 0xff,
				},
			},
			image.Pt(-1, -1),
			0.5,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0xff, 0x7f, 0x7f, 0xff, 0x00, 0xff, 0x00, 0xff,
					0x7f, 0x7f, 0x7f, 0xff, 0x20, 0x20, 0x20, 0x7f,
				},
			},
		},
		{
			"Overlay 2x2 2x2 0.5 no intersection",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 1),
				Stride: 2 * 4,
				Pix: []uint8{
					0xff, 0x00, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff,
					0x00, 0x00, 0xff, 0xff, 0x20, 0x20, 0x20, 0x00,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 1),
				Stride: 2 * 4,
				Pix: []uint8{
					0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00,
					0xff, 0xff, 0x00, 0xff, 0x20, 0x20, 0x20, 0xff,
				},
			},
			image.Pt(-10, 10),
			0.5,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0xff, 0x00, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff,
					0x00, 0x00, 0xff, 0xff, 0x20, 0x20, 0x20, 0x00,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Overlay(tc.src1, tc.src2, tc.p, tc.a)
			if !compareNRGBA(got, tc.want, 0) {
				t.Fatalf("got result %#v want %#v", got, tc.want)
			}
		})
	}
}

func BenchmarkOverlay(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Overlay(testdataBranchesJPG, testdataFlowersSmallPNG, image.Pt(100, 100), 0.5)
	}
}

func TestOverlayCenter(t *testing.T) {
	testCases := []struct {
		name string
		src1 image.Image
		src2 image.Image
		a    float64
		want *image.NRGBA
	}{
		{
			"OverlayCenter 2x3 2x1",
			&image.NRGBA{
				Rect:   image.Rect(-1, -1, 1, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x10, 0x10, 0x10, 0xff, 0x10, 0x10, 0x10, 0xff,
					0x10, 0x10, 0x10, 0xff, 0x10, 0x10, 0x10, 0xff,
					0x10, 0x10, 0x10, 0xff, 0x10, 0x10, 0x10, 0xff,
				},
			},
			&image.NRGBA{
				Rect:   image.Rect(1, 1, 3, 2),
				Stride: 2 * 4,
				Pix: []uint8{
					0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80,
				},
			},
			0.5,
			&image.NRGBA{
				Rect:   image.Rect(0, 0, 2, 3),
				Stride: 2 * 4,
				Pix: []uint8{
					0x10, 0x10, 0x10, 0xff, 0x10, 0x10, 0x10, 0xff,
					0x2c, 0x2c, 0x2c, 0xff, 0x2c, 0x2c, 0x2c, 0xff,
					0x10, 0x10, 0x10, 0xff, 0x10, 0x10, 0x10, 0xff,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := OverlayCenter(tc.src1, tc.src2, 0.5)
			if !compareNRGBA(got, tc.want, 0) {
				t.Fatalf("got result %#v want %#v", got, tc.want)
			}
		})
	}
}
