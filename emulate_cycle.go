package main

import (
	"math/rand"
)

func (chip *Chip8) emulateCycle(opcode uint16) {
	if opcode == 0 {
		opcode = uint16(chip.memory[chip.pc])<<8 | uint16(chip.memory[chip.pc+1])
	}
	instruction := opcode & 0xF000

	switch instruction {
	case 0x0000:
		switch opcode & 0x00FF {
		case 0x00E0:
			// 00E0: Clear screen
			chip.clearDisplay()
		case 0x00EE:
			// 00EE: Return from subroutine
			chip.returnFromSubroutine()
		}
	case 0x1000:
		// 1nnn: Jump to address
		chip.jumpToAddress(opcode & 0x0FFF)
	case 0x2000:
		// 2nnn: Call Subroutine
		chip.callSubroutine(opcode & 0x0FFF)
	case 0x3000:
		// 3xkk: Skip next if Vx == kk
		chip.skipIfEqual(opcode&0x0F00, opcode&0x00FF)
	case 0x4000:
		// 4xkk: Skip next if Vx != kk
		chip.skipIfNotEqual(opcode&0x0F00, opcode&0x00FF)
	case 0x5000:
		// 5xy0: Skip next if Vx == Vy
		chip.skipIfRegistersEqual(opcode&0x0F00, opcode&0x00F0)
	case 0x6000:
		// 6xkk: Set Vx = kk
		chip.setRegister(opcode&0x0F00, opcode&0x00FF)
	case 0x7000:
		// 7xkk: Set Vx = Vx + kk
		chip.addtoRegister(opcode&0x0F00, opcode&0x00FF)
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000:
			// 8xy0: Set Vx = Vy
			chip.storeRegister(opcode&0x0F00, opcode&0x00F0)
		case 0x0001:
			// 8xy1: Set Vx = Vx OR Vy
			chip.orRegisters(opcode&0x0F00, opcode&0x00F0)
		case 0x0002:
			// 8xy2: Set Vx = Vx AND Vy
			chip.andRegisters(opcode&0x0F00, opcode&0x00F0)
		case 0x0003:
			// 8xy3: Set Vx = Vx XOR Vy
			chip.xorRegisters(opcode&0x0F00, opcode&0x00F0)
		case 0x0004:
			// 8xy4: Set Vx = Vx + Vy, set VF = carry
			chip.addRegistersWithCarry(opcode&0x0F00, opcode&0x00F0)
		case 0x0005:
			// 8xy5: Set Vx = Vx - Vy, set VF = NOT borrow
			chip.subtractRegistersWithBorrowVxVy(opcode&0x0F00, opcode&0x00F0)
		case 0x0006:
			// 8xy6: Set Vx = SHR 1 Vy
			chip.shiftRight(opcode&0x0F00, opcode&0x00F0)
		case 0x0007:
			// 8xy7: Set Vx = Vy - Vx, set VF = NOT borrow
			chip.subtractRegistersWithBorrowVyVx(opcode&0x0F00, opcode&0x00F0)
		case 0x000E:
			// 8xyE: Set Vx = SHL 1 Vy
			chip.shiftLeft(opcode&0x0F00, opcode&0x00F0)
		}
	case 0x9000:
		// 9xy0: Skip next instruction if Vx != Vy
		chip.skipIfRegistersNotEqual(opcode&0x0F00, opcode&0x00F0)
	case 0xA000:
		// Annn: Set I = nnn
		chip.setIndexRegister(opcode & 0x0FFF)
	case 0xB000:
		// Bnnn: Jump to location nnn + V0
		chip.jumpToAddressPlusV0(opcode & 0x0FFF)
	case 0xC000:
		// Cxkk: Vx = random byte AND kk
		// TODO: Implement random
		chip.randomWithMask(opcode&0x0F00, opcode&0x00FF)
	case 0xD000:
		// Dxyn: Display n-byte sprint starting at memory location I (Vx, Vy), set VF = collision
		chip.drawSprite(opcode&0x0F00, opcode&0x00F0, opcode&0x000F)
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E:
			// Ex9E: Skip next instruction if key with value of Vx is pressed
			chip.skipIfKeyPressed(opcode & 0x0F00)
		case 0x00A1:
			// ExA1: Skip next instruction if key with value of Vx is not pressed
			chip.skipIfKeyNotPressed(opcode & 0x0F00)
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x0007:
			// Fx07: Set Vx = delay timer value
			chip.getDelayTimer(opcode & 0x0F00)
		case 0x000A:
			// Fx0A: Wait for a key press, store the value of the key in Vx
			chip.waitForKeyPress(opcode & 0x0F00)
		case 0x0015:
			// Fx15: Set delay timer = Vx
			chip.setDelayTimer(opcode & 0x0F00)
		case 0x0018:
			// Fx18: Set sound timer = Vx
			chip.setSoundTimer(opcode & 0x0F00)
		case 0x001E:
			// Fx1E: Set I = I + Vx
			chip.addToIndexRegister(opcode & 0x0F00)
		case 0x0029:
			// Fx29: Set I = location of sprite for digit Vx
			chip.setIndexToSpriteAddress(opcode & 0x0F00)
		case 0x0033:
			// Fx33: Store BCD representation of Vx in memory locations I, I+1 and I+2
			chip.storeBcd(opcode & 0x0F00)
		case 0x0055:
			// Fx55: Store registers V0 through Vx in memory starting at location I
			chip.storeRegistersToMemory(opcode & 0x0F00)
		case 0x0065:
			// Fx65: Read registers v0 through Vx from memory starting at location I
			chip.readRegistersFromMemory(opcode & 0x0F00)
		}

	}
}

