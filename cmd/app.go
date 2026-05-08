package main

import (
	"4g-proxy/config"
	"4g-proxy/internal/models"
	"4g-proxy/internal/routes"
	"4g-proxy/internal/sctp"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	configFile   string
	listenAddr   string
	listenPort   int
	mmeAddr      string
	mmePort      int
	apiPort      int
	verbose      bool
	showHelp     bool
	showVersion  bool
)

const version = "1.0.0"

func init() {
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.StringVar(&configFile, "c", "", "Path to configuration file (shorthand)")
	flag.StringVar(&listenAddr, "listen", "", "Listen address for eNB connections")
	flag.IntVar(&listenPort, "port", 0, "Listen port for eNB connections (default: 36412)")
	flag.StringVar(&mmeAddr, "mme", "", "MME address (e.g., 127.0.0.1:36412)")
	flag.IntVar(&mmePort, "mme-port", 0, "MME port")
	flag.IntVar(&apiPort, "api-port", 0, "HTTP API port (default: 8080)")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&verbose, "v", false, "Enable verbose logging (shorthand)")
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showHelp, "h", false, "Show help (shorthand)")
	flag.BoolVar(&showVersion, "version", false, "Show version")
}

func main() {
	flag.Parse()

	if showHelp {
		printUsage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("4g-proxy version %s\n", version)
		os.Exit(0)
	}

	// Load configuration
	cfg := config.DefaultConfig()
	if configFile != "" {
		var err error
		cfg, err = config.Load(configFile)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	}

	// Load environment variables (takes precedence over config file)
	cfg.LoadFromEnv()

	// Override config with command line flags (takes highest precedence)
	if listenAddr != "" {
		cfg.Proxy.ListenAddress = listenAddr
	}
	if listenPort != 0 {
		cfg.Proxy.ListenPort = listenPort
	}
	if mmeAddr != "" {
		cfg.MME.Address = mmeAddr
	}
	if mmePort != 0 {
		cfg.MME.Port = mmePort
	}
	if apiPort != 0 {
		cfg.API.Port = apiPort
	}
	if verbose {
		cfg.Logging.Verbose = true
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	log.Println("===========================================")
	log.Println("       4G LTE S1AP Proxy")
	log.Printf("       Version: %s", version)
	log.Println("===========================================")

	// Create proxy state
	state := models.NewProxyState()

	// Apply delay configuration from config
	applyDelayConfig(cfg, state)

	// Create and start SCTP proxy
	proxy := sctp.NewProxy(cfg.ProxyEndpoint(), cfg.MMEEndpoint(), state)
	if err := proxy.Start(); err != nil {
		log.Fatalf("Failed to start proxy: %v", err)
	}

	// Start HTTP API if enabled
	if cfg.API.Enabled {
		router := routes.SetupRouter(state)
		go func() {
			log.Printf("HTTP API listening on %s", cfg.APIEndpoint())
			if err := router.Run(cfg.APIEndpoint()); err != nil {
				log.Printf("HTTP API error: %v", err)
			}
		}()
	}

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Proxy started. Press Ctrl+C to stop.")

	<-sigCh
	log.Println("\nShutting down...")

	proxy.Stop()
	log.Println("Goodbye!")
}

func applyDelayConfig(cfg *config.Config, state *models.ProxyState) {
	state.DelayConfig.SetDelayByName("attach", cfg.Delay.Attach)
	state.DelayConfig.SetDelayByName("detach", cfg.Delay.Detach)
	state.DelayConfig.SetDelayByName("tau", cfg.Delay.TAU)
	state.DelayConfig.SetDelayByName("serviceRequest", cfg.Delay.ServiceRequest)
	state.DelayConfig.SetDelayByName("ueContextRelease", cfg.Delay.UEContextRelease)
	state.DelayConfig.SetDelayByName("pdnConnectivity", cfg.Delay.PDNConnectivity)
	state.DelayConfig.SetDelayByName("handover", cfg.Delay.Handover)
	state.DelayConfig.SetDelayByName("handoverRequired", cfg.Delay.HandoverRequired)
	state.DelayConfig.SetDelayByName("handoverNotify", cfg.Delay.HandoverNotify)
	state.DelayConfig.SetDelayByName("reset", cfg.Delay.Reset)
	state.DelayConfig.SetDelayByName("paging", cfg.Delay.Paging)
	state.DelayConfig.SetDelayByName("default", cfg.Delay.Default)

	// Log configured delays
	delays := state.DelayConfig.GetAll()
	hasDelay := false
	for _, v := range delays {
		if v > 0 {
			hasDelay = true
			break
		}
	}
	if hasDelay {
		log.Println("Delay configuration:")
		for name, ms := range delays {
			if ms > 0 {
				log.Printf("  %s: %dms", name, ms)
			}
		}
	}
}

func printUsage() {
	fmt.Println("4G LTE S1AP Proxy")
	fmt.Println()
	fmt.Println("A transparent proxy for S1AP protocol between eNB and MME with")
	fmt.Println("message inspection, selective message dropping, and delay capabilities.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  4g-proxy [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -c, --config FILE      Path to configuration file")
	fmt.Println("  --listen ADDR          Listen address for eNB connections (default: 0.0.0.0)")
	fmt.Println("  --port PORT            Listen port (default: 36412)")
	fmt.Println("  --mme ADDR             MME address")
	fmt.Println("  --mme-port PORT        MME port")
	fmt.Println("  --api-port PORT        HTTP API port (default: 8080)")
	fmt.Println("  -v, --verbose          Enable verbose logging")
	fmt.Println("  -h, --help             Show this help message")
	fmt.Println("  --version              Show version")
	fmt.Println()
	fmt.Println("Environment Variables (for Kubernetes):")
	fmt.Println("  PROXY_LISTEN_ADDRESS       Listen address")
	fmt.Println("  PROXY_LISTEN_PORT          Listen port")
	fmt.Println("  MME_ADDRESS                MME address")
	fmt.Println("  MME_PORT                   MME port")
	fmt.Println("  API_ENABLED                Enable HTTP API (true/false)")
	fmt.Println("  API_PORT                   HTTP API port")
	fmt.Println("  DELAY_ATTACH_MS            Delay for Attach messages (ms)")
	fmt.Println("  DELAY_DETACH_MS            Delay for Detach messages (ms)")
	fmt.Println("  DELAY_TAU_MS               Delay for TAU messages (ms)")
	fmt.Println("  DELAY_SERVICE_REQUEST_MS   Delay for Service Request messages (ms)")
	fmt.Println("  DELAY_UE_CONTEXT_RELEASE_MS Delay for UE Context Release messages (ms)")
	fmt.Println("  DELAY_PDN_CONNECTIVITY_MS  Delay for PDN Connectivity messages (ms)")
	fmt.Println("  DELAY_HANDOVER_MS          Delay for Handover messages (ms)")
	fmt.Println("  DELAY_HANDOVER_REQUIRED_MS Delay for HandoverRequired messages (ms)")
	fmt.Println("  DELAY_HANDOVER_NOTIFY_MS   Delay for HandoverNotify messages (ms)")
	fmt.Println("  DELAY_RESET_MS             Delay for Reset messages (ms)")
	fmt.Println("  DELAY_PAGING_MS            Delay for Paging messages (ms)")
	fmt.Println("  DELAY_DEFAULT_MS           Default delay for other messages (ms)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Run with default settings")
	fmt.Println("  4g-proxy")
	fmt.Println()
	fmt.Println("  # Run with custom MME address")
	fmt.Println("  4g-proxy --mme 192.168.1.100 --mme-port 36412")
	fmt.Println()
	fmt.Println("  # Run with configuration file")
	fmt.Println("  4g-proxy -c config/config.yaml")
	fmt.Println()
	fmt.Println("  # Run with delay via environment variables")
	fmt.Println("  DELAY_ATTACH_MS=1000 DELAY_TAU_MS=500 4g-proxy")
	fmt.Println()
	fmt.Println("HTTP API Endpoints:")
	fmt.Println("  GET  /health                    Health check")
	fmt.Println("  GET  /api/v1/status             Get proxy status")
	fmt.Println("  GET  /api/v1/stats              Get statistics")
	fmt.Println("  GET  /api/v1/drop               Get drop flags")
	fmt.Println("  PUT  /api/v1/drop               Set multiple drop flags")
	fmt.Println("  DELETE /api/v1/drop             Reset all drop flags")
	fmt.Println("  GET  /api/v1/delay              Get delay settings")
	fmt.Println("  PUT  /api/v1/delay              Set multiple delays")
	fmt.Println("  PUT  /api/v1/delay/:type        Set delay for signal type")
	fmt.Println("  DELETE /api/v1/delay            Reset all delays")
	fmt.Println()
}
