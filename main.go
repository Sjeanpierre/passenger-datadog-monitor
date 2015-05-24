package main

import (
	"bytes"
	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"
	"encoding/xml"
	"fmt"
	"io"
	_ "os"
	"os/exec"
	"sort"
	"time"
    "log"
)

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

type stats struct {
	min int
	len int
	avg int
	max int
	sum int
}

func summerizeStats(statsArray *[]int) stats {
	var processedStats stats
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

func parsePassengerXML(xmlData *io.Reader) (passengerStatus, error){
    var ParsedPassengerXML passengerStatus
    dec := xml.NewDecoder(*xmlData)
    dec.CharsetReader = charset.NewReader
    err := dec.Decode(&ParsedPassengerXML)
    if err != nil {
        return passengerStatus{}, err
    }
    return ParsedPassengerXML, nil
}

func chartProcessed(passengerDetails *passengerStatus) {
	fmt.Println("|=====Processed====|")
	var processed []int
	processes := passengerDetails.Processes
	for _, processStats := range processes {
		processed = append(processed, processStats.Processed)
	}
	stats := summerizeStats(&processed)
	fmt.Println("Total processed", stats.sum)   //sum processed by threads
	fmt.Println("Average processed", stats.avg) //average processed by threads
	fmt.Println("Minimum processed", stats.min)
	fmt.Println("Maximum processed", stats.max)
}

func chartMemory(passengerDetails *passengerStatus) {
	fmt.Println("|=====Memory====|")
	var memory []int
	processes := passengerDetails.Processes
	for _, processStats := range processes {
		memory = append(memory, processStats.Memory)
	}
	stats := summerizeStats(&memory)
	fmt.Println("Total memory", stats.sum/1024)   //sum processed by threads
	fmt.Println("Average memory", stats.avg/1024) //average processed by threads
	fmt.Println("Minimum memory", stats.min/1024)
	fmt.Println("Maximum memory", stats.max/1024)
}

func chartPendingRequest(passengerDetails *passengerStatus) {
    fmt.Println("|=====Queue Depth====|")
	fmt.Println("Queue Depth", passengerDetails.QueuedCount)
}

func chartPoolUse(passengerDetails *passengerStatus) {
    fmt.Println("|=====Pool Usage====|")
	fmt.Println("Used Pool", passengerDetails.ProcessCount)
	fmt.Println("Max Pool", passengerDetails.PoolMax)
}

func chartProcessUptime(passengerDetails *passengerStatus) {
	fmt.Println("|=====Process uptime====|")
	var upTimes []int
	processes := passengerDetails.Processes
	for _, processStats := range processes {
		SpawnedNano := time.Unix(0, int64(processStats.SpawnTime*1000))
		uptime := time.Since(SpawnedNano)
		upTimes = append(upTimes, int(uptime.Minutes()))
	}
	stats := summerizeStats(&upTimes)
	fmt.Println("Average uptime", stats.avg, "min")
	fmt.Println("Minimum uptime", stats.min, "min")
	fmt.Println("Maximum uptime", stats.max, "min")
}

func main() {
	xmlData, err := retrievePassengerStats()
	if err != nil {
		log.Fatal("Error getting passenger data:", err)
	}
	PassengerStatusData, err := parsePassengerXML(&xmlData)
    if err != nil {
        log.Fatal("Error parsing passenger data:", err)
    }
	chartProcessed(&PassengerStatusData)
	chartMemory(&PassengerStatusData)
	chartPendingRequest(&PassengerStatusData)
	chartPoolUse(&PassengerStatusData)
	chartProcessUptime(&PassengerStatusData)
}
