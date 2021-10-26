package cpu

import (
	"fmt"
	"strings"
)

const ( //Op DataType
	D8 = iota
	D16
	R8
	UNSET
)

type OpCodeBuilder struct {
	validated    bool
	label        string
	opcode       int
	ops          []Op	
	lastDataType int	
}

func opcodesForValues(start int, step int, values ...string) map[int]string {
	m := make(map[int]string)
	s := start

	for _, v := range values {
		m[s] = v
		s += step
	}

	return m
}

func NewOpCodeBuilder(opcode int, lb string) OpCodeBuilder {
	ob := OpCodeBuilder{
		opcode: opcode,
		label: lb,
	}

	if len(OemBug) < 2 {
		OemBug = [2]IntRegistryFunc{
			ALU.funcs[NewAluFunctionKey("INC", D16)], 
			ALU.funcs[NewAluFunctionKey("DEC", D16)],
		}
	}

	return ob
}

func NewOpcodes() {

	InitializeArguments()

	OemBug = [2]IntRegistryFunc{
		ALU.funcs[NewAluFunctionKey("INC", D16)],
		ALU.funcs[NewAluFunctionKey("DEC", D16)],
	}

	opcodes := make(map[int]OpCodeBuilder)
	extOpcodes := make(map[int]OpCodeBuilder)

	regCmd(opcodes, 0x00, "NOP", "")

	//1
	for k, v := range opcodesForValues(0x01, 0x10, "BC", "DE", "HL", "SP") {
		regLoad(opcodes, k, v, "A")
	}

	//2
	for k, v := range opcodesForValues(0x02, 0x10, "(BC)", "(DE)") {
		regLoad(opcodes, k, v, "A")
	}

	//3
	for k, v := range opcodesForValues(0x03, 0x10, "BC", "DE", "HL", "SP") {
		o := regCmd(opcodes, k, "INC {}", v)
		o.Load(v)
		o.Alu("INC")
		o.Store(v)
	}

	//4
	for k, v := range opcodesForValues(0x04, 0x08, "B", "C", "D", "E", "H", "L", "(HL)", "A") {
		regCmd(opcodes, k, "INC {}", v) //Load().Alu().Store()
	}

	/*just to not alert error*/
	fmt.Printf(string(extOpcodes[0].label))

}

//opcode = key, label = value
func regCmd(cmds map[int]OpCodeBuilder, opCode int, label string, replaceValue string) OpCodeBuilder {

	if replaceValue != "" {
		label = strings.ReplaceAll(label, "{}", replaceValue)
	}

	//check if opcode already indexed
	if val, ok := cmds[opCode]; ok {
		panic(fmt.Sprintf("Opcode %x already exists: %s", opCode, val.label))
	}

	builder := OpCodeBuilder{
		opcode:    opCode,
		label:     label,
		validated: true,
	}
	cmds[opCode] = builder
	return builder
}

func regLoad(cmds map[int]OpCodeBuilder, key int, target string, source string) {
	//to do
}

func (o *OpCodeBuilder) Load(src string) {
	arg := GetArgument(src)
	o.lastDataType = arg.DataType
	o.ops = append(o.ops, NewLoadOp(arg))
}

func (o *OpCodeBuilder) Alu(op string) {
	//to do alu functions
}

func (o *OpCodeBuilder) Store(t string) {
	arg := GetArgument(t)

	if o.lastDataType == D16 && arg.Label == s_P_a16 {

	}
}

