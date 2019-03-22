// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

// +build docker

package docker

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types"

	"github.com/DataDog/datadog-agent/pkg/util/containers"
	"github.com/DataDog/datadog-agent/pkg/util/containers/metrics"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

var healthRe = regexp.MustCompile(`\(health: (\w+)\)`)

// ContainerListConfig allows to pass listing options
type ContainerListConfig struct {
	IncludeExited  bool
	FlagExcluded   bool
	PopulateLimits bool
}

// Containers gets a list of all containers on the current node using a mix of
// the Docker APIs and cgroups stats. We attempt to limit syscalls where possible.
func (d *DockerUtil) ListContainers(cfg *ContainerListConfig) ([]*containers.Container, error) {
	if !d.useCgroups {
		cfg.PopulateLimits = true
	}
	cList, err := d.dockerContainers(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not get docker containers: %s", err)
	}

	if !d.useCgroups {
		return cList, d.updateContainerMetricsNoCgroups(cList)
	}

	cgByContainer, err := metrics.ScrapeAllCgroups()
	if err != nil {
		return nil, fmt.Errorf("could not get cgroups: %s", err)
	}

	for _, container := range cList {
		if container.State != containers.ContainerRunningState || container.Excluded {
			continue
		}
		cgroup, ok := cgByContainer[container.ID]
		if !ok {
			log.Debugf("No matching cgroups for container %s, skipping", container.ID[:12])
			continue
		}
		container.SetCgroups(cgroup)
		err = container.FillCgroupLimits()
		if err != nil {
			log.Debugf("Cannot get limits for container %s: %s, skipping", container.ID[:12], err)
			continue
		}
	}

	err = d.UpdateContainerMetrics(cList)
	return cList, err
}

// UpdateContainerMetrics updates cgroup / network performance metrics for
// a provided list of Container objects
func (d *DockerUtil) UpdateContainerMetrics(cList []*containers.Container) error {
	if !d.useCgroups {
		return d.updateContainerMetricsNoCgroups(cList)
	}

	for _, container := range cList {
		if container.State != containers.ContainerRunningState || container.Excluded {
			continue
		}

		err := container.FillCgroupMetrics()
		if err != nil {
			log.Debugf("Cannot get metrics for container %s: %s", container.ID[:12], err)
			continue
		}

		if d.cfg.CollectNetwork {
			d.Lock()
			networks := d.networkMappings[container.ID]
			d.Unlock()

			nwByIface := make(map[string]string)
			for _, nw := range networks {
				nwByIface[nw.iface] = nw.dockerName
			}

			err = container.FillNetworkMetrics(nwByIface)
			if err != nil {
				log.Debugf("Cannot get network stats for container %s: %s", container.ID, err)
				continue
			}
		}
	}
	return nil
}

func (d *DockerUtil) updateContainerMetricsNoCgroups(cList []*containers.Container) error {
	for _, container := range cList {
		if container.State != containers.ContainerRunningState || container.Excluded {
			continue
		}
		stats, err := d.getDockerStats(container.ID)
		if err != nil {
			log.Debugf("Cannot get metrics for container %s: %s", container.ID[:12], err)
			continue
		}

		err = fillContainerStats(container, stats)
		if err != nil {
			log.Debugf("Cannot get metrics for container %s: %s", container.ID[:12], err)
			continue
		}
	}
	return nil
}

