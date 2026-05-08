package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	apiURL     string
	showHelp   bool
)

func init() {
	flag.StringVar(&apiURL, "url", "http://localhost:8080", "Proxy API URL")
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showHelp, "h", false, "Show help (shorthand)")
}

func main() {
	flag.Parse()

	if showHelp || flag.NArg() == 0 {
		printUsage()
		os.Exit(0)
	}

	args := flag.Args()
	command := args[0]

	switch command {
	case "status":
		getStatus()
	case "stats":
		getStats()
	case "list":
		listDropFlags()
	case "delays":
		listDelays()
	case "drop":
		if len(args) < 2 {
			fmt.Println("Error: signal type required")
			fmt.Println("Usage: dropctl drop <signal-type>")
			os.Exit(1)
		}
		setDrop(args[1], true)
	case "allow":
		if len(args) < 2 {
			fmt.Println("Error: signal type required")
			fmt.Println("Usage: dropctl allow <signal-type>")
			os.Exit(1)
		}
		setDrop(args[1], false)
	case "delay":
		if len(args) < 3 {
			fmt.Println("Error: signal type and delay (ms) required")
			fmt.Println("Usage: dropctl delay <signal-type> <milliseconds>")
			os.Exit(1)
		}
		delayMs, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			fmt.Printf("Error: invalid delay value: %s\n", args[2])
			os.Exit(1)
		}
		setDelay(args[1], delayMs)
	case "reset":
		resetDropFlags()
	case "reset-delays":
		resetDelays()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("dropctl - 4G Proxy Drop & Delay Control CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  dropctl [options] <command> [args]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --url URL      Proxy API URL (default: http://localhost:8080)")
	fmt.Println("  -h, --help     Show this help message")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  status                         Show proxy status")
	fmt.Println("  stats                          Show statistics")
	fmt.Println("  list                           List current drop flags")
	fmt.Println("  delays                         List current delay settings")
	fmt.Println("  drop <signal-type>             Enable dropping of signal type")
	fmt.Println("  allow <signal-type>            Disable dropping of signal type")
	fmt.Println("  delay <signal-type> <ms>       Set delay for signal type (milliseconds)")
	fmt.Println("  reset                          Reset all drop flags")
	fmt.Println("  reset-delays                   Reset all delays to zero")
	fmt.Println()
	fmt.Println("Signal Types:")
	fmt.Println("  attach               Attach Request/Accept/Complete/Reject")
	fmt.Println("  detach               Detach Request/Accept")
	fmt.Println("  tau                  Tracking Area Update")
	fmt.Println("  service-request      Service Request")
	fmt.Println("  ue-context-release   UE Context Release")
	fmt.Println("  pdn-connectivity     PDN Connectivity")
	fmt.Println("  handover             All Handover messages")
	fmt.Println("  handover-required    HandoverRequired (Source eNB -> MME)")
	fmt.Println("  handover-notify      HandoverNotify (Target eNB -> MME)")
	fmt.Println("  reset                Reset")
	fmt.Println("  paging               Paging")
	fmt.Println("  default              Default (for unspecified types)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Drop all Attach messages")
	fmt.Println("  dropctl drop attach")
	fmt.Println()
	fmt.Println("  # Allow Attach messages again")
	fmt.Println("  dropctl allow attach")
	fmt.Println()
	fmt.Println("  # Add 1000ms delay to Attach messages")
	fmt.Println("  dropctl delay attach 1000")
	fmt.Println()
	fmt.Println("  # Add 500ms delay to all TAU messages")
	fmt.Println("  dropctl delay tau 500")
	fmt.Println()
	fmt.Println("  # Add 2000ms delay to HandoverRequired (simulates core lag)")
	fmt.Println("  dropctl delay handover-required 2000")
	fmt.Println()
	fmt.Println("  # Add 1500ms delay to HandoverNotify")
	fmt.Println("  dropctl delay handover-notify 1500")
	fmt.Println()
	fmt.Println("  # Show current drop flags")
	fmt.Println("  dropctl list")
	fmt.Println()
	fmt.Println("  # Show current delay settings")
	fmt.Println("  dropctl delays")
	fmt.Println()
	fmt.Println("  # Reset all delays")
	fmt.Println("  dropctl reset-delays")
}

