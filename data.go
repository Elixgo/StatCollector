package main

import (
	"encoding/json"
)

type Data struct {
	Config *Config
	MemoryData *MemoryData
	CpuData *CpuData
	DiskData *DiskData
}

func NewData(config *Config) (*Data, error) {
	var data Data
	data.Config = config

	memoryData, err := NewMemoryData()
	memoryData.UpdateData()
	if err != nil {
		return nil, err
	}
	data.MemoryData = memoryData

	cpuData := CpuData{}
	err = cpuData.GetCpuInfo()
	if err != nil {
		return nil, err
	}
	cpuData.UpdateCpuData()
	data.CpuData = &cpuData

	diskData := NewDiskData()
	data.DiskData = diskData
	err = data.DiskData.Update()
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (d Data) UpdateData() {
	d.CpuData.GetCpuUsage()
	d.CpuData.GetCpuTimeData()
	d.MemoryData.UpdateData()
	d.DiskData.Update()
}

func (d Data) ToJsonString() (string, error) {
	jsonData, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}