func (chip *Chip8) clearDisplay() {
	// Bludgeion the screen
	chip.screen = [32][64]uint8{}
}

func (chip *Chip8) returnFromSubroutine() {
	chip.pc = chip.stack[len(chip.stack)-1]
	chip.stack = chip.stack[:len(chip.stack)-1]
}

func (chip *Chip8) jumpToAddress(address uint16) {
	chip.pc = address
}

func (chip *Chip8) callSubroutine(address uint16) {
	chip.stack = append(chip.stack, chip.pc)
	chip.pc = address
}

func (chip *Chip8) skipIfEqual(register uint16, value uint16) {
	vIndex := register >> 8
	comparisonValue := value >> 4

	vValue := chip.V[vIndex]
	originalPc := chip.pc

	registerEqual := vValue == uint8(value)
	if registerEqual {
		chip.pc += 2
	}

	LogOperation("skipIfEqual", []LogDetail{
		{Pairs: []LogPair{{"vIndex", vIndex}, {"vValue", vValue}}},
		{Pairs: []LogPair{{"comparisonValue", comparisonValue}}},
		{Pairs: []LogPair{{"registerEqual", registerEqual}}},
		{Pairs: []LogPair{{"originalPc", originalPc}, {"PC", chip.pc}}},
	})
}

func (chip *Chip8) skipIfNotEqual(register uint16, value uint16) {
	vIndex := register >> 8
	comparisonValue := value >> 4

	vValue := chip.V[vIndex]
	originalPc := chip.pc

	registerNotEqual := vValue != uint8(value)
	if registerNotEqual {
		chip.pc += 2
	}

	LogOperation("skipIfNotEqual", []LogDetail{
		{Pairs: []LogPair{{"vIndex", vIndex}, {"vValue", vValue}}},
		{Pairs: []LogPair{{"comparisonValue", comparisonValue}}},
		{Pairs: []LogPair{{"registerNotEqual", registerNotEqual}}},
		{Pairs: []LogPair{{"originalPc", originalPc}, {"PC", chip.pc}}},
	})
}

func (chip *Chip8) skipIfRegistersEqual(registerX uint16, registerY uint16) {
	vxIndex := registerX >> 8
	vyIndex := registerY >> 4

	originalVx := chip.V[vxIndex]
	originalVy := chip.V[vyIndex]
	originalPc := chip.pc

	registersEqual := chip.V[vxIndex] == chip.V[vyIndex]
	if registersEqual {
		chip.pc += 2
	}

	LogOperation("skipIfRegistersEqual", []LogDetail{
		{Pairs: []LogPair{{"vxIndex", vxIndex}, {"originalVx", originalVx}}},
		{Pairs: []LogPair{{"vyIndex", vyIndex}, {"originalVy", originalVy}}},
		{Pairs: []LogPair{{"registersEqual", registersEqual}}},
		{Pairs: []LogPair{{"originalPc", originalPc}, {"PC", chip.pc}}},
	})
}

