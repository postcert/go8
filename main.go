package main

import (
	"fmt"
	"io"
	"os"
)

type Chip8 struct {
	memory     [4096]uint8   // 4K memory
	V          [16]uint8     // V0-VF registers
	I          uint16        // Index register
	pc         uint16        // Program counter
	screen     [32][64]uint8 // Screen
	delayTimer uint8
	soundTimer uint8
	stack      []uint16
	keys       [16]uint8 // Hex keypad

	waitingForKeyPress bool
	keyRegister        uint8
}

func NewChip8() *Chip8 {
	chip := &Chip8{}
	chip.initialize()
	return chip
}

func (chip *Chip8) initialize() {
	chip.memory = [4096]uint8{}
	chip.V = [16]uint8{}
	chip.I = 0
	chip.pc = 0x200 // Set PC to starting address: 0x200
	chip.screen = [32][64]uint8{}
	chip.delayTimer = 0
	chip.soundTimer = 0
	chip.stack = []uint16{}
	chip.keys = [16]uint8{}

	chip.waitingForKeyPress = false
	chip.keyRegister = 0

	// Load fontset
	chip.loadFontset()
}

func (chip *Chip8) loadFontset() {
	fontset := []uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	copy(chip.memory[:], fontset)
}

func (chip *Chip8) loadRom(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening ROM file: %s", err)
	}
	defer file.Close()

	startAddress := 0x200
	bytesRead, err := io.ReadFull(file, chip.memory[startAddress:])
	if err != nil {
		return fmt.Errorf("error reading ROM file: %s", err)
	}

	fmt.Printf("Loaded %d bytes into memory\n", bytesRead)

	return nil
}

func main() {
	chip := NewChip8()
	err := chip.loadRom("roms/test.rom")
	if err != nil {
		fmt.Println("Failed to load ROM: ", err)
		os.Exit(1)
	}
}
