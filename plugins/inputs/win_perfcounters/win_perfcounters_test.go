// +build windows

package win_perfcounters

import (
	"errors"
	//"fmt"
	"testing"

	//"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/require"
)

func TestWinPerfcountersConfigGet1(t *testing.T) {
	validmetrics := itemList{}

	var instances = make([]string, 1)
	var counters = make([]string, 1)
	var perfobjects = make([]perfobject, 1)

	objectname := "Processor Information"
	instances[0] = "_Total"
	counters[0] = "% Processor Time"

	var measurement string = "test"
	var warnonmissing bool = false
	var failonmissing bool = true
	var includetotal bool = false

	PerfObject := perfobject{
		ObjectName:    objectname,
		Instances:     instances,
		Counters:      counters,
		Measurement:   measurement,
		WarnOnMissing: warnonmissing,
		FailOnMissing: failonmissing,
		IncludeTotal:  includetotal,
	}

	perfobjects[0] = PerfObject

	m := Win_PerfCounters{PrintValid: false, Object: perfobjects}

	err := m.ParseConfig(&validmetrics)
	require.NoError(t, err)
}

func TestWinPerfcountersConfigGet2(t *testing.T) {
	metrics := itemList{}

	var instances = make([]string, 1)
	var counters = make([]string, 1)
	var perfobjects = make([]perfobject, 1)

	objectname := "Processor Information"
	instances[0] = "_Total"
	counters[0] = "% Processor Time"

	var measurement string = "test"
	var warnonmissing bool = false
	var failonmissing bool = true
	var includetotal bool = false

	PerfObject := perfobject{
		ObjectName:    objectname,
		Instances:     instances,
		Counters:      counters,
		Measurement:   measurement,
		WarnOnMissing: warnonmissing,
		FailOnMissing: failonmissing,
		IncludeTotal:  includetotal,
	}

	perfobjects[0] = PerfObject

	m := Win_PerfCounters{PrintValid: false, Object: perfobjects}

	err := m.ParseConfig(&metrics)
	require.NoError(t, err)

	if len(metrics.items) == 1 {
		require.NoError(t, nil)
	} else if len(metrics.items) == 0 {
		var errorstring1 string = "No results returned from the query: " + string(len(metrics.items))
		err2 := errors.New(errorstring1)
		require.NoError(t, err2)
	} else if len(metrics.items) > 1 {
		var errorstring1 string = "Too many results returned from the query: " + string(len(metrics.items))
		err2 := errors.New(errorstring1)
		require.NoError(t, err2)
	}
}

func TestWinPerfcountersConfigGet3(t *testing.T) {
	metrics := itemList{}

	var instances = make([]string, 1)
	var counters = make([]string, 2)
	var perfobjects = make([]perfobject, 1)

	objectname := "Processor Information"
	instances[0] = "_Total"
	counters[0] = "% Processor Time"
	counters[1] = "% Idle Time"

	var measurement string = "test"
	var warnonmissing bool = false
	var failonmissing bool = true
	var includetotal bool = false

	PerfObject := perfobject{
		ObjectName:    objectname,
		Instances:     instances,
		Counters:      counters,
		Measurement:   measurement,
		WarnOnMissing: warnonmissing,
		FailOnMissing: failonmissing,
		IncludeTotal:  includetotal,
	}

	perfobjects[0] = PerfObject

	m := Win_PerfCounters{PrintValid: false, Object: perfobjects}

	err := m.ParseConfig(&metrics)
	require.NoError(t, err)

	if len(metrics.items) == 2 {
		require.NoError(t, nil)
	} else if len(metrics.items) < 2 {

		var errorstring1 string = "Too few results returned from the query. " + string(len(metrics.items))
		err2 := errors.New(errorstring1)
		require.NoError(t, err2)
	} else if len(metrics.items) > 2 {

		var errorstring1 string = "Too many results returned from the query: " + string(len(metrics.items))
		err2 := errors.New(errorstring1)
		require.NoError(t, err2)
	}
}

func TestWinPerfcountersConfigGet4(t *testing.T) {
	metrics := itemList{}

	var instances = make([]string, 2)
	var counters = make([]string, 1)
	var perfobjects = make([]perfobject, 1)

	objectname := "Processor Information"
	instances[0] = "_Total"
	instances[1] = "0"
	counters[0] = "% Processor Time"

	var measurement string = "test"
	var warnonmissing bool = false
	var failonmissing bool = true
	var includetotal bool = false

	PerfObject := perfobject{
		ObjectName:    objectname,
		Instances:     instances,
		Counters:      counters,
		Measurement:   measurement,
		WarnOnMissing: warnonmissing,
		FailOnMissing: failonmissing,
		IncludeTotal:  includetotal,
	}

	perfobjects[0] = PerfObject

	m := Win_PerfCounters{PrintValid: false, Object: perfobjects}

	err := m.ParseConfig(&metrics)
	require.NoError(t, err)

	if len(metrics.items) == 2 {
		require.NoError(t, nil)
	} else if len(metrics.items) < 2 {

		var errorstring1 string = "Too few results returned from the query: " + string(len(metrics.items))
		err2 := errors.New(errorstring1)
		require.NoError(t, err2)
	} else if len(metrics.items) > 2 {

		var errorstring1 string = "Too many results returned from the query: " + string(len(metrics.items))
		err2 := errors.New(errorstring1)
		require.NoError(t, err2)
	}
}

