package cpu

import (
	"fmt"
	"testing"

	"github.com/theleao/goingboy/core"
	"github.com/theleao/goingboy/cpu"
	"github.com/theleao/goingboy/gpu"
)

type CpuTest struct {
	cpu cpu.Cpu
	mem core.AddressSpace
}

func TestTiming(t *testing.T) {
	m := core.NewRam(0x00, 0x10000)

	ct := CpuTest{
		cpu: cpu.NewCpu(&m, core.NewInterrupter(false), gpu.Gpu{}, &gpu.NullDisplay{}, core.SpeedMode{}),
		mem: &m,
	}

	//RET
	x := ct.assertTiming(16, 0xc9, 0, 0)
	if x != nil {
		t.Errorf(x.Error())
	}

	//RETI
	x = ct.assertTiming(16, 0xd9, 0, 0)
	if x != nil {
		t.Errorf(x.Error())
	}

	//RET NZ
	ct.cpu.Regs.Flags.SetZ(false)
	x = ct.assertTiming(20, 0xc0, 0, 0)
	if x != nil {
		t.Errorf(x.Error())
	}

	//RET NZ
	ct.cpu.Regs.Flags.SetZ(true)
	x = ct.assertTiming(8, 0xc0, 0, 0)
	if x != nil {
		t.Errorf(x.Error())
	}

	//CALL a16
	x = ct.assertTiming(24, 0xcd, 0, 0)
	if x != nil {
		t.Errorf(x.Error())
	}

	//PUSH BC
	x = ct.assertTiming(16, 0xc5)
	if x != nil {
		t.Errorf(x.Error())
	}

	//POP AF
	x = ct.assertTiming(12, 0xf1)
	if x != nil {
		t.Errorf(x.Error())
	}

	//SUB A, d8
	x = ct.assertTiming(8, 0xd6, 0)
	if x != nil {
		t.Errorf(x.Error())
	}

	//JR nc, r8
	ct.cpu.Regs.Flags.SetC(true)
	x = ct.assertTiming(8, 0x30, 0)
	if x != nil {
		t.Errorf(x.Error())
	}

	//JR nc, r8
	ct.cpu.Regs.Flags.SetC(false)
	x = ct.assertTiming(12, 0x30, 0)
	if x != nil {
		t.Errorf(x.Error())
	}

	//JP nc, a16
	ct.cpu.Regs.Flags.SetC(true)
	x = ct.assertTiming(12, 0xd2, 0)
	if x != nil {
		t.Errorf(x.Error())
	}

	//JP nc, a16
	ct.cpu.Regs.Flags.SetC(false)
	x = ct.assertTiming(16, 0xd2, 0)
	if x != nil {
		t.Errorf(x.Error())
	}

	//JP a16
	x = ct.assertTiming(16, 0xc3, 0, 0)
	if x != nil {
		t.Errorf(x.Error())
	}

	//XOR a
	x = ct.assertTiming(4, 0xaf)
	if x != nil {
		t.Errorf(x.Error())
	}

	//LD (ff00+05),A
	x = ct.assertTiming(12, 0xe0, 0x05)
	if x != nil {
		t.Errorf(x.Error())
	}

	//LD A,(ff00+05)
	x = ct.assertTiming(12, 0xf0, 0x05)
	if x != nil {
		t.Errorf(x.Error())
	}

	//OR
	x = ct.assertTiming(4, 0xb7, 0x05)
	if x != nil {
		t.Errorf(x.Error())
	}

	//LDA A,E
	x = ct.assertTiming(4, 0x7b)
	if x != nil {
		t.Errorf(x.Error())
	}

	//SUB A,d8
	x = ct.assertTiming(8, 0xd6, 0)
	if x != nil {
		t.Errorf(x.Error())
	}

	//RL D
	x = ct.assertTiming(8, 0xcb, 0x12)
	if x != nil {
		t.Errorf(x.Error())
	}

	//ADD A
	x = ct.assertTiming(4, 0x87)
	if x != nil {
		t.Errorf(x.Error())
	}

	//DI
	x = ct.assertTiming(4, 0xf3)
	if x != nil {
		t.Errorf(x.Error())
	}

	//LD (HL-),A
	x = ct.assertTiming(8, 0x32)
	if x != nil {
		t.Errorf(x.Error())
	}

	//LD (HL),d8
	x = ct.assertTiming(12, 0x36)
	if x != nil {
		t.Errorf(x.Error())
	}

	//LD (a16),A
	x = ct.assertTiming(16, 0xea, 0, 0)
	if x != nil {
		t.Errorf(x.Error())
	}

	//ADD HL,BC
	x = ct.assertTiming(8, 0x09)
	if x != nil {
		t.Errorf(x.Error())
	}

	//RST 00H
	x = ct.assertTiming(16, 0xc7)
	if x != nil {
		t.Errorf(x.Error())
	}

	//LDA A,51
	x = ct.assertTiming(8, 0x3e, 0x51)
	if x != nil {
		t.Errorf(x.Error())
	}

	//RRA
	x = ct.assertTiming(4, 0x1f)
	if x != nil {
		t.Errorf(x.Error())
	}

	//ADC A,01
	x = ct.assertTiming(8, 0xce, 0x01)
	if x != nil {
		t.Errorf(x.Error())
	}

	//NOP
	x = ct.assertTiming(4, 0)
	if x != nil {
		t.Errorf(x.Error())
	}
}

func (c *CpuTest) assertTiming(expectedTiming int, opcodes ...int) error {
	for i := 0; i < len(opcodes); i++ {
		c.mem.SetByte(0x100+i, opcodes[i])
	}
	c.cpu.ClearState()
	c.cpu.Regs.Pc = 0x100

	ticks := 0
	var opcode cpu.Opcode = cpu.Opcode{}

	for {
		c.cpu.Tick()
		if opcode.IsEmpty() && !c.cpu.CrrOpCode.IsEmpty() {
			opcode = c.cpu.CrrOpCode
		}
		ticks++

		if c.cpu.State == cpu.OPCODE && ticks >= 4 {
			break
		}
	}

	if opcode.IsEmpty() {
		if expectedTiming != ticks {
			return fmt.Errorf("invalid timing value for %s for Expected: %d Ticks: %d", hexArray(opcodes), expectedTiming, ticks)
		} else {
			return nil
		}
	} else {
		if expectedTiming != ticks {
			return fmt.Errorf("invalid timing value for [%s] for Expected: %d Ticks: %d", opcode.GetString(), expectedTiming, ticks)
		} else {
			return nil
		}
	}

}

func hexArray(data []int) string {
	var b string = "["
	ln := len(data)
	for i := 0; i < ln; i++ {
		b += fmt.Sprintf("%02x", data[i])
		if i < ln-1 {
			b += " "
		}
	}
	b += "]"
	return b
}
