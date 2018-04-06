// This utility takes a alertmanager.yml configuration file and a set of
// Prometheus labels as normally found on an alert, and prints the
// Alertmanager receiver that would have been notified of a alert with the
// given label.  The goal is to build some form of unit testing of routes
// in the Alertmanager configuration.

package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

// Import the Alertmanager dispatch library to use the exact route processing
// code that Alertmanager does,
import (
	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/dispatch"
	"github.com/prometheus/common/model"
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(),
		"Evalute what receiver(s) are used when an alert with the given labels\n")
	fmt.Fprintf(flag.CommandLine.Output(),
		"is passed through the routing tree found in the Alertmanager config.\n\n")
	fmt.Fprintf(flag.CommandLine.Output(),
		"\t%s [options] label=value [label=value ...]\n\n", os.Args[0])
	flag.PrintDefaults()

	os.Exit(1)
}

func main() {
	var ctree config.Config
	flag.Usage = usage
	configFile := flag.String("f", "alertmanager.yml",
		"Alertmanager configuration file")
	expectedReceivers := flag.String("e", "",
		"Comma seperated list of expected alert receivers")
	flag.Parse()

	if flag.NArg() == 0 {
		usage()
	}

	// Read the configuration file
	blob, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Error reading %s: %s", *configFile, err.Error())
	}
	if err := yaml.Unmarshal([]byte(blob), &ctree); err != nil {
		log.Fatalf("Error unmarhsalling YAML: %s", err.Error())
	}

	// Parse label=value command line argument pairs
	lset := model.LabelSet{}
	for _, arg := range flag.Args() {
		fields := strings.SplitN(arg, "=", 2)
		if len(fields) == 2 {
			lset[model.LabelName(fields[0])] = model.LabelValue(fields[1])
		} else {
			log.Fatalf("Invalid label=value pair: %s", arg)
		}
	}
	if err := lset.Validate(); err != nil {
		log.Fatalf("Bad label set: %s", err.Error())
	}

	// Build the route tree
	routeTree := dispatch.NewRoute(ctree.Route, nil)

	// Okay, what does the routing say?
	routes := routeTree.Match(lset)
	results := []string{}
	for _, route := range routes {
		results = append(results, route.RouteOpts.Receiver)
	}
	sort.Strings(results)
	for _, receiver := range results {
		log.Printf("Executed receiver: %s", receiver)
	}

	// Compare to expected results
	if *expectedReceivers == "" {
		os.Exit(0)
	}
	expected := strings.Split(*expectedReceivers, ",")
	sort.Strings(expected)
	if len(expected) != len(results) {
		log.Printf("FAILED: number of receivers does not match number of expected")
		os.Exit(2)
	}
	for i := range expected {
		if expected[i] != results[i] {
			log.Printf("FAILED: Receiver(s) do not match expected.")
			os.Exit(2)
		}
	}
}
