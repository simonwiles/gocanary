package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
)

const AppVersion = "1.1"

func main() {

	fs := ff.NewFlagSet("gocanary")
	var (
		version    = fs.BoolLong("version", "prints current app version")
		host       = fs.String('h', "host", "localhost", "Port to run the server on")
		port       = fs.Uint('p', "port", 8930, "Port to run the server on")
		alertExprs = fs.StringSetLong("alert-when", "alert rule (repeatable)")
		_          = fs.String('c', "config", "", "Path to config file")
	)

	if err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("GC"),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
	); err != nil {
		fmt.Printf("%s\n", ffhelp.Flags(fs))
		fmt.Printf("err=%v\n", err)
		os.Exit(0)
	}

	if *version {
		fmt.Printf("gocanary %s\n", AppVersion)
		os.Exit(0)
	}

	alertRules, err := compileAlerts(alertExprs)
	if err != nil {
		fmt.Printf("Failed to compile alert rules, with error:\n%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Listening on %s:%d...\n", *host, *port)
	http.Handle("/", buildResponseHandler(alertRules))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil))
}
