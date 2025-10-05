package vm

import "fmt"

type VM struct {
	Stack   *Stack
	Mem     *Memory
	Prog    []Instr
	IP      int
	halted  bool
	pcIndex map[int]int
	debug   bool
}

func NewVM(prog []Instr, debug bool) *VM {
	pcIdx := make(map[int]int, len(prog))
	for i, ins := range prog {
		pcIdx[ins.PC] = i
	}
	return &VM{
		Stack:   NewStack(),
		Mem:     NewMemory(),
		Prog:    prog,
		pcIndex: pcIdx,
		debug:   debug,
	}
}

func (m *VM) Run() error {
	for !m.halted && m.IP < len(m.Prog) {
		ins := m.Prog[m.IP]
		// Usar una condición para imprimir solo si debug es true
		if m.debug {
			fmt.Printf("PC %d  %-16s arg=%q  stack=%d\n", ins.PC, ins.Op, ins.Arg, m.Stack.Len())
		}
		if err := m.exec(ins); err != nil {
			return fmt.Errorf("PC %d (%s): %w", ins.PC, ins.Op, err)
		}
	}
	return nil
}
