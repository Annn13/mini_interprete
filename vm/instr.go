package vm

type Instr struct {
	PC     int
	Op     string
	Arg    string
	HasArg bool
}
