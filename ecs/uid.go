package ecs

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

const (
	CEpoch         = 1474802888000
	CWorkerIdBits  = 10
	CSenquenceBits = 12

	CWorkerIdShift  = 12
	CTimeStampShift = 22

	CSequenceMask = 0xfff
	CMaxWorker    = 0x3ff
)
var uidGenerator *UIDGenerator

func init()  {
	uidGenerator = NewUIDGenerator()
}

func NextUID()int64 {
	var uid int64
	var err error
	for{
		uid,err =uidGenerator.NextId()
		if err==nil {
			return uid
		}
	}
}

type UIDGenerator struct {
	rand          int64
	lastTimeStamp int64
	sequence      int64
	lock          *sync.Mutex
}

func NewUIDGenerator() (iw *UIDGenerator) {
	iw = new(UIDGenerator)
	iw.rand = rand.Int63n(getMaxWorkerId())
	iw.lastTimeStamp = -1
	iw.sequence = 0
	iw.lock = new(sync.Mutex)
	return iw
}

func getMaxWorkerId() int64 {
	return -1 ^ -1<<CWorkerIdBits
}

func getSequenceMask() int64 {
	return -1 ^ -1<<CSenquenceBits
}

func (iw *UIDGenerator) timeGen() int64 {
	return time.Now().UnixNano() / 1000 / 1000
}

func (iw *UIDGenerator) timeReGen(last int64) int64 {
	ts := time.Now().UnixNano() / 1000 / 1000
	for {
		if ts <= last {
			ts = iw.timeGen()
		} else {
			break
		}
	}
	return ts
}

func (iw *UIDGenerator) NextId() (ts int64, err error) {
	iw.lock.Lock()
	defer iw.lock.Unlock()
	ts = iw.timeGen()
	if ts == iw.lastTimeStamp {
		iw.sequence = (iw.sequence + 1) & CSequenceMask
		if iw.sequence == 0 {
			ts = iw.timeReGen(ts)
		}
	} else {
		iw.sequence = 0
	}

	if ts < iw.lastTimeStamp {
		err = errors.New("clock moved backwards, refuse gen id")
		return 0, err
	}
	iw.lastTimeStamp = ts
	ts = (ts-CEpoch)<<CTimeStampShift | iw.rand<<CWorkerIdShift | iw.sequence
	return ts, nil
}

func ParseId(id int64) (t time.Time, ts int64, workerId int64, seq int64) {
	seq = id & CSequenceMask
	workerId = (id >> CWorkerIdShift) & CMaxWorker
	ts = (id >> CTimeStampShift) + CEpoch
	t = time.Unix(ts/1000, (ts%1000)*1000000)
	return
}