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

	// Now run the external command to display the image on the eink screen
	cmd := exec.Command("../IT8951/IT8951", "0", "0", "screenshot.bmp")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func fixedSizeScreenshot(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.EmulateViewport(800, 600),
		chromedp.Sleep(1 * time.Second),
		chromedp.CaptureScreenshot(res),
	}
}
