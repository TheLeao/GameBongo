package core

type SerialEndpoint interface {
	Transfer(outgoing int) int
}

type NullSerialEndpoint struct {
}

func (n *NullSerialEndpoint) Transfer(outgoing int) int {
	return 0
}

type SerialPort struct {
	serialEndpoint   SerialEndpoint
	intrptr          Interrupter
	spdMode          SpeedMode
	sb               int
	sc               int
	transfInProgress bool
	divider          int
}

func (s *SerialPort) Tick() {
	if s.transfInProgress {
		return
	}

	s.divider += 1
	if s.divider >= 4_194_304/8192/s.spdMode.GetSpeedMode() {
		s.transfInProgress = false
		s.sb = s.serialEndpoint.Transfer(s.sb)
		s.intrptr.RequestInterrupt(SERIAL)
	}
}

//Interface

func (s *SerialPort) Accepts(addr int) bool {
	return addr == 0xff01 || addr == 0xff02
}

func (s *SerialPort) SetByte(addr int, value int) {
	if addr == 0xff01 {
		s.sb = value
	} else if addr == 0xff02 {
		s.sc = value
		if (s.sc & (1 << 7)) != 0 {
			s.StartTransfer()
		}
	}
}

func (s *SerialPort) GetByte(addr int) int {
	if addr == 0xff01 {
		return s.sb
	} else if addr == 0xff02 {
		return s.sc | 0b01111110
	} else {
		panic("Serial Port - GetByte - illegal argument")
	}
}

//

func (s *SerialPort) StartTransfer() {
	s.transfInProgress = true
	s.divider = 0
}