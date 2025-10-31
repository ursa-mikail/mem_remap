package main

import (
    "fmt"
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
    // Example: Temporary patch for debugging or testing
    data := []byte{0x10, 0x20, 0x30, 0x40, 0x50}
    fmt.Printf("Original data: % x\n", data)
    
    // Create offset tracker
    tracker := &offsetMemoryAddressing{
        offsets: make(map[int64]byte),
    }
    
    // Apply temporary modifications
    tracker.modifyBytes(data, 1, 0x99)  // Change byte at offset 1
    tracker.modifyBytes(data, 3, 0x77)  // Change byte at offset 3
    
    fmt.Printf("After modifications: % x\n", data)
    fmt.Printf("Tracked offsets: %v\n", tracker.offsets)
    
    // Restore original values
    err := tracker.removeMemoryAddress(data)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("After restoration: % x\n", data)
    fmt.Printf("Remaining tracked offsets: %v\n", tracker.offsets)
}

/*
Original data: 10 20 30 40 50
After modifications: 10 99 30 77 50
Tracked offsets: map[1:32 3:64]

"""
map[1:32 3:64] means:

Offset 1 maps to byte value 32 (0x20 in hex)
Offset 3 maps to byte value 64 (0x40 in hex)
"""

After restoration: 10 20 30 40 50
Remaining tracked offsets: map[]
*/