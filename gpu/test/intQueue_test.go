package gpu

import (
	"github.com/theleao/goingboy/core"
	"github.com/theleao/goingboy/gpu"
	"testing"
)

type testQueue struct {
	queue gpu.DmgPixelQueue
}

func newTestQueue() testQueue {
	var mr []core.MemRegisterType
	gpuRegs := gpu.GpuRegisters()
	for _, r := range gpuRegs {
		mr = append(mr, core.MemRegisterType{
			Addr: r,
		})
	}

	memRegs := core.NewMemoryRegisters(mr...)
	bgp, _ := gpu.GetGpuRegister(gpu.BGP)
	memRegs.Put(bgp, 0b11100100)

	return testQueue{
		queue: gpu.DmgPixelQueue{
			Display: &gpu.NullDisplay{},
			Regs:    memRegs,
		},
	}
}

func TestEnqueue(t *testing.T) {
	q := newTestQueue()
	var x [8]int
	z := zip(0b11001001, 0b11110000, false, x)
	q.queue.Enqueue8Pixels(z[:], 0)

	list := []int{3, 3, 2, 2, 1, 0, 0, 1}
	queueAsList := arrayQueueAsList(q.queue.Pixels)

	v := true
	for _, i := range list {
		if queueAsList[i] != i {
			v = false
			break
		}
	}

	if v {
		t.Errorf("Lists are not equal")
	}
}

//Fetcher Zip
func zip(dt1 int, dt2 int, reverse bool, pxLine [8]int) [8]int {
	for i := 7; i >= 0; i-- {
		mask := 1 << i

		p := 0
		if dt2&mask != 0 {
			p = 2
		}
		if dt1&mask != 0 {
			p += 1
		}

		if reverse {
			pxLine[i] = p
		} else {
			pxLine[7-i] = p
		}
	}

	return pxLine
}

func arrayQueueAsList(q gpu.IntQueue) []int {
	s := len(q)
	l := make([]int, s)
	for i := 0; i < s; i++ {
		l = append(l, q[i])
	}
	return l
}
