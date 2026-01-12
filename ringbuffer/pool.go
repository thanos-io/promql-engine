// Copyright (c) The Thanos Community Authors.
// Licensed under the Apache License 2.0.

package ringbuffer

// BufferPool manages a pool of reusable ring buffers for memory efficiency.
// Buffers are pre-allocated and accessed via round-robin indexing.
type BufferPool struct {
	// Pre-allocated buffers for deterministic behavior
	buffers []Buffer
	size    int
}

// NewBufferPool creates a new buffer pool with the specified size.
// The factory function is called to create new buffers.
func NewBufferPool(size int, factory func() Buffer) *BufferPool {
	if size <= 0 {
		size = 1
	}
	
	buffers := make([]Buffer, size)
	for i := range buffers {
		buffers[i] = factory()
	}
	
	return &BufferPool{
		buffers: buffers,
		size:    size,
	}
}

// GetBuffer returns a buffer for the given index.
// Uses modulo to map any index to the pool size.
func (p *BufferPool) GetBuffer(index int) Buffer {
	return p.buffers[index%p.size]
}

// Size returns the number of buffers in the pool.
func (p *BufferPool) Size() int {
	return p.size
}
