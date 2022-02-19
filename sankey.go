package main

import (
	"context"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/types"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	sankeyNode = []opts.SankeyNode{
		{Name: "Original Track - Potato"},
		{Name: "One Semitone up"},
		{Name: "Same Semitone"},
		{Name: "One Semitone down"},
		{Name: "song1"},
		{Name: "song2"},
		{Name: "song3"},
		{Name: "song4"},
		{Name: "song5"},
		{Name: "song6"},
		{Name: "song7"},
		{Name: "song8"},
		{Name: "song9"},
	}

	sankeyLink = []opts.SankeyLink{
		{Source: "Original Track - Potato", Target: "One Semitone up", Value: 8},
		{Source: "Original Track - Potato", Target: "Same Semitone", Value: 7},
		{Source: "Original Track - Potato", Target: "One Semitone down", Value: 6},
		{Source: "One Semitone up", Target: "song1", Value: 3},
		{Source: "One Semitone up", Target: "song2", Value: 2},
		{Source: "One Semitone up", Target: "song3", Value: 1},
		{Source: "Same Semitone", Target: "song4", Value: 3},
		{Source: "Same Semitone", Target: "song5", Value: 2},
		{Source: "Same Semitone", Target: "song6", Value: 1},
		{Source: "One Semitone down", Target: "song7", Value: 3},
		{Source: "One Semitone down", Target: "song8", Value: 2},
		{Source: "One Semitone down", Target: "song9", Value: 1},
	}
)

func sankeyBase() *charts.Sankey {
	sankey := charts.NewSankey()
	sankey.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:      "800px",
			Height:     "600px",
			Theme:      types.ThemeChalk,
			AssetsHost: "https://robyeates.github.io/go-echarts-assets/",
		}),
	)

	sankey.AddSeries("sankey", sankeyNode, sankeyLink,
		charts.WithLineStyleOpts(opts.LineStyle{
			Opacity: 0.8,
		}),
		charts.WithLabelOpts(opts.Label{
			Show:  true,
			Color: "white",
		}))
	return sankey
}

func Examples() {
	page := components.NewPage()
	page.AddCharts(
		sankeyBase(),
	)

	f, err := os.Create("sankey.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)

	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate(filepath.Join(mydir, f.Name())),
		chromedp.FullScreenshot(&buf, 100),
	); err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("fullScreenshotDark.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote fullScreenshot.jpeg")
	img, err := readImage("fullScreenshotDark.png")
	if err != nil {
		log.Fatal(err)
	}

	// I've hard-coded a crop rectangle, start (0,0), end (100, 100).
	img, err = cropImage(img, image.Rect(8, 8, 774, 608))
	if err != nil {
		log.Fatal(err)
	}

	writeImage(img, "pic-cropped.png")
}

func main() {
	Examples()
}

// readImage reads a image file from disk. We're assuming the file will be png
// format.
func readImage(name string) (image.Image, error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	// image.Decode requires that you import the right image package. We've
	// imported "image/png", so Decode will work for png files. If we needed to
	// decode jpeg files then we would need to import "image/jpeg".
	//
	// Ignored return value is image format name.
	img, _, err := image.Decode(fd)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// cropImage takes an image and crops it to the specified rectangle.
func cropImage(img image.Image, crop image.Rectangle) (image.Image, error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	// img is an Image interface. This checks if the underlying value has a
	// method called SubImage. If it does, then we can use SubImage to crop the
	// image.
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(crop), nil
}

// writeImage writes an Image back to the disk.
func writeImage(img image.Image, name string) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return png.Encode(fd, img)
}
