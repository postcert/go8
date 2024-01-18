package main

import "math/rand"

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
	if chip.V[register>>8] == uint8(value) {
		chip.pc += 2
	}
}

func (chip *Chip8) skipIfNotEqual(opcode uint16, value uint16) {
	if chip.V[opcode>>8] != uint8(value) {
		chip.pc += 2
	}
}

func (chip *Chip8) skipIfRegistersEqual(registerX uint16, registerY uint16) {
	if chip.V[registerX>>8] == chip.V[registerY>>4] {
		chip.pc += 2
	}
}

func (chip *Chip8) setRegister(registerX uint16, value uint16) {
	chip.V[registerX>>8] = uint8(value)
}

func (chip *Chip8) addtoRegister(register uint16, value uint16) {
	chip.V[register>>8] += uint8(value)
}

func (chip *Chip8) storeRegister(registerX uint16, registerY uint16) {
	chip.V[registerX>>8] = chip.V[registerY>>4]
}

func (chip *Chip8) orRegisters(registerX uint16, registerY uint16) {
	chip.V[registerX>>8] |= chip.V[registerY>>4]
}

func (chip *Chip8) andRegisters(registerX uint16, registerY uint16) {
	chip.V[registerX>>8] &= chip.V[registerY>>4]
}

func (chip *Chip8) xorRegisters(registerX uint16, registerY uint16) {
	chip.V[registerX>>8] ^= chip.V[registerY>>4]
}

func (chip *Chip8) addRegistersWithCarry(opcode uint16, value uint16) {
	// TODO: 255 + 1 = 0, not 256 becuase of type. Failing overflow test below
	temp := uint16(chip.V[opcode>>8]) + uint16(chip.V[opcode>>4])
	if temp > 255 {
		chip.V[0xF] = 1
	} else {
		chip.V[0xF] = 0
	}

	chip.V[opcode>>8] = uint8(temp)
}

func (chip *Chip8) subtractRegistersWithBorrowVxVy(registerX uint16, registerY uint16) {
	if chip.V[registerX>>8] < chip.V[registerY>>4] {
		chip.V[0xF] = 0
	} else {
		chip.V[0xF] = 1
	}

	chip.V[registerX>>8] -= chip.V[registerY>>4]
}

func (chip *Chip8) shiftRight(registerX uint16, registerY uint16) {
	chip.V[registerX>>8] = chip.V[registerY>>4] >> 1
}

func (chip *Chip8) subtractRegistersWithBorrowVyVx(registerX uint16, registerY uint16) {
	if chip.V[registerY>>4] < chip.V[registerX>>8] {
		chip.V[0xF] = 0
	} else {
		chip.V[0xF] = 1
	}

	chip.V[registerY>>4] -= chip.V[registerX>>8]
}

func (chip *Chip8) shiftLeft(registerX uint16, registerY uint16) {
	chip.V[registerX>>8] = chip.V[registerY>>4] << 1
}

func (chip *Chip8) skipIfRegistersNotEqual(registerX uint16, registerY uint16) {
	if chip.V[registerX>>8] != chip.V[registerY>>4] {
		chip.pc += 2
	}
}

func (chip *Chip8) setIndexRegister(address uint16) {
	chip.I = address
}

func (chip *Chip8) jumpToAddressPlusV0(value uint16) {
	chip.pc = uint16(chip.V[0]) + value
}

func (chip *Chip8) randomWithMask(register uint16, mask uint16) {
	chip.V[register>>8] = uint8(mask) & uint8(rand.Intn(255))
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
}

func (chip *Chip8) skipIfKeyPressed(opcode uint16) {
	if chip.keys[chip.V[opcode>>8]] != 0 {
		chip.pc += 2
	}
}

func (chip *Chip8) skipIfKeyNotPressed(opcode uint16) {
	if chip.keys[chip.V[opcode>>8]] == 0 {
		chip.pc += 2
	}
}

func (chip *Chip8) getDelayTimer(register uint16) {
	chip.V[register>>8] = chip.delayTimer
}

func (chip *Chip8) waitForKeyPress(opcode uint16) {
	chip.waitingForKeyPress = true
	chip.keyRegister = uint8(opcode >> 8)
}

func (chip *Chip8) setDelayTimer(register uint16) {
	chip.delayTimer = chip.V[register>>8]
}

func (chip *Chip8) setSoundTimer(register uint16) {
	chip.soundTimer = chip.V[register>>8]
}

func (chip *Chip8) addToIndexRegister(register uint16) {
	chip.I += uint16(chip.V[register>>8])
}

func (chip *Chip8) setIndexToSpriteAddress(register uint16) {
	chip.I = uint16(chip.V[register>>8]) * 5
}

func (chip *Chip8) storeBcd(register uint16) {
	chip.memory[chip.I] = chip.V[register>>8] / 100
	chip.memory[chip.I+1] = (chip.V[register>>8] / 10) % 10
	chip.memory[chip.I+2] = (chip.V[register>>8] % 100) % 10
}

func (chip *Chip8) storeRegistersToMemory(upperRegister uint16) {
	for i := uint16(0); i <= upperRegister>>8; i++ {
		chip.memory[chip.I+i] = chip.V[i]
	}
}

func (chip *Chip8) readRegistersFromMemory(upperRegister uint16) {
	for i := uint16(0); i <= upperRegister>>8; i++ {
		chip.V[i] = chip.memory[chip.I+i]
	}
}