// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"bytes"
	"context"
	"image/png"
	"log"
	"math"
	"os"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"golang.org/x/image/bmp"
)

func main() {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	// Capture screenshot of a page at a particular browser size
	if err := chromedp.Run(ctx, fixedSizeScreenshot(`http://solar.local/solar`, &buf)); err != nil {
		log.Fatal(err)
	}
	// if err := ioutil.WriteFile("screenshot.png", buf, 0644); err != nil {
	// 	log.Fatal(err)
	// }
	// Convert from png to bmp
	reader := bytes.NewReader(buf)
	image, err := png.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	// Create a file for writing
	file, err := os.Create("screenshot.bmp")
	if err != nil {
		log.Fatal(err)
	}
	err = bmp.Encode(file, image)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func fixedSizeScreenshot(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.EmulateViewport(800, 600),
		chromedp.Sleep(1 * time.Second),
		chromedp.CaptureScreenshot(res),
	}
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(sel, chromedp.ByID),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.ByID),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Liberally copied from puppeteer's source.
//
// Note: this will override the viewport emulation settings.
func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