func TestWinPerfcountersConfigGet5(t *testing.T) {
	metrics := itemList{}

	var instances = make([]string, 2)
	var counters = make([]string, 2)
	var perfobjects = make([]perfobject, 1)

	objectname := "Processor Information"
	instances[0] = "_Total"
	instances[1] = "0"
	counters[0] = "% Processor Time"
	counters[1] = "% Idle Time"

	var measurement string = "test"
	var warnonmissing bool = false
	var failonmissing bool = true
	var includetotal bool = false

	PerfObject := perfobject{
		ObjectName:    objectname,
		Instances:     instances,
		Counters:      counters,
		Measurement:   measurement,
		WarnOnMissing: warnonmissing,
		FailOnMissing: failonmissing,
		IncludeTotal:  includetotal,
	}

	perfobjects[0] = PerfObject

	m := Win_PerfCounters{PrintValid: false, Object: perfobjects}

	err := m.ParseConfig(&metrics)
	require.NoError(t, err)

	if len(metrics.items) == 4 {
		require.NoError(t, nil)
	} else if len(metrics.items) < 4 {
		var errorstring1 string = "Too few results returned from the query: " + string(len(metrics.items))
		err2 := errors.New(errorstring1)
		require.NoError(t, err2)
	} else if len(metrics.items) > 4 {
		var errorstring1 string = "Too many results returned from the query: " + string(len(metrics.items))
		err2 := errors.New(errorstring1)
		require.NoError(t, err2)
	}
}

func TestWinPerfcountersConfigError1(t *testing.T) {
	metrics := itemList{}

	var instances = make([]string, 1)
	var counters = make([]string, 1)
	var perfobjects = make([]perfobject, 1)

	objectname := "Processor InformationERROR"
	instances[0] = "_Total"
	counters[0] = "% Processor Time"

	var measurement string = "test"
	var warnonmissing bool = false
	var failonmissing bool = true
	var includetotal bool = false

	PerfObject := perfobject{
		ObjectName:    objectname,
		Instances:     instances,
		Counters:      counters,
		Measurement:   measurement,
		WarnOnMissing: warnonmissing,
		FailOnMissing: failonmissing,
		IncludeTotal:  includetotal,
	}

	perfobjects[0] = PerfObject

	m := Win_PerfCounters{PrintValid: false, Object: perfobjects}

	err := m.ParseConfig(&metrics)
	require.Error(t, err)
}

func TestWinPerfcountersConfigError2(t *testing.T) {
	metrics := itemList{}

	var instances = make([]string, 1)
	var counters = make([]string, 1)
	var perfobjects = make([]perfobject, 1)

	objectname := "Processor"
	instances[0] = "SuperERROR"
	counters[0] = "% C1 Time"

	var measurement string = "test"
	var warnonmissing bool = false
	var failonmissing bool = true
	var includetotal bool = false

	PerfObject := perfobject{
		ObjectName:    objectname,
		Instances:     instances,
		Counters:      counters,
		Measurement:   measurement,
		WarnOnMissing: warnonmissing,
		FailOnMissing: failonmissing,
		IncludeTotal:  includetotal,
	}

	perfobjects[0] = PerfObject

	m := Win_PerfCounters{PrintValid: false, Object: perfobjects}

	err := m.ParseConfig(&metrics)
	require.Error(t, err)
}

func TestWinPerfcountersConfigError3(t *testing.T) {
	metrics := itemList{}

	var instances = make([]string, 1)
	var counters = make([]string, 1)
	var perfobjects = make([]perfobject, 1)

	objectname := "Processor Information"
	instances[0] = "_Total"
	counters[0] = "% Processor TimeERROR"

	var measurement string = "test"
	var warnonmissing bool = false
	var failonmissing bool = true
	var includetotal bool = false

	PerfObject := perfobject{
		ObjectName:    objectname,
		Instances:     instances,
		Counters:      counters,
		Measurement:   measurement,
		WarnOnMissing: warnonmissing,
		FailOnMissing: failonmissing,
		IncludeTotal:  includetotal,
	}

	perfobjects[0] = PerfObject

	m := Win_PerfCounters{PrintValid: false, Object: perfobjects}

	err := m.ParseConfig(&metrics)
	require.Error(t, err)
}

// Broken, working on resolving
//func TestWinPerfcountersCollect1(t *testing.T) {

//	var instances = make([]string, 1)
//	var counters = make([]string, 2)
//	var perfobjects = make([]perfobject, 1)

//	objectname := "Processor Information"
//	instances[0] = "_Total"
//	counters[0] = "% Processor Time"
//	counters[1] = "% Idle Time"

//	var measurement string = "test"
//	var warnonmissing bool = false
//	var failonmissing bool = true
//	var includetotal bool = false
//	var testmode bool = true

//	PerfObject := perfobject{
//		ObjectName:    objectname,
//		Instances:     instances,
//		Counters:      counters,
//		Measurement:   measurement,
//		WarnOnMissing: warnonmissing,
//		FailOnMissing: failonmissing,
//		IncludeTotal:  includetotal,
//		TestMode:      testmode,
//	}

//	perfobjects[0] = PerfObject

//	m := Win_PerfCounters{PrintValid: false, Object: perfobjects}
//	var acc testutil.Accumulator
//	err := m.Gather(&acc)

//	//err := m.ParseConfig(&metrics)
//	require.NoError(t, err)
//	fmt.Printf("Map: %v\n", acc)
//	acc.AssertContainsTaggedFields(t,
//		"% Idle Time",
//		map[string]interface{}{
//			"% Idle Time": int(4),
//		},
//		map[string]string{
//			"instance": "_Total",
//		},
//	)
//}

