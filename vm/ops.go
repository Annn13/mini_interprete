package vm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (m *VM) exec(ins Instr) error {
	switch ins.Op {
	case "LOAD_CONST":
		if !ins.HasArg {
			return errors.New("falta argumento")
		}
		v, err := parseLiteral(ins.Arg)
		if err != nil {
			return err
		}
		m.Stack.Push(v)

		if m.debug {
			fmt.Printf("DEBUG: Pila después de LOAD_CONST: %s\n", m.Stack.String())
		}

	case "LOAD_FAST":
		if !ins.HasArg {
			return errors.New("falta nombre de variable")
		}
		v, ok := m.Mem.Get(ins.Arg)
		if !ok {
			return fmt.Errorf("variable no definida: %s", ins.Arg)
		}
		m.Stack.Push(v)

		if m.debug {
			fmt.Println("DEBUG: LOAD_FAST, Pila:", m.Stack.String())
		}

	case "STORE_FAST":
		if !ins.HasArg {
			return errors.New("falta nombre de variable")
		}
		v, ok := m.Stack.Pop()
		if !ok {
			return errors.New("underflow en STORE_FAST")
		}
		m.Mem.Set(ins.Arg, v)

		if m.debug {
			fmt.Println("DEBUG: STORE_FAST, Memoria:", m.Mem.String())
			fmt.Println("DEBUG: STORE_FAST, Pila:", m.Stack.String())
		}

	case "LOAD_GLOBAL":
		// Solo "print"
		if !ins.HasArg {
			return errors.New("falta nombre global")
		}
		if strings.ToLower(ins.Arg) != "print" {
			return fmt.Errorf("global no soportado: %s", ins.Arg)
		}

		m.Stack.Push(nativePrint{})

		//PARA PROBAR
		if m.debug {
			fmt.Printf("DEBUG: LOAD_GLOBAL cargó la referencia a 'print'. Pila: %s\n", m.Stack.String())
		}

	case "CALL_FUNCTION":
		if !ins.HasArg {
			return errors.New("falta cantidad")
		}
		n, err := strconv.Atoi(ins.Arg)
		if err != nil || n < 0 {
			return errors.New("aridad inválida")
		}

		// sacar n argumentos (LIFO → guardamos y luego invertimos para orden natural)
		args := make([]any, 0, n)
		for i := 0; i < n; i++ {
			v, ok := m.Stack.Pop()
			if !ok {
				return errors.New("underflow en CALL_FUNCTION (args)")
			}
			args = append(args, v)
		}
		// referencia a función (abajo de los args)
		fn, ok := m.Stack.Pop()
		if !ok {
			return errors.New("underflow en CALL_FUNCTION (funcref)")
		}

		// solo print nativo
		if _, isPrint := fn.(nativePrint); !isPrint {
			return errors.New("CALL_FUNCTION: solo se permite print")
		}
		// imprimir en orden de izquierda a derecha
		for i := len(args) - 1; i >= 0; i-- {
			if i != len(args)-1 {
				fmt.Print(" ")
			}
			fmt.Print(args[i])
		}
		fmt.Println()

	case "COMPARE_OP":
		if !ins.HasArg {
			return errors.New("falta operador (<,<=,==,!=,>,>=)")
		}
		b, ok := m.Stack.Pop()
		if !ok {
			return errors.New("underflow")
		}
		a, ok := m.Stack.Pop()
		if !ok {
			return errors.New("underflow")
		}
		res, err := compare(a, b, ins.Arg)
		if err != nil {
			return err
		}
		m.Stack.Push(res) // usa bool Go, o si prefieres 0/1

		if m.debug {
			fmt.Printf("DEBUG: COMPARE_OP: %v %s %v -> %v. Pila: %s\n", a, ins.Arg, b, res, m.Stack.String())
		}

	case "BINARY_ADD":
		if err := m.binOp(add); err != nil {
			return err
		}

		if m.debug {
			fmt.Println("DEBUG: BINARY_ADD:", m.Stack.String())
		}

	case "BINARY_SUBSTRACT":
		if err := m.binOp(sub); err != nil {
			return err
		}
		if m.debug {
			fmt.Println("DEBUG: BINARY_SUBSTRACT:", m.Stack.String())
		}

	case "BINARY_MULTIPLY":
		if err := m.binOp(mul); err != nil {
			return err
		}

		if m.debug {
			fmt.Println("DEBUG: BINARY_MULTIPLY:", m.Stack.String())
		}

	case "BINARY_DIVIDE":
		if err := m.binOp(div); err != nil {
			return err
		} // división entera

		if m.debug {
			fmt.Println("DEBUG: BINARY_DIVIDE:", m.Stack.String())
		}

	case "BINARY_AND":
		if err := m.binOp(andOp); err != nil {
			return err
		}

		if m.debug {
			fmt.Println("DEBUG: BINARY_AND, Pila:", m.Stack.String())
		}

	case "BINARY_OR":
		if err := m.binOp(orOp); err != nil {
			return err
		}
		if m.debug {
			fmt.Println("DEBUG: BINARY_OR:", m.Stack.String())
		}

	case "BINARY_MODULO":
		if err := m.binOp(mod); err != nil {
			return err
		}

		if m.debug {
			fmt.Println("DEBUG: BINARY_MODULO:", m.Stack.String())
		}

	case "STORE_SUBSCR":
		// array[index] = value
		idx, ok := m.Stack.Pop()
		if !ok {
			return errors.New("underflow value")
		}
		arr, ok := m.Stack.Pop()
		if !ok {
			return errors.New("underflow index")
		}
		val, ok := m.Stack.Pop()
		if !ok {
			return errors.New("underflow array")
		}

		if err := storeSubscr(arr, idx, val); err != nil {
			return err
		}

		if m.debug {
			fmt.Println("DEBUG: STORE_SUBSCR:", m.Stack.String())
		}

	case "BINARY_SUBSCR":
		// push array[index]
		idx, ok := m.Stack.Pop()
		if !ok {
			return errors.New("underflow index")
		}
		arr, ok := m.Stack.Pop()
		if !ok {
			return errors.New("underflow array")
		}
		el, err := getSubscr(arr, idx)
		if err != nil {
			return err
		}
		m.Stack.Push(el)

		if m.debug {
			fmt.Println("DEBUG: BINARY_SUBSCR:", m.Stack.String())
		}

	case "BUILD_LIST":
		if !ins.HasArg {
			return errors.New("falta cantidad de elementos")
		}
		n, err := strconv.Atoi(ins.Arg)
		if err != nil || n < 0 {
			return errors.New("cantidad inválida")
		}
		tmp := make([]any, 0, n)
		for i := 0; i < n; i++ {
			v, ok := m.Stack.Pop()
			if !ok {
				return errors.New("underflow en BUILD_LIST")
			}
			tmp = append(tmp, v)
		}
		// revertir para orden humano (elem1, elem2, ..., elemN)
		for i, j := 0, len(tmp)-1; i < j; i, j = i+1, j-1 {
			tmp[i], tmp[j] = tmp[j], tmp[i]
		}
		m.Stack.Push(tmp) // []any

		if m.debug {
			fmt.Println("DEBUG: BUILD_LIST:", m.Stack.String())
		}

	case "JUMP_ABSOLUTE":
		if !ins.HasArg {
			return errors.New("falta target")
		}
		target, err := strconv.Atoi(ins.Arg)
		if err != nil {
			return errors.New("target inválido")
		}
		idx, ok := m.pcIndex[target]
		if !ok {
			return fmt.Errorf("target no existe: %d", target)
		}
		m.IP = idx
		return nil // importante: NO auto-incrementar

	case "JUMP_IF_TRUE":
		if !ins.HasArg {
			return errors.New("falta target")
		}
		target, err := strconv.Atoi(ins.Arg)
		if err != nil {
			return errors.New("target inválido")
		}
		val, ok := m.Stack.Pop()
		if !ok {
			return errors.New("underflow en JUMP_IF_TRUE")
		}
		if isTruthy(val) {
			idx, ok := m.pcIndex[target]
			if !ok {
				return fmt.Errorf("target no existe: %d", target)
			}
			m.IP = idx
			return nil
		}

	case "JUMP_IF_FALSE":
		if !ins.HasArg {
			return errors.New("falta target")
		}
		target, err := strconv.Atoi(ins.Arg)
		if err != nil {
			return errors.New("target inválido")
		}
		val, ok := m.Stack.Pop()
		if !ok {
			return errors.New("underflow en JUMP_IF_FALSE")
		}
		if !isTruthy(val) {
			idx, ok := m.pcIndex[target]
			if !ok {
				return fmt.Errorf("target no existe: %d", target)
			}
			m.IP = idx
			return nil
		}

	case "END":
		m.halted = true

	default:
		return fmt.Errorf("opcode no soportado: %s", ins.Op)
	}

	// avance normal
	m.IP++
	return nil
}

