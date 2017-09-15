package sandbox

import (
	"path/filepath"
	"os/exec"
)

const (
	cgCPUPathPrefix    = "/sys/fs/cgroup/cpu/"
	cgPidPathPrefix    = "/sys/fs/cgroup/pids/"
	cgMemoryPathPrefix = "/sys/fs/cgroup/memory/"
)

func InitCGroup(pid, containerID, memory string) error {
	if err := cpuCGroup(pid, containerID); err != nil {
		return err
	}

	if err := pidCGroup(pid, containerID); err != nil {
		return err
	}

	if err := memoryCGroup(pid, containerID, memory); err != nil {
		return err
	}

	return nil
}

// https://www.kernel.org/doc/Documentation/scheduler/sched-bwc.txt
func cpuCGroup(pid, containerID string) error {
	cgCPUPath := filepath.Join(cgCPUPathPrefix, containerID)

	// add sub cgroup system
	mkdirCmd := exec.Command("/usr/bin/mkdir", cgCPUPath)
	if err := mkdirCmd.Run(); err != nil {
		return err
	}

	// add current pid to cgroup cpu
	taskCmd := exec.Command("/usr/bin/echo", pid, ">", filepath.Join(cgCPUPath, "/tasks"))
	if err := taskCmd.Run(); err != nil {
		return err
	}

	// limit a group to 2% of 1 CPU
	quotaCmd := exec.Command("/usr/bin/echo", "2000", ">", filepath.Join(cgCPUPath, "/cpu.cfs_quota_us"))
	if err := quotaCmd.Run(); err != nil {
		return err
	}

	return nil
}

// https://www.kernel.org/doc/Documentation/cgroup-v1/pids.txt
func pidCGroup(pid, containerID string) error {
	cgPidPath := filepath.Join(cgPidPathPrefix, containerID)

	// add sub cgroup system
	mkdirCmd := exec.Command("/usr/bin/mkdir", cgPidPath)
	if err := mkdirCmd.Run(); err != nil {
		return err
	}

	// add current pid to cgroup pids
	taskCmd := exec.Command("/usr/bin/echo", pid, ">", filepath.Join(cgPidPath, "/cgroup.procs"))
	if err := taskCmd.Run(); err != nil {
		return err
	}

	// max pids limitation
	quotaCmd := exec.Command("/usr/bin/echo", "2", ">", filepath.Join(cgPidPath, "/pids.max"))
	if err := quotaCmd.Run(); err != nil {
		return err
	}

	return nil
}

// https://www.kernel.org/doc/Documentation/cgroup-v1/memory.txt
func memoryCGroup(pid, containerID, memory string) error {
	cgMemoryPath := filepath.Join(cgMemoryPathPrefix, containerID)

	// add sub cgroup system
	mkdirCmd := exec.Command("/usr/bin/mkdir", cgMemoryPath)
	if err := mkdirCmd.Run(); err != nil {
		return err
	}

	// add current pid to cgroup memory
	taskCmd := exec.Command("/usr/bin/echo", pid, ">", filepath.Join(cgMemoryPath, "/tasks"))
	if err := taskCmd.Run(); err != nil {
		return err
	}

	// set memory usage limitation
	memoryCmd := exec.Command("/usr/bin/echo", memory+"m", ">", filepath.Join(cgMemoryPath, "/memory.limit_in_bytes"))
	if err := memoryCmd.Run(); err != nil {
		return err
	}

	swappinessCmd := exec.Command("/usr/bin/echo", "0", ">", filepath.Join(cgMemoryPath, "/memory.swappiness"))
	if err := swappinessCmd.Run(); err != nil {
		return err
	}

	kernelMemoryCmd := exec.Command("/usr/bin/echo", memory+"m", ">", filepath.Join(cgMemoryPath, "/memory.kmem.limit_in_bytes"))
	if err := kernelMemoryCmd.Run(); err != nil {
		return err
	}

	return nil
}

func CleanupCGroup(containerID string) error {
	cleanCPUCommand := exec.Command("rmdir", filepath.Join(cgCPUPathPrefix, containerID))
	if err := cleanCPUCommand.Run(); err != nil {
		return err
	}

	cleanPidCommand := exec.Command("rmdir", filepath.Join(cgPidPathPrefix, containerID))
	if err := cleanPidCommand.Run(); err != nil {
		return err
	}

	cleanMemoryCommand := exec.Command("rmdir", filepath.Join(cgMemoryPathPrefix, containerID))
	if err := cleanMemoryCommand.Run(); err != nil {
		return err
	}

	return nil
}
