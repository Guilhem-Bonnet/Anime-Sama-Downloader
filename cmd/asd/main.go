package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	baseURL := flag.String("server", envOr("ASD_SERVER_URL", "http://127.0.0.1:8080"), "URL du serveur (ex: http://127.0.0.1:8080)")
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout HTTP")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: asd [health|version]")
		os.Exit(2)
	}

	client := &http.Client{Timeout: *timeout}

	switch args[0] {
	case "health":
		run(client, *baseURL+"/api/v1/health")
	case "version":
		run(client, *baseURL+"/api/v1/version")
	default:
		fmt.Fprintln(os.Stderr, "Commande inconnue:", args[0])
		os.Exit(2)
	}
}

func run(client *http.Client, url string) {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	var pretty any
	if err := json.Unmarshal(b, &pretty); err == nil {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(pretty)
		if resp.StatusCode >= 400 {
			os.Exit(1)
		}
		return
	}

	os.Stdout.Write(b)
	os.Stdout.Write([]byte("\n"))
	if resp.StatusCode >= 400 {
		os.Exit(1)
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
