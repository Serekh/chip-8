package chip8

import (
	"chip-8/pkg/utils"
	"log"
	"math"
	"math/rand"
)

const (
	registersCount uint8  = 16
	programAddress uint16 = 0x0200
)

// Emulator CHIP-8 virtual machine
type Emulator struct {
	display    *Display
	ram        *Memory
	pc         uint16                // Program Counter
	registers  [registersCount]uint8 // V0 to VF
	i          uint16                // Index
	stack      *Stack
	delayTimer uint8
	soundTimer uint8
	input      *Input
	speaker    *utils.Beeper
}

func NewChip8Emulator(beeper *utils.Beeper) *Emulator {
	return &Emulator{
		display: NewDisplay(),
		ram:     NewMemory(),
		pc:      programAddress,
		stack:   NewStack(),
		input:   NewInput(),
		speaker: beeper,
	}
}

func (emulator *Emulator) Reset() {
	emulator.display.Clear()
	emulator.ram.Clear()
	emulator.pc = programAddress
	for i := uint8(0); i < registersCount; i++ {
		emulator.registers[i] = 0x0
	}
	//emulator.registers.
	emulator.i = 0x0
	emulator.stack.Clear()
	emulator.delayTimer = 0
	emulator.soundTimer = 0
	emulator.input.Clear()
}

func (emulator *Emulator) GetDisplayScreen() [ScreenW * ScreenH]uint8 {
	return emulator.display.GetScreen()
}

func (emulator *Emulator) KeyDown(key uint8) {
	emulator.input.KeyPressed(key, true)
}

func (emulator *Emulator) KeyUp(key uint8) {
	emulator.input.KeyPressed(key, false)
}

func (emulator *Emulator) Load(data []uint8) {
	emulator.ram.Copy(programAddress, data)
}

func (emulator *Emulator) Tick() {
	// Fetch
	opcode := emulator.fetch()
	// Decode & Execute
	emulator.execute(opcode)
}

func (emulator *Emulator) TickTimers() {
	if emulator.delayTimer > 0 {
		emulator.delayTimer -= 1
	}

	if emulator.soundTimer > 0 {
		if emulator.soundTimer == 1 {
			emulator.speaker.Play()
		}
		emulator.soundTimer -= 1
	}
}

func (emulator *Emulator) fetch() uint16 {
	opcode := (uint16(emulator.ram.Read(emulator.pc)) << 8) | uint16(emulator.ram.Read(emulator.pc+1))
	emulator.pc += 2
	return opcode
}