func (chip *Chip8) setRegister(register uint16, value uint16) {
	vIndex := register >> 8

	originalV := chip.V[vIndex]

	chip.V[vIndex] = uint8(value)

	LogOperation("setRegister", []LogDetail{
		{Pairs: []LogPair{{"vIndex", vIndex}, {"originalV", originalV}}},
		{Pairs: []LogPair{{"value", value}}},
		{Pairs: []LogPair{{"result", chip.V[vIndex]}}},
	})
}

func (chip *Chip8) addtoRegister(register uint16, value uint16) {
	vIndex := register >> 8

	originalV := chip.V[vIndex]

	chip.V[vIndex] += uint8(value)

	LogOperation("addtoRegister", []LogDetail{
		{Pairs: []LogPair{{"vIndex", vIndex}, {"originalV", originalV}}},
		{Pairs: []LogPair{{"value", value}}},
		{Pairs: []LogPair{{"result", chip.V[vIndex]}}},
	})
}

func (chip *Chip8) storeRegister(registerX uint16, registerY uint16) {
	vxIndex := registerX >> 8
	vyIndex := registerY >> 4

	originalVx := chip.V[vxIndex]
	originalVy := chip.V[vyIndex]

	chip.V[vxIndex] = chip.V[vyIndex]

	LogOperation("storeRegister", []LogDetail{
		{Pairs: []LogPair{{"vxIndex", vxIndex}, {"originalVx", originalVx}}},
		{Pairs: []LogPair{{"vyIndex", vyIndex}, {"originalVy", originalVy}}},
		{Pairs: []LogPair{{"result", chip.V[vxIndex]}}},
	})
}

func (chip *Chip8) orRegisters(registerX uint16, registerY uint16) {
	vxIndex := registerX >> 8
	vyIndex := registerY >> 4

	originalVx := chip.V[vxIndex]
	originalVy := chip.V[vyIndex]

	chip.V[vxIndex] |= chip.V[vyIndex]

	LogOperation("orRegisters", []LogDetail{
		{Pairs: []LogPair{{"vxIndex", vxIndex}, {"originalVx", originalVx}}},
		{Pairs: []LogPair{{"vyIndex", vyIndex}, {"originalVy", originalVy}}},
		{Pairs: []LogPair{{"result", chip.V[vxIndex]}}},
	})
}

func (chip *Chip8) andRegisters(registerX uint16, registerY uint16) {
	vxIndex := registerX >> 8
	vyIndex := registerY >> 4

	originalVx := chip.V[vxIndex]
	originalVy := chip.V[vyIndex]

	chip.V[vxIndex] &= chip.V[vyIndex]

	LogOperation("andRegisters", []LogDetail{
		{Pairs: []LogPair{{"vxIndex", vxIndex}, {"originalVx", originalVx}}},
		{Pairs: []LogPair{{"vyIndex", vyIndex}, {"originalVy", originalVy}}},
		{Pairs: []LogPair{{"result", chip.V[vxIndex]}}},
	})
}

func (chip *Chip8) xorRegisters(registerX uint16, registerY uint16) {
	vxIndex := registerX >> 8
	vyIndex := registerY >> 4

	originalVx := chip.V[vxIndex]
	originalVy := chip.V[vyIndex]

	chip.V[vxIndex] ^= chip.V[vyIndex]

	LogOperation("xorRegisters", []LogDetail{
		{Pairs: []LogPair{{"vxIndex", vxIndex}, {"originalVx", originalVx}}},
		{Pairs: []LogPair{{"vyIndex", vyIndex}, {"originalVy", originalVy}}},
		{Pairs: []LogPair{{"result", chip.V[vxIndex]}}},
	})
}

func (chip *Chip8) addRegistersWithCarry(registerX uint16, registerY uint16) {
	vxIndex := registerX >> 8
	vyIndex := registerY >> 4

	originalVx := chip.V[vxIndex]
	originalVy := chip.V[vyIndex]

	additionTemp := uint16(originalVx) + uint16(originalVy)
	if additionTemp > 255 {
		chip.V[0xF] = 1
	} else {
		chip.V[0xF] = 0
	}

	chip.V[vxIndex] = uint8(additionTemp)

	LogOperation("addRegistersWithCarry", []LogDetail{
		{Pairs: []LogPair{{"vxIndex", vxIndex}, {"originalVx", originalVx}}},
		{Pairs: []LogPair{{"vyIndex", vyIndex}, {"originalVy", originalVy}}},
		{Pairs: []LogPair{{"result", chip.V[vxIndex]}, {"carry", chip.V[0xF] == 0}}},
	})
}

