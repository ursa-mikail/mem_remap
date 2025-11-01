# image_scrambler

Automatically creates an images folder if it doesn't exist

Looks for images/sample.jpg as the input image

Scrambles the image using your memory addressing system

Saves the scrambled version as images/scrambled_sample.jpg

Recovers the original image using removeMemoryAddress

Saves the recovered version as images/recovered_sample.jpg

## How it works:
1. Image Loading: The program loads JPEG or PNG images and converts them to raw byte arrays.

2. Scrambling:
- Creates a random permutation of all byte indices
- Uses your offsetMemoryAddressing to track original byte values while applying the scrambled mapping
- Each byte is moved to a new position based on the scrambled indices

3. Recovery:
- Uses removeMemoryAddress to restore all bytes to their original positions
- The tracker map contains the information needed to reverse the scrambling

4. Image Saving: The scrambled and recovered images are saved as new files.

The scrambling effect will make the original image appear as random noise, but the recovery process will perfectly restore it using your memory addressing system.

# Pixel Shuffler

A Go program that scrambles and unscrambles image pixels by shuffling their positions while preserving pixel values.

## How It Works

- **Pixel Shuffling**: Creates a random permutation of pixel positions and moves pixels to new locations
- **Value Preservation**: Pixel colors (RGB values) remain exactly the same - only positions change
- **Reversible**: Uses a reverse mapping to perfectly restore the original image
- **Memory Tracking**: Implements byte-level tracking using `offsetMemoryAddressing` system

## Usage

1. Place your input image as `images/sample.jpg`
2. Run the program:
   ```bash
   go run main.go
   ```
   