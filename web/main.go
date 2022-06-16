// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

// Image must be PNG 800x600
func displayImage(reader io.Reader) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "image.png")
	if err != nil {
		return err
	}
	_, err = io.Copy(part, reader)
	if err != nil {
		return err
	}
	writer.Close()
	// TODO: Make the URL below configurable
	response, err := http.Post("http://eink-web-api:8080/image", writer.FormDataContentType(), body)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("unexpected response from eink web api service: %s", string(b))
	}
	return err
}

func screenshotAndDisplay(ctx context.Context, url string) error {
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
	reader := bytes.NewReader(buf)
	return displayImage(reader)
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
	file, err := os.Open("finn.png")
	if err != nil {
		log.Fatal(err)
	}
	err = displayImage(file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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