func (chip *Chip8) subtractRegistersWithBorrowVxVy(registerX uint16, registerY uint16) {
	vxIndex := registerX >> 8
	vyIndex := registerY >> 4

	originalVx := chip.V[vxIndex]
	originalVy := chip.V[vyIndex]

	if originalVx >= originalVy {
		chip.V[0xF] = 1
	} else {
		chip.V[0xF] = 0
	}

	chip.V[vxIndex] = originalVx - originalVy

	LogOperation("subtractRegistersWithBorrowVxVy", []LogDetail{
		{Pairs: []LogPair{{"vxIndex", vxIndex}, {"originalVx", originalVx}}},
		{Pairs: []LogPair{{"vyIndex", vyIndex}, {"originalVy", originalVy}}},
		{Pairs: []LogPair{{"result", chip.V[vxIndex]}, {"borrow", chip.V[0xF] == 0}}},
	})
}

func (chip *Chip8) shiftRight(registerX uint16, registerY uint16) {
	vxIndex := registerX >> 8
	vyIndex := registerY >> 4

	originalVx := chip.V[vxIndex]
	originalVy := chip.V[vyIndex]

	chip.V[vxIndex] = chip.V[vyIndex] >> 1

	LogOperation("shiftRight", []LogDetail{
		{Pairs: []LogPair{{"vxIndex", vxIndex}, {"originalVx", originalVx}}},
		{Pairs: []LogPair{{"vyIndex", vyIndex}, {"originalVy", originalVy}}},
		{Pairs: []LogPair{{"result", chip.V[vxIndex]}}},
	})
}

func (chip *Chip8) subtractRegistersWithBorrowVyVx(registerX uint16, registerY uint16) {
	vxIndex := registerX >> 8
	vyIndex := registerY >> 4

	originalVx := chip.V[vxIndex]
	originalVy := chip.V[vyIndex]

	if originalVy < originalVx {
		chip.V[0xF] = 0
	} else {
		chip.V[0xF] = 1
	}

	chip.V[vyIndex] -= chip.V[vxIndex]

	LogOperation("subtractRegistersWithBorrowVyVx", []LogDetail{
		{Pairs: []LogPair{{"vyIndex", vyIndex}, {"originalVy", originalVy}}},
		{Pairs: []LogPair{{"vxIndex", vxIndex}, {"originalVx", originalVx}}},
		{Pairs: []LogPair{{"result", chip.V[vyIndex]}, {"borrow", chip.V[0xF] == 0}}},
	})
}

func (chip *Chip8) shiftLeft(registerX uint16, registerY uint16) {
	vxIndex := registerX >> 8
	vyIndex := registerY >> 4

	originalVx := chip.V[vxIndex]
	originalVy := chip.V[vyIndex]

	chip.V[vxIndex] = chip.V[vyIndex] << 1

	LogOperation("shiftRight", []LogDetail{
		{Pairs: []LogPair{{"vxIndex", vxIndex}, {"originalVx", originalVx}}},
		{Pairs: []LogPair{{"vyIndex", vyIndex}, {"originalVy", originalVy}}},
		{Pairs: []LogPair{{"result", chip.V[vxIndex]}}},
	})
}

func (chip *Chip8) skipIfRegistersNotEqual(registerX uint16, registerY uint16) {
	vxIndex := registerX >> 8
	vyIndex := registerY >> 4

	originalVx := chip.V[vxIndex]
	originalVy := chip.V[vyIndex]
	originalPc := chip.pc

	registersNotEqual := chip.V[vxIndex] != chip.V[vyIndex]
	if registersNotEqual {
		chip.pc += 2
	}

	LogOperation("skipIfRegistersEqual", []LogDetail{
		{Pairs: []LogPair{{"vxIndex", vxIndex}, {"originalVx", originalVx}}},
		{Pairs: []LogPair{{"vyIndex", vyIndex}, {"originalVy", originalVy}}},
		{Pairs: []LogPair{{"registersNotEqual", registersNotEqual}}},
		{Pairs: []LogPair{{"originalPc", originalPc}, {"PC", chip.pc}}},
	})
}

