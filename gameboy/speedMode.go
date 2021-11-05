package gameboy

type SpeedMode struct {
	CurrentSpeed    bool
	PrepSpeedSwitch bool
}

const (
	MINSPEED = 0x007E // 126
)

func (*SpeedMode) Accepts(addr int) bool {
	return addr == 0xff4d
}

func (s *SpeedMode) SetByte(addr int, value int) {
	s.PrepSpeedSwitch = value != 0
}

func (s *SpeedMode) GetByte(addr int) int {
	var byteInt int = MINSPEED //126
	if s.CurrentSpeed {
		byteInt += 128
	}
	if s.PrepSpeedSwitch {
		byteInt += 1
	}

	return byteInt
}

func (s *SpeedMode) OnStop() bool {
	if s.PrepSpeedSwitch {
		return false
	}

	s.CurrentSpeed = !s.CurrentSpeed
	s.PrepSpeedSwitch = false
	return true
}

func (s *SpeedMode) GetSpeedMode() int {
	if s.CurrentSpeed {
		return 2
	} else {
		return 1
	}
}