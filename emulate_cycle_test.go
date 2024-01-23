package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Reenable test suite functionality when neotest supports it
// https://github.com/nvim-neotest/neotest-go/pull/67

// type EmulateCycleTestSuite struct {
// 	suite.Suite
// 	chip *Chip8
// }

// func (suite *EmulateCycleTestSuite) SetupTest() {
//   suite.chip = NewChip8()
// }

func TestClearDisplay(t *testing.T) {
	chip := NewChip8()
	chip.screen[0][0] = 1

	chip.clearDisplay()

	assert.Equal(t, uint8(0), chip.screen[0][0])
}

func TestEmulateCycleClearDisplay(t *testing.T) {
	chip := NewChip8()
	chip.screen[0][0] = 1

	chip.emulateCycle(0x00E0)

	assert.Equal(t, uint8(0), chip.screen[0][0])
}

func TestReturnFromSubroutine(t *testing.T) {
	chip := NewChip8()
	chip.stack = []uint16{0x200, 0x300}
	chip.pc = 0x400
	chip.returnFromSubroutine()
	assert.Equal(t, uint16(0x300), chip.pc)
	assert.Equal(t, []uint16{0x200}, chip.stack)
}

func TestEmulateCycleReturnFromSubroutine(t *testing.T) {
	chip := NewChip8()
	chip.stack = []uint16{0x200, 0x300}
	chip.pc = 0x400
	chip.emulateCycle(0x00EE)
	assert.Equal(t, uint16(0x300), chip.pc)
	assert.Equal(t, []uint16{0x200}, chip.stack)
}

func TestJumpToAddress(t *testing.T) {
	chip := NewChip8()
	chip.jumpToAddress(0x300)
	assert.Equal(t, uint16(0x300), chip.pc)
}

func TestEmulateCycleJumpToAddress(t *testing.T) {
	chip := NewChip8()
	chip.emulateCycle(0x1300)
	assert.Equal(t, uint16(0x300), chip.pc)
}

func TestCallSubroutine(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.callSubroutine(0x300)
	assert.Equal(t, uint16(0x300), chip.pc)
	assert.Equal(t, []uint16{0x200}, chip.stack)
}

func TestEmulateCycleCallSubroutine(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.emulateCycle(0x2300)
	assert.Equal(t, uint16(0x300), chip.pc)
	assert.Equal(t, []uint16{0x200}, chip.stack)
}

func TestSkipIfEqual_Success(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.V[0] = 0x1
	chip.skipIfEqual(0x000, 0x1)
	assert.Equal(t, uint16(0x202), chip.pc)
}

func TestEmulateCycleSkipIfEqual_Success(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.V[0] = 0x1
	chip.emulateCycle(0x3001)
	assert.Equal(t, uint16(0x202), chip.pc)
}

func TestSkipIfEqual_Failure(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.V[0] = 0x1
	chip.skipIfEqual(0x000, 0x2)
	assert.Equal(t, uint16(0x200), chip.pc)
}

func TestSkipIfNotEqual_Success(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.V[0] = 0x1
	chip.skipIfNotEqual(0x000, 0x2)
	assert.Equal(t, uint16(0x202), chip.pc)
}

func TestEmulateCycleSkipIfNotEqual_Success(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.V[0] = 0x1
	chip.emulateCycle(0x4002)
	assert.Equal(t, uint16(0x202), chip.pc)
}

func TestSkipIfNotEqual_Failure(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.V[0] = 0x1
	chip.skipIfNotEqual(0x000, 0x1)
	assert.Equal(t, uint16(0x200), chip.pc)
}

func TestSkipIfRegistersEqual_Success(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.V[0] = 0x1
	chip.V[1] = 0x1
	chip.skipIfRegistersEqual(0x0, 0x10)
	assert.Equal(t, uint16(0x202), chip.pc)
}

func TestEmulateCycleSkipIfRegistersEqual_Success(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.V[0] = 0x1
	chip.V[1] = 0x1
	chip.emulateCycle(0x5010)
	assert.Equal(t, uint16(0x202), chip.pc)
}

