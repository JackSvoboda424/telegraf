package traceroute

import (
	_ "embed"

	"testing"
	"github.com/stretchr/testify/require"
	"github.com/influxdata/telegraf/testutil"
)

func TestGather(t *testing.T) {
	type test struct {
		P *Traceroute
	}
	
	tests := []test{
		{
			P: &Traceroute{
					Url:           "google.com",
					Max_Hops:       30,
					Timeout:        1000,
				},
		},	
		{
			P: &Traceroute{
					Url:           "localhost",
					Max_Hops:       30,
					Timeout:        1000,
				},
		},	
	}
    
	for _, tt := range tests {
		var acc testutil.Accumulator
		require.NoError(t, tt.P.Init())
		require.NoError(t, acc.GatherError(tt.P.Gather))
		require.True(t, acc.HasField("traceroute", "hops"))
		require.InDelta(t, 15, acc.Metrics[0].Fields["hops"], 15)
    }
}
