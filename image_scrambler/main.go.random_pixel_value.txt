package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"time"
)

// offsetMemoryAddressing tracks byte changes at specific memory offsets
type offsetMemoryAddressing struct {
	offsets map[int64]byte
}

func (o *offsetMemoryAddressing) removeMemoryAddress(mem []byte) error {
	if len(mem) == 0 {
		return fmt.Errorf("mem block is empty")
	}

	for offset, originalByte := range o.offsets {
		if offset < 0 || offset >= int64(len(mem)) {
			return fmt.Errorf("offset %d is out of bounds", offset)
		}
		mem[offset] = originalByte
	}

	o.offsets = make(map[int64]byte)
	return nil
}

// Helper method to modify bytes and track changes
func (o *offsetMemoryAddressing) modifyBytes(mem []byte, offset int64, newValue byte) error {
	if offset < 0 || offset >= int64(len(mem)) {
		return fmt.Errorf("offset %d is out of bounds", offset)
	}

	// Store original byte if not already tracked
	if _, exists := o.offsets[offset]; !exists {
		o.offsets[offset] = mem[offset]
	}

	// Apply modification
	mem[offset] = newValue
	return nil
}

// loadImage loads an image from file
func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Determine image format based on extension
	if len(filename) > 4 {
		ext := filename[len(filename)-4:]
		switch ext {
		case ".jpg", ".jpeg":
			return jpeg.Decode(file)
		case ".png":
			return png.Decode(file)
		}
	}
	return nil, fmt.Errorf("unsupported image format: %s", filename)
}

// saveImage saves an image to file
func saveImage(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if len(filename) > 4 {
		ext := filename[len(filename)-4:]
		switch ext {
		case ".jpg", ".jpeg":
			return jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
		case ".png":
			return png.Encode(file, img)
		}
	}
	return fmt.Errorf("unsupported image format: %s", filename)
}

// imageToBytes converts image to byte array
func imageToBytes(img image.Image) []byte {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Each pixel has 4 bytes (RGBA)
	data := make([]byte, width*height*4)

	index := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			data[index] = byte(r >> 8)
			data[index+1] = byte(g >> 8)
			data[index+2] = byte(b >> 8)
			data[index+3] = byte(a >> 8)
			index += 4
		}
	}
	return data
}

// bytesToImage converts byte array back to image
func bytesToImage(data []byte, width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	index := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := data[index]
			g := data[index+1]
			b := data[index+2]
			a := data[index+3]
			img.Set(x, y, color.RGBA{r, g, b, a})
			index += 4
		}
	}
	return img
}

// scrambleImage scrambles the image using the offsetMemoryAddressing
func scrambleImage(img image.Image) (image.Image, *offsetMemoryAddressing, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Convert image to bytes
	data := imageToBytes(img)

	// Create tracker and scramble the data
	tracker := &offsetMemoryAddressing{
		offsets: make(map[int64]byte),
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a permutation of indices
	indices := make([]int64, len(data))
	for i := range indices {
		indices[i] = int64(i)
	}

	// Shuffle indices using Fisher-Yates algorithm
	for i := len(indices) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		indices[i], indices[j] = indices[j], indices[i]
	}

	// Apply scrambling by swapping bytes according to scrambled indices
	scrambledData := make([]byte, len(data))
	copy(scrambledData, data)

	for i := 0; i < len(data); i++ {
		targetIdx := indices[i]
		err := tracker.modifyBytes(scrambledData, targetIdx, data[i])
		if err != nil {
			return nil, nil, err
		}
	}

	// Convert scrambled bytes back to image
	scrambledImg := bytesToImage(scrambledData, width, height)

	return scrambledImg, tracker, nil
}

// recoverImage recovers the original image from scrambled data
func recoverImage(scrambledImg image.Image, tracker *offsetMemoryAddressing) (image.Image, error) {
	bounds := scrambledImg.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Convert scrambled image to bytes
	scrambledData := imageToBytes(scrambledImg)

	// Recover original data
	err := tracker.removeMemoryAddress(scrambledData)
	if err != nil {
		return nil, err
	}

	// Convert recovered bytes back to image
	recoveredImg := bytesToImage(scrambledData, width, height)

	return recoveredImg, nil
}

func main() {
	// Static image path - change this to your image path
	imagePath := "images/sample.jpg"

	// Create images directory if it doesn't exist
	if _, err := os.Stat("images"); os.IsNotExist(err) {
		err := os.Mkdir("images", 0755)
		if err != nil {
			fmt.Printf("Error creating images directory: %v\n", err)
			return
		}
		fmt.Println("Created 'images' directory")
		fmt.Printf("Please place your image file as '%s'\n", imagePath)
		return
	}

	// Check if image exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		fmt.Printf("Image file '%s' not found!\n", imagePath)
		fmt.Println("Please place an image file (JPEG or PNG) in the 'images' folder as 'sample.jpg'")
		return
	}

	// Load the image
	fmt.Printf("Loading image: %s\n", imagePath)
	img, err := loadImage(imagePath)
	if err != nil {
		fmt.Printf("Error loading image: %v\n", err)
		return
	}

	bounds := img.Bounds()
	fmt.Printf("Image loaded: %dx%d pixels\n", bounds.Dx(), bounds.Dy())

	// Scramble the image
	fmt.Println("\nScrambling image...")
	scrambledImg, tracker, err := scrambleImage(img)
	if err != nil {
		fmt.Printf("Error scrambling image: %v\n", err)
		return
	}

	// Save scrambled image
	scrambledFile := "images/scrambled_sample.jpg"
	err = saveImage(scrambledFile, scrambledImg)
	if err != nil {
		fmt.Printf("Error saving scrambled image: %v\n", err)
		return
	}
	fmt.Printf("Scrambled image saved as: %s\n", scrambledFile)
	fmt.Printf("Number of modified bytes tracked: %d\n", len(tracker.offsets))

	// Recover the image
	fmt.Println("\nRecovering image...")
	recoveredImg, err := recoverImage(scrambledImg, tracker)
	if err != nil {
		fmt.Printf("Error recovering image: %v\n", err)
		return
	}

	// Save recovered image
	recoveredFile := "images/recovered_sample.jpg"
	err = saveImage(recoveredFile, recoveredImg)
	if err != nil {
		fmt.Printf("Error saving recovered image: %v\n", err)
		return
	}
	fmt.Printf("Recovered image saved as: %s\n", recoveredFile)

	fmt.Println("\nProcess completed successfully!")
	fmt.Println("Check the 'images' folder for:")
	fmt.Println("  - scrambled_sample.jpg (scrambled image)")
	fmt.Println("  - recovered_sample.jpg (recovered image)")
}
