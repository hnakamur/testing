package imagestream

import (
	"fmt"
	"image"
	"io"
	"net/http"

	"golang.org/x/net/context"
)

// PixelFormat defines the data structure of the image.
type PixelFormat int

// Pixel formats definitions:
const (
	YCbCr444 PixelFormat = iota // planer YCbCr 4:4:4, 24bpp
	YCbCr422                    // planer YCbCr 4:2:2, 16bpp
	YCbCr420                    // planer YCbCr 4:2:0, 12bpp
	YCbCr440                    // planer YCbCr 4:4:0
	YCbCr411                    // planer YCbCr 4:1:1, 12bpp
	YCbCr410                    // planer YCbCr 4:1:0, 9bpp
	RGB24                       // packed RGB 8:8:8, 24bpp
	RGBA32                      // packed RGBA 8:8:8:, 32bpp
	Gray8                       // 8bpp grayscale
)

// Config holds configurations of the image planes.
type Config struct {
	// width and height (in number of elements, i.e., pixels) of each plane.
	PlaneWidth, PlaneHeight []int
	// Stride holds strides of each plane. It defines number of bytes between
	// adjacent scan lines.
	Stride []int
	// PixelFormat is pixel format of the image.
	PixelFormat PixelFormat
}

// ScanLineReader is the interface to read scan lines of the image planes.
type ScanLineReader interface {
	// Config returns the configurations of the image planes. It holds the pixel
	// format of the image and structural parameters of each plane.
	Config() *Config

	// ReadScanLines reads scan lines into the p for each image component. it
	// reads up to len(p[0]) at once and returns number of bytes read and any
	// error. When the image is YCbCr 4:2:0,
	//   ReadScanLines(p) returns
	//      p[0]: read scan lines of the Y component.
	//      p[1]: read scan lines of the Cb component.
	//      p[2]: read scan lines of the Cr component.
	//
	// ReadScanLines returns an error in the same manner of io.Reader.
	// When ReadScanLines reads to the end of image, it returns io.EOF.
	ReadScanLines(p [][]byte) (n int, err error)
}

// Filter is the interface of filter to reading scan lines of the image planes.
// It reads from src ScanLineReader, filter the data, and writes to dest
// ScanLineReader. If the filter can not support src ScanLineReader, it returns
// an error.
type Filter interface {
	Filter(ctx context.Context, src ScanLineReader) (dest ScanLineReader, err error)
}

// FilterFunc is a function that implements Filter interface.
type FilterFunc func(ctx context.Context, src ScanLineReader) (dest ScanLineReader, err error)

// Filter calls its own implementation.
func (f FilterFunc) Filter(ctx context.Context, src ScanLineReader) (dest ScanLineReader, err error) {
	dest, err = f(ctx, src)
	return
}

// Source is the interface to read image from the any source.
//
// Open returns the reader which reads scan lines from the image source. If it
// failed to open the any resource or read the data, return an error.
type Source interface {
	Open(ctx context.Context) (image ScanLineReader, err error)
}

// SourceFilter implements filter interface. It can use the first of filter
// chain as a filter.
type SourceFilter struct {
	src Source
}

// NewSourceFilter returns a filter to read scan lines from the source s.
func NewSourceFilter(s Source) *SourceFilter {
	return &SourceFilter{s}
}

// Filter returns a reader to read scan lines from the source. Note that it
// do not read from the src reader, the reader is just ignored.
func (f *SourceFilter) Filter(ctx context.Context, _ ScanLineReader) (dest ScanLineReader, err error) {
	dest, err = f.src.Open(ctx)
	return
}

// HTTPEncoder is the interface to write encoded image to HTTP response.
//
// HTTPEncode reads from src and writes encoded image to w. If it does not
// finish writing image, it returns an error.
type HTTPEncoder interface {
	HTTPEncode(ctx context.Context, w http.ResponseWriter, src ScanLineReader) error
}

// FileEncoder is the interface to write encoded image into the file.
//
// Write reads image from src and writes encoded to w.
type FileEncoder interface {
	Write(ctx context.Context, w io.Writer, src ScanLineReader) error
}

// FilterChain is a filter which can chain multiple filters.
type FilterChain struct {
	chain FilterFunc
}

// Append appends the filter f into end of the filter chain.
func (c *FilterChain) Append(f Filter) {
	c.chain = func(ctx context.Context, src ScanLineReader) (dest ScanLineReader, err error) {
		chainDest, err := c.chain(ctx, src)
		if err != nil {
			return chainDest, err
		}
		dest, err = f.Filter(ctx, chainDest)
		return
	}
}

func NewImageSource(src image.Image) Source {
	switch s := src.(type) {
	case *image.RGBA:
		return &imageRGBASource{s}
	default:
		return nil
	}
}

type imageRGBASource struct {
	img *image.RGBA
}

func (s *imageRGBASource) Open(ctx context.Context) (image ScanLineReader, err error) {
	return &imageRGBAReader{s.img, 0}, nil
}

type imageRGBAReader struct {
	img *image.RGBA
	off int
}

func (r *imageRGBAReader) Config() *Config {
	w, h := r.img.Bounds().Dx(), r.img.Bounds().Dy()
	return &Config{
		PlaneWidth:  []int{w},
		PlaneHeight: []int{h},
		Stride:      []int{r.img.Stride},
		PixelFormat: RGBA32,
	}
}

func (r *imageRGBAReader) ReadScanLines(p [][]byte) (n int, err error) {
	if r.off >= len(r.img.Pix) {
		if len(p) == 0 {
			return
		}
		return 0, io.EOF
	}
	n = copy(p[0], r.img.Pix[r.off:])
	r.off += n
	return
}

type imageYCbCrReader struct {
	img *image.YCbCr
	off int
}

func (r *imageYCbCrReader) Config() *Config {
	w, h := r.img.Bounds().Dx(), r.img.Bounds().Dy()
	return &Config{
		PlaneWidth:  []int{w},
		PlaneHeight: []int{h},
		Stride:      []int{r.img.YStride, r.img.CStride, r.img.CStride},
		PixelFormat: RGBA32,
	}
}

// WriteImage reads from the src and returns an image.
func WriteImage(src ScanLineReader) (image.Image, error) {
	config := src.Config()
	switch config.PixelFormat {
	case RGBA32:
		return encodeToRGBA(src)
	default:
		return nil, fmt.Errorf("unsupported pixel format: %v", config.PixelFormat)
	}
}

func encodeToRGBA(src ScanLineReader) (*image.RGBA, error) {
	w, h := src.Config().PlaneWidth[0], src.Config().PlaneHeight[0]
	dest := image.NewRGBA(image.Rect(0, 0, w, h))
	off := 0
	for {
		n, err := src.ReadScanLines([][]byte{dest.Pix[off:]})
		if err != nil {
			if err == io.EOF {
				break
			}
			return dest, err
		}
		off += n
	}
	return dest, nil
}
