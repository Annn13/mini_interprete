package main

import (
	"fmt"
	"microinterprete/vm"
	"strings"
)

func main() {

	programPathWithFlag := "programs/ejemplo_profe.txt -d"

	parts := strings.Split(programPathWithFlag, " ")
	programPath := parts[0]
	isDebugging := false

	for _, part := range parts[1:] {
		if part == "-d" {
			isDebugging = true
			break
		}
	}

	prog, err := vm.ParseProgram(programPath)
	if err != nil {
		fmt.Println("Parser error:", err)
		return
	}

	m := vm.NewVM(prog, isDebugging)
	if err := m.Run(); err != nil {
		fmt.Println("Runtime error:", err)
	}
}
