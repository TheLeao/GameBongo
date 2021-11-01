package cpu

type SpeedMode struct {
	currentSpeed    bool
	prepSpeedSwitch bool
}

const (
	MINSPEED = 0x007E // 126
)

func (*SpeedMode) accepts(addr int) bool {
	return addr == 0xff4d
}

func (s *SpeedMode) setByte(addr int, value int) {
	s.prepSpeedSwitch = value != 0
}

func (s *SpeedMode) getByte(addr int) int {
	var byteInt int = MINSPEED //126
	if s.currentSpeed {
		byteInt += 128
	}
	if s.prepSpeedSwitch {
		byteInt += 1
	}

	return byteInt
}

func (s *SpeedMode) onStop() bool {
	if s.prepSpeedSwitch {
		return false
	}

	s.currentSpeed = !s.currentSpeed
	s.prepSpeedSwitch = false
	return true
}

func (s *SpeedMode) GetSpeedMode() int {
	if s.currentSpeed {
		return 2
	} else {
		return 1
	}
}