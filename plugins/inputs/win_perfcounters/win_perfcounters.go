// +build windows

package win_perfcounters

import (
	"fmt"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"os"
	"os/signal"

	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/lxn/win"
)

var sampleConfig string = `#  By default this plugin gathers IO 
#  Uncomment examples below or write your own as you see fit. If the system being polled for data does not have the Object at startup of the Telegraf agent, it will not be gathered.
#  More examples can be found at: https://gist.github.com/TheFlyingCorpse/d1c4def0906ce8d35430
#  # Settings:
#  #PrintValid = false # Print All matching performance counters
#  [[inputs.win_perfcounters.object]]
#    # Useful if the server hosts a DFS Namespace or is a Domain Controller
#    ObjectName = "DFS Namespace Service Referrals"
#    Counters = ["Requests Processed","Requests Failed","Avg. Response Time"]
#    Measurement = "win_dfsn"
#  [[inputs.win_perfcounters.object]] 
#    # Useful if the server hosts a DFS Replication folder or is a Domain Controller
#    ObjectName = "DFS Replication Service Volumes"
#    Counters = ["Data Lookups","Database Commits"]
#    Measurement = "win_dfsr"
#  [[inputs.win_perfcounters.object]]
#    ObjectName = "DirectoryServices"
#    Counters = ["Base Searches/sec","Database adds/sec","Database deletes/sec","Database modifys/sec","Database recycles/sec","LDAP Client Sessions"]
#    Measurement = "win_ad" # Set an alternative measurement to win_perfcounters if wanted.
#    #Instances = [""] # Gathers all instances by default, specify to only gather these
#    #IncludeTotal = false # Set this to get the instance _Total back if its included in the instances either via * or as the single option
#    #WarnOnMissing = false # Print out when the performance counter is missing in some way, either of object, counter or instance.
#  [[inputs.win_perfcounters.object]]
#    ObjectName = "Security System-Wide Statistics"
#    Counters = ["NTLM Authentications","Kerberos Authentications","Digest Authentications"]
#    Measurement = "win_ad"
#  [[inputs.win_perfcounters.object]]
#    # HTTP Service request queues in the Kernel before being handed over to User Mode.
#    ObjectName = "HTTP Service Request Queues"
#    Counters = ["CurrentQueueSize","RejectedRequests"]
#    Measurement = "win_http_queues"
#  [[inputs.win_perfcounters.object]]
#    # Processor usage, alternative to native, reports on a per core.
#    ObjectName = "Processor"
#    Counters = ["% Idle Time", "% Interrupt Time", "% Privileged Time", "% User Time", "% Processor Time"]
#    Measurement = "win_cpu"
#  [[inputs.win_perfcounters.object]]
#    # Disk times and queues
#    ObjectName = "LogicalDisk"
#    Counters = ["% Idle Time", "% Disk Read Time", "% Disk Write Time", "% User Time", "Current Disk Queue Length"]
#    Measurement = "win_disk"
#  [[inputs.win_perfcounters.object]]
#    # Example query where the Instance portion must be removed to get data back, such as from the Memory object.
#    ObjectName = "Memory"
#    Counters = ["Available Bytes","Cache Bytes"]
#    Instances = ["------"] # Use 6 x - to remove the Instance bit from the query.
#    Measurement = "win_mem"
#
`

// Flag to break out of loops if set to false.
var execute bool = true

type Win_PerfCounters struct {
	Interval   int
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

func (m *itemList) AddItem(query string, objectName string, counter string, instance string, measurement string, include_total bool) {
	if m.items == nil {
		m.items = make(map[int]*item)
	}

	var handle win.PDH_HQUERY
	var counterHandle win.PDH_HCOUNTER
	ret := win.PdhOpenQuery(0, 0, &handle)
	ret = win.PdhAddEnglishCounter(handle, query, 0, &counterHandle)

	_ = ret

	temp := &item{query, objectName, counter, instance, measurement, include_total, false, handle, counterHandle}
	index := len(m.items)
	m.items[index] = temp
}

func (s *Win_PerfCounters) Description() string {
	return "Input plugin to query Performance Counters on Windows operating systems"
}

func (s *Win_PerfCounters) SampleConfig() string {
	return sampleConfig
}

func (s *Win_PerfCounters) ParseConfig(metrics *itemList) {
	var query string

	for _, PerfObject := range s.Object {
		for _, counter := range PerfObject.Counters {
			for _, instance := range PerfObject.Instances {
				objectname := PerfObject.ObjectName

				if instance == "------" {
					query = "\\" + objectname + "\\" + counter
				} else if instance == "" {
					query = "\\" + objectname + "(*)\\" + counter
				} else {
					query = "\\" + objectname + "(" + instance + ")\\" + counter
				}

				var exists uint32 = win.PdhValidatePath(query)

				if exists == win.ERROR_SUCCESS {
					if s.PrintValid {
						fmt.Printf("Valid: %s\n", query)
					}
					metrics.AddItem(query, objectname, counter, instance, PerfObject.Measurement, PerfObject.IncludeTotal)
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
}

func (s *Win_PerfCounters) Cleanup(metrics *itemList) {
	// Cleanup

	for _, metric := range metrics.items {
		ret := win.PdhCloseQuery(metric.handle)
		_ = ret
	}
}

func (s *Win_PerfCounters) Gather(acc inputs.Accumulator) error {
	metrics := itemList{}

	if execute {
		s.ParseConfig(&metrics)
	}

	// When interrupt is called
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() error {
		<-c
		execute = false
		s.Cleanup(&metrics)
		return nil
	}()

	var interval int
	if s.Interval == 0 {
		interval = 10000
	} else {
		interval = s.Interval
	}

	var bufSize uint32
	var bufCount uint32
	var size uint32 = uint32(unsafe.Sizeof(win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE{}))
	var emptyBuf [1]win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE // need at least 1 addressable null ptr.

	for execute {
		if len(metrics.items) == 0 {
			break
		}

		// For iterate over the known metrics and get the samples.
		for _, metric := range metrics.items {
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
						} else if !strings.Contains(s, "_Total") {
							add = true
						}

						if add {
							fields := make(map[string]interface{})
							tags := make(map[string]string)
							tags["instance"] = s
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
		time.Sleep(time.Duration(interval) * time.Millisecond)

	}

	return nil
}

func init() {
	inputs.Add("win_perfcounters", func() inputs.Input { return &Win_PerfCounters{} })
}
