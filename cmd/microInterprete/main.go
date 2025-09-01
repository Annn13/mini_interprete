package main

import (
	"fmt"
	"microinterprete/vm"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: microinterprete <ruta_programa.txt>")
		return
	}

	prog, err := vm.ParseProgram(os.Args[1])
	if err != nil {
		fmt.Println("Parser error:", err)
		return
	}

	m := vm.NewVM(prog)
	if err := m.Run(); err != nil {
		fmt.Println("Runtime error:", err)
	}
}
