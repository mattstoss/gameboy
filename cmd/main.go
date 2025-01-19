package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

const (
	MemorySize = 0x10000 // 64 KB addressable memory space

	EntryPoint = 0x100 // The initial program counter value
)

type (
	CPU struct {
		// registers
		A  uint8
		B  uint8
		C  uint8
		D  uint8
		E  uint8
		H  uint8
		L  uint8
		SP uint16 // stack pointer
		PC uint16 // program counter

		// flags
		Flags struct {
			Z bool // zero
			N bool // subtract
			H bool // half carry
			C bool // carry
		}
	}

	Gameboy struct {
		Memory []byte
		CPU    *CPU
	}
)

func (cpu *CPU) GetHL() uint16 {
	return uint16(cpu.H)<<8 | uint16(cpu.L)
}

func (cpu *CPU) SetHL(value uint16) {
	cpu.H = uint8(value >> 8)
	cpu.L = uint8(value & 0xFF)
}

func NewGameboy(romPath string) (Gameboy, error) {
	rom, err := os.ReadFile(romPath)
	if err != nil {
		return Gameboy{}, fmt.Errorf("failed to read rom: %w", err)
	}

	memory := make([]uint8, MemorySize)
	copy(memory, rom)

	return Gameboy{
		Memory: memory,
		CPU:    &CPU{PC: EntryPoint},
	}, nil
}

func (g *Gameboy) Run() error {
	for {
		pc := g.CPU.PC
		cpu := g.CPU
		opcode := g.Memory[pc]
		operand := binary.LittleEndian.Uint16(g.Memory[pc+1 : pc+3])

		switch opcode {
		case 0x00: // nop
			cpu.PC += 1
		case 0xc3: // jp d16
			cpu.PC = operand
		case 0xaf: // xor A,A
			cpu.PC += 1
			cpu.A = 0
		case 0x21: // ld HL,d16
			cpu.PC += 3
			cpu.SetHL(operand)
		default:
			return fmt.Errorf("unsupported opcode: 0x%x", opcode)
		}
	}
}

func run() error {
	gameboy, err := NewGameboy("testdata/tetris.gb")
	if err != nil {
		return err
	}

	return gameboy.Run()
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
