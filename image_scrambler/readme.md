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