func TestSkipIfRegistersEqual_Failure(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.V[0] = 0x1
	chip.V[1] = 0x2
	chip.skipIfRegistersEqual(0x000, 0x10)
	assert.Equal(t, uint16(0x200), chip.pc)
}

func TestSetRegister(t *testing.T) {
	chip := NewChip8()
	chip.setRegister(0x0, 0x1)
	assert.Equal(t, uint8(0x1), chip.V[0])
}

func TestEmulateCycleSetRegister(t *testing.T) {
	chip := NewChip8()
	chip.emulateCycle(0x6001)
	assert.Equal(t, uint8(0x1), chip.V[0])
}

func TestAddToRegister(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.addtoRegister(0x0, 0x1)
	assert.Equal(t, uint8(0x2), chip.V[0])
}

func TestEmulateCycleAddToRegister(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.emulateCycle(0x7001)
	assert.Equal(t, uint8(0x2), chip.V[0])
}

func TestStoreRegister(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0x2
	chip.storeRegister(0x0, 0x10)
	assert.Equal(t, uint8(0x2), chip.V[0])
}

func TestEmulateCycleStoreRegister(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0x2
	chip.storeRegister(0x0, 0x10)
	assert.Equal(t, uint8(0x2), chip.V[0])
}

func TestOrRegisters(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0x2
	chip.orRegisters(0x0, 0x10)
	assert.Equal(t, uint8(0x3), chip.V[0])
}

func TestEmulateCycleOrRegisters(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0x2
	chip.emulateCycle(0x8011)
	assert.Equal(t, uint8(0x3), chip.V[0])
}

func TestAndRegisters(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0x2
	chip.andRegisters(0x0, 0x10)
	assert.Equal(t, uint8(0x0), chip.V[0])
}

func TestEmulateCycleAndRegisters(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0x2
	chip.emulateCycle(0x8012)
	assert.Equal(t, uint8(0x0), chip.V[0])
}

func TestXorRegisters(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0x3
	chip.xorRegisters(0x0, 0x10)
	assert.Equal(t, uint8(0x2), chip.V[0])
}

func TestEmulateCycleXorRegisters(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0x3
	chip.emulateCycle(0x8013)
	assert.Equal(t, uint8(0x2), chip.V[0])
}

func TestAddRegistersWithCarry_Success(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0xFF
	chip.addRegistersWithCarry(0x10, 0x10)
	assert.Equal(t, uint8(0x00), chip.V[0])
	assert.Equal(t, uint8(0x1), chip.V[0xF])
}

func TestEmulateCycleAddRegistersWithCarry_Success(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0xFE
	chip.emulateCycle(0x8014)
	assert.Equal(t, uint8(0xFF), chip.V[0])
	assert.Equal(t, uint8(0x0), chip.V[0xF])
}

func TestSubtractRegistersWithBorrowVxVy_NoBorrow(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x2
	chip.V[1] = 0x1
	chip.subtractRegistersWithBorrowVxVy(0x0, 0x10)
	assert.Equal(t, uint8(0x1), chip.V[0])
	assert.Equal(t, uint8(0x1), chip.V[0xF])
}

func TestEmulateCycleSubtractRegistersWithBorrowVxVy_NoBorrow(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x2
	chip.V[1] = 0x1
	chip.emulateCycle(0x8015)
	assert.Equal(t, uint8(0x1), chip.V[0])
	assert.Equal(t, uint8(0x1), chip.V[0xF])
}

func TestSubtractRegistersWithBorrowVxVy_Borrow(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0x2
	chip.subtractRegistersWithBorrowVxVy(0x0, 0x10)
	assert.Equal(t, uint8(0xFF), chip.V[0])
	assert.Equal(t, uint8(0x0), chip.V[0xF])
}

func TestEmulateCycleSubtractRegistersWithBorrowVxVy_Borrow(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0x2
	chip.emulateCycle(0x8015)
	assert.Equal(t, uint8(0xFF), chip.V[0])
	assert.Equal(t, uint8(0x0), chip.V[0xF])
}

