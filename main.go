package main

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/krzysztofzaucha/protobuf-sandbox/internal"
	"os"
	"plugin"
	"runtime"
	"strings"
)

const (
	retries int = 3
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	help := flag.Bool("help", false, "print help")
	configFilePath := flag.String("config", "config.json", "configuration file (e.g. config.json)")
	module := flag.String("module", "", "module to be used (e.g. producer, consumer, etc...)")

	flag.Parse()

	// print help
	if *help == true {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// load configuration
	config, err := internal.LoadConfig(*configFilePath)
	if err != nil {
		panic(err)
	}

	// ensure module is provided, if not print help
	if *module == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	symName := loadPlugin(*module)

	configureWith(symName, config)

	// execute plugin logic
	var executor internal.Executor

	executor, ok := symName.(internal.Executor)
	if !ok {
		panic("plugin is not an executor")
	}

	err = executor.Execute()
	if err != nil {
		panic(err)
	}
}

func generateSymbolName(module string) string {
	return strings.ReplaceAll(strings.Title(strings.ReplaceAll(module, "-", " ")), " ", "")
}

func loadPlugin(module string) plugin.Symbol {
	// locate and load the plugin
	plug, err := plugin.Open("bin/" + module + ".so")
	if err != nil {
		panic(err)
	}

	symName, err := plug.Lookup(generateSymbolName(module))
	if err != nil {
		panic(err)
	}

	return symName
}

func configureWith(symbol plugin.Symbol, config *internal.Config) {
	// configure Config
	if _, ok := symbol.(internal.ConfigConfigurator); ok {
		symbol.(internal.ConfigConfigurator).WithConfig(config)
	}
}
