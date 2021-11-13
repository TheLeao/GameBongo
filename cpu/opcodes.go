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

func (o Opcode) IsEmpty() bool {
	return o.value == 0 && o.label == "" && o.ops == nil && o.length == 0
}

func (o Opcode) GetString() string {
	return fmt.Sprintf("%02x %s", o.value, o.label)
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

func NewOpcodes() ([]Opcode, []Opcode) {

	//opcodes := make(map[int]OpCodeBuilder)
	//extOpcodes := make(map[int]OpCodeBuilder)

	opcodes := make([]OpCodeBuilder, 0x100)
	extOpcodes := make([]OpCodeBuilder, 0x100)

	regCmd(opcodes, 0x00, "NOP")

	//1
	for k, v := range opcodesForValues(0x01, 0x10, "BC", "DE", "HL", "SP") {
		opcodes[k] = regLoad(opcodes, k, v, "d16")
	}

	//2
	for k, v := range opcodesForValues(0x02, 0x10, "(BC)", "(DE)") {
		opcodes[k] = regLoad(opcodes, k, v, "A")
	}

	//3
	for k, v := range opcodesForValues(0x03, 0x10, "BC", "DE", "HL", "SP") {
		o := regCmd(opcodes, k, "INC "+v)
		o.Load(v)
		o.Alu3("INC")
		o.Store(v)
		opcodes[k] = o
	}

	//4
	for k, v := range opcodesForValues(0x04, 0x08, "B", "C", "D", "E", "H", "L", "(HL)", "A") {
		o := regCmd(opcodes, k, "INC "+v)
		o.Load(v)
		o.Alu3("INC")
		o.Store(v)
		opcodes[k] = o
	}

	//5
	for k, v := range opcodesForValues(0x05, 0x08, "B", "C", "D", "E", "H", "L", "(HL)", "A") {
		o := regCmd(opcodes, k, "DEC "+v)
		o.Load(v)
		o.Alu3("DEC")
		o.Store(v)
		opcodes[k] = o
	}

	//6
	for k, v := range opcodesForValues(0x06, 0x08, "B", "C", "D", "E", "H", "L", "(HL)", "A") {
		opcodes[k] = regLoad(opcodes, k, v, s_d8)
	}

	//7
	for k, v := range opcodesForValues(0x07, 0x08, "RLC", "RRC", "RL", "RR") {
		o := regCmd(opcodes, k, v+"A")
		o.Load("A")
		o.Alu3(v)
		o.ClearZ()
		o.Store("A")
		opcodes[k] = o
	}

	opcodes[0x08] = regLoad(opcodes, 0x08, s_P_a16, s_SP)

	//8
	for k, v := range opcodesForValues(0x09, 0x10, "BC", "DE", "HL", "SP") {
		o := regCmd(opcodes, k, "ADD HL,"+v)
		o.Load("HL")
		o.Alu1("ADD", v)
		o.Store("HL")
		opcodes[k] = o
	}

	//9
	for k, v := range opcodesForValues(0x0a, 0x10, "(BC)", "(DE)") {
		opcodes[k] = regLoad(opcodes, k, "A", v)
	}

	//10
	for k, v := range opcodesForValues(0x0b, 0x10, "BC", "DE", "HL", "SP") {
		o := regCmd(opcodes, k, "DEC "+v)
		o.Load(v)
		o.Alu3("DEC")
		o.Store(v)
		opcodes[k] = o
	}

	opcodes[0x10] = regCmd(opcodes, 0x10, "STOP")

	x := regCmd(opcodes, 0x18, "JR r8")
	x.Load(s_PC)
	x.Alu3("DEC")
	x.Store(s_PC)
	opcodes[0x18] = x

	//11
	for k, v := range opcodesForValues(0x20, 0x08, "NZ", "Z", "NC", "C") {
		o := regCmd(opcodes, k, fmt.Sprintf("JR %s,r8", v))
		o.Load(s_PC)
		o.ProceedIf(v)
		o.Alu1("ADD", s_r8)
		o.Store(s_PC)
		opcodes[k] = o
	}

	x = regCmd(opcodes, 0x22, "LD (HL+),A")
	x.CopyByte("(HL)", "A")
	x.AluHL("INC")
	opcodes[0x22] = x

	x = regCmd(opcodes, 0x2a, "LD A,(HL+)")
	x.CopyByte("A", "(HL)")
	x.AluHL("INC")
	opcodes[0x2a] = x

	x = regCmd(opcodes, 0x27, "DAA")
	x.Load("A")
	x.Alu3("DAA")
	x.Store("A")
	opcodes[0x27] = x

	x = regCmd(opcodes, 0x2f, "CPL")
	x.Load("A")
	x.Alu3("CPL")
	x.Store("A")
	opcodes[0x2f] = x

	x = regCmd(opcodes, 0x32, "LD (HL-),A")
	x.CopyByte("(HL)", "A")
	x.AluHL("DEC")
	opcodes[0x32] = x

	x = regCmd(opcodes, 0x3a, "LD A,(HL-)")
	x.CopyByte("A", "(HL)")
	x.AluHL("DEC")
	opcodes[0x3a] = x

	x = regCmd(opcodes, 0x37, "SCF")
	x.Load("A")
	x.Alu3("SCF")
	x.Store("A")
	opcodes[0x37] = x

	x = regCmd(opcodes, 0x3f, "CCF")
	x.Load("A")
	x.Alu3("CCF")
	x.Store("A")
	opcodes[0x3f] = x

	//12
	for k, v := range opcodesForValues(0x40, 0x08, "B", "C", "D", "E", "H", "L", "(HL)", "A") {
		for key, val := range opcodesForValues(k, 0x01, "B", "C", "D", "E", "H", "L", "(HL)", "A") {
			if key == 0x76 {
				continue
			}

			opcodes[key] = regLoad(opcodes, key, v, val)
		}
	}

	opcodes[0x76] = regCmd(opcodes, 0x76, "HALT")

	//13
	for k, v := range opcodesForValues(0x80, 0x08, "ADD", "ADC", "SUB", "SBC", "AND", "XOR", "OR", "CP") {
		for ik, iv := range opcodesForValues(k, 0x01, "B", "C", "D", "E", "H", "L", "(HL)", "A") {
			o := regCmd(opcodes, ik, v+" "+iv)
			o.Load("A")
			o.Alu1(v, iv)
			o.Store("A")
			opcodes[ik] = o
		}
	}

	//14
	for k, v := range opcodesForValues(0xc0, 0x08, "NZ", "Z", "NC", "C") {
		o := regCmd(opcodes, k, "RET "+v)
		o.ExtraCycle()
		o.ProceedIf(v)
		o.Pop()
		o.ForceFinish()
		o.Store(s_PC)
		opcodes[k] = o
	}

	//15
	for k, v := range opcodesForValues(0xc1, 0x10, "BC", "DE", "HL", "AF") {
		o := regCmd(opcodes, k, "POP "+v)
		o.Pop()
		o.Store(v)
		opcodes[k] = o
	}

	//16
	for k, v := range opcodesForValues(0xc2, 0x08, "NZ", "Z", "NC", "C") {
		o := regCmd(opcodes, k, "JP "+v+",a16")
		o.Load(s_a16)
		o.ProceedIf(v)
		o.Store(s_PC)
		o.ExtraCycle()
		opcodes[k] = o
	}

	x = regCmd(opcodes, 0xc3, "JP a16")
	x.Load(s_a16)
	x.Store(s_PC)
	x.ExtraCycle()
	opcodes[0xc3] = x

	//17
	for k, v := range opcodesForValues(0xc4, 0x08, "NZ", "Z", "NC", "C") {
		o := regCmd(opcodes, k, "CALL "+v+",a16")
		o.ProceedIf(v)
		o.ExtraCycle()
		o.Load(s_PC)
		o.Push()
		o.Load(s_a16)
		opcodes[k] = o
	}

	//18
	for k, v := range opcodesForValues(0xc5, 0x10, "BC", "DE", "HL", "AF") {
		o := regCmd(opcodes, k, "PUSH "+v)
		o.ExtraCycle()
		o.Load(v)
		o.Push()
		opcodes[k] = o
	}

	//19
	for k, v := range opcodesForValues(0xc6, 0x08, "ADD", "ADC", "SUB", "SBC", "AND", "XOR", "OR", "CP") {
		o := regCmd(opcodes, k, v+" d8")
		o.Load("A")
		o.Alu1(v, s_d8)
		o.Store("A")
		opcodes[k] = o
	}

	//20
	j := 0x00
	for i := 0xc7; i <= 0xf7; i += 0x10 {
		o := regCmd(opcodes, i, fmt.Sprintf("RST %02xH", j))
		o.Load(s_PC)
		o.Push()
		o.ForceFinish()
		o.LoadWord(j)
		o.Store(s_PC)
		opcodes[i] = o
		j += 0x10
	}

	x = regCmd(opcodes, 0xc9, "RET")
	x.Pop()
	x.ForceFinish()
	x.Store(s_PC)
	opcodes[0xc9] = x

	x = regCmd(opcodes, 0xcd, "CALL a16")
	x.Load(s_PC)
	x.ExtraCycle()
	x.Push()
	x.Load(s_a16)
	x.Store(s_PC)
	opcodes[0xcd] = x

	//21
	j = 0x08
	for i := 0xcf; i <= 0xff; i += 0x10 {
		o := regCmd(opcodes, i, "RST %02xH")
		o.Load(s_PC)
		o.Push()
		o.ForceFinish()
		o.LoadWord(j)
		o.Store(s_PC)
		opcodes[i] = o
		j += 0x10
	}

	x = regCmd(opcodes, 0xd9, "RETI")
	x.Pop()
	x.ForceFinish()
	x.Store(s_PC)
	x.SwitchInterrupts(true, false)
	opcodes[0xd9] = x

	opcodes[0xe2] = regLoad(opcodes, 0xe2, s_P_C, "A")
	opcodes[0xf2] = regLoad(opcodes, 0xf2, "A", s_P_C)

	x = regCmd(opcodes, 0xe9, "JP (HL)")
	x.Load(s_HL)
	x.Store(s_PC)
	opcodes[0xe9] = x

	x = regCmd(opcodes, 0xe0, "LDH (a8),A")
	x.CopyByte("(a8)", "A")
	opcodes[0xe0] = x

	x = regCmd(opcodes, 0xf0, "LDH A,(a8)")
	x.CopyByte("A", "(a8)")
	opcodes[0xf0] = x

	x = regCmd(opcodes, 0xe8, "ADD SP,r8")
	x.Load(s_SP)
	x.Alu1("ADD_SP", "r8")
	x.ExtraCycle()
	x.Store(s_SP)
	opcodes[0xe8] = x

	x = regCmd(opcodes, 0xf8, "LD HL,SP+r8")
	x.Load(s_SP)
	x.Alu1("ADD_SP", "r8")
	x.Store(s_HL)
	opcodes[0xf8] = x

	opcodes[0xea] = regLoad(opcodes, 0xea, s_P_a16, "A")
	opcodes[0xfa] = regLoad(opcodes, 0xfa, "A", s_P_a16)

	x = regCmd(opcodes, 0xf3, "DI")
	x.SwitchInterrupts(false, true)
	opcodes[0xf3] = x

	x = regCmd(opcodes, 0xfb, "EI")
	x.SwitchInterrupts(true, true)
	opcodes[0xfb] = x

	x = regLoad(opcodes, 0xf9, s_SP, s_HL)
	x.ExtraCycle()
	opcodes[0xf9] = x

	//22
	for k, v := range opcodesForValues(0x00, 0x08, "RLC", "RRC", "RL", "RR", "SLA", "SRA", "SWAP", "SRL") {
		for k2, v2 := range opcodesForValues(k, 0x01, "B", "C", "D", "E", "H", "L", "(HL)", "A") {
			o := regCmd(extOpcodes, k2, v+" "+v2)
			o.Load(v2)
			o.Alu3(v)
			o.Store(v2)
			extOpcodes[k2] = o
		}
	}

	//23
	for k, v := range opcodesForValues(0x40, 0x40, "BIT", "RES", "SET") {
		for b := 0; b < 0x08; b++ {
			for k2, v2 := range opcodesForValues(k+(b*0x08), 0x01, "B", "C", "D", "E", "H", "L", "(HL)", "A") {
				if "BIT" == v && s_HL == v2 {
					o := regCmd(extOpcodes, k2, fmt.Sprintf("BIT %d,(HL)", b))
					o.BitHL(b)
					extOpcodes[k2] = o
				} else {
					o := regCmd(extOpcodes, k2, fmt.Sprintf("%s %d,%s", v, b, v2))
					o.Load(v2)
					o.Alu2(v, b)
					o.Store(v2)
					extOpcodes[k2] = o
				}
			}
		}
	}

	var commands []Opcode
	var extCommands []Opcode

	for _, ob := range opcodes {
		commands = append(commands, ob.NewOpcode())
	}

	for _, ob := range extOpcodes {
		extCommands = append(extCommands, ob.NewOpcode())
	}

	return commands, extCommands
}

//opcode = key, label = value
func regCmd(cmds []OpCodeBuilder, opCode int, label string) OpCodeBuilder {
	//check if opcode already indexed
	if cmds[opCode].label != "" {
		panic(fmt.Sprintf("Opcode %x already exists: %s", opCode, cmds[opCode].label))
	}

	builder := OpCodeBuilder{
		opcode:    opCode,
		label:     label,
		validated: true,
	}
	cmds[opCode] = builder
	return builder
}

func regLoad(cmds []OpCodeBuilder, key int, target string, source string) OpCodeBuilder {
	o := regCmd(cmds, key, "LD "+target+","+source)
	o.CopyByte(target, source)
	return o
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

	if fn == nil {
		panic("smart panic")
	}

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
	o.ops = append(o.ops, NewAluHlOp(ALU.GetFunction(operation, D16), operation))
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

func (o *OpCodeBuilder) LoadWord(v int) {
	o.lastDataType = D16
	o.ops = append(o.ops, NewLoadWordOp(v))
}

// create Opcode from OpCodeBuilder
func (o *OpCodeBuilder) NewOpcode() Opcode {
	opcode := Opcode{
		label: o.label,
		value: o.opcode,
		ops:   o.ops,
	}

	if len(opcode.ops) <= 0 {
		opcode.length = 0
	} else {
		var max, temp int
		for _, opcop := range opcode.ops {
			if opcop.OperandLength() > temp {
				temp = opcop.OperandLength()
				max = temp
			}
		}
		opcode.length = max
	}

	return opcode
}

func (o *OpCodeBuilder) GetString() string {
	return o.label
}