func TestShiftRight(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x2
	chip.V[1] = 0x2
	chip.shiftRight(0x0, 0x10)
	assert.Equal(t, uint8(0x1), chip.V[0])
}

func TestEmulateCycleShiftRight(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x2
	chip.V[1] = 0x2
	chip.emulateCycle(0x8016)
	assert.Equal(t, uint8(0x1), chip.V[0])
}

func TestSubtractRegistersWithBorrowVyVx_NoBorrow(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0x2
	chip.subtractRegistersWithBorrowVyVx(0x0, 0x10)
	assert.Equal(t, uint8(0x1), chip.V[1])
	assert.Equal(t, uint8(0x1), chip.V[0xF])
}

func TestEmulateCycleSubtractRegistersWithBorrowVyVx_NoBorrow(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0x2
	chip.emulateCycle(0x8017)
	assert.Equal(t, uint8(0x1), chip.V[1])
	assert.Equal(t, uint8(0x1), chip.V[0xF])
}

func TestSubtractRegistersWithBorrowVyVx_Borrow(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x2
	chip.V[1] = 0x1
	chip.subtractRegistersWithBorrowVyVx(0x0, 0x10)
	assert.Equal(t, uint8(0xFF), chip.V[1])
	assert.Equal(t, uint8(0x0), chip.V[0xF])
}

func TestEmulateCycleSubtractRegistersWithBorrowVyVx_Borrow(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x2
	chip.V[1] = 0x1
	chip.emulateCycle(0x8017)
	assert.Equal(t, uint8(0xFF), chip.V[1])
	assert.Equal(t, uint8(0x0), chip.V[0xF])
}

func TestShiftLeft(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0x1
	chip.shiftLeft(0x0, 0x10)
	assert.Equal(t, uint8(0x2), chip.V[0])
}

func TestEmulateCycleShiftLeft(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.V[1] = 0x1
	chip.emulateCycle(0x801E)
	assert.Equal(t, uint8(0x2), chip.V[0])
}

func TestSkipIfRegistersNotEqual_Success(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.V[0] = 0x1
	chip.skipIfRegistersNotEqual(0x0, 0x20)
	assert.Equal(t, uint16(0x202), chip.pc)
}

func TestEmulateCycleSkipIfRegistersNotEqual_Success(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.V[0] = 0x1
	chip.emulateCycle(0x9020)
	assert.Equal(t, uint16(0x202), chip.pc)
}

func TestSkipIfRegistersNotEqual_Failure(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.V[0] = 0x1
	chip.V[1] = 0x1
	chip.skipIfRegistersNotEqual(0x0, 0x10)
	assert.Equal(t, uint16(0x200), chip.pc)
}

func TestEmulateCycleSkipIfRegistersNotEqual_Failure(t *testing.T) {
	chip := NewChip8()
	chip.pc = 0x200
	chip.V[0] = 0x1
	chip.V[1] = 0x1
	chip.emulateCycle(0x9010)
	assert.Equal(t, uint16(0x200), chip.pc)
}

func TestSetIndexRegister(t *testing.T) {
	chip := NewChip8()
	chip.setIndexRegister(0x200)
	assert.Equal(t, uint16(0x200), chip.I)
}

func TestEmulateCycleSetIndexRegister(t *testing.T) {
	chip := NewChip8()
	chip.emulateCycle(0xA200)
	assert.Equal(t, uint16(0x200), chip.I)
}

func TestJumpToAddressPlusV0(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.jumpToAddressPlusV0(0x200)
	assert.Equal(t, uint16(0x201), chip.pc)
}

func TestEmulateCycleJumpToAddressPlusV0(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.emulateCycle(0xB200)
	assert.Equal(t, uint16(0x201), chip.pc)
}

func TestRandomWithMask(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x0
	chip.randomWithMask(0x0, 0xF)
	// Mask is lower 4bits, ensure upper stays 0
	// Also check that the result is < 16 (4bit max)
	assert.Equal(t, uint8(0x0), chip.V[0]>>4)
	assert.Less(t, int(chip.V[0]), 16)
}

