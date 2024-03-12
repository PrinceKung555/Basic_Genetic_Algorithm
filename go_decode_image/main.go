package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	// Read the Base64 encoded string from the file
	base64String, err := readBase64FromFile("base64_encoded.txt")
	if err != nil {
		fmt.Println("Error reading Base64 from file:", err)
		return
	}
	// Decode the base64 string
	decoded, err := base64.StdEncoding.DecodeString(strings.Replace(base64String, "data:image/jpeg;base64,", "", -1))
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		return
	}

	// Decode the image data
	img, _, err := image.Decode(strings.NewReader(string(decoded)))
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// Create a file to save the decoded image
	outFile, err := os.Create("decoded_image.png")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	// Encode the decoded image data to PNG and save it to the file
	err = png.Encode(outFile, img)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}

	fmt.Println("Image successfully decoded and saved as decoded_image.png")
}

func readBase64FromFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
