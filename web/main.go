// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"bytes"
	"context"
	"fmt"
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

func screenshotAndDisplay(ctx context.Context, url string) error {
	// First for debugging purposes get the first bit of text from the URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	fmt.Printf("%s", body)

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
		// Hacky hack - wait for two seconds for the animation to finish
		chromedp.Sleep(2 * time.Second),
		chromedp.CaptureScreenshot(res),
	}
}

func main() {
	url := os.Getenv("URL")
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	for {
		err := screenshotAndDisplay(ctx, url)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Sleeping for thirty seconds...")
		time.Sleep(30 * time.Second)
	}
}
