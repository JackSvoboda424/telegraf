package traceroute

import (
	_ "embed"

	"fmt"
 	"os/exec"
 	"runtime"
	"strings"
	"strconv"
	"errors"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

//go:embed sample.conf
var sampleConfig string

type Traceroute struct {
	Url 		string 	`toml:"url"`
	Max_Hops 	int		`toml:"max_hops"`

	// Windows only setting that adjusts the wait timeout for each reply (in ms)
	Timeout 	int   	`toml:"timeout"`
}

func (*Traceroute) SampleConfig() string {
	return sampleConfig
}

func (s *Traceroute) Init() error {
	if s.Url == "" {
		return errors.New("Need to provide a URL")
	}

	// The interval cannot be below 0.2 seconds, matching ping implementation: https://linux.die.net/man/8/ping
	if s.Timeout < 10 {
		return errors.New("Timeouts below 10ms are invalid")
	} 

	if s.Max_Hops < 1 {
		return errors.New("Max Hops must be greater than 0")
	} 

	return nil
}

func (s *Traceroute) Gather(acc telegraf.Accumulator) error {

    target := s.Url
	if target == "" {
		return fmt.Errorf("Traceroute target url cannot be empty")
	}
	timeout := 1000

	if s.Timeout > 0 {
		timeout = s.Timeout
	}

	max_hops := 30
	if s.Max_Hops > 0 && s.Max_Hops <= 30 {
		max_hops = s.Max_Hops
	}
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tracert", "-w", strconv.Itoa(timeout), "-h", strconv.Itoa(max_hops), target)
	} else {
		cmd = exec.Command("traceroute", target, "-m", strconv.Itoa(max_hops))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		if runtime.GOOS == "windows" {
			if strings.Contains(string(output), "Unable to resolve target") {
				return errors.New(string(output))
			}

			return errors.New(err.Error())
		} else {
			if strings.Contains(string(output), "Name or service not known") {
				error_msg := "Unable to resolve target system name " + target
				return errors.New(error_msg)
			}

			return errors.New(err.Error())
		}
	} else{
		result := string(output)

		lines := strings.Split(result, "\n")
		hops := 0
		if runtime.GOOS == "windows" {
			hops = len(lines) - 7

		} else {
			hops = len(lines) - 1
		}


		fields := make(map[string]interface{})
		fields["hops"] = hops

		tags := make(map[string]string)
		tags["url"] = target
		tags["timeout"] = strconv.Itoa(timeout)

		acc.AddFields("traceroute", fields, tags)

		return nil
	}

	
}

func init() {
	inputs.Add("traceroute", func() telegraf.Input { 
		s := &Traceroute{
			Url:		"127.0.0.1",
			Max_Hops:	30,
			Timeout:	1000,
		}
		return s
	})
}
