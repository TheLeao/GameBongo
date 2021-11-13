package cpu

import (
	"fmt"
	"github.com/theleao/goingboy/core"
	"github.com/theleao/goingboy/gpu"
)

type Cpu struct {
	clockCycle  int
	haltBugMode bool
	State       int
	CrrOpCode   Opcode
	opCode1     int
	opCode2     int
	operand     [2]int
	ops         []Op
	oprndIndex  int
	opIndex     int
	opCntxt     int
	Addrs       core.AddressSpace
	intrpt      core.Interrupter
	speedMode   core.SpeedMode
	intrFlag    int
	intrEnabled int
	gpu         gpu.Gpu
	Regs        Registers
	display     gpu.Display
	rqstIrq     int
	cmds        []Opcode
	extCmds     []Opcode
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

func NewCpu(addr core.AddressSpace, intrptr core.Interrupter, g gpu.Gpu, d gpu.Display, spd core.SpeedMode) Cpu {

	InitializeAlu()
	InitializeArguments()
	opCmds, opExtCmds := NewOpcodes()

	return Cpu{
		Addrs:     addr,
		intrpt:    intrptr,
		rqstIrq:   -1,
		cmds:      opCmds,
		extCmds:   opExtCmds,
		display:   d,
		gpu:       g,
		speedMode: spd,
	}
}

func (c *Cpu) Tick() {

	c.clockCycle++
	speed := c.speedMode.GetSpeedMode()

	if c.clockCycle >= (4 / speed) {
		c.clockCycle = 0
	} else {
		return
	}

	//checking interruptions
	if c.State == OPCODE || c.State == HALTED || c.State == STOPPED {
		//finish this
		if c.intrpt.Ime && (c.intrpt.InterruptEnabled != 0 && c.intrpt.InterruptFlag != 0) {
			if c.State == STOPPED {
				c.display.EnableLcd()
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
		if (c.intrFlag & c.intrEnabled) != 0 {
			c.State = OPCODE
		}
	}

	if c.State == HALTED || c.State == STOPPED {
		return
	}

	memoryAccessed := false

	for {
		pc := c.Regs.Pc

		switch c.State {
		case OPCODE:
			c.ClearState()
			c.opCode1 = c.Addrs.GetByte(pc)
			memoryAccessed = true
			if c.opCode1 == 0xcb {
				c.State = EXT_OPCODE
			} else if c.opCode1 == 0x10 {
				c.CrrOpCode = c.cmds[c.opCode1]
				c.State = EXT_OPCODE
			} else {
				c.State = OPERAND
				c.CrrOpCode = c.cmds[c.opCode1]
				if c.CrrOpCode.label == "" {
					panic(fmt.Sprintf("No command for 0x%02x", c.opCode1)) //get from java....
				}
			}

			if !c.haltBugMode {
				c.Regs.incrementPC()
			} else {
				c.haltBugMode = false
			}
		case EXT_OPCODE:
			if memoryAccessed {
				return
			}

			memoryAccessed = true
			c.opCode2 = c.Addrs.GetByte(pc)
			c.CrrOpCode = c.extCmds[c.opCode2]

			c.State = OPERAND
			c.Regs.incrementPC()

		case OPERAND:
			for {
				if c.oprndIndex >= c.CrrOpCode.length {
					break
				}

				if memoryAccessed {
					return
				}

				memoryAccessed = true
				c.operand[c.oprndIndex] = c.Addrs.GetByte(pc)
				c.oprndIndex++
				c.Regs.incrementPC()
			}

			c.ops = c.CrrOpCode.ops
			c.State = RUNNING

		case RUNNING:
			if c.opCode1 == 0x10 {
				if c.speedMode.OnStop() {
					c.State = OPCODE
				} else {
					c.State = STOPPED
					c.display.DisableLcd()
				}
			} else if c.opCode1 == 0x76 {
				if c.intrpt.IsHaltBug() {
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
				hasCorruption, corruptionType := op.CausesOemBug(&c.Regs, c.opCntxt)
				if hasCorruption {
					if !c.gpu.Lcdc.IsLcdEnabled() {
						return
					}

					//GPU Stat register
					stat := c.Addrs.GetByte(0xff41)
					if (stat&0b11) == gpu.OAMSEARCH && c.gpu.TicksInLine < 79 {
						gpu.CorruptOam(&c.Addrs, corruptionType, c.gpu.TicksInLine)
					}
				}

				c.opCntxt = op.Execute(&c.Regs, c.Addrs, c.operand[:], c.opCntxt)
				op.SwitchInterrupts(&c.intrpt)

				if !op.Proceed(c.Regs) {
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
			c.intrpt.ClearInterrupt(c.rqstIrq)
			c.intrpt.DisableInterrupts(false)
		}
	case IRQ_PUSH_1:
		c.Regs.decrementSP()
		c.Addrs.SetByte(c.Regs.Sp, (c.Regs.Pc&0xff00)>>8)
		c.State = IRQ_PUSH_2
	case IRQ_PUSH_2:
		c.Regs.decrementSP()
		c.Addrs.SetByte(c.Regs.Sp, c.Regs.Pc&0x00ff)
		c.State = IRQ_JUMP
	case IRQ_JUMP:
		c.Regs.Pc = c.rqstIrq
		c.rqstIrq = -1
		c.State = OPCODE
	}
}

func (c *Cpu) ClearState() {
	c.opCode1 = 0
	c.opCode2 = 0
	c.CrrOpCode = Opcode{}
	c.ops = nil

	c.operand[0] = 0
	c.operand[1] = 0
	c.oprndIndex = 0
	c.opIndex = 0
	c.opCntxt = 0
	c.intrFlag = 0
	c.intrEnabled = 0
}
