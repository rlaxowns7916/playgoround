package kill

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

type Opts struct {
	PID  int
	Name string
	Port int
}

func Cmd() *cobra.Command {
	opts := &Opts{}
	cmd := &cobra.Command{
		Use:   "kill",
		Short: "Kill process by selector",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			f := cmd.Flags()
			if f.Changed("pid") && opts.PID <= 0 {
				return fmt.Errorf("--pid must be > 0")
			}
			if f.Changed("port") && (opts.Port <= 0 || opts.Port > 65535) {
				return fmt.Errorf("--port must be in 1..65535")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			f := cmd.Flags()
			switch {
			case f.Changed("pid"):
				return killByPID(cmd, opts.PID)
			case f.Changed("name"):
				return killByName(cmd, opts.Name)
			case f.Changed("port"):
				return killByPort(cmd, opts.Port)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&opts.PID, "pid", "P", 0, "process ID")
	cmd.Flags().StringVarP(&opts.Name, "name", "n", "", "process name (substring match)")
	cmd.Flags().IntVarP(&opts.Port, "port", "p", 0, "listening port")

	cmd.MarkFlagsMutuallyExclusive("pid", "name", "port")
	cmd.MarkFlagsOneRequired("pid", "name", "port")

	return cmd
}

func killByPID(cmd *cobra.Command, pid int) error {
	if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
		if errors.Is(err, syscall.ESRCH) {
			return fmt.Errorf("process %d does not exist", pid)
		}
		if errors.Is(err, syscall.EPERM) {
			return fmt.Errorf("permission denied to kill process %d", pid)
		}
		return fmt.Errorf("failed to kill process %d: %w", pid, err)
	}
	cmd.Printf("Successfully killed process %d\n", pid)
	return nil
}

func killByName(cobraCmd *cobra.Command, name string) error {
	cmd := exec.Command("pgrep", "-f", name)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return fmt.Errorf("no process found matching name: %s", name)
		}
		return fmt.Errorf("failed to search processes: %w", err)
	}

	pidStrs := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(pidStrs) == 0 {
		return fmt.Errorf("no process found matching name: %s", name)
	}

	var killed []int
	var failed []string
	for _, pidStr := range pidStrs {
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			if errors.Is(err, syscall.ESRCH) {
				failed = append(failed, fmt.Sprintf("%d (not exist)", pid))
			} else if errors.Is(err, syscall.EPERM) {
				failed = append(failed, fmt.Sprintf("%d (permission denied)", pid))
			} else {
				failed = append(failed, fmt.Sprintf("%d (%v)", pid, err))
			}
		} else {
			killed = append(killed, pid)
		}
	}

	if len(killed) == 0 {
		return fmt.Errorf("failed to kill any processes matching '%s': %s", name, strings.Join(failed, ", "))
	}

	cobraCmd.Printf("Successfully killed %d process(es): %v\n", len(killed), killed)
	if len(failed) > 0 {
		cobraCmd.Printf("Failed to kill: %s\n", strings.Join(failed, ", "))
	}
	return nil
}

func killByPort(cobraCmd *cobra.Command, port int) error {
	cmd := exec.Command("lsof", "-ti", fmt.Sprintf(":%d", port))
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return fmt.Errorf("no process found listening on port %d", port)
		}
		return fmt.Errorf("failed to find process on port %d: %w", port, err)
	}

	pidStr := strings.TrimSpace(string(output))
	lines := strings.Split(pidStr, "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no process found listening on port %d", port)
	}

	var killed []int
	var failed []string
	for _, line := range lines {
		pid, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			continue
		}
		if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
			if errors.Is(err, syscall.ESRCH) {
				failed = append(failed, fmt.Sprintf("%d (not exist)", pid))
			} else if errors.Is(err, syscall.EPERM) {
				failed = append(failed, fmt.Sprintf("%d (permission denied)", pid))
			} else {
				failed = append(failed, fmt.Sprintf("%d (%v)", pid, err))
			}
		} else {
			killed = append(killed, pid)
		}
	}

	if len(killed) == 0 {
		return fmt.Errorf("failed to kill any process on port %d: %s", port, strings.Join(failed, ", "))
	}

	cobraCmd.Printf("Successfully killed %d process(es) on port %d: %v\n", len(killed), port, killed)
	if len(failed) > 0 {
		cobraCmd.Printf("Failed to kill: %s\n", strings.Join(failed, ", "))
	}
	return nil
}
