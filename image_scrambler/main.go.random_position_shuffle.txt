package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"math/rand"
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

func main() {
	// Load sample image from /images directory
	imagePath := filepath.Join("images", "sample.jpg")
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Printf("Error loading image: %v\n", err)
		fmt.Printf("Please make sure %s exists\n", imagePath)
		return
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Printf("Error decoding image: %v\n", err)
		return
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	totalPixels := width * height
	
	fmt.Printf("Original image: %dx%d pixels, %d total pixels\n", width, height, totalPixels)
	
	// Convert image to RGBA for easier pixel manipulation
	originalImg := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			originalImg.Set(x, y, img.At(x, y))
		}
	}

	// Create a permutation (shuffle) of all pixel indices
	rand.Seed(time.Now().UnixNano())
	shuffledIndices := make([]int, totalPixels)
	for i := 0; i < totalPixels; i++ {
		shuffledIndices[i] = i
	}
	
	// Shuffle the indices using Fisher-Yates algorithm
	for i := totalPixels - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffledIndices[i], shuffledIndices[j] = shuffledIndices[j], shuffledIndices[i]
	}

	fmt.Printf("Created shuffled indices for %d pixels\n", len(shuffledIndices))

	// Show some shuffle examples
	fmt.Println("\nShuffle examples:")
	for i := 0; i < 5; i++ {
		origX, origY := i%width, i/width
		shuffledIdx := shuffledIndices[i]
		shuffledX, shuffledY := shuffledIdx%width, shuffledIdx/width
		
		// Get original pixel values
		origR, origG, origB, _ := originalImg.At(origX, origY).RGBA()
		fmt.Printf("  Pixel at (%d,%d) [R:%d G:%d B:%d] -> moves to (%d,%d)\n", 
			origX, origY, origR>>8, origG>>8, origB>>8, shuffledX, shuffledY)
	}

	// Create scrambled image - start with empty image
	scrambledImg := image.NewRGBA(bounds)
	tracker := &offsetMemoryAddressing{
		offsets: make(map[int64]byte),
	}

	fmt.Println("\nApplying pixel shuffle...")
	bytesPerPixel := 4 // RGBA

	// Apply the shuffle: copy each pixel from original position to shuffled position
	for originalIdx, shuffledIdx := range shuffledIndices {
		// Calculate coordinates
		origX := originalIdx % width
		origY := originalIdx / width
		shuffledX := shuffledIdx % width
		shuffledY := shuffledIdx / width

		// Get the pixel value from ORIGINAL image
		r, g, b, a := originalImg.At(origX, origY).RGBA()
		
		// Calculate byte offsets in scrambled image
		scrambledOffset := (shuffledY*width + shuffledX) * bytesPerPixel
		
		// Store the ORIGINAL pixel value at the NEW position
		// Since scrambledImg starts empty, we're setting pixels for the first time
		tracker.modifyBytes(scrambledImg.Pix, int64(scrambledOffset), byte(r>>8))
		tracker.modifyBytes(scrambledImg.Pix, int64(scrambledOffset+1), byte(g>>8))
		tracker.modifyBytes(scrambledImg.Pix, int64(scrambledOffset+2), byte(b>>8))
		tracker.modifyBytes(scrambledImg.Pix, int64(scrambledOffset+3), byte(a>>8))
	}

	fmt.Printf("Tracked %d byte modifications during shuffling\n", len(tracker.offsets))

	// Save scrambled image
	scrambledPath := filepath.Join("images", "scrambled.jpg")
	scrambledFile, err := os.Create(scrambledPath)
	if err != nil {
		fmt.Printf("Error creating scrambled image: %v\n", err)
		return
	}
	defer scrambledFile.Close()
	
	jpeg.Encode(scrambledFile, scrambledImg, &jpeg.Options{Quality: 90})
	fmt.Printf("Saved shuffled image as: %s\n", scrambledPath)

	// Verify that pixel VALUES are preserved in shuffled image
	fmt.Println("\nVerifying pixel values after shuffling:")
	for i := 0; i < 3; i++ {
		origX, origY := i%width, i/width
		shuffledIdx := shuffledIndices[i]
		shuffledX, shuffledY := shuffledIdx%width, shuffledIdx/width
		
		// Get original pixel value
		origR, origG, origB, _ := originalImg.At(origX, origY).RGBA()
		// Get the value at the shuffled position in scrambled image
		shuffledR, shuffledG, shuffledB, _ := scrambledImg.At(shuffledX, shuffledY).RGBA()
		
		valuesMatch := origR == shuffledR && origG == shuffledG && origB == shuffledB
		fmt.Printf("  Pixel from (%d,%d) [R:%d] moved to (%d,%d) [R:%d] - Values match: %t\n", 
			origX, origY, origR>>8, shuffledX, shuffledY, shuffledR>>8, valuesMatch)
	}

	// Now restore using removeMemoryAddress
	fmt.Println("\nRestoring using removeMemoryAddress...")
	
	// Create a copy of scrambled image for restoration
	scrambledCopy := image.NewRGBA(bounds)
	copy(scrambledCopy.Pix, scrambledImg.Pix)
	
	// This should restore the scrambled image to its original empty state
	// But wait - this won't work because we want to restore the ORIGINAL image, not empty!
	// Let's think differently...
	
	// Instead, let's create a new image and use the reverse mapping
	fmt.Println("\nCreating reverse shuffle mapping...")
	reverseShuffle := make([]int, totalPixels)
	for originalIdx, shuffledIdx := range shuffledIndices {
		reverseShuffle[shuffledIdx] = originalIdx
	}

	// Create restored image by applying reverse shuffle
	restoredImg := image.NewRGBA(bounds)
	tracker2 := &offsetMemoryAddressing{
		offsets: make(map[int64]byte),
	}

	fmt.Println("Applying reverse shuffle...")
	for shuffledIdx, originalIdx := range reverseShuffle {
		// Calculate coordinates
		shuffledX := shuffledIdx % width
		shuffledY := shuffledIdx / width
		origX := originalIdx % width
		origY := originalIdx / width

		// Get pixel from scrambled image
		r, g, b, a := scrambledImg.At(shuffledX, shuffledY).RGBA()
		
		// Calculate byte offsets in restored image
		restoredOffset := (origY*width + origX) * bytesPerPixel
		
		// Store pixel back at original position
		tracker2.modifyBytes(restoredImg.Pix, int64(restoredOffset), byte(r>>8))
		tracker2.modifyBytes(restoredImg.Pix, int64(restoredOffset+1), byte(g>>8))
		tracker2.modifyBytes(restoredImg.Pix, int64(restoredOffset+2), byte(b>>8))
		tracker2.modifyBytes(restoredImg.Pix, int64(restoredOffset+3), byte(a>>8))
	}

	// Save restored image
	restoredPath := filepath.Join("images", "restored.jpg")
	restoredFile, err := os.Create(restoredPath)
	if err != nil {
		fmt.Printf("Error creating restored image: %v\n", err)
		return
	}
	defer restoredFile.Close()
	
	jpeg.Encode(restoredFile, restoredImg, &jpeg.Options{Quality: 90})
	fmt.Printf("Saved restored image as: %s\n", restoredPath)

	// Verify restoration by comparing with original
	fmt.Println("\nVerifying restoration...")
	matches := 0
	for i := 0; i < totalPixels && i < 100; i++ {
		x := i % width
		y := i / width
		
		origR, origG, origB, origA := originalImg.At(x, y).RGBA()
		restoredR, restoredG, restoredB, restoredA := restoredImg.At(x, y).RGBA()
		
		if origR == restoredR && origG == restoredG && origB == restoredB && origA == restoredA {
			matches++
		}
	}
	
	fmt.Printf("Pixel matching: %d/100 pixels match original\n", matches)
	
	if matches == 100 {
		fmt.Println("✓ Perfect restoration achieved!")
	} else {
		fmt.Println("⚠ Some pixels don't match")
	}

	// Demonstrate removeMemoryAddress on the scrambled image
	fmt.Println("\n=== Demonstrating removeMemoryAddress on scrambled image ===")
	err = tracker.removeMemoryAddress(scrambledImg.Pix)
	if err != nil {
		fmt.Printf("Error during removeMemoryAddress: %v\n", err)
		return
	}
	fmt.Printf("After removeMemoryAddress - scrambled image bytes were reset\n")
	fmt.Printf("Remaining tracked offsets: %d\n", len(tracker.offsets))
}