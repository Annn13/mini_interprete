package vm

import "fmt"

type VM struct {
	Stack   *Stack
	Mem     *Memory
	Prog    []Instr
	IP      int            // índice de instrucción actual (0..len-1)
	halted  bool
	pcIndex map[int]int    // PC textual → índice en Prog
}

func NewVM(prog []Instr) *VM {
	pcIdx := make(map[int]int, len(prog))
	for i, ins := range prog {
		pcIdx[ins.PC] = i
	}
	return &VM{
		Stack:   NewStack(),
		Mem:     NewMemory(),
		Prog:    prog,
		pcIndex: pcIdx,
	}
}

func (m *VM) Run() error {
	for !m.halted && m.IP < len(m.Prog) {
		ins := m.Prog[m.IP]
		fmt.Printf("PC %d  %-16s arg=%q  stack=%d\n", ins.PC, ins.Op, ins.Arg, m.Stack.Len())
		if err := m.exec(ins); err != nil {
			return fmt.Errorf("PC %d (%s): %w", ins.PC, ins.Op, err)
		}
	}
	return nil
}
