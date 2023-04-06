package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gabriel-vasile/mimetype"
	"fmt"
	"io/ioutil"
	"image/jpeg"
	"image/png"
	"image/gif"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"github.com/Kagami/go-avif"
	"github.com/chai2010/webp"
	"log"
	"bytes"
	"image"
	"strconv"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	var supportedTypes = map[string]bool {
		"image/jpeg": true,
		"image/png": true,
		"image/gif": true,
		"image/bmp": true,
		"image/tiff": true,
		"image/webp": true,
	}

	app.Post("/upload", func(c *fiber.Ctx) error {
	    file, err := c.FormFile("file")
		mode := c.FormValue("mode")
		
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		intMode, err := strconv.Atoi(mode)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		buffer, err := file.Open()

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}		
		
		defer buffer.Close()

		data, err := ioutil.ReadAll(buffer)
		if err != nil {
			fmt.Println(err)
		}

		contentType := mimetype.Detect(data).String()

		if !supportedTypes[contentType] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   "not supported image format",
			})
		} 

		var img image.Image
		var buf bytes.Buffer

		switch contentType {
			case "image/jpeg":
				img, err = jpeg.Decode(bytes.NewReader(data))
			case "image/png":
				img, err = png.Decode(bytes.NewReader(data))
			case "image/gif":
				img, err = gif.Decode(bytes.NewReader(data))
			case "image/webp":
				img, err = webp.Decode(bytes.NewReader(data))
			case "image/bmp":
				img, err = bmp.Decode(bytes.NewReader(data))
			case "image/tiff":
				img, err = tiff.Decode(bytes.NewReader(data))
			default:
				fmt.Println("image error")
		}

		if err != nil {
			fmt.Println("error: ")
			fmt.Println(err.Error())
		}

		if intMode == 0 {
			if err = webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: 50}); err != nil {
				log.Println(err)
			}
		}

		if intMode == 1 {
			if err = avif.Encode(&buf, img, &avif.Options{Quality: 25}); err != nil {
				log.Println(err)
			}
		}

		c.Set("Content-Type", contentType)

		return c.Send(buf.Bytes())
	})

	app.Listen(":3000")
}
