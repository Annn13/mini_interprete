package vm

import (
	"bufio" //leer linea por linea
	"fmt"
	"os"      //manejo de archivos
	"strconv" //conversiones
	"strings" //manipulación de texto
)

func ParseProgram(path string) ([]Instr, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var prog []Instr
	sc := bufio.NewScanner(f)
	expectedPC := 0 //verifica que los PCs sean consecutivos
	lineNo := 0     //numero de linea actual

	for sc.Scan() {
		lineNo++
		raw := sc.Text() //texto de la linea actual

		// cortar comentarios
		if i := strings.Index(raw, "#"); i >= 0 {
			raw = raw[:i]
		}
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}

		toks := strings.Fields(raw) // maneja tabs y múltiples espacios
		if len(toks) < 2 {

			return nil, fmt.Errorf("l%d: se esperan al menos [PC OPCODE]", lineNo)
		}

		pc, err := strconv.Atoi(toks[0])
		if err != nil {

			return nil, fmt.Errorf("l%d: pc inválido: %v", lineNo, toks[0])
		}
		if pc != expectedPC { //debe ser consecutivo

			return nil, fmt.Errorf("l%d: pc %d no coincide con esperado %d", lineNo, pc, expectedPC)
		}
		expectedPC++

		op := strings.ToUpper(toks[1])
		var arg string
		hasArg := false
		if len(toks) > 2 {
			arg = strings.Join(toks[2:], " ") // todo lo que sigue es UN parámetro
			hasArg = true
		}

		prog = append(prog, Instr{PC: pc, Op: op, Arg: arg, HasArg: hasArg})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return prog, nil
}
