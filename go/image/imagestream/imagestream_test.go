package imagestream_test

import (
	"image"
	"log"
	"testing"

	"golang.org/x/net/context"

	"github.com/harukasan/testing/go/image/imagestream"
)

func makeTestRGBA() *image.RGBA {
	s := image.NewRGBA(image.Rect(0, 0, 5, 5))
	for i := 0; i < 25; i++ {
		s.Pix[i*4+0] = uint8(i * 10)
		s.Pix[i*4+1] = uint8(i * 10)
		s.Pix[i*4+2] = uint8(i * 10)
		s.Pix[i*4+3] = 0
	}

	return s
}

func TestRGBAtoRGBA(t *testing.T) {
	s := makeTestRGBA()
	srcFilter := imagestream.NewSourceFilter(imagestream.NewImageSource(s))
	dest, err := srcFilter.Filter(nil)
	if err != nil {
		t.Fatalf("got error on %v", err)
	}
	d, err := imagestream.WriteImage(context.Background(), dest)
	if err != nil {
		t.Fatalf("got error on %v", err)
	}
	log.Println(s)
	log.Println(d)
}
