package cpu

import (
	"fmt"

	"github.com/theleao/gamebongo/gameboy"
	"github.com/theleao/gamebongo/gpu"
)

type Cpu struct {
	clockCycle  int
	haltBugMode bool
	State       int
	crrOpCode   Opcode
	opCode1     int
	opCode2     int
	operand     [2]int
	ops         []Op
	oprndIndex  int
	opIndex     int
	opCntxt     int
	Addrs       gameboy.AddressSpace
	intrpt      Interrupter
	speedMode   SpeedMode
	intrFlag    int
	intrEnabled int
	gpu         gpu.Gpu
	reg         Registers
	display     gpu.Display
	rqstIrq  	int
	cmds 	    []Opcode
	extCmds 	[]Opcode
}

type Opcode struct {
	value  int
	label  string
	ops    []Op
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
		rqstIrq: -1,
	}
}

func NewCpuTest() Cpu {
	return Cpu{
		crrOpCode: Opcode{
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
	if c.State == OPCODE || c.State == HALTED || c.State == STOPPED {
		//finish this
		if c.intrpt.ime && (c.intrpt.interruptEnabled != 0 && c.intrpt.interruptFlag != 0) {
			if c.State == STOPPED {
				// c#: _display.Enabled = true;
			}

			c.State = IRQ_READ_IF
		}
	}

	switch c.State {
	case IRQ_READ_IF:
	case IRQ_READ_IE:
	case IRQ_PUSH_1:
	case IRQ_PUSH_2:
	case IRQ_JUMP:
		c.handleInterrupt()
		return
	case HALTED:
		if c.intrpt.interruptEnabled != 0 && c.intrpt.interruptFlag != 0 {
			//continue switch
			c.State = OPCODE
		}
	}

	if c.State == HALTED || c.State == STOPPED {
		return
	}

	memoryAccessed := false

	for {
		var pc int = 0 //Registers.PC

		switch c.State {
		case OPCODE:
			c.clearState()
			c.opCode1 = c.Addrs.GetByte(pc)
			memoryAccessed = true
			if c.opCode1 == 0xcb {
				c.State = EXT_OPCODE
			} else if c.opCode1 == 0x10 {
				c.crrOpCode = c.cmds[c.opCode1] //opcodes java:Opcodes.COMMANDS.get(opcode1);
				c.State = EXT_OPCODE
			} else {
				c.State = OPERAND
				c.crrOpCode = c.cmds[c.opCode1] //opcodes java:Opcodes.COMMANDS.get(opcode1);
				if c.crrOpCode == nil {
					panic(fmt.Sprintf("No command for OpCode 1 : %x", c.opCode1))
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
				c.crrOpCode = c.extCmds[c.opCode2] 
				//_opcodes.ExtCommands[_opcode2];
			}
			if c.crrOpCode == nil {
				panic(fmt.Sprintf("No command for OpCode 2 : %x", c.opCode2))
			}

			c.State = OPERAND
			c.reg.incrementPC()

		case OPERAND:
			for ok := true; ok; ok = (c.oprndIndex < c.crrOpCode.length) {
				if memoryAccessed {
					return
				}

				c.oprndIndex++
				c.operand[c.oprndIndex] = c.Addrs.GetByte(pc)
				c.reg.incrementPC()
			}

			c.ops = c.crrOpCode.ops
			c.State = RUNNING

		case RUNNING:
			if c.opCode1 == 0x10 {
				if c.speedMode.onStop() {
					c.State = OPCODE
				} else {
					c.State = STOPPED
					c.display.DisableLcd()
				}
			} else if c.opCode1 == 0x76 {
				if c.intrpt.isHaltBug() {
					c.State = OPCODE
					c.haltBugMode = true
					return
				} else {
					c.State = HALTED
					return
				}
			}

			if c.opIndex < len(c.ops) {
				var op Op = c.ops[c.opIndex]
				var opMemoryAccessed bool = op.ReadsMemory() || op.WritesMemory()

				if opMemoryAccessed && memoryAccessed {
					return
				}
				c.opIndex++

				//handle sprite bug
				hasCorruption, corruptionType := op.CausesOemBug(c.reg, c.opCntxt)
				if hasCorruption {
					if !c.gpu.Lcdc.Enabled {
						return
					}

					//GPU Stat register
					stat := c.Addrs.GetByte(0xff41)
					if (stat&0b11) == gpu.OAMSEARCH && c.gpu.TicksInLine < 79 {
						gpu.CorruptOam(&c.Addrs, corruptionType, c.gpu.TicksInLine)
					}
				}

				c.opCntxt = op.Execute(c.reg, c.Addrs, c.operand, c.opCntxt)
				op.SwitchInterrupts(c.intrpt)

				if !op.Proceed(c.reg) {
					c.opIndex = len(c.ops)
					break
				}

				if op.ForceFinishCycle() {
					return
				}

				if opMemoryAccessed {
					memoryAccessed = true
				}
			}

			if c.opIndex >= len(c.ops) {
				c.State = OPCODE
				c.oprndIndex = 0
				c.intrpt.OnInstructionFinished()
				return
			}
			break

		case HALTED:
		case STOPPED:
			return
		}
	}
}

func getSpeed() int {
	//Speed mode
	return 0
}

func (c *Cpu) handleInterrupt() {
	//TO do
	switch c.State {
	case IRQ_READ_IF:
		c.intrFlag = c.Addrs.GetByte(0xff0f)
		c.State = IRQ_READ_IE
	case IRQ_READ_IE:
		c.intrEnabled = c.Addrs.GetByte(0xffff)
		c.rqstIrq = -1

		for i := 0; i < 5; i++ {
			if (c.intrFlag & c.intrEnabled & (1 << i)) != 0 {
				c.rqstIrq = i
				break
			}
		}

		if c.rqstIrq == -1 {
			c.State = OPCODE
		} else {
			c.State = IRQ_PUSH_1
			c.intrpt.clearInterrupt(c.rqstIrq)
			c.intrpt.disableInterrupts(false)
		}
	case IRQ_PUSH_1:
		c.reg.decrementSP()
		c.Addrs.SetByte(c.reg.Sp, (c.reg.Pc & 0xff00) >> 8)
		c.State = IRQ_PUSH_2
	case IRQ_PUSH_2:
		c.reg.decrementSP()
		c.Addrs.SetByte(c.reg.Sp, c.reg.Pc & 0x00ff)
		c.State = IRQ_JUMP
	case IRQ_JUMP:
		c.reg.Pc = c.rqstIrq
		c.rqstIrq = -1
		c.State = OPCODE
	}
}

func (c *Cpu) clearState() {
	c.opCode1 = 0
	c.opCode2 = 0
	c.crrOpCode = new(Opcode)
	c.ops = nil

	c.operand[0] = 0
	c.operand[1] = 0
	c.oprndIndex = 0
	c.opIndex = 0
	c.opCntxt = 0
	c.intrFlag = 0
	c.intrEnabled = 0
}
