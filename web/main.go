// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"bytes"
	"context"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/image/bmp"
)

// path is relative to the eink-driver directory. Must be
// an 800x600 bmp image
func displayBmp(path string) error {
	// If the command doesn't complete in 10 seconds it will get killed
	// This is a workaround for seeing the IT8951 occasionaly hanging
	// using 100% cpu. Obviously it would be good to figure out why this
	// is actually happening
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "IT8951/IT8951", "0", "0", path)
	return cmd.Run()
}

func screenshot(ctx context.Context, url string, path string) error {
	// First for debugging purposes get the first bit of text from the URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	var buf []byte
	// Capture screenshot of a page at a particular browser size
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
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	err = bmp.Encode(file, image)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func screenshotAndDisplay(ctx context.Context, url string) error {
	err := screenshot(ctx, url, "screenshot.bmp")
	if err != nil {
		return err
	}
	return displayBmp("screenshot.bmp")
}

func fixedSizeScreenshot(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.EmulateViewport(800, 600),
		// Hacky hack - wait for two seconds for the animation to finish
		chromedp.Sleep(2 * time.Second),
		chromedp.CaptureScreenshot(res),
	}
}

func main() {
	// First things first. Show a picture of Finn to show
	// that we're starting up.
	log.Println("Showing startup image...")
	err := displayBmp("finn.bmp")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Sleeping for thirty seconds...")
	time.Sleep(30 * time.Second)

	url := os.Getenv("URL")
	tz := os.Getenv("TZ")

	// Set timezone of the browser
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Env("TZ="+tz),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create context
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	for {
		err = screenshotAndDisplay(ctx, url)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Sleeping for thirty seconds...")
		time.Sleep(30 * time.Second)
	}
}
