// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"bytes"
	"context"
	"image/png"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/image/bmp"
)

func screenshotAndDisplay(ctx context.Context, url string) error {
	var buf []byte
	// Capture screenshot of a page at a particular browser size
	// Note we can't use multicast DNS to use the nice name solar.local because
	// it doesn't work inside a docker container. So, using a hardcoded IP for the time
	// being. This IP at least is made to be static on the dhcp server (router)
	if err := chromedp.Run(ctx, fixedSizeScreenshot(url, &buf)); err != nil {
		return err
	}
	// Convert from png to bmp
	reader := bytes.NewReader(buf)
	image, err := png.Decode(reader)
	if err != nil {
		return err
	}
	// Create a file for writing
	file, err := os.Create("screenshot.bmp")
	if err != nil {
		return err
	}
	err = bmp.Encode(file, image)
	if err != nil {
		return err
	}
	file.Close()

	// Now run the external command to display the image on the eink screen
	cmd := exec.Command("IT8951/IT8951", "0", "0", "screenshot.bmp")
	return cmd.Run()
}

func fixedSizeScreenshot(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.EmulateViewport(800, 600),
		chromedp.Sleep(1 * time.Second),
		chromedp.CaptureScreenshot(res),
	}
}

func main() {
	url := os.Getenv("URL")
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := screenshotAndDisplay(ctx, url)
	if err != nil {
		log.Fatal(err)
	}
}
