package cpu

import (
	"fmt"

	"github.com/theleao/gamebongo/gameboy"
	"github.com/theleao/gamebongo/gpu"
)

type Cpu struct {
	clockCycle   int
	haltBugMode  bool
	state        int
	crrOpCode    *Opcode
	opCode1      int
	opCode2      int
	operand      [2]int
	ops          []int
	operandIndex int
	opIndex      int
	opCntxt      int
	Addrs        gameboy.AddressSpace
	intrpt       Interrupter
	speedMode    SpeedMode
	intrFlag     int
	intrEnabled  int
	gpu          gpu.Gpu
}

type Opcode struct {
	value  int
	label  string
	ops    []int
	length int
}

const (
	OPCODE = iota
	EXT_OPCODE
	OPERAND
	RUNNING
	IRQ_READ_IF
	IRQ_READ_IE
	IRQ_PUSH_1
	IRQ_PUSH_2
	IRQ_JUMP
	STOPPED
	HALTED
)

func NewCpu(addr gameboy.AddressSpace, intrptr Interrupter) Cpu {
	return Cpu{
		Addrs:   addr,
		intrpt:  intrptr,
		opCode1: 10,
	}
}

func NewCpuTest() Cpu {
	return Cpu{
		crrOpCode: &Opcode{
			value: 99,
			label: "Moscau",
		},
		speedMode: SpeedMode{
			currentSpeed:    true,
			prepSpeedSwitch: true,
		},
	}
}

func LittleTest() {
	c := NewCpuTest()

	fmt.Println(c.crrOpCode.label)
	fmt.Println("Changing")
	c.crrOpCode.label = "Lalalal"
	fmt.Println(c.crrOpCode.label)

	c.speedMode.currentSpeed = false
	fmt.Println(c.speedMode.currentSpeed)
}

func (c *Cpu) Tick() {

	c.clockCycle++
	speed := getSpeed()

	if c.clockCycle >= (4 / speed) {
		c.clockCycle = 0
	} else {
		return
	}

	//checking interruptions
	if c.state == OPCODE || c.state == HALTED || c.state == STOPPED {
		//finish this
		if c.intrpt.ime && (c.intrpt.interruptEnabled != 0 && c.intrpt.interruptFlag != 0) {
			if c.state == STOPPED {
				// c#: _display.Enabled = true;
			}

			c.state = IRQ_READ_IF
		}
	}

	switch c.state {
	case IRQ_READ_IF:
	case IRQ_READ_IE:
	case IRQ_PUSH_1:
	case IRQ_PUSH_2:
	case IRQ_JUMP:
		handleInterrupt()
		return
	case HALTED:
		if c.intrpt.interruptEnabled != 0 && c.intrpt.interruptFlag != 0 {
			//continue switch
			c.state = OPCODE
		}
	}

	if c.state == HALTED || c.state == STOPPED {
		return
	}

	memoryAccessed := false

	for {
		var pc int = 0 //Registers.PC

		switch c.state {
		case OPCODE:
			c.clearState()
			c.opCode1 = c.Addrs.GetByte(pc)
			memoryAccessed = true
			if c.opCode1 == 0xcb {
				c.state = EXT_OPCODE
			} else if c.opCode1 == 0x10 {
				//c.crrOpCode = nil //opcodes java:Opcodes.COMMANDS.get(opcode1);
				c.state = EXT_OPCODE
			} else {
				c.state = OPERAND
				c.crrOpCode = nil //opcodes java:Opcodes.COMMANDS.get(opcode1);
				if c.crrOpCode == nil {
					panic(nil) //--exception
				}
			}

			if !c.haltBugMode {
				//java:registers.incrementPC()
			} else {
				c.haltBugMode = false
			}
		case EXT_OPCODE:
			if memoryAccessed {
				return
			}
			memoryAccessed = true
			c.opCode2 = c.Addrs.GetByte(pc)

			if c.crrOpCode == nil {
				c.crrOpCode = nil //_opcodes.ExtCommands[_opcode2];
			}
			if c.crrOpCode == nil {
				panic(nil) //exception "No command for %0xcb 0x%02x"
			}
		}

	}
}

func getSpeed() int {
	//Speed mode
	return 0
}

func handleInterrupt() {
	//TO do
}

func (c *Cpu) clearState() {
	c.opCode1 = 0
	c.opCode2 = 0
	c.crrOpCode = nil
	c.ops = nil

	c.operand[0] = 0
	c.operand[1] = 0
	c.operandIndex = 0
	c.opIndex = 0
	c.opCntxt = 0
	c.intrFlag = 0
	c.intrEnabled = 0
}