// dockerContainers returns the running container list from the docker API
func (d *DockerUtil) dockerContainers(cfg *ContainerListConfig) ([]*containers.Container, error) {
	if cfg == nil {
		return nil, errors.New("configuration is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), d.queryTimeout)
	defer cancel()
	cList, err := d.cli.ContainerList(ctx, types.ContainerListOptions{All: cfg.IncludeExited})
	if err != nil {
		return nil, fmt.Errorf("error listing containers: %s", err)
	}
	ret := make([]*containers.Container, 0, len(cList))
	for _, c := range cList {
		if d.cfg.CollectNetwork && c.State == containers.ContainerRunningState {
			// FIXME: We might need to invalidate this cache if a containers networks are changed live.
			d.Lock()
			if _, ok := d.networkMappings[c.ID]; !ok {
				i, err := d.Inspect(c.ID, false)
				if err != nil {
					d.Unlock()
					log.Debugf("Error inspecting container %s: %s", c.ID, err)
					continue
				}
				d.networkMappings[c.ID] = findDockerNetworks(c.ID, i.State.Pid, c)
			}
			d.Unlock()
		}

		image, err := d.ResolveImageName(c.Image)
		if err != nil {
			log.Warnf("Can't resolve image name %s: %s", c.Image, err)
		}

		excluded := d.cfg.filter.IsExcluded(c.Names[0], image)
		if excluded && !cfg.FlagExcluded {
			continue
		}

		entityID := ContainerIDToEntityName(c.ID)
		container := &containers.Container{
			Type:     "Docker",
			ID:       c.ID,
			EntityID: entityID,
			Name:     c.Names[0],
			Image:    image,
			ImageID:  c.ImageID,
			Created:  c.Created,
			State:    c.State,
			Excluded: excluded,
			Health:   parseContainerHealth(c.Status),
		}

		if cfg.PopulateLimits {
			i, err := d.Inspect(c.ID, false)
			if err == nil && i.HostConfig != nil {
				container.MemLimit = uint64(i.HostConfig.Memory)
				container.SoftMemLimit = uint64(i.HostConfig.MemoryReservation)
				if (i.HostConfig.CPUPeriod > 0) && (i.HostConfig.CPUQuota > 0) {
					container.CPULimit = (float64(i.HostConfig.CPUQuota) / float64(i.HostConfig.CPUPeriod)) * 100.0
				} else {
					container.CPULimit = 100.0
				}
			}
		}

		ret = append(ret, container)
	}

	// Resolve docker networks after we've processed all containers so all
	// routing maps are available.
	if d.cfg.CollectNetwork {
		d.Lock()
		resolveDockerNetworks(d.networkMappings)
		d.Unlock()
	}

	if d.lastInvalidate.Add(invalidationInterval).After(time.Now()) {
		d.cleanupCaches(cList)
	}

	return ret, nil
}

// Parse the health out of a container status. The format is either:
//  - 'Up 5 seconds (health: starting)'
//  - 'Up about an hour'
func parseContainerHealth(status string) string {
	// Avoid allocations in most cases by just checking for '('
	if strings.IndexByte(status, '(') == -1 {
		return ""
	}
	all := healthRe.FindAllStringSubmatch(status, -1)
	if len(all) < 1 || len(all[0]) < 2 {
		return ""
	}
	return all[0][1]
}

// cleanupCaches removes cache entries for unknown containers and images
func (d *DockerUtil) cleanupCaches(containers []types.Container) {
	liveContainers := make(map[string]struct{})
	liveImages := make(map[string]struct{})
	for _, c := range containers {
		liveContainers[c.ID] = struct{}{}
		liveImages[c.Image] = struct{}{}
	}
	d.Lock()
	for cid := range d.networkMappings {
		if _, ok := liveContainers[cid]; !ok {
			delete(d.networkMappings, cid)
		}
	}
	for image := range d.imageNameBySha {
		if _, ok := liveImages[image]; !ok {
			delete(d.imageNameBySha, image)
		}
	}
	d.Unlock()
}

// fillContainerStats uses docker stats info to populate the container metrics
// For now, only fields used in the docker check are populated
func fillContainerStats(c *containers.Container, s *types.StatsJSON) error {
	if c == nil || s == nil {
		return errors.New("nil pointer input")
	}

	// CPU stats
	c.CPU = &metrics.CgroupTimesStat{
		System:     s.CPUStats.CPUUsage.UsageInKernelmode,
		User:       s.CPUStats.CPUUsage.UsageInUsermode,
		UsageTotal: float64(s.CPUStats.CPUUsage.TotalUsage),
	}
	c.CPUNrThrottled = s.CPUStats.ThrottlingData.ThrottledPeriods

	// Mem stats
	c.Memory = &metrics.CgroupMemStat{}
	for k, v := range s.MemoryStats.Stats {
		switch k {
		case "cache":
			c.Memory.Cache = v
		case "swap":
			c.Memory.Swap = v
			c.Memory.SwapPresent = true
		case "rss":
			c.Memory.RSS = v
		case "rss_huge":
			c.Memory.RSSHuge = v
		case "mapped_file":
			c.Memory.MappedFile = v
		case "pgpgin":
			c.Memory.Pgpgin = v
		case "pgpgout":
			c.Memory.Pgpgout = v
		case "pgfault":
			c.Memory.Pgfault = v
		case "pgmajfault":
			c.Memory.Pgmajfault = v
		case "inactive_anon":
			c.Memory.InactiveAnon = v
		case "active_anon":
			c.Memory.ActiveAnon = v
		case "inactive_file":
			c.Memory.InactiveFile = v
		case "active_file":
			c.Memory.ActiveFile = v
		case "unevictable":
			c.Memory.Unevictable = v
		case "hierarchical_memory_limit":
			c.Memory.HierarchicalMemoryLimit = v
		case "hierarchical_memsw_limit":
			c.Memory.HierarchicalMemSWLimit = v
		case "total_cache":
			c.Memory.TotalCache = v
		case "total_rss":
			c.Memory.TotalRSS = v
		case "total_rssHuge":
			c.Memory.TotalRSSHuge = v
		case "total_mapped_file":
			c.Memory.TotalMappedFile = v
		case "total_pgpgin":
			c.Memory.TotalPgpgIn = v
		case "total_pgpgout":
			c.Memory.TotalPgpgOut = v
		case "total_pgfault":
			c.Memory.TotalPgFault = v
		case "total_pgmajfault":
			c.Memory.TotalPgMajFault = v
		case "total_inactive_anon":
			c.Memory.TotalInactiveAnon = v
		case "total_active_anon":
			c.Memory.TotalActiveAnon = v
		case "total_inactive_file":
			c.Memory.TotalInactiveFile = v
		case "total_active_file":
			c.Memory.TotalActiveFile = v
		case "total_unevictable":
			c.Memory.TotalUnevictable = v
		}
	}
	c.MemFailCnt = s.MemoryStats.Failcnt

	// Net stats
	for net, netstats := range s.Networks {
		n := &metrics.InterfaceNetStats{
			NetworkName: net,
			BytesSent:   netstats.TxBytes,
			BytesRcvd:   netstats.RxBytes,
			PacketsSent: netstats.TxPackets,
			PacketsRcvd: netstats.RxPackets,
		}
		c.Network = append(c.Network, n)
	}

	// IO stats, only sum is exposed
	c.IO = sumBlkioStatEntry(s.BlkioStats.IoServiceBytesRecursive)

	// Thread stats
	c.ThreadCount = s.PidsStats.Current
	c.ThreadLimit = s.PidsStats.Limit

	return nil
}

func sumBlkioStatEntry(entries []types.BlkioStatEntry) *metrics.CgroupIOStat {
	s := &metrics.CgroupIOStat{}
	for _, e := range entries {
		switch e.Op {
		case "Read":
			s.ReadBytes += e.Value
		case "Write":
			s.WriteBytes += e.Value
		}
	}
	return s
}