func TestEmulateCycleRandomWithMask(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x0
	chip.emulateCycle(0xC00F)
	// Mask is lower 4bits, ensure upper stays 0
	// Also check that the result is < 16 (4bit max)
	assert.Equal(t, uint8(0x0), chip.V[0]>>4)
	assert.Less(t, int(chip.V[0]), 16)
}

// TODO: DrawSprite tests

// TODO: KeyPress tests

func TestGetDelayTimer(t *testing.T) {
	chip := NewChip8()
	chip.delayTimer = 0x1
	chip.getDelayTimer(0x0)
	assert.Equal(t, uint8(0x1), chip.V[0])
}

func TestEmulateCycleGetDelayTimer(t *testing.T) {
	chip := NewChip8()
	chip.delayTimer = 0x1
	chip.emulateCycle(0xF007)
	assert.Equal(t, uint8(0x1), chip.V[0])
}

func TestSetDelayTimer(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.setDelayTimer(0x0)
	assert.Equal(t, uint8(0x1), chip.delayTimer)
}

func TestEmulateCycleSetDelayTimer(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.emulateCycle(0xF015)
	assert.Equal(t, uint8(0x1), chip.delayTimer)
}

func TestSetSoundTimer(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.setSoundTimer(0x0)
	assert.Equal(t, uint8(0x1), chip.soundTimer)
}

func TestEmulateCycleSetSoundTimer(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.emulateCycle(0xF018)
	assert.Equal(t, uint8(0x1), chip.soundTimer)
}

func TestAddToIndexRegister(t *testing.T) {
	chip := NewChip8()
	chip.I = 0x1
	chip.V[0] = 0x1
	chip.addToIndexRegister(0x0)
	assert.Equal(t, uint16(0x2), chip.I)
}

func TestEmulateCycleAddToIndexRegister(t *testing.T) {
	chip := NewChip8()
	chip.I = 0x1
	chip.V[0] = 0x1
	chip.emulateCycle(0xF01E)
	assert.Equal(t, uint16(0x2), chip.I)
}

func TestSetIndexToSpriteAddress(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.setIndexToSpriteAddress(0x0)
	assert.Equal(t, uint16(0x5), chip.I)
}

func TestEmulateCycleSetIndexToSpriteAddress(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x1
	chip.emulateCycle(0xF029)
	chip.setIndexToSpriteAddress(0x0)
	assert.Equal(t, uint16(0x5), chip.I)
}

func TestStoreBcd(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x7F
	chip.I = 0x400
	chip.storeBcd(0x0)
	assert.Equal(t, uint8(0x1), chip.memory[0x400])
	assert.Equal(t, uint8(0x2), chip.memory[0x401])
	assert.Equal(t, uint8(0x7), chip.memory[0x402])
}

func TestEmulateCycleStoreBcd(t *testing.T) {
	chip := NewChip8()
	chip.V[0] = 0x7F
	chip.I = 0x400
	chip.emulateCycle(0xF033)
	assert.Equal(t, uint8(0x1), chip.memory[0x400])
	assert.Equal(t, uint8(0x2), chip.memory[0x401])
	assert.Equal(t, uint8(0x7), chip.memory[0x402])
}

func TestReadRegistersFromMemory(t *testing.T) {
	chip := NewChip8()
	chip.I = 0x400
	chip.memory[0x400] = 0x1
	chip.memory[0x401] = 0x2
	chip.memory[0x402] = 0x3
	chip.readRegistersFromMemory(0x300)
	assert.Equal(t, uint8(0x1), chip.V[0])
	assert.Equal(t, uint8(0x2), chip.V[1])
	assert.Equal(t, uint8(0x3), chip.V[2])
}

func TestEmulateCycleReadRegistersFromMemory(t *testing.T) {
	chip := NewChip8()
	chip.I = 0x400
	chip.memory[0x400] = 0x1
	chip.memory[0x401] = 0x2
	chip.memory[0x402] = 0x3
	chip.emulateCycle(0xF365)
	assert.Equal(t, uint8(0x1), chip.V[0])
	assert.Equal(t, uint8(0x2), chip.V[1])
	assert.Equal(t, uint8(0x3), chip.V[2])
}
