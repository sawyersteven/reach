package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	NoColor      bool `short:"c" long:"nocolor" description:"Print output without colors."`
	Timeout      int  `long:"timeout" default:"15" description:"HTTP request timeout in seconds"`
	MaxRedirects int  `long:"maxredirects" default:"20" description:"Maximum redirects to follow"`
	Help         bool `long:"help" description:"Display usage instructions"`
	Version      bool `long:"version" description:"Display version and license information"`
}

var trace = &httptrace.ClientTrace{
	PutIdleConn: func(err error) {
		if err != nil {
			handleTransportError("Could Not Finish Connection", err)
		} else {
			printTransportProgress("Connection finished")
		}
	},

	GotFirstResponseByte: func() {
		printTransportProgress("Receiving Response")
	},

	Got100Continue: func() {
		printTransportProgress("Received 100 Response - Waiting... ")
	},

	DNSStart: func(i httptrace.DNSStartInfo) {
		printTransportProgress("Starting DNS Lookup")
	},

	DNSDone: func(i httptrace.DNSDoneInfo) {
		if i.Err != nil {
			handleTransportError("DNS Lookup Failed", i.Err)
		} else {
			printTransportProgress("DNS Lookup Complete")
		}
	},

	ConnectStart: func(network, addr string) {
		printTransportProgress("Connection Started")
	},

	ConnectDone: func(network, addr string, err error) {
		if err != nil {
			handleTransportError("Connection Failed", err)
		} else {
			printTransportProgress("Connected - waiting for response...")
		}
	},

	TLSHandshakeDone: func(state tls.ConnectionState, err error) {
		if err != nil {
			handleTransportError("TLS Handshake Failed", err)
		} else {
			printTransportProgress("TLS Handshake Complete.")
		}
	},
}

func main() {
	args := parseArgs()

	if opts.Help {
		printHelp()
		os.Exit(0)
	} else if opts.Version {
		printVersion()
		os.Exit(0)
	}

	getURL(args[0])
}

func parseArgs() []string {
	p := flags.NewParser(&opts, 0)

	args, err := p.Parse()

	if err != nil || len(args) < 1 {
		printHelp()
		os.Exit(0)
	}
	return args
}

func getURL(targetURL string) {
	/* Attempts to get url while explicitly handling redirects

	Makes a custom Transport (roundtripper) to specify timeout

	Validates URL and sends request, which uses trace to follow
		step-by-step through the request process

	If a redirect response is returned it is printed then followed
	*/

	var roundTripper = &http.Transport{
		ResponseHeaderTimeout: time.Duration(opts.Timeout) * time.Second,
	}

	var nextRequestURL string

	if !strings.HasPrefix(targetURL, "https://") && !strings.HasPrefix(targetURL, "http://") {
		nextRequestURL = "http://" + targetURL
	} else {
		nextRequestURL = targetURL
	}

	for i := 0; i < opts.MaxRedirects; i++ {

		if !verifyURL(nextRequestURL) {
			handleTransportError("Invalid URL", errors.New("Unable to parse '"+nextRequestURL+"'"))
		}

		request, err := http.NewRequest("HEAD", nextRequestURL, nil)
		if err != nil {
			panic(err)
		}

		request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))

		response, err := roundTripper.RoundTrip(request)

		if err != nil {
			switch err := err.(type) {
			case net.Error:
				if err.Timeout() {
					handleTransportError("Request Failed", errors.New("Request timed out before a response was received."))
				}
			default:
				handleTransportError("Request Failed", err)

			}
			os.Exit(0)
		}
		printResponseInfo(response)

		if response.StatusCode/100 == 3 {
			nextRequestURL = response.Header.Get("Location")
			fmt.Println()
		} else {
			break
		}
	}
	fmt.Println()
}

func verifyURL(u string) bool {
	parsed, err := url.Parse(u)
	if err != nil ||
		parsed.Scheme == "" ||
		parsed.Host == "" ||
		(parsed.Scheme != "http" && parsed.Scheme != "https") {
		return false
	}
	return true
}

func printClear() {
	// If this doesn't work just make it longer
	fmt.Print("\r                                   \r")
}

func printTransportProgress(name string) {
	printClear()
	fmt.Printf("%v", name)
}

func handleTransportError(name string, err error) {
	printClear()
	if opts.NoColor {
		fmt.Printf("%s : %s", name, err.Error())
	} else {
		fmt.Printf("\x1b[91m%s:\x1b[0m %s", name, err.Error())
	}
	os.Exit(0)
}

func printResponseInfo(response *http.Response) {
	var bg int
	var fg int
	var add string

	printClear()
	if opts.NoColor {
		fmt.Printf("%s %s", response.Status, add)
	} else {
		switch response.StatusCode / 100 {
		case 2:
			bg = 102 // Green
			fg = 30  // Black
		case 3:
			bg = 106 // Blue
			fg = 30  // Black
			add = fmt.Sprintf("-> %s", response.Header.Get("Location"))
		case 4:
			bg = 101 // Red
			fg = 30  // Black
		case 5:
			bg = 105 // Purple
			fg = 30  // Black
		}

		statusDescription := fmt.Sprintf("\x1b[107m\x1b[30m %s \x1b[0m", response.Status[4:]) // Black on White

		fmt.Printf("\x1b[%dm\x1b[%dm %d \x1b[0m%s %s", bg, fg, response.StatusCode, statusDescription, add)
	}
}

func printHelp() {
	fmt.Println(`Usage: reach [OPTIONS] URL

Options:
  -c, --nocolor               Print output without colors
  --maxredirects=REDIRECTS    Maximum redirects to follow [default: 20]
  --timeout=SECONDS           HTTP request timeout in seconds [default: 15]
  --help                      Display this help message
  --version                   Display version and license info
  `)
}

func printVersion() {
	fmt.Println(`reach 0.1
Copyright (C) 2016 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Written by Steven Sawyer.`)
}
