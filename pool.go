// Package pool provides optimized memory pooling for Lux projects
// Inspired by VictoriaMetrics memory management strategies
package pool

import (
	"sync"
	"unsafe"
)

// ByteSlicePool provides optimized byte slice pooling
// Similar to VictoriaMetrics' approach for reducing allocations
var ByteSlicePool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 1024)
	},
}

// GetByteSlice gets a byte slice from the pool
func GetByteSlice() []byte {
	return ByteSlicePool.Get().([]byte)
}

// PutByteSlice returns a byte slice to the pool
func PutByteSlice(b []byte) {
	if cap(b) > 1024*1024 { // Don't pool large slices
		return
	}
	ByteSlicePool.Put(b[:0]) // Reset to zero length but keep capacity
}

// StringPool provides string pooling for common strings
// Uses unsafe conversion to avoid allocations
var StringPool = sync.Pool{
	New: func() interface{} {
		return ""
	},
}

// GetString gets a string from the pool
func GetString() string {
	return StringPool.Get().(string)
}

// PutString returns a string to the pool
func PutString(s string) {
	if len(s) > 1024 { // Don't pool large strings
		return
	}
	StringPool.Put(s)
}

// InterfaceSlicePool provides pooling for interface slices
// Useful for JSON marshaling/unmarshaling
var InterfaceSlicePool = sync.Pool{
	New: func() interface{} {
		return make([]interface{}, 0, 16)
	},
}

// GetInterfaceSlice gets an interface slice from the pool
func GetInterfaceSlice() []interface{} {
	return InterfaceSlicePool.Get().([]interface{})
}

// PutInterfaceSlice returns an interface slice to the pool
func PutInterfaceSlice(s []interface{}) {
	if cap(s) > 1024 { // Don't pool large slices
		return
	}
	InterfaceSlicePool.Put(s[:0])
}

// MapPool provides pooling for maps
// Specialized for common map types
var MapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]interface{}, 16)
	},
}

// GetMap gets a map from the pool
func GetMap() map[string]interface{} {
	return MapPool.Get().(map[string]interface{})
}

// PutMap returns a map to the pool
func PutMap(m map[string]interface{}) {
	if len(m) > 1024 { // Don't pool large maps
		return
	}
	// Clear the map before returning to pool
	for k := range m {
		delete(m, k)
	}
	MapPool.Put(m)
}

// FastBuffer provides a zero-allocation buffer
// Similar to VictoriaMetrics' buffer optimizations
type FastBuffer struct {
	buf []byte
}

// NewFastBuffer creates a new FastBuffer
func NewFastBuffer() *FastBuffer {
	return &FastBuffer{
		buf: GetByteSlice(),
	}
}

// Write implements io.Writer interface
func (b *FastBuffer) Write(p []byte) (n int, err error) {
	b.buf = append(b.buf, p...)
	return len(p), nil
}

// Bytes returns the underlying byte slice
func (b *FastBuffer) Bytes() []byte {
	return b.buf
}

// String returns the buffer as a string without allocation
func (b *FastBuffer) String() string {
	return *(*string)(unsafe.Pointer(&b.buf))
}

// Reset clears the buffer
func (b *FastBuffer) Reset() {
	b.buf = b.buf[:0]
}

// Free returns the buffer to the pool
func (b *FastBuffer) Free() {
	PutByteSlice(b.buf)
	b.buf = nil
}

// GetFastBuffer gets a FastBuffer from the pool
var fastBufferPool = sync.Pool{
	New: func() interface{} {
		return NewFastBuffer()
	},
}

func GetFastBuffer() *FastBuffer {
	return fastBufferPool.Get().(*FastBuffer)
}

// PutFastBuffer returns a FastBuffer to the pool
func PutFastBuffer(b *FastBuffer) {
	b.Reset()
	fastBufferPool.Put(b)
}