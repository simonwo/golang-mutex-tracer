package muxtracer

// supports Mutex and RWMutex, timers with nanosecond precision

import (
	"sync"
	"sync/atomic"
)

type Mutex struct {
	lock sync.Mutex

	// internal trace fields
	threshold        uint64 // 0 when disabled, else threshold in nanoseconds
	beginAwaitLock   uint64 // start time in unix nanoseconds from start waiting for lock
	beginAwaitUnlock uint64 // start time in unix nanoseconds from start waiting for unlock
	lockObtained     uint64 // once we've entered the lock in unix nanoseconds
}

func (m *Mutex) Lock() {
	tracingThreshold := m.isTracing()
	if tracingThreshold != 0 {
		m.traceBeginAwaitLock()
	}

	// actual lock
	m.lock.Lock()

	if tracingThreshold != 0 {
		m.traceEndAwaitLock(tracingThreshold)
	}
}

func (m *Mutex) Unlock() {
	tracingThreshold := m.isTracing()
	if tracingThreshold != 0 {
		m.traceBeginAwaitUnlock()
	}

	// unlock
	m.lock.Unlock()

	if tracingThreshold != 0 {
		m.traceEndAwaitUnlock(tracingThreshold)
	}
}

type TraceLocker interface {
	EnableTracer()
	DisableTracer()
	EnableTracerWithOpts(o Opts)
}

func (m *Mutex) EnableTracer() {
	m.EnableTracerWithOpts(obtainGlobalOpts())
}

func (m *Mutex) EnableTracerWithOpts(o Opts) {
	atomic.StoreUint64(&m.threshold, uint64(o.Threshold.Nanoseconds()))
}

func (m *Mutex) DisableTracer() {
	atomic.StoreUint64(&m.threshold, 0)
}