func (chip *Chip8) setIndexRegister(address uint16) {
	prevIndex := chip

	chip.I = address

	LogOperation("setIndexRegister", []LogDetail{
		{Pairs: []LogPair{{"prevIndex", prevIndex}, {"address", address}}},
		{Pairs: []LogPair{{"newIndex", chip.I}}},
	})
}

func (chip *Chip8) jumpToAddressPlusV0(value uint16) {
	prevPc := chip.pc
	vValue := chip.V[0]

	chip.pc = uint16(vValue) + value

	LogOperation("jumpToAddressPlusV0", []LogDetail{
		{Pairs: []LogPair{{"V0", vValue}, {"value", value}}},
		{Pairs: []LogPair{{"prevPC", prevPc}, {"newPC", chip.pc}}},
	})
}

func (chip *Chip8) randomWithMask(register uint16, mask uint16) {
	vIndex := register >> 4
	vValue := chip.V[vIndex]

	randomVal := uint8(rand.Intn(255))

	chip.V[vIndex] = uint8(mask) & uint8(rand.Intn(255))

	LogOperation("randomWithMask", []LogDetail{
		{Pairs: []LogPair{{"vIndex", vIndex}, {"vValue", vValue}}},
		{Pairs: []LogPair{{"mask", mask}, {"randomVal", randomVal}}},
		{Pairs: []LogPair{{"result", chip.V[vIndex]}}},
	})
}

func (chip *Chip8) drawSprite(registerX uint16, registerY uint16, height uint16) {
	x := uint16(chip.V[registerX>>8])
	y := uint16(chip.V[registerY>>4])
	chip.V[0xF] = 0

	for yline := uint16(0); yline < height; yline++ {
		pixel := chip.memory[chip.I+yline]
		for xline := uint16(0); xline < 8; xline++ {
			if (pixel & (0x80 >> xline)) != 0 {
				if chip.screen[x+xline][y+yline] == 1 {
					chip.V[0xF] = 1
				}
				chip.screen[x+xline][y+yline] ^= 1
			}
		}
	}

	LogOperation("drawSprite", []LogDetail{
		{Pairs: []LogPair{{"xLoc", x}, {"yLoc", y}}},
		{Pairs: []LogPair{{"height", height}, {"collision", chip.V[0xF] == 1}}},
	})
}

func (chip *Chip8) skipIfKeyPressed(opcode uint16) {
	vIndex := opcode >> 8
	keyRegister := chip.V[vIndex]
	prevPc := chip.pc

	keyPressed := chip.keys[keyRegister] != 0
	if keyPressed {
		chip.pc += 2
	}

	LogOperation("skipIfKeyPressed", []LogDetail{
		{Pairs: []LogPair{{"key", vIndex}, {"keyPressed", keyPressed}}},
		{Pairs: []LogPair{{"prevPC", prevPc}, {"newPC", chip.pc}}},
	})
}

func (chip *Chip8) skipIfKeyNotPressed(opcode uint16) {
	vIndex := opcode >> 8
	keyRegister := chip.V[vIndex]
	prevPc := chip.pc

	keyNotPressed := chip.keys[keyRegister] == 0
	if keyNotPressed {
		chip.pc += 2
	}

	LogOperation("skipIfKeyNotPressed", []LogDetail{
		{Pairs: []LogPair{{"key", vIndex}, {"keyNotPressed", keyNotPressed}}},
		{Pairs: []LogPair{{"prevPC", prevPc}, {"newPC", chip.pc}}},
	})
}

func (chip *Chip8) getDelayTimer(register uint16) {
	vIndex := register >> 8
	originalV := chip.V[vIndex]

	chip.V[vIndex] = chip.delayTimer

	LogOperation("getDelayTimer", []LogDetail{
		{Pairs: []LogPair{{"vIndex", vIndex}, {"originalV", originalV}}},
		{Pairs: []LogPair{{"newV", chip.V[vIndex]}}},
	})
}

func (chip *Chip8) waitForKeyPress(opcode uint16) {
	chip.waitingForKeyPress = true
	keyRegister := uint8(opcode >> 8)

	chip.keyRegister = keyRegister

	LogOperation("waitingForKeyPress", []LogDetail{
		{Pairs: []LogPair{{"keyRegister", keyRegister}}},
	})
}