func getStatus() {
	resp, err := http.Get(apiURL + "/api/v1/status")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	prettyPrint(result)
}

func getStats() {
	resp, err := http.Get(apiURL + "/api/v1/stats")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	fmt.Println("Statistics:")
	fmt.Println("-----------")
	for key, value := range result {
		fmt.Printf("  %s: %v\n", key, value)
	}
}

func listDropFlags() {
	resp, err := http.Get(apiURL + "/api/v1/drop")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var flags map[string]bool
	json.Unmarshal(body, &flags)

	fmt.Println("Drop Flags:")
	fmt.Println("-----------")
	for signalType, drop := range flags {
		status := "ALLOW"
		if drop {
			status = "DROP"
		}
		fmt.Printf("  %-20s %s\n", signalType+":", status)
	}
}

func setDrop(signalType string, drop bool) {
	// Normalize signal type name
	signalType = normalizeSignalType(signalType)

	url := fmt.Sprintf("%s/api/v1/drop/%s", apiURL, signalType)

	reqBody, _ := json.Marshal(map[string]bool{"drop": drop})

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if resp.StatusCode == http.StatusOK {
		action := "will be forwarded"
		if drop {
			action = "will be dropped"
		}
		fmt.Printf("%s messages %s\n", signalType, action)
	} else {
		fmt.Printf("Error: %v\n", result["error"])
		os.Exit(1)
	}
}

func resetDropFlags() {
	req, err := http.NewRequest(http.MethodDelete, apiURL+"/api/v1/drop", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Println("All drop flags have been reset")
}

func listDelays() {
	resp, err := http.Get(apiURL + "/api/v1/delay")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var delays map[string]float64
	json.Unmarshal(body, &delays)

	fmt.Println("Delay Settings (milliseconds):")
	fmt.Println("------------------------------")
	for signalType, ms := range delays {
		if ms > 0 {
			fmt.Printf("  %-20s %dms\n", signalType+":", int64(ms))
		} else {
			fmt.Printf("  %-20s (none)\n", signalType+":")
		}
	}
}

func setDelay(signalType string, delayMs int64) {
	// Normalize signal type name
	signalType = normalizeSignalType(signalType)

	url := fmt.Sprintf("%s/api/v1/delay/%s", apiURL, signalType)

	reqBody, _ := json.Marshal(map[string]int64{"delayMs": delayMs})

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if resp.StatusCode == http.StatusOK {
		if delayMs > 0 {
			fmt.Printf("%s messages will be delayed by %dms\n", signalType, delayMs)
		} else {
			fmt.Printf("%s messages will not be delayed\n", signalType)
		}
	} else {
		fmt.Printf("Error: %v\n", result["error"])
		os.Exit(1)
	}
}

func resetDelays() {
	req, err := http.NewRequest(http.MethodDelete, apiURL+"/api/v1/delay", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Println("All delays have been reset to zero")
}

func normalizeSignalType(s string) string {
	// Handle hyphenated names
	s = strings.ReplaceAll(s, "-", "")

	switch strings.ToLower(s) {
	case "attach":
		return "attach"
	case "detach":
		return "detach"
	case "tau", "trackingareaupdate":
		return "tau"
	case "servicerequest", "service":
		return "serviceRequest"
	case "uecontextrelease", "release":
		return "ueContextRelease"
	case "pdnconnectivity", "pdn":
		return "pdnConnectivity"
	case "handover", "ho":
		return "handover"
	case "handoverrequired", "horequired":
		return "handoverRequired"
	case "handovernotify", "honotify":
		return "handoverNotify"
	case "reset":
		return "reset"
	case "paging":
		return "paging"
	default:
		return s
	}
}

func prettyPrint(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}
