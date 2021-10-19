package cpu

type SpeedMode struct {
	currentSpeed    bool
	prepSpeedSwitch bool
}

const (
	MINSPEED = 0x007E // 126
)

func (*SpeedMode) Accepts(addr int) bool {
	return addr == 0xff4d
}

func (s *SpeedMode) SetByte(addr int, value int) {
	s.prepSpeedSwitch = value != 0
}

func (s *SpeedMode) GetByte(addr int) int {
	var byteInt int = MINSPEED //126
	if s.currentSpeed {
		byteInt += 128
	}
	if s.prepSpeedSwitch {
		byteInt += 1
	}

	return byteInt
}
