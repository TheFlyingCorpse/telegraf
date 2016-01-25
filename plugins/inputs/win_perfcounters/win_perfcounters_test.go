package win_perfcounters

import (
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/require"
)

func TestWinPerfcountersGet(t *testing.T) {
	var instances []string
	var counters []string

	objectname := "Processor Information"
	instances[0] = "_Total"
	counters[0] = "% Processor Time"

	var measurement string = "none"
	var warnonmissing bool = false
	var failonmissing bool = true
	var includetotal bool = false

	p := perfobject{ObjectName: objectname, Instances: instances, Measurement: "none", WarnOnMissing: warnonmissing, FailOnMissing: failonmissing, IncludeTotal: includetotal}

	m := Win_PerfCounters{PrintValid: true, Object: &s}

	var acc testutil.Accumulator
	err := m.Gather(&acc)
	require.NoError(t, err)
}

func TestWinPerfcountersError1(t *testing.T) {

	var instances []string
	var counters []string

	objectname := "Processor InformationERROR"
	instances[0] = "_Total"
	counters[0] = "% Processor Time"

	var measurement string = "none"
	var warnonmissing bool = false
	var failonmissing bool = true
	var includetotal bool = false

	p := perfobject{ObjectName: objectname, Instances: instances, Measurement: "none", WarnOnMissing: warnonmissing, FailOnMissing: failonmissing, IncludeTotal: includetotal}

	m := Win_PerfCounters{PrintValid: true, Object: &s}

	var acc testutil.Accumulator
	err := m.Gather(&acc)
	require.Error(t, err)
}

func TestWinPerfcountersError2(t *testing.T) {
	var instances []string
	var counters []string

	objectname := "Processor Information"
	instances[0] = "_TotalERROR"
	counters[0] = "% Processor Time"

	var measurement string = "none"
	var warnonmissing bool = false
	var failonmissing bool = true
	var includetotal bool = false

	p := perfobject{ObjectName: objectname, Instances: instances, Measurement: "none", WarnOnMissing: warnonmissing, FailOnMissing: failonmissing, IncludeTotal: includetotal}

	m := Win_PerfCounters{PrintValid: true, Object: &s}

	var acc testutil.Accumulator
	err := m.Gather(&acc)
	require.Error(t, err)
}

func TestWinPerfcountersError3(t *testing.T) {
	var instances []string
	var counters []string

	objectname := "Processor Information"
	instances[0] = "_Total"
	counters[0] = "% Processor TimeERROR"

	var measurement string = "none"
	var warnonmissing bool = false
	var failonmissing bool = true
	var includetotal bool = false

	p := perfobject{ObjectName: objectname, Instances: instances, Measurement: "none", WarnOnMissing: warnonmissing, FailOnMissing: failonmissing, IncludeTotal: includetotal}

	m := Win_PerfCounters{PrintValid: true, Object: &s}

	var acc testutil.Accumulator
	err := m.Gather(&acc)
	require.Error(t, err)
}
