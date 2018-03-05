package loopback

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/prometheus/common/log"
	"github.com/r3boot/anycast-operator/pkg/utils"
)

type LoopbackInterfaceConfig struct {
	Interface string
}

type LoopbackInterface struct {
	config *LoopbackInterfaceConfig
	cmd    string
}

func NewLoopback(cfg *LoopbackInterfaceConfig) (*LoopbackInterface, error) {
	var err error

	intf := &LoopbackInterface{
		config: cfg,
	}

	intf.cmd, err = exec.LookPath("ip")
	if err != nil {
		return nil, fmt.Errorf("NewLoopback: %v", err)
	}

	return intf, nil
}

func (l *LoopbackInterface) ip(args ...string) ([]string, error) {
	var (
		stdoutBuf bytes.Buffer
		stderrBuf bytes.Buffer
	)

	cmd := exec.Command(l.cmd, args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	log.Debugf("Running '%s %s'", l.cmd, strings.Join(args, " "))

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("LoopbackInterface.ip: cmd.Run: %v", err)
	}

	return strings.Split(stdoutBuf.String(), "\n"), nil
}

func (l *LoopbackInterface) GetAnycastIPs() ([]string, error) {
	result, err := l.ip("addr", "show", "dev", l.config.Interface)
	if err != nil {
		return nil, fmt.Errorf("LoopbackInterface.GetAnycastIPs: %v", err)
	}

	allAnycastIPs := []string{}
	for _, line := range result {
		if !strings.Contains(line, "scope global") {
			continue
		}

		ipAddr := strings.Split(strings.Split(line, " ")[5], "/")[0]

		allAnycastIPs = append(allAnycastIPs, ipAddr)
	}

	return allAnycastIPs, nil
}

func (l *LoopbackInterface) AddIPAddress(addr string) error {
	addr = utils.AddCIDR(addr)

	_, err := l.ip("addr", "add", addr, "dev", l.config.Interface)
	if err != nil {
		return fmt.Errorf("LoopbackInterface.AddIPAddress: %v", err)
	}

	return nil
}

func (l *LoopbackInterface) RemoveIPAddress(addr string) error {
	addr = utils.AddCIDR(addr)

	_, err := l.ip("addr", "del", addr, "dev", l.config.Interface)
	if err != nil {
		return fmt.Errorf("LoopbackInterface.RemoveIPAddress: %v", err)
	}

	return nil
}
