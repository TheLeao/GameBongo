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

type Opcode struct {
	value  int
	label  string
	ops    []Op
	length int
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
	return ob
}

func NewOpcodes() {

	InitializeArguments()

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
		o.ops = append(o.ops, NewStoreA160Op1(arg))
		o.ops = append(o.ops, NewStoreA160Op2(arg))
	} else if o.lastDataType == arg.DataType {
		o.ops = append(o.ops, NewStoreLastDataTypeOp(arg))
	} else {
		panic(fmt.Sprintf("Can't write %d to %s", o.lastDataType, t))
	}
}

func (o *OpCodeBuilder) ProceedIf(cond string) {
	o.ops = append(o.ops, NewProceedIfOp(cond))
}

func (o *OpCodeBuilder) Push() {
	dec := ALU.GetFunction("DEC", D16)

	o.ops = append(o.ops, NewPushOp1(dec))
	o.ops = append(o.ops, NewPushOp2(dec))
}

func (o *OpCodeBuilder) Pop() {
	inc := ALU.GetFunction("INC", D16)
	o.lastDataType = D16
	o.ops = append(o.ops, NewPopOp1(inc))
	o.ops = append(o.ops, NewPopOp2(inc))
}

func (o *OpCodeBuilder) ExtraCycle() {
	o.ops = append(o.ops, NewExtraCycleOp())
}

func (o *OpCodeBuilder) Alu1(operation string, arg string) {
	a := GetArgument(arg)
	fn := ALU.GetBiIntFunction(operation, o.lastDataType, a.DataType)
	o.ops = append(o.ops, NewAluOp1(fn, a, operation, o.lastDataType))

	if o.lastDataType == D16 {
		o.ExtraCycle()
	}
}

func (o *OpCodeBuilder) Alu2(operation string, d8Val int) {
	f := ALU.GetBiIntFunction(operation, o.lastDataType, D8)
	o.ops = append(o.ops, NewAluOp2(f, operation, d8Val))

	if o.lastDataType == D16 {
		o.ExtraCycle()
	}
}

func (o *OpCodeBuilder) Alu3(operation string) {
	f := ALU.GetFunction(operation, o.lastDataType)
	o.ops = append(o.ops, NewAluOp3(f, operation, o.lastDataType))

	if o.lastDataType == D16 {
		o.ExtraCycle()
	}
}

func (o *OpCodeBuilder) AluHL(operation string) {
	o.Load("HL")
	o.ops = append(o.ops, NewAluHlOp(ALU.GetFunction(operation, D16)))
	o.Store("HL")
}

func (o *OpCodeBuilder) BitHL(bit int) {
	o.ops = append(o.ops, NewBitHlOp(bit))
}

func (o *OpCodeBuilder) ClearZ() {
	o.ops = append(o.ops, NewClearZOp())
}

func (o *OpCodeBuilder) SwitchInterrupts(enable bool, withDelay bool) {
	o.ops = append(o.ops, NewSwitchInterruptsOp(enable, withDelay))
}

func (o *OpCodeBuilder) ForceFinish() {
	o.ops = append(o.ops, NewForceFinishOp())
}

//Todo create Opcode from OpCodeBuilder
func (o *OpCodeBuilder) BuildOpcode() Opcode {
	return Opcode{}
}

func (o *OpCodeBuilder) GetString() string {
	return o.label
}