// ---------- helpers de ejecución ----------

type nativePrint struct{}

func (m *VM) binOp(op func(a, b any) (any, error)) error {
	// Saca el segundo operando (el que está en el tope)
	b, ok := m.Stack.Pop()
	if !ok {
		return errors.New("underflow")
	}
	// Saca el primer operando
	a, ok := m.Stack.Pop()
	if !ok {
		return errors.New("underflow")
	}

	// Aplica la operación andOp, orOp, etc.
	res, err := op(a, b)
	if err != nil {
		return err
	}

	// Empuja el resultado a la pila
	m.Stack.Push(res)

	return nil
}

func parseLiteral(arg string) (any, error) {
	arg = strings.TrimSpace(arg)

	// boolean
	if strings.ToLower(arg) == "true" {
		return true, nil
	}
	if strings.ToLower(arg) == "false" {
		return false, nil
	}

	// literal de lista
	if strings.HasPrefix(arg, "[") && strings.HasSuffix(arg, "]") {
		// Eliminar los corchetes
		inner := strings.Trim(arg, "[]")

		// MANEJO DE LISTAS VACÍAS
		if inner == "" {
			return []any{}, nil
		}

		// Dividir por comas y parsear cada elemento
		parts := strings.Split(inner, ",")
		result := make([]any, len(parts))
		for i, p := range parts {
			v, err := parseLiteral(strings.TrimSpace(p))
			if err != nil {
				return nil, err
			}
			result[i] = v
		}
		return result, nil
	}

	// string
	if strings.HasPrefix(arg, "\"") && strings.HasSuffix(arg, "\"") && len(arg) >= 2 {
		return strings.Trim(arg, "\""), nil
	}
	// char
	if strings.HasPrefix(arg, "'") && strings.HasSuffix(arg, "'") && len(arg) == 3 {
		return rune(arg[1]), nil
	}
	// int
	if i, err := strconv.Atoi(arg); err == nil {
		return i, nil
	}
	// float
	if f, err := strconv.ParseFloat(arg, 64); err == nil {
		return f, nil
	}
	return nil, fmt.Errorf("literal inválido: %s", arg)
}

