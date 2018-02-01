package main

import (
	"github.com/shirou/gopsutil/disk"
)

type DiskData struct {
	DiskStats  []DiskStats
	Partitions []Partition
	Usage      []Usage
}

func NewDiskData() *DiskData {
	return &DiskData{}
}

func (d *DiskData) Update() error {
	diskStats, err := disk.IOCounters()
	if err != nil {
		return err
	}
	d.DiskStats = ioCountersStatsToDiskStats(diskStats)

	partitions, err := disk.Partitions(true)
	if err != nil {
		return err
	}
	d.Partitions = partitionStatToPartition(partitions)

	var usage []Usage
	for _, value := range d.Partitions {
		var u Usage
		u, err = getUsageByDrive(value.MountPoint)
		if err != nil {
			return err
		}
		usage = append(usage, u)
	}
	d.Usage = usage
	return nil
}

type DiskStats struct {
	Disk             string
	ReadCount        uint64
	MergedReadCount  uint64
	WriteCount       uint64
	MergedWriteCount uint64
	ReadBytes        uint64
	WriteBytes       uint64
	ReadTime         uint64
	WriteTime        uint64
	IopsIntProgress  uint64
	IoTime           uint64
	WeightedIo       uint64
}

type Partition struct {
	Device     string
	MountPoint string
	Fstype     string
	Opts       string
}

type Usage struct {
	Path        string
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}

func getUsageByDrive(drive string) (Usage, error) {
	usage, err := disk.Usage(drive + "/")
	if err != nil {
		return Usage{}, nil
	}
	return usageStatToUsage(*usage), nil
}

func usageStatToUsage(stat disk.UsageStat) Usage {
	return Usage{
		Path:        stat.Path,
		Total:       stat.Total,
		Free:        stat.Total,
		Used:        stat.Used,
		UsedPercent: stat.UsedPercent,
	}
}

func partitionStatToPartition(stats []disk.PartitionStat) []Partition {
	var partitions []Partition
	for _, value := range stats {
		var partition Partition
		partition.Device = value.Device
		partition.MountPoint = value.Mountpoint
		partition.Fstype = value.Fstype
		partition.Opts = value.Opts

		partitions = append(partitions, partition)
	}

	return partitions
}

func ioCountersStatsToDiskStats(data map[string]disk.IOCountersStat) []DiskStats {
	var diskStats []DiskStats
	for key, value := range data {
		var diskStat DiskStats
		diskStat.Disk = key
		diskStat.ReadCount = value.ReadCount
		diskStat.MergedReadCount = value.MergedReadCount
		diskStat.WriteCount = value.WriteCount
		diskStat.MergedWriteCount = value.MergedWriteCount
		diskStat.ReadBytes = value.ReadBytes
		diskStat.WriteBytes = value.WriteBytes
		diskStat.ReadTime = value.ReadTime
		diskStat.WriteTime = value.WriteTime
		diskStat.IopsIntProgress = value.IopsInProgress
		diskStat.IoTime = value.IoTime
		diskStat.WeightedIo = value.WeightedIO

		diskStats = append(diskStats, diskStat)
	}
	return diskStats
}
