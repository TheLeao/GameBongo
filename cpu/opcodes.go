package cpu

import "fmt"

type OpCodeBuilder struct {
	valid bool
	label string //change ?^
	opcode Opcode 
}

func opcodesForValues(start int, step int, values...string) map[int]string {
	m := make(map[int]string)
	s := start

	for _, v := range values {
		m[s] = v
		s += step
	}

	return m
}

func NewOpCodeBuilder() *OpCodeBuilder {
	return &OpCodeBuilder {

	}
}

func NewOpcodes() {

	opcodes := make([]OpCodeBuilder, 0x100)
	extOpcodes := make([]OpCodeBuilder, 0x100)

	regCmd(opcodes, 0x00, "NOP")
	
	for k, v := range opcodesForValues(0x01, 0x10, "BC", "DE", "HL", "SP") {
		
	}


}

//opcode = key, label = value
func regCmd(cmds []OpCodeBuilder, opCode int, label string) {
	if !cmds[opCode].valid {
		panic(fmt.Sprintf("Opcode %x already exists: %s", opCode, cmds[opCode].label))
	}
	
	//builder := OpCodeBuilder {}
}

func regLoad(opcodes []OpCodeBuilder, key int, target string, source string) {
	//to do
}