func inferType(v any) string {
	switch v.(type) {
	case int:
		return "int"
	case float64:
		return "float"
	case string:
		return "string"
	case rune:
		return "char"
	case []any:
		return "list"
	default:
		return "unknown"
	}
}

func asNums(a, b any) (float64, float64, string, error) {
	switch x := a.(type) {
	case int:
		switch y := b.(type) {
		case int:
			return float64(x), float64(y), "int", nil
		case float64:
			return float64(x), y, "float", nil
		}
	case float64:
		switch y := b.(type) {
		case int:
			return x, float64(y), "float", nil
		case float64:
			return x, y, "float", nil
		}
	}
	return 0, 0, "", errors.New("operandos no numéricos")
}

func add(a, b any) (any, error) {

	if sa, ok := a.(string); ok {
		if sb, ok := b.(string); ok {
			return sa + sb, nil
		}
	}
	x, y, kind, err := asNums(a, b)
	if err != nil {
		return nil, err
	}
	if kind == "int" {
		return int(x + y), nil
	}
	return x + y, nil
}
func sub(a, b any) (any, error) {
	x, y, kind, err := asNums(a, b)
	if err != nil {
		return nil, err
	}
	if kind == "int" {
		return int(x - y), nil
	}
	return x - y, nil
}
func mul(a, b any) (any, error) {
	x, y, kind, err := asNums(a, b)
	if err != nil {
		return nil, err
	}
	if kind == "int" {
		return int(x * y), nil
	}
	return x * y, nil
}
func div(a, b any) (any, error) {
	x, y, _, err := asNums(a, b)
	if err != nil {
		return nil, err
	}
	if y == 0 {
		return nil, errors.New("división por cero")
	}
	// división entera
	return int(x / y), nil
}
func mod(a, b any) (any, error) {
	ai, aok := a.(int)
	bi, bok := b.(int)
	if !aok || !bok {
		return nil, errors.New("módulo requiere int")
	}
	if bi == 0 {
		return nil, errors.New("módulo por cero")
	}
	return ai % bi, nil
}

func isTruthy(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case int:
		return t != 0
	case float64:
		return t != 0.0
	case string:
		return t != ""
	case []any:
		return len(t) > 0
	case nil:
		return false
	default:
		return true
	}
}

func compare(a, b any, op string) (bool, error) {
	// numéricos
	if x, y, _, err := asNums(a, b); err == nil {
		switch op {
		case "<":
			return x < y, nil
		case "<=":
			return x <= y, nil
		case "==":
			return x == y, nil
		case "!=":
			return x != y, nil
		case ">":
			return x > y, nil
		case ">=":
			return x >= y, nil
		}
	}
	// strings con == y !=
	if sa, ok := a.(string); ok {
		if sb, ok := b.(string); ok {
			switch op {
			case "==":
				return sa == sb, nil
			case "!=":
				return sa != sb, nil
			}
		}
	}
	return false, fmt.Errorf("compare no soportado para op %s", op)
}

// AND/OR "lógicos" por truthiness, devuelven bool
func andOp(a, b any) (any, error) { return isTruthy(a) && isTruthy(b), nil }
func orOp(a, b any) (any, error)  { return isTruthy(a) || isTruthy(b), nil }

// listas
func getSubscr(arr any, idx any) (any, error) {
	list, ok := arr.([]any)
	if !ok {
		return nil, errors.New("BINARY_SUBSCR: no es lista")
	}
	i, ok := idx.(int)
	if !ok {
		return nil, errors.New("BINARY_SUBSCR: índice no int")
	}
	if i < 0 || i >= len(list) {
		return nil, errors.New("índice fuera de rango")
	}
	return list[i], nil
}
func storeSubscr(arr any, idx any, val any) error {
	list, ok := arr.([]any)
	if !ok {
		return errors.New("STORE_SUBSCR: no es lista")
	}
	i, ok := idx.(int)
	if !ok {
		return errors.New("STORE_SUBSCR: índice no int")
	}
	if i < 0 || i >= len(list) {
		return errors.New("índice fuera de rango")
	}
	list[i] = val
	return nil
}
