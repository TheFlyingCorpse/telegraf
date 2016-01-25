// +build windows

package win_perfcounters

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"os"
	"os/signal"

	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/lxn/win"
)

var sampleConfig string = `  # By default this plugin returns basic CPU and Disk statistics. See the README file for more examples.
  # Uncomment examples below or write your own as you see fit. If the system being polled for data does not have the Object at startup of the Telegraf agent, it will not be gathered.
  # Settings:
  #PrintValid = false # Print All matching performance counters
  [[inputs.win_perfcounters.object]]
    # Processor usage, alternative to native, reports on a per core.
    ObjectName = "Processor"
    Instances = ["*"]
    Counters = ["% Idle Time", "% Interrupt Time", "% Privileged Time", "% User Time", "% Processor Time"]
    Measurement = "win_cpu"
    #IncludeTotal=false #Set to true to include _Total instance when querying for all (*).
  [[inputs.win_perfcounters.object]]
    # Disk times and queues
    ObjectName = "LogicalDisk"
    Instances = ["*"]
    Counters = ["% Idle Time", "% Disk Time","% Disk Read Time", "% Disk Write Time", "% User Time", "Current Disk Queue Length"]
    Measurement = "win_disk"
    #IncludeTotal=false #Set to true to include _Total instance when querying for all (*).
    #WarnOnMissing = false # Print out when the performance counter is missing, either of object, counter or instance.
  [[inputs.win_perfcounters.object]]
    ObjectName = "System"
    Counters = ["Context Switches/sec","System Calls/sec"]
    Instances = ["------"]
    Measurement = "win_system"
    #IncludeTotal=false #Set to true to include _Total instance when querying for all (*).
#  [[inputs.win_perfcounters.object]]
#    # Example query where the Instance portion must be removed to get data back, such as from the Memory object.
#    ObjectName = "Memory"
#    Counters = ["Available Bytes","Cache Faults/sec","Demand Zero Faults/sec","Page Faults/sec","Pages/sec","Transition Faults/sec","Pool Nonpaged Bytes","Pool Paged Bytes"]
#    Instances = ["------"] # Use 6 x - to remove the Instance bit from the query.
#    Measurement = "win_mem"
#    #IncludeTotal=false #Set to true to include _Total instance when querying for all (*).
`

var gItemList = make(map[int]*item)
var configParsed bool

type Win_PerfCounters struct {
	PrintValid bool
	Object     []perfobject
}

type perfobject struct {
	ObjectName    string
	Counters      []string
	Instances     []string
	Measurement   string
	WarnOnMissing bool
	IncludeTotal  bool
}

// Parsed configuration ends up here after it has been validated for valid Performance Counter paths
type itemList struct {
	items map[int]*item
}

type item struct {
	query         string
	objectName    string
	counter       string
	instance      string
	measurement   string
	include_total bool
	result        bool
	handle        win.PDH_HQUERY
	counterHandle win.PDH_HCOUNTER
}

func AddItem(query string, objectName string, counter string, instance string, measurement string, include_total bool) {

	var handle win.PDH_HQUERY
	var counterHandle win.PDH_HCOUNTER
	ret := win.PdhOpenQuery(0, 0, &handle)
	ret = win.PdhAddEnglishCounter(handle, query, 0, &counterHandle)

	_ = ret

	temp := &item{query, objectName, counter, instance, measurement, include_total, false, handle, counterHandle}
	index := len(gItemList)
	gItemList[index] = temp
}

func (m *Win_PerfCounters) Description() string {
	return "Input plugin to query Performance Counters on Windows operating systems"
}

func (m *Win_PerfCounters) SampleConfig() string {
	return sampleConfig
}

