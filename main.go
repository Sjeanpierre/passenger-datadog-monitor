package main

import (
	"bytes"
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"encoding/xml"
	"fmt"
	"github.com/PagerDuty/godspeed"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"
	"time"
)

var print bool

type passengerStatus struct {
	XMLName      xml.Name  `xml:"info"`
	ProcessCount int       `xml:"process_count"`
	PoolMax      int       `xml:"max"`
	PoolCurrent  int       `xml:"capacity_used"`
	QueuedCount  int       `xml:"get_wait_list_size"`
	Processes    []process `xml:"supergroups>supergroup>group>processes>process"`
}

type process struct {
	CurrentSessions int `xml:"sessions"`
	Processed       int `xml:"processed"`
	SpawnTime       int `xml:"spawn_end_time"`
	CPU             int `xml:"cpu"`
	Memory          int `xml:"real_memory"`
}

//Stats is used to store stats
type Stats struct {
	min int
	len int
	avg int
	max int
	sum int
}

func summerizeStats(statsArray *[]int) Stats {
	var processedStats Stats
	sum, count := 0, len(*statsArray)
	sort.Sort(sort.IntSlice(*statsArray))

	for _, stat := range *statsArray {
		sum += stat
	}
	sortedStats := *statsArray
	processedStats.min = sortedStats[0]
	processedStats.len = count
	processedStats.avg = sum / count
	processedStats.max = sortedStats[len(sortedStats)-1]
	processedStats.sum = sum

	return processedStats
}

func retrievePassengerStats() (io.Reader, error) {
	out, err := exec.Command("passenger-status", "--show=xml").Output()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return bytes.NewReader(out), nil
}
func parsePassengerXML(xmlData *io.Reader) (passengerStatus, error) {
	var ParsedPassengerXML passengerStatus
	dec := xml.NewDecoder(*xmlData)
	dec.CharsetReader = charset.NewReader
	err := dec.Decode(&ParsedPassengerXML)
	if err != nil {
		return passengerStatus{}, err
	}
	return ParsedPassengerXML, nil
}

func floatMyInt(value int) float64{
    return float64(value)
}

func processed(passengerDetails *passengerStatus) Stats {
	var processed []int
	processes := passengerDetails.Processes
	for _, processStats := range processes {
		processed = append(processed, processStats.Processed)
	}
	return summerizeStats(&processed)
}
func memory(passengerDetails *passengerStatus) Stats {
	var memory []int
	processes := passengerDetails.Processes
	for _, processStats := range processes {
		memory = append(memory, processStats.Memory)
	}
	return summerizeStats(&memory)
}
func processUptime(passengerDetails *passengerStatus) Stats {

	var upTimes []int
	processes := passengerDetails.Processes
	for _, processStats := range processes {
		SpawnedNano := time.Unix(0, int64(processStats.SpawnTime*1000))
		uptime := time.Since(SpawnedNano)
		upTimes = append(upTimes, int(uptime.Minutes()))
	}
	return summerizeStats(&upTimes)
}

func chartPendingRequest(passengerDetails *passengerStatus, DogStatsD *godspeed.Godspeed) {
	if print {
		fmt.Println("|=====Queue Depth====|")
		fmt.Println("Queue Depth", passengerDetails.QueuedCount)
	}
    DogStatsD.Gauge("passenger.queue.depth", floatMyInt(passengerDetails.QueuedCount),nil)
}
func chartPoolUse(passengerDetails *passengerStatus, DogStatsD *godspeed.Godspeed) {
	if print {
		fmt.Println("|=====Pool Usage====|")
		fmt.Println("Used Pool", passengerDetails.ProcessCount)
		fmt.Println("Max Pool", passengerDetails.PoolMax)
	}
        DogStatsD.Gauge("passenger.pool.used", floatMyInt(passengerDetails.ProcessCount),nil)
        DogStatsD.Gauge("passenger.pool.max", floatMyInt(passengerDetails.PoolMax),nil)
}
func chartProcessed(passengerDetails *passengerStatus, DogStatsD *godspeed.Godspeed) {
	stats := processed(passengerDetails)
	if print {
		fmt.Println("|=====Processed====|")
		fmt.Println("Total processed", stats.sum)   //sum processed by threads
		fmt.Println("Average processed", stats.avg) //average processed by threads
		fmt.Println("Minimum processed", stats.min)
		fmt.Println("Maximum processed", stats.max)
	}
    DogStatsD.Gauge("passenger.processed.total",floatMyInt(stats.sum),nil)
    DogStatsD.Gauge("passenger.processed.avg",floatMyInt(stats.avg),nil)
    DogStatsD.Gauge("passenger.processed.min",floatMyInt(stats.min),nil)
    DogStatsD.Gauge("passenger.processed.max",floatMyInt(stats.max),nil)

}
func chartMemory(passengerDetails *passengerStatus, DogStatsD *godspeed.Godspeed) {
	stats := memory(passengerDetails)
	if print {
		fmt.Println("|=====Memory====|")
		fmt.Println("Total memory", stats.sum/1024)
		fmt.Println("Average memory", stats.avg/1024)
		fmt.Println("Minimum memory", stats.min/1024)
		fmt.Println("Maximum memory", stats.max/1024)
	}
    DogStatsD.Gauge("passenger.memory.total",floatMyInt(stats.sum/1024),nil)
    DogStatsD.Gauge("passenger.memory.avg",floatMyInt(stats.avg/1024),nil)
    DogStatsD.Gauge("passenger.memory.min",floatMyInt(stats.min/1024),nil)
    DogStatsD.Gauge("passenger.memory.max",floatMyInt(stats.max/1024),nil)
}
func chartProcessUptime(passengerDetails *passengerStatus, DogStatsD *godspeed.Godspeed) {
	stats := processUptime(passengerDetails)
	if print {
		fmt.Println("|=====Process uptime====|")
		fmt.Println("Average uptime", stats.avg, "min")
		fmt.Println("Minimum uptime", stats.min, "min")
		fmt.Println("Maximum uptime", stats.max, "min")
	}
    DogStatsD.Gauge("passenger.uptime.avg",floatMyInt(stats.avg),nil)
    DogStatsD.Gauge("passenger.uptime.min",floatMyInt(stats.min),nil)
    DogStatsD.Gauge("passenger.uptime.max",floatMyInt(stats.max),nil)
}

func main() {
	if len(os.Args[1:]) > 0 {
		if os.Args[1] == "print" {
			print = true
		}
	}
    for {
        xmlData, err := retrievePassengerStats()
        if err != nil {
            log.Fatal("Error getting passenger data:", err)
        }
        PassengerStatusData, err := parsePassengerXML(&xmlData)
        if err != nil {
            log.Fatal("Error parsing passenger data:", err)
        }
        DogStatsD, err := godspeed.NewDefault()
        if err != nil {
            log.Fatal("Error establishing StatsD connection", err)
        }
        defer DogStatsD.Conn.Close()

        chartProcessed(&PassengerStatusData, DogStatsD)
        chartMemory(&PassengerStatusData, DogStatsD)
        chartPendingRequest(&PassengerStatusData, DogStatsD)
        chartPoolUse(&PassengerStatusData, DogStatsD)
        chartProcessUptime(&PassengerStatusData, DogStatsD)
        time.Sleep(10 * time.Second)
    }
}
