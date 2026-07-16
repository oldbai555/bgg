package logic

import (
	"context"
	"runtime"
	"time"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/zeromicro/go-zero/core/logx"
)

type MonitorStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMonitorStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MonitorStatusLogic {
	return &MonitorStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// MonitorStatus 读的是 iam-rpc 进程自己所在机器/容器的资源，拆分前是和 gateway 同一进程，
// 现在是有意的语义变化（见 iam.proto 里 MonitorStatusResponse 的注释）。
func (l *MonitorStatusLogic) MonitorStatus(in *iam.Empty) (*iam.MonitorStatusResponse, error) {
	resp := &iam.MonitorStatusResponse{}

	if percentages, cores, err := getCPUInfo(l.ctx); err != nil {
		l.Errorf("获取CPU信息失败: %v", err)
		resp.CpuCores = int32(runtime.NumCPU())
	} else {
		resp.CpuUsage = percentages
		resp.CpuCores = int32(cores)
	}

	if vmStat, err := mem.VirtualMemoryWithContext(l.ctx); err != nil {
		l.Errorf("获取内存信息失败: %v", err)
	} else {
		resp.MemoryTotal = vmStat.Total
		resp.MemoryUsed = vmStat.Used
		resp.MemoryAvailable = vmStat.Available
		resp.MemoryUsage = vmStat.UsedPercent
	}

	diskStat, err := disk.UsageWithContext(l.ctx, "/")
	if err != nil {
		diskStat, err = disk.UsageWithContext(l.ctx, "C:")
	}
	if err != nil {
		l.Errorf("获取磁盘信息失败: %v", err)
	} else {
		resp.DiskTotal = diskStat.Total
		resp.DiskUsed = diskStat.Used
		resp.DiskAvailable = diskStat.Free
		resp.DiskUsage = diskStat.UsedPercent
	}

	if ioCounters, err := net.IOCountersWithContext(l.ctx, false); err != nil {
		l.Errorf("获取网络信息失败: %v", err)
	} else {
		for _, counter := range ioCounters {
			resp.NetworkBytesSent += counter.BytesSent
			resp.NetworkBytesRecv += counter.BytesRecv
			resp.NetworkPacketsSent += counter.PacketsSent
			resp.NetworkPacketsRecv += counter.PacketsRecv
		}
	}

	return resp, nil
}

func getCPUInfo(ctx context.Context) (usage float64, cores int, err error) {
	percentages, err := cpu.PercentWithContext(ctx, 1*time.Second, false)
	if err != nil {
		return 0, 0, err
	}
	if len(percentages) > 0 {
		usage = percentages[0]
	}
	cores, err = cpu.Counts(true)
	if err != nil {
		cores = runtime.NumCPU()
	}
	return usage, cores, nil
}
