package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/devSLAVUS/yagometrix22/internal/agent/collector"
	"github.com/devSLAVUS/yagometrix22/internal/agent/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	address := flag.String("a", cfg.Address, "server http://ip:port")
	reportInterval := flag.Int("r", cfg.ReportInterval, "agent report interval in seconds")
	pullInterval := flag.Int("p", cfg.PullInterval, "agent poll interval in seconds")
	flag.Parse()

	if *address != cfg.Address {
		cfg.Address = *address
	}
	if *pullInterval != cfg.PullInterval {
		cfg.PullInterval = *pullInterval
	}
	if *reportInterval != cfg.ReportInterval {
		cfg.ReportInterval = *reportInterval
	}

	agent := collector.NewAgent(
		"http://"+cfg.Address,
		time.Duration(cfg.PullInterval)*time.Second,
		time.Duration(cfg.ReportInterval)*time.Second,
	)

	fmt.Printf("Starting agent on %s at %s\n", "http://"+cfg.Address, time.Now())
	fmt.Printf("Poll Interval: %d seconds, Report Interval: %d seconds\n", cfg.PullInterval, cfg.ReportInterval)

	agent.RunSleep()
}
