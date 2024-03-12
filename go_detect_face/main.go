package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	pigo "github.com/esimov/pigo/core"
)

func main() {
	// Load cascade file
	cascadeFile, err := os.ReadFile("C:/Users/User/go/pkg/mod/github.com/esimov/pigo@v1.4.6/cascade/facefinder")
	if err != nil {
		log.Fatalf("Error reading the cascade file: %v", err)
	}
	// Load image
	src, err := pigo.GetImage("C:/Users/User/Pictures/Screenshots/Screenshot 2023-12-02 202334.png")
	if err != nil {
		log.Fatalf("Cannot open the image file: %v", err)
	}
	// Convert image to *image.RGBA
	rgbaSrc := image.NewRGBA(src.Bounds())
	draw.Draw(rgbaSrc, src.Bounds(), src, image.Point{}, draw.Src)
	// Convert image to grayscale
	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y
	// Define cascade parameters
	cParams := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     1000,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,

		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}
	// Initialize pigo
	pigo := pigo.NewPigo()
	// Unpack the cascade file
	classifier, err := pigo.Unpack(cascadeFile)
	if err != nil {
		log.Fatalf("Error reading the cascade file: %s", err)
	}
	// cascade rotation angle. 0.0 is 0 radians and 1.0 is 2*pi radians
	angle := 0.0
	// Run the classifier over the obtained leaf nodes and return the detection results.
	dets := classifier.RunCascade(cParams, angle)
	// Calculate the intersection over union (IoU) of two clusters.
	dets = classifier.ClusterDetections(dets, 0.2)
	// Draw rectangles around detected faces
	outputs := drawRectangles(rgbaSrc, dets, src)
	// Save the image
	for k, output := range outputs {
		saveImage(output, fmt.Sprint(k)+".png")
	}
	// Encode the image to Base64
	encodedImage, err := encodeImageToBase64(outputs[2])
	if err != nil {
		log.Fatalf("Error encoding image to Base64: %v", err)
	}
	filePath := "D:/New folder/base64_encoded.txt"
	err = os.WriteFile(filePath, []byte(encodedImage), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func drawRectangles(img *image.RGBA, dets []pigo.Detection, src image.Image) []image.Image {
	// Define the color for drawing rectangles (red in this case)
	rectColor := color.RGBA{255, 0, 0, 255}
	// Define the thickness for the rectangle outlines
	outlineThickness := 1
	// Draw rectangles around detected faces
	outputs := []image.Image{}
	for _, det := range dets {
		x := det.Col - int(det.Scale/2) - 50
		y := det.Row - int(det.Scale/2) - 170
		w := int(det.Scale) + 100
		h := int(det.Scale) + 300
		// Draw the top horizontal line
		drawHorizontalLine(img, x, y, w, outlineThickness, rectColor)
		// Draw the bottom horizontal line
		drawHorizontalLine(img, x, y+h, w, outlineThickness, rectColor)
		// Draw the left vertical line
		drawVerticalLine(img, x, y, h, outlineThickness, rectColor)
		// Draw the right vertical line
		drawVerticalLine(img, x+w, y, h, outlineThickness, rectColor)
		//bounds := src.Bounds()
		//width := bounds.Dx()
		//height := bounds.Dy()
		cropSize := image.Rect(x, y, x+w, y+h)
		croppedImage := src.(SubImager).SubImage(cropSize)
		outputs = append(outputs, croppedImage)
	}
	return outputs
}

func drawHorizontalLine(img *image.RGBA, x, y, w, thickness int, color color.RGBA) {
	for i := 0; i < w; i++ {
		for j := 0; j < thickness; j++ {
			img.Set(x+i, y+j, color)
		}
	}
}

func drawVerticalLine(img *image.RGBA, x, y, h, thickness int, color color.RGBA) {
	for i := 0; i < h; i++ {
		for j := 0; j < thickness; j++ {
			img.Set(x+j, y+i, color)
		}
	}
}

func saveImage(img image.Image, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating image file: %v", err)
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		log.Fatalf("Error encoding image: %v", err)
	}
	fmt.Printf("Image with detected faces saved as %s\n", filename)
}

func encodeImageToBase64(img image.Image) (string, error) {
	// Create a buffer to store the Base64-encoded image
	var buf bytes.Buffer
	// Encode the image to Base64 and write it to the buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("error encoding image to PNG: %v", err)
	}
	// Encode the buffer to Base64 and return the encoded string
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}
