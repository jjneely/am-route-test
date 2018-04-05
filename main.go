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
	"strings"
)

// Import the Alertmanager dispatch library to use the exact route processing
// code that Alertmanager does,
import (
	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/dispatch"
	"github.com/prometheus/common/model"
)

func main() {
	var ctree config.Config
	configFile := flag.String("-f", "alertmanager.yml",
		"Alertmanager configuration file")
	flag.Parse()

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
	for _, route := range routes {
		fmt.Printf("%s\n", route.RouteOpts.Receiver)
	}
}
