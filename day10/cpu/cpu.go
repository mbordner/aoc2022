package cpu

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type CPU struct {
	cycle     int
	registers map[string]int
	states    []string
}

func (cpu *CPU) String() string {
	return fmt.Sprintf("cycle: %d, registers: %s", cpu.cycle, cpu.CurrentRegisterState())
}

func NewCPU() *CPU {
	cpu := new(CPU)
	cpu.Reset()
	return cpu
}

func (cpu *CPU) Reset() {
	cpu.cycle = 0
	cpu.registers = make(map[string]int)
	cpu.registers["X"] = 1
	cpu.states = make([]string, 0, 500)
}

func (cpu *CPU) GetState(c int) (string, error) {
	if c >= 0 && c < len(cpu.states) {
		return cpu.states[c], nil
	}
	return "", errors.New("invalid cycle")
}

func (cpu *CPU) CurrentRegisterState() string {
	s, _ := json.Marshal(cpu.registers)
	return string(s)
}

func (cpu *CPU) GetValueFromState(register string, c int) (int, error) {
	s, err := cpu.GetState(c - 1)
	if err != nil {
		return -1, err
	}
	var data map[string]int
	_ = json.Unmarshal([]byte(s), &data)
	if v, exists := data[register]; exists {
		return v, nil
	}
	return -2, errors.New("unknown register")
}

func (cpu *CPU) Process(instruction string) {
	if strings.HasPrefix(instruction, "noop") {
		cpu.tick()
	} else if strings.HasPrefix(instruction, "addx ") {
		cpu.tick()
		cpu.tick()
		v, _ := strconv.Atoi(instruction[5:])
		cpu.registers["X"] += v
	}
}

func (cpu *CPU) tick() {
	s := cpu.CurrentRegisterState()
	cpu.states = append(cpu.states, string(s))
	cpu.cycle++
}
