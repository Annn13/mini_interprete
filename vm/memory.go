package vm

type Var struct {
	Nombre string
	Tipo   string // "int", "float", "string", "char", "list"
	Valor  any
}

type Memory struct{ m map[string]Var }

func NewMemory() *Memory { return &Memory{m: map[string]Var{}} }

func (mem *Memory) Set(name string, val any) {
	mem.m[name] = Var{Nombre: name, Tipo: inferType(val), Valor: val}
}

func (mem *Memory) Get(name string) (any, bool) {
	v, ok := mem.m[name]
	if !ok {
		return nil, false
	}
	return v.Valor, true
}