func (chip *Chip8) setDelayTimer(register uint16) {
	prevDelayTimer := chip.delayTimer

	vIndex := register >> 8
	vValue := chip.V[vIndex]

	chip.delayTimer = vValue

	LogOperation("setDelayTimer", []LogDetail{
		{Pairs: []LogPair{{"vIndex", vIndex}, {"vValue", vValue}}},
		{Pairs: []LogPair{{"prevDelayTimer", prevDelayTimer}, {"newDelayTimer", chip.delayTimer}}},
	})
}

func (chip *Chip8) setSoundTimer(register uint16) {
	prevSoundTimer := chip.soundTimer

	vIndex := register >> 8
	vValue := chip.V[vIndex]

	chip.soundTimer = vValue

	LogOperation("setSoundTimer", []LogDetail{
		{Pairs: []LogPair{{"vIndex", vIndex}, {"vValue", vValue}}},
		{Pairs: []LogPair{{"prevSoundTimer", prevSoundTimer}, {"newSoundTimer", chip.soundTimer}}},
	})
}

func (chip *Chip8) addToIndexRegister(register uint16) {
	prevIndex := chip.I

	vIndex := register >> 8
	vValue := chip.V[vIndex]

	chip.I += uint16(vValue)

	LogOperation("addToIndexRegister", []LogDetail{
		{Pairs: []LogPair{{"vIndex", vIndex}, {"vValue", vValue}}},
		{Pairs: []LogPair{{"prevIndex", prevIndex}, {"newIndex", chip.I}}},
	})
}

func (chip *Chip8) setIndexToSpriteAddress(register uint16) {
	prevIndex := chip.I

	vIndex := register >> 8
	vValue := chip.V[vIndex]

	chip.I = uint16(vValue) * 5

	LogOperation("setIndexToSpriteAddress", []LogDetail{
		{Pairs: []LogPair{{"vIndex", vIndex}, {"vValue", vValue}}},
		{Pairs: []LogPair{{"prevIndex", prevIndex}, {"newIndex", chip.I}}},
	})
}

func (chip *Chip8) storeBcd(register uint16) {
	vIndex := register >> 8
	vValue := chip.V[vIndex]

	chip.memory[chip.I] = vValue / 100
	chip.memory[chip.I+1] = (vValue / 10) % 10
	chip.memory[chip.I+2] = (vValue % 100) % 10

	LogOperation("storeBcd", []LogDetail{
		{Pairs: []LogPair{{"vIndex", vIndex}, {"vValue", vValue}}},
		{Pairs: []LogPair{{"BCD 100's", chip.memory[chip.I]}}},
		{Pairs: []LogPair{{"BCD 10's", chip.memory[chip.I+1]}}},
		{Pairs: []LogPair{{"BCD 1's", chip.memory[chip.I+2]}}},
	})
}

func (chip *Chip8) storeRegistersToMemory(upperRegister uint16) {
	upperRegisterIndex := upperRegister >> 8

	var details []LogDetail
	for i := uint16(0); i <= upperRegisterIndex; i++ {
		details = append(details, LogDetail{Pairs: []LogPair{{"index", chip.memory[chip.I+i]}, {"register", i}}})

		prevValue := chip.memory[chip.I+i]

		chip.memory[chip.I+i] = chip.V[i]
		details = append(details, LogDetail{Pairs: []LogPair{{"prevValue", prevValue}, {"newValue", chip.memory[chip.I+i]}}})
	}

	header := []LogDetail{
		{Pairs: []LogPair{{"upperRegisterIndex", upperRegisterIndex}, {"indexRegister", chip.I}}},
	}

	header = append(header, details...)
	LogOperation("storeRegistersToMemory", header)
}

func (chip *Chip8) readRegistersFromMemory(upperRegister uint16) {
	upperRegisterIndex := upperRegister >> 8

	var details []LogDetail
	for i := uint16(0); i <= upperRegisterIndex; i++ {
		details = append(details, LogDetail{Pairs: []LogPair{{"register", i}, {"index", chip.memory[chip.I+i]}}})

		prevValue := chip.V[i]

		chip.V[i] = chip.memory[chip.I+i]
		details = append(details, LogDetail{Pairs: []LogPair{{"prevValue", prevValue}, {"newValue", chip.V[i]}}})
	}

	header := []LogDetail{
		{Pairs: []LogPair{{"upperRegisterIndex", upperRegisterIndex}, {"indexRegister", chip.I}}},
	}

	header = append(header, details...)
	LogOperation("readRegistersFromMemory", header)
}
