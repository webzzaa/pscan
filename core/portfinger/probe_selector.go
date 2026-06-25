package portfinger

import (
	"sort"
	"strconv"
	"strings"
)

// PortInRange 检查端口是否在指定的端口范围字符串内
// 端口范围格式: "21,22,80,1000-2000,8080"
func PortInRange(port int, portsStr string) bool {
	if portsStr == "" {
		return false
	}

	parts := strings.Split(portsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 检查是否是范围 (如 "1000-2000")
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) == 2 {
				start, err1 := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
				end, err2 := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
				if err1 == nil && err2 == nil && port >= start && port <= end {
					return true
				}
			}
		} else {
			// 单个端口
			p, err := strconv.Atoi(part)
			if err == nil && p == port {
				return true
			}
		}
	}
	return false
}

// GetProbesForPort 获取适用于指定端口的所有探测器
// 根据 Probe.Ports 字段筛选，并按 Rarity 从低到高排序
func (v *VScan) GetProbesForPort(port int) []*Probe {
	var result []*Probe

	for i := range v.Probes {
		probe := &v.Probes[i]
		// 跳过 UDP 探测器
		if probe.Protocol == "udp" {
			continue
		}
		// 检查端口是否在探测器的 ports 范围内
		if PortInRange(port, probe.Ports) {
			result = append(result, probe)
		}
	}

	// 按 Rarity 从低到高排序 (rarity 越低越优先)
	sort.Slice(result, func(i, j int) bool {
		// rarity 为 0 表示未设置，视为最低优先级 (放最后)
		ri, rj := result[i].Rarity, result[j].Rarity
		if ri == 0 {
			ri = 10
		}
		if rj == 0 {
			rj = 10
		}
		return ri < rj
	})

	return result
}

// GetSSLProbesForPort 获取适用于指定端口的 SSL 探测器
func (v *VScan) GetSSLProbesForPort(port int) []*Probe {
	var result []*Probe

	for i := range v.Probes {
		probe := &v.Probes[i]
		// 检查端口是否在探测器的 sslports 范围内
		if PortInRange(port, probe.SSLPorts) {
			result = append(result, probe)
		}
	}

	// 按 Rarity 排序
	sort.Slice(result, func(i, j int) bool {
		ri, rj := result[i].Rarity, result[j].Rarity
		if ri == 0 {
			ri = 10
		}
		if rj == 0 {
			rj = 10
		}
		return ri < rj
	})

	return result
}

// GetAllProbesSortedByRarity 获取所有 TCP 探测器，按 Rarity 排序
func (v *VScan) GetAllProbesSortedByRarity() []*Probe {
	result := make([]*Probe, 0, len(v.Probes))

	for i := range v.Probes {
		probe := &v.Probes[i]
		if probe.Protocol != "udp" {
			result = append(result, probe)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		ri, rj := result[i].Rarity, result[j].Rarity
		if ri == 0 {
			ri = 10
		}
		if rj == 0 {
			rj = 10
		}
		return ri < rj
	})

	return result
}

// FilterProbesByIntensity 根据 intensity 过滤探测器
// intensity 范围 1-9，默认 7
func FilterProbesByIntensity(probes []*Probe, intensity int) []*Probe {
	if intensity <= 0 {
		intensity = 7
	}
	if intensity > 9 {
		intensity = 9
	}

	var result []*Probe
	for _, probe := range probes {
		// rarity 为 0 表示未设置，视为 1 (最常用)
		rarity := probe.Rarity
		if rarity == 0 {
			rarity = 1
		}
		if rarity <= intensity {
			result = append(result, probe)
		}
	}
	return result
}
