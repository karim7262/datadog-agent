// +build windows

package checks

import (
	//"runtime"
	"time"

	"github.com/DataDog/datadog-agent/pkg/process/model"
	"github.com/DataDog/gopsutil/cpu"
	"github.com/DataDog/gopsutil/process"
)

func formatUser(fp *process.FilledProcess) *model.ProcessUser {
	return &model.ProcessUser{
		Name: fp.Username,
	}
}

func formatCPU(fp *process.FilledProcess, t2, t1, syst2, syst1 cpu.TimesStat) *model.CPUStat {
	//numCPU := float64(runtime.NumCPU())
	//deltaSys := float64(t2.Timestamp - t1.Timestamp)
	// under windows, utime & stime are number of 100-ns increments.  The elapsed time
	// is in nanoseconds.
	return &model.CPUStat{
		LastCpu:    t2.CPU,
		TotalPct:   float32(fp.CpuTime.User + fp.CpuTime.System),
		UserPct:    float32(fp.CpuTime.User),
		SystemPct:  float32(fp.CpuTime.System),
		NumThreads: fp.NumThreads,
		Cpus:       []*model.SingleCPUStat{},
		Nice:       fp.Nice,
		UserTime:   int64(0),
		SystemTime: int64(0),
	}
}

func calculatePct(deltaProc, deltaTime, numCPU float64) float32 {
	if deltaTime == 0 {
		return 0
	}

	// Calculates utilization split across all CPUs. A busy-loop process
	// on a 2-CPU-core system would be reported as 50% instead of 100%.
	overalPct := (deltaProc / deltaTime) * 100

	// In cases where we get values that don't make sense, clamp to (100% * number of CPUS)
	if overalPct > (numCPU * 100) {
		overalPct = numCPU * 100
	}
	return float32(overalPct)
}

func formatIO(fp *process.FilledProcess, lastIO *process.IOCountersStat, before time.Time) *model.IOStat {
	// This will be nill for Mac
	if fp.IOStat == nil {
		return &model.IOStat{}
	}

	diff := time.Now().Unix() - before.Unix()
	if before.IsZero() || diff <= 0 {
		return &model.IOStat{}
	}
	// Reading 0 as a counter means the file could not be opened due to permissions. We distinguish this from a real 0 in rates.
	var readRate float32
	readRate = -1
	if fp.IOStat.ReadCount != 0 {
		readRate = float32(fp.IOStat.ReadCount)/1000
	}
	var writeRate float32
	writeRate = -1
	if fp.IOStat.WriteCount != 0 {
		writeRate = float32(fp.IOStat.WriteCount)/1000
	}
	var readBytesRate float32
	readBytesRate = -1
	if fp.IOStat.ReadBytes != 0 {
		readBytesRate = float32(fp.IOStat.ReadBytes)/1000
	}
	var writeBytesRate float32
	writeBytesRate = -1
	if fp.IOStat.WriteBytes != 0 {
		writeBytesRate = float32(fp.IOStat.WriteBytes)/1000
	}
	return &model.IOStat{
		ReadRate:       readRate,
		WriteRate:      writeRate,
		ReadBytesRate:  readBytesRate,
		WriteBytesRate: writeBytesRate,
	}
}