func (emulator *Emulator) execute(opcode uint16) {
	// Decode
	digit1 := (opcode & 0xF000) >> 12
	digit2 := (opcode & 0x0F00) >> 8
	digit3 := (opcode & 0x00F0) >> 4
	digit4 := opcode & 0x000F

	switch {
	case digit1 == 0 && digit2 == 0 && digit3 == 0 && digit4 == 0: // 0000 - NOP
		return
	case digit1 == 0 && digit2 == 0 && digit3 == 0xE && digit4 == 0: // 00E0 - Clear screen
		emulator.display.Clear()
		break
	case digit1 == 0 && digit2 == 0 && digit3 == 0xE && digit4 == 0xE: // 00EE - Return from Subroutine
		emulator.pc = emulator.stack.Pop()
		break
	case digit1 == 1: // 1NNN - Jump
		emulator.pc = opcode & 0xFFF
		break
	case digit1 == 2: // 2NNN - Call Subroutine
		emulator.stack.Push(emulator.pc)
		emulator.pc = opcode & 0xFFF
		break
	case digit1 == 3: // 3XNN - Skip next if VX == NN
		nn := uint8(opcode & 0xFF)
		if emulator.registers[digit2] == nn {
			emulator.pc += 2
		}
		break
	case digit1 == 4: // 4XNN - Skip next if VX != NN
		nn := uint8(opcode & 0xFF)
		if emulator.registers[digit2] != nn {
			emulator.pc += 2
		}
		break
	case digit1 == 5: // 5XY0 - Skip next if VX == VY
		if emulator.registers[digit2] == emulator.registers[digit3] {
			emulator.pc += 2
		}
		break
	case digit1 == 6: // 6XNN - VX = NN
		emulator.registers[digit2] = uint8(opcode & 0xFF)
		break
	case digit1 == 7: // 7XNN - VX += NN (carry flag is not changed)
		emulator.registers[digit2] += uint8(opcode & 0xFF)
		break
	case digit1 == 8: // 8XY0 - VX = VY
		emulator.registers[digit2] = emulator.registers[digit3]
		break
	case digit1 == 8 && digit4 == 1: // 8XY1 -	VX |= VY
		emulator.registers[digit2] |= emulator.registers[digit3]
		break
	case digit1 == 8 && digit4 == 2: // 8XY2 - VX &= VY
		emulator.registers[digit2] &= emulator.registers[digit3]
		break
	case digit1 == 8 && digit4 == 3: // 8XY3 - VX ^= VY
		emulator.registers[digit2] ^= emulator.registers[digit3]
		break
	case digit1 == 8 && digit4 == 4: // 8XY4 - VX += VY (VF is set to 1 when there's a carry, and to 0 when there is not)
		emulator.registers[digit2], emulator.registers[0xF] = utils.Add8(emulator.registers[digit2], emulator.registers[digit3], 0)
		break
	case digit1 == 8 && digit4 == 5: // 8XY5 - VX -= VY (VF is set to 0 when there's a borrow, and 1 when there is not)
		emulator.registers[digit2], emulator.registers[0xF] = utils.Sub8(emulator.registers[digit2], emulator.registers[digit3], 0)
		emulator.registers[0xF] ^= 0x1
		break
	case digit1 == 8 && digit4 == 6: // 8XY6 - VX »= 1 (Stores the least significant bit of VX in VF and then shifts VX to the right by 1)
		emulator.registers[0xF] = emulator.registers[digit2] & 0x1
		emulator.registers[digit2] >>= 1
		break
	case digit1 == 8 && digit4 == 7: // 8XY7 - VX = VY - VX (VF is set to 0 when there's a borrow, and 1 when there is not)
		emulator.registers[digit2], emulator.registers[0xF] = utils.Sub8(emulator.registers[digit3], emulator.registers[digit2], 0)
		emulator.registers[0xF] ^= 0x1
		break
	case digit1 == 8 && digit4 == 0xE: // 8XYE - VX «= 1 (Stores the most significant bit of VX in VF and then shifts VX to the left by 1)
		emulator.registers[0xF] = (emulator.registers[digit2] >> 7) & 0x1
		emulator.registers[digit2] <<= 1
		break
	case digit1 == 9 && digit4 == 0: // 9XY0 - Skip if VX != VY
		if emulator.registers[digit2] != emulator.registers[digit3] {
			emulator.pc += 2
		}
		break
	case digit1 == 0xA: // ANNN - I = NNN
		emulator.i = opcode & 0xFFF
		break
	case digit1 == 0xB: // BNNN - PC = V0 + NNN (Jump to V0 + NNN)
		emulator.pc = uint16(emulator.registers[0x0]) + (opcode & 0xFFF)
		break
	case digit1 == 0xC: // CXNN - VX = rand() & NN
		emulator.registers[digit2] = uint8(rand.Intn(math.MaxUint8+1)) & uint8(opcode&0xFF)
		break
	case digit1 == 0xD: // DXYN - Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels
		// Get coords for our sprite
		xCord := uint16(emulator.registers[digit2])
		yCord := uint16(emulator.registers[digit3])
		// VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that does not happen
		emulator.registers[0xF] = 0
		// Iterate over each row of our sprite
		for yLine := uint16(0); yLine < digit4; yLine++ {
			// Determine which memory address our row's data is stored
			addr := emulator.i + yLine
			pixels := emulator.ram.Read(addr)
			// Iterate over each column in our row
			for xLine := uint16(0); xLine < 8; xLine++ {
				// Use a mask to fetch current pixel's bit. Only flip if a 1
				if (pixels & (0b1000_0000 >> xLine)) != 0 {
					// Check if we're about to flip the pixel and set
					emulator.registers[0xF] |= emulator.display.Draw(xCord+xLine, yCord+yLine)
				}
			}
		}
		break
	case digit1 == 0xE && digit3 == 9 && digit4 == 0xE: // EX9E - Skips the next instruction if the key stored in VX is pressed
		if emulator.input.IsKeyPressed(emulator.registers[digit2]) {
			emulator.pc += 2
		}
		break
	case digit1 == 0xE && digit3 == 0xA && digit4 == 1: // EXA1 - Skips the next instruction if the key stored in VX is not pressed
		if !emulator.input.IsKeyPressed(emulator.registers[digit2]) {
			emulator.pc += 2
		}
		break
	case digit1 == 0xF && digit3 == 0 && digit4 == 7: // FX07 - VX = DT (Sets VX to the value of the delay timer)
		emulator.registers[digit2] = emulator.delayTimer
		break
	case digit1 == 0xF && digit3 == 0 && digit4 == 0xA: // FX0A - Wait for Key Press (blocking operation, all instruction halted until next key event)
		key, pressed := emulator.input.IsAnyKeyPressed()
		if pressed {
			emulator.registers[digit2] = key
		} else {
			// Redo opcode
			emulator.pc -= 2
		}
		break
	case digit1 == 0xF && digit3 == 1 && digit4 == 5: // FX15 - DT = VX (Sets the delay timer to VX)
		emulator.delayTimer = emulator.registers[digit2]
		break
	case digit1 == 0xF && digit3 == 1 && digit4 == 8: // FX18 - ST = VX (Sets the sound timer to VX)
		emulator.soundTimer = emulator.registers[digit2]
		break
	case digit1 == 0xF && digit3 == 1 && digit4 == 0xE: // FX1E - I += VX (VF is not affected)
		emulator.i += uint16(emulator.registers[digit2])
		break
	case digit1 == 0xF && digit3 == 2 && digit4 == 9: // FX29 - I = sprite_addr[VX] (Sets I to the location of the sprite for the character in VX)
		emulator.i = uint16(emulator.registers[digit2]) * 5
		break
	case digit1 == 0xF && digit3 == 3 && digit4 == 3: // FX33 - I = BCD of VX (Stores the binary-coded decimal representation of VX, with the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2)
		// Fetch the hundreds digit by dividing by 100 and tossing the decimal
		emulator.ram.Write(emulator.i, emulator.registers[digit2]/100)
		// Fetch the tens digit by dividing by 10, tossing the ones digit and the decimal
		emulator.ram.Write(emulator.i+1, emulator.registers[digit2]/10%10)
		// Fetch the ones digit by tossing the hundreds and the tens
		emulator.ram.Write(emulator.i+2, emulator.registers[digit2]%10)
		break
	case digit1 == 0xF && digit3 == 5 && digit4 == 5: // FX55 - Store V0 - VX into I (Stores from V0 to VX [included], in memory, starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified)
		for idx := uint16(0); idx < digit2+1; idx++ {
			emulator.ram.Write(emulator.i+idx, emulator.registers[idx])
		}
		break
	case digit1 == 0xF && digit3 == 6 && digit4 == 5: // FX65 - Load I into V0 - VX (Fills from V0 to VX [included] with values from memory, starting at address I. The offset from I is increased by 1 for each value read, but I itself is left unmodified)
		for idx := uint16(0); idx < digit2+1; idx++ {
			emulator.registers[idx] = emulator.ram.Read(emulator.i + idx)
		}
		break
	default:
		log.Fatalf("Unimplemented opcode: %v", opcode)
	}
}
