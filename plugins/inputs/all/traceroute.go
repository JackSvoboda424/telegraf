//go:build !custom || inputs || inputs.traceroute

package all

import _ "github.com/influxdata/telegraf/plugins/inputs/traceroute" // register plugin
