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
