package atomic

import "sync/atomic"

// implement atomic<bool>
type Boolean uint32

func (b *Boolean) Get() bool {
	return atomic.LoadUint32((*uint32)(b)) != 0
}

func (b *Boolean) Set(v bool) {
	if v {
		atomic.StoreUint32((*uint32)(b), 1)
	} else {
		atomic.StoreUint32((*uint32)(b), 0)
	}
}
