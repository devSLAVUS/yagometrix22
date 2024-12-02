package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type Agent struct {
	PullCount      int64
	PullInterval   time.Duration
	ReportInterval time.Duration
	ServerAddress  string
	Metrics        map[string]float64
}

func NewAgent(serverAddress string, pullInterval, reportInterval time.Duration) *Agent {
	return &Agent{
		ServerAddress:  serverAddress,
		PullInterval:   pullInterval,
		ReportInterval: reportInterval,
		Metrics:        make(map[string]float64),
	}
}

func (a *Agent) collectMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	a.Metrics["Alloc"] = float64(m.Alloc)
	a.Metrics["BuckHashSys"] = float64(m.BuckHashSys)
	a.Metrics["Frees"] = float64(m.Frees)
	a.Metrics["GCCPUFraction"] = m.GCCPUFraction
	a.Metrics["GCSys"] = float64(m.GCSys)
	a.Metrics["HeapAlloc"] = float64(m.HeapAlloc)
	a.Metrics["HeapIdle"] = float64(m.HeapIdle)
	a.Metrics["HeapInuse"] = float64(m.HeapInuse)
	a.Metrics["HeapObjects"] = float64(m.HeapObjects)
	a.Metrics["HeapReleased"] = float64(m.HeapReleased)
	a.Metrics["HeapSys"] = float64(m.HeapSys)
	a.Metrics["LastGC"] = float64(m.LastGC)
	a.Metrics["Lookups"] = float64(m.Lookups)
	a.Metrics["MCacheInuse"] = float64(m.MCacheInuse)
	a.Metrics["MCacheSys"] = float64(m.MCacheSys)
	a.Metrics["MSpanInuse"] = float64(m.MSpanInuse)
	a.Metrics["MSpanSys"] = float64(m.MSpanSys)
	a.Metrics["Mallocs"] = float64(m.Mallocs)
	a.Metrics["NextGC"] = float64(m.NextGC)
	a.Metrics["NumForcedGC"] = float64(m.NumForcedGC)
	a.Metrics["NumGC"] = float64(m.NumGC)
	a.Metrics["OtherSys"] = float64(m.OtherSys)
	a.Metrics["PauseTotalNs"] = float64(m.PauseTotalNs)
	a.Metrics["StackInuse"] = float64(m.StackInuse)
	a.Metrics["StackSys"] = float64(m.StackSys)
	a.Metrics["Sys"] = float64(m.Sys)
	a.Metrics["TotalAlloc"] = float64(m.TotalAlloc)

	a.PullCount++
	a.Metrics["PullCount"] = float64(a.PullCount)
	a.Metrics["RandomValue"] = rand.Float64()
}

func (a *Agent) sendMetric(metricType MetricType, metricName string, metricValue float64) {
	url := fmt.Sprintf("%s/update/%s/%s/%v", a.ServerAddress, metricType, metricName, metricValue)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Server returned non-OK status: %s\n", resp.Status)
	}
}

func (a *Agent) run() {
	tPull := time.NewTicker(a.PullInterval)
	tReport := time.NewTicker(a.ReportInterval)
	defer tPull.Stop()
	defer tReport.Stop()

	for {
		select {
		case <-tPull.C:
			a.collectMetrics()

		case <-tReport.C:
			for name, value := range a.Metrics {
				var metricType MetricType
				if name == "PullCount" {
					metricType = Counter
				} else {
					metricType = Gauge
				}
				a.sendMetric(metricType, name, value)
			}
		}
	}
}

func (a *Agent) runSleep() {

	for {
		time.Sleep(a.PullInterval)
		a.collectMetrics()

		time.Sleep(a.ReportInterval - a.PullInterval)
		for name, value := range a.Metrics {
			var metricType MetricType
			if name == "PullCount" {
				metricType = Counter
			} else {
				metricType = Gauge
			}
			a.sendMetric(metricType, name, value)
		}
	}
}

func main() {
	a := flag.String("a", "http://localhost:8080", "server http://ip:port")
	r := flag.Int("r", 10, "agent report interval")
	p := flag.Int("p", 2, "agent pull interval")
	flag.Parse()
	agent := NewAgent(*a, time.Duration(*p)*time.Second, time.Duration(*r)*time.Second)
	fmt.Println("Start agent", time.Now())
	agent.runSleep()
}
