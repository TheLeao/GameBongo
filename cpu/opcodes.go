package cpu

import (
	"fmt"
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
		o := regCmd(opcodes, k, "INC " + v)
		o.Load(v)
		o.Alu3("INC")
		o.Store(v)
	}

	//4
	for k, v := range opcodesForValues(0x04, 0x08, "B", "C", "D", "E", "H", "L", "(HL)", "A") {
		o := regCmd(opcodes, k, "INC " + v)
		o.Load(v)
		o.Alu3("INC")
		o.Store(v)
	}

	//5
	for k, v := range opcodesForValues(0x05, 0x08, "B", "C", "D", "E", "H", "L", "(HL)", "A") {
		o := regCmd(opcodes, k, "DEC " + v)
		o.Load(v)
		o.Alu3("DEC")
		o.Store(v)
	}

	//6
	for k, v := range opcodesForValues(0x06, 0x08, "B", "C", "D", "E", "H", "L", "(HL)", "A") {
		regLoad(opcodes, k, v, s_d8)
	}

	//7
	for k, v := range opcodesForValues(0x07, 0x08, "RLC", "RRC", "RL", "RR") {
		o := regCmd(opcodes, k, v + "A")
		o.Load("A")
		o.Alu3(v)
		o.ClearZ()
		o.Store("A")
	}

	regLoad(opcodes, 0x08, s_P_a16, s_SP)

	//8
	for k, v := range opcodesForValues(0x09, 0x10, "BC", "DE", "HL", "SP") {
		o := regCmd(opcodes, k, "ADD HL," + v)
		o.Load("HL")
		o.Alu1("ADD", v)
		o.Store("HL")
	}

	//9
	for k, v := range opcodesForValues(0x0a, 0x10, "(BC)", "(DE)") {
		regLoad(opcodes, k, "A", v)
	}

	//10
	for k, v := range opcodesForValues(0x0b, 0x10, "BC", "DE", "HL", "SP") {
		o := regCmd(opcodes, k, "DEC " + v)
		o.Load(v)
		o.Alu3("DEC")
		o.Store(v)
	}

	regCmd(opcodes, 0x10, "STOP")
	x := regCmd(opcodes, 0x18, "JR r8")
	x.Load(s_PC)
	x.Alu3("DEC")
	x.Store(s_PC)

	//11
	for k, v := range opcodesForValues(0x20, 0x08, "NZ", "Z", "NC", "C") {
		o := regCmd(opcodes, k, fmt.Sprintf("JR %s,r8", v))
		o.Load(s_PC)
		o.ProceedIf(v)
		o.Alu1("ADD", s_r8)
		o.Store(s_PC)
	}

	x = regCmd(opcodes, 0x22, "LD (HL+),A")
	x.CopyByte("(HL)", "A")
	x.AluHL("INC")

	x = regCmd(opcodes, 0x2a, "LD A,(HL+)")
	x.CopyByte("A", "(HL)")
	x.AluHL("INC")

	x = regCmd(opcodes, 0x27, "DAA")
	x.Load("A")
	x.Alu3("DAA")
	x.Store("A")

	x = regCmd(opcodes, 0x2f, "CPL")
	x.Load("A")
	x.Alu3("CPL")
	x.Store("A")

	x = regCmd(opcodes, 0x32, "LD (HL-),A")
	x.CopyByte("(HL)", "A")
	x.AluHL("DEC")

	x = regCmd(opcodes, 0x3a, "LD A,(HL-)")
	x.CopyByte("A", "(HL)")
	x.AluHL("DEC")

	x = regCmd(opcodes, 0x37, "SCF")
	x.Load("A")
	x.Alu3("SCF")
	x.Store("A")

	x = regCmd(opcodes, 0x3f, "CCF")
	x.Load("A")
	x.Alu3("CCF")
	x.Store("A")

	//12
	for k, v := range opcodesForValues(0x40, 0x08, "B", "C", "D", "E", "H", "L", "(HL)", "A") {
		for key, val := range opcodesForValues(k, 0x01, "B", "C", "D", "E", "H", "L", "(HL)", "A"){
			if key == 0x76 {
				continue
			}

			regLoad(opcodes, key, v, val)
		}
	}

	regCmd(opcodes, 0x76, "HALT")

	//13
	for k, v := range opcodesForValues(0x80, 0x08, "ADD", "ADC", "SUB", "SBC", "AND", "XOR", "OR", "CP") {
		for ik, iv := range opcodesForValues(k, 0x01, "B", "C", "D", "E", "H", "L", "(HL)", "A") {
			o := regCmd(opcodes, ik, v + " " + iv)
			o.Load("A")
			o.Alu1(v, iv)
			o.Store("A")
		}
	}

	//14
	for k, v := range opcodesForValues(0xc0, 0x08, "NZ", "Z", "NC", "C") {
		o := regCmd(opcodes, k, "RET " + v)
		o.ExtraCycle()
		o.ProceedIf(v)
		o.Pop()
		o.ForceFinish()
		o.Store(s_PC)
	}

	//15
	for k, v := range opcodesForValues(0xc1, 0x10, "BC", "DE", "HL", "AF") {
		o := regCmd(opcodes, k, "POP " + v)
		o.Pop()
		o.Store(v)
	}

	//16
	


	/*just to not alert error*/
	fmt.Printf(string(extOpcodes[0].label))
}

//opcode = key, label = value
func regCmd(cmds map[int]OpCodeBuilder, opCode int, label string) OpCodeBuilder {
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

func (o *OpCodeBuilder) CopyByte(t string, src string) {
	o.Load(src)
	o.Store(t)
}

func (o *OpCodeBuilder) Load(src string) {
	arg := GetArgument(src)
	o.lastDataType = arg.DataType
	o.ops = append(o.ops, NewLoadOp(arg))
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