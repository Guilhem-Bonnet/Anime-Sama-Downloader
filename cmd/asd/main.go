package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Guilhem-Bonnet/Anime-Sama-Downloader/internal/domain"
)

func main() {
	baseURL := flag.String("server", envOr("ASD_SERVER_URL", "http://127.0.0.1:8080"), "URL du serveur (ex: http://127.0.0.1:8080)")
	timeout := flag.Duration("timeout", 10*time.Second, "Timeout HTTP")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(2)
	}

	client := &http.Client{Timeout: *timeout}

	switch args[0] {
	case "health":
		runGET(client, *baseURL+"/api/v1/health")
	case "version":
		runGET(client, *baseURL+"/api/v1/version")
	case "openapi":
		runGET(client, *baseURL+"/api/v1/openapi.json")
	case "settings":
		runSettings(client, *baseURL, args[1:])
	default:
		fmt.Fprintln(os.Stderr, "Commande inconnue:", args[0])
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: asd [flags] <command> ...")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  health")
	fmt.Fprintln(os.Stderr, "  version")
	fmt.Fprintln(os.Stderr, "  openapi")
	fmt.Fprintln(os.Stderr, "  settings get")
	fmt.Fprintln(os.Stderr, "  settings set [--destination PATH] [--output-naming legacy|media-server] [--separate-lang true|false] [--max-workers N] [--max-concurrent-downloads N]")
}

func runGET(client *http.Client, url string) {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	printResponse(resp.StatusCode, b)
}

func runJSON(client *http.Client, method string, url string, body any) {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Erreur:", err)
			os.Exit(1)
		}
		r = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur:", err)
		os.Exit(1)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	printResponse(resp.StatusCode, b)
}

func printResponse(status int, body []byte) {
	var pretty any
	if err := json.Unmarshal(body, &pretty); err == nil {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(pretty)
		if status >= 400 {
			os.Exit(1)
		}
		return
	}

	os.Stdout.Write(body)
	os.Stdout.Write([]byte("\n"))
	if status >= 400 {
		os.Exit(1)
	}
}

func runSettings(client *http.Client, baseURL string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: asd settings [get|set]")
		os.Exit(2)
	}

	switch args[0] {
	case "get":
		runGET(client, baseURL+"/api/v1/settings")
	case "set":
		settings := getSettings(client, baseURL)

		fs := flag.NewFlagSet("settings set", flag.ContinueOnError)
		fs.SetOutput(io.Discard)

		destination := newOptionalStringFlag("")
		outputNaming := newOptionalStringFlag("")
		separateLang := newOptionalBoolFlag(false)
		maxWorkers := newOptionalIntFlag(0)
		maxConcurrentDownloads := newOptionalIntFlag(0)

		fs.Var(destination, "destination", "Chemin racine destination")
		fs.Var(outputNaming, "output-naming", "Mode de nommage: legacy|media-server")
		fs.Var(separateLang, "separate-lang", "Séparer les langues (true/false)")
		fs.Var(maxWorkers, "max-workers", "Nombre de workers")
		fs.Var(maxConcurrentDownloads, "max-concurrent-downloads", "Téléchargements concurrents")

		if err := fs.Parse(args[1:]); err != nil {
			fmt.Fprintln(os.Stderr, "Erreur:", err)
			os.Exit(2)
		}

		if destination.set {
			settings.Destination = destination.value
		}
		if outputNaming.set {
			settings.OutputNamingMode = domain.OutputNamingMode(strings.TrimSpace(outputNaming.value))
		}
		if separateLang.set {
			settings.SeparateLang = separateLang.value
		}
		if maxWorkers.set {
			settings.MaxWorkers = maxWorkers.value
		}
		if maxConcurrentDownloads.set {
			settings.MaxConcurrentDownloads = maxConcurrentDownloads.value
		}

		runJSON(client, http.MethodPut, baseURL+"/api/v1/settings", settings)
	default:
		fmt.Fprintln(os.Stderr, "Commande inconnue:", args[0])
		os.Exit(2)
	}
}

func getSettings(client *http.Client, baseURL string) domain.Settings {
	resp, err := client.Get(baseURL + "/api/v1/settings")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		printResponse(resp.StatusCode, b)
		os.Exit(1)
	}
	var s domain.Settings
	if err := json.Unmarshal(b, &s); err != nil {
		fmt.Fprintln(os.Stderr, "Erreur:", err)
		os.Exit(1)
	}
	return s
}

type optionalStringFlag struct {
	set   bool
	value string
}

func newOptionalStringFlag(def string) *optionalStringFlag {
	return &optionalStringFlag{value: def}
}

func (f *optionalStringFlag) String() string { return f.value }
func (f *optionalStringFlag) Set(v string) error {
	f.set = true
	f.value = v
	return nil
}

type optionalIntFlag struct {
	set   bool
	value int
}

func newOptionalIntFlag(def int) *optionalIntFlag { return &optionalIntFlag{value: def} }
func (f *optionalIntFlag) String() string         { return fmt.Sprintf("%d", f.value) }
func (f *optionalIntFlag) Set(v string) error {
	var parsed int
	_, err := fmt.Sscanf(v, "%d", &parsed)
	if err != nil {
		return fmt.Errorf("invalid int: %s", v)
	}
	f.set = true
	f.value = parsed
	return nil
}

type optionalBoolFlag struct {
	set   bool
	value bool
}

func newOptionalBoolFlag(def bool) *optionalBoolFlag { return &optionalBoolFlag{value: def} }
func (f *optionalBoolFlag) String() string           { return fmt.Sprintf("%t", f.value) }
func (f *optionalBoolFlag) Set(v string) error {
	vv := strings.TrimSpace(strings.ToLower(v))
	switch vv {
	case "true", "1", "yes", "y", "on":
		f.value = true
	case "false", "0", "no", "n", "off":
		f.value = false
	default:
		return fmt.Errorf("invalid bool: %s", v)
	}
	f.set = true
	return nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