func (m *Win_PerfCounters) ParseConfig(metrics *itemList) {
	var query string

	for _, PerfObject := range m.Object {
		for _, counter := range PerfObject.Counters {
			for _, instance := range PerfObject.Instances {
				objectname := PerfObject.ObjectName

				if instance == "------" {
					query = "\\" + objectname + "\\" + counter
				} else {
					query = "\\" + objectname + "(" + instance + ")\\" + counter
				}

				var exists uint32 = win.PdhValidatePath(query)

				if exists == win.ERROR_SUCCESS {
					if m.PrintValid {
						fmt.Printf("Valid: %s\n", query)
					}
					AddItem(query, objectname, counter, instance, PerfObject.Measurement, PerfObject.IncludeTotal)
				} else if exists == 3221228472 { // win.PDH_CSTATUS_NO_OBJECT
					if PerfObject.WarnOnMissing {
						fmt.Printf("Performance Object '%s' does not exist in query: %s\n", objectname, query)
					}
				} else if exists == 3221228473 { //win.PDH_CSTATUS_NO_COUNTER
					if PerfObject.WarnOnMissing {
						fmt.Printf("Counter '%s' does not exist in query: %s\n", counter, query)
					}
				} else if exists == 2147485649 { //win.PDH_CSTATUS_NO_INSTANCE
					if PerfObject.WarnOnMissing {
						fmt.Printf("Instance '%s' does not exist in query: %s\n", instance, query)
					}
				} else {
					fmt.Printf("Invalid result: %v, query: %s\n", exists, query)
				}
			}
		}
	}
	configParsed = true
}

func (m *Win_PerfCounters) Cleanup(metrics *itemList) {
	// Cleanup

	for _, metric := range metrics.items {
		ret := win.PdhCloseQuery(metric.handle)
		_ = ret
	}
}

func (m *Win_PerfCounters) Gather(acc inputs.Accumulator) error {
	metrics := itemList{}

	// We only need to parse the config during the init, it uses the global variable after.
	if configParsed == false {
		m.ParseConfig(&metrics)
	}

	// When interrupt or terminate is called.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() error {
		<-c
		m.Cleanup(&metrics)
		return nil
	}()

	var bufSize uint32
	var bufCount uint32
	var size uint32 = uint32(unsafe.Sizeof(win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE{}))
	var emptyBuf [1]win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE // need at least 1 addressable null ptr.

	// For iterate over the known metrics and get the samples.
	for _, metric := range gItemList {
		// collect
		ret := win.PdhCollectQueryData(metric.handle)
		if ret == win.ERROR_SUCCESS {
			ret = win.PdhGetFormattedCounterArrayDouble(metric.counterHandle, &bufSize, &bufCount, &emptyBuf[0]) // uses null ptr here according to MSDN.
			if ret == win.PDH_MORE_DATA {
				filledBuf := make([]win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE, bufCount*size)
				ret = win.PdhGetFormattedCounterArrayDouble(metric.counterHandle, &bufSize, &bufCount, &filledBuf[0])
				for i := 0; i < int(bufCount); i++ {
					c := filledBuf[i]
					var s string = win.UTF16PtrToString(c.SzName)

					var add bool

					// If IncludeTotal is set, include all.
					if metric.include_total {
						add = true
					} else if metric.instance == "*" && !strings.Contains(s, "_Total") { // Catch if set to * and that it is not a '*_Total*' instance.
						add = true
					} else if metric.instance == s { // Catch if we set it to total or some form of it
						add = true
					} else if metric.instance == "------" {
						add = true
					}

					if add {
						fields := make(map[string]interface{})
						tags := make(map[string]string)
						if s != "" {
							tags["instance"] = s
						}
						tags["objectname"] = metric.objectName
						fields[string(metric.counter)] = c.FmtValue.DoubleValue

						var measurement string
						if metric.measurement == "" {
							measurement = "win_perfcounters"
						} else {
							measurement = metric.measurement
						}
						acc.AddFields(measurement, fields, tags)
					}
				}

				filledBuf = nil
				// Need to at least set bufSize to zero, because if not, the function will not
				// return PDH_MORE_DATA and will not set the bufSize.
				bufCount = 0
				bufSize = 0
			}

		}
	}

	return nil
}

func init() {
	inputs.Add("win_perfcounters", func() inputs.Input { return &Win_PerfCounters{} })
}
