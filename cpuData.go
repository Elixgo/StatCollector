package main

import (
	"github.com/shirou/gopsutil/cpu"
	"errors"
	"sync"
)

const ERROR_CPU_INFO_NOT_INITIALIZED = "Cpu Info isn't initialized"

type CpuData struct {
	sync.Mutex
	CpuInfo []CpuInfo
	initializedInfo bool

	CpuUsage []float64
	CpuTimeData CpuTimeData
}

func (c CpuData) CpuCount() int {
	return len(c.CpuInfo)
}

func (c CpuData) UpdateCpuData() {
	c.GetCpuUsage()
	c.GetCpuTimeData()
}

type CpuInfo struct {
	VendorId string
	Family string
	PhysicalId string
	Cores int32
	ModelName string
	Mhz float64
}

type CpuTime struct {
	Cpu string
	User float64
	System float64
	Idle float64
	Nice float64
	IOWait float64
	Irq float64
	Softirq float64
	Steal float64
	Guest float64
	GuestNice float64
	Stolen float64
}

type CpuTimeData struct {
	Total CpuTime
	TotalCpuOne CpuTime
	TotalCpuTwo CpuTime
	CoreData []CpuTime `json:"coreData"`
}

func (c *CpuData) GetCpuInfo() error {
	c.Lock()
	defer c.Unlock()

	cpuData, err := cpu.Info()
	if err != nil {
		return err
	}

	var cpus []CpuInfo
	for _, value := range cpuData {
		var cpu CpuInfo
		cpu.VendorId = value.VendorID
		cpu.Family = value.Family
		cpu.PhysicalId = value.PhysicalID
		cpu.Cores = value.Cores
		cpu.ModelName = value.ModelName
		cpu.Mhz = value.Mhz

		cpus = append(cpus, cpu)
	}
	c.CpuInfo = cpus
	c.initializedInfo = true
	return nil
}

func (c *CpuData) GetCpuUsage() error {
	c.Lock()
	defer c.Unlock()

	if !c.initializedInfo {
		return errors.New(ERROR_CPU_INFO_NOT_INITIALIZED)
	}

	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		return err
	}
	c.CpuUsage = cpuUsage
	return nil
}

func timesStatToCpuTime(stat cpu.TimesStat) CpuTime {
	return CpuTime{
		Cpu: stat.CPU,
		User: stat.User,
		System: stat.System,
		Idle: stat.Idle,
		Nice: stat.Nice,
		IOWait: stat.Iowait,
		Irq: stat.Irq,
		Softirq: stat.Softirq,
		Steal: stat.Steal,
		Guest: stat.Guest,
		GuestNice: stat.GuestNice,
		Stolen: stat.Stolen,
	}
}

func (c *CpuData) GetCpuTimeData() error {
	c.Lock()
	defer c.Unlock()
	if !c.initializedInfo {
		return errors.New(ERROR_CPU_INFO_NOT_INITIALIZED)
	}

	cpuTimes, err := cpu.Times(true)
	if err != nil {
		return err
	}

	var TotalTime CpuTime
	var TotalTimeCpuOne CpuTime
	var TotalTimeCpuTwo CpuTime

	var coreData []CpuTime

	for _, value := range cpuTimes {
		if value.CPU == "_Total" {
			TotalTime = timesStatToCpuTime(value)
		} else if value.CPU == "0,_Total" {
			TotalTimeCpuOne = timesStatToCpuTime(value)
		} else if value.CPU == "1,_Total" {
			TotalTimeCpuTwo = timesStatToCpuTime(value)
		} else {
			coreData = append(coreData, timesStatToCpuTime(value))
		}
	}

	c.CpuTimeData.Total = TotalTime
	c.CpuTimeData.TotalCpuOne = TotalTimeCpuOne
	c.CpuTimeData.TotalCpuTwo = TotalTimeCpuTwo
	c.CpuTimeData.CoreData = coreData
	return nil
}
