package main

import (
	"sync"
	"github.com/shirou/gopsutil/mem"
	"errors"
)

const ERROR_MEMORY_INFO_NOT_INITIALIZED = "Cpu Info isn't initialized"

type MemoryData struct {
	sync.Mutex
	TotalMemory uint64
	AvailableMemory uint64
	UsedMemory uint64
	UsedPercent float64

	initializedInfo bool
}

func NewMemoryData() (*MemoryData, error) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	memoryData := MemoryData{}
	memoryData.TotalMemory = memory.Total
	memoryData.AvailableMemory = memory.Available
	memoryData.UsedMemory = memory.Used
	memoryData.UsedPercent = memory.UsedPercent
	memoryData.initializedInfo = true

	return &memoryData, nil
}

func (m *MemoryData) UpdateData() error {
	m.Lock()
	defer m.Unlock()

	if !m.initializedInfo {
		return errors.New(ERROR_MEMORY_INFO_NOT_INITIALIZED)
	}

	memory, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	m.AvailableMemory = memory.Available
	m.UsedMemory = memory.Used
	m.UsedPercent = memory.UsedPercent
	return nil
}