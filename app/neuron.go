// This package will help one to config neuron as per the entries made in the configuration file.
// This will help in enabling API and UI as well.
package neuron

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	conf "neuron/app/config"
	rou "neuron/app/routers"
	err "neuron/error"
	neulog "neuron/logger"
	"sync"
)

type Config struct {
	EnableAPI  bool
	ConfigPath string
}

var (
	ConfResponse = conf.ConfigResponse{}
)

// All one has to do is to call this function to config neuron.
func (c *Config) ConfigureNeuron() error {

	//Initializing log first before anyother thing
	logerr := neulog.Init()
	if logerr != nil {
		return logerr
	}

	//configuring neuron to prepare it for operations
	config, conferr := conf.ConfigNeuron(c.ConfigPath)
	if conferr != nil {
		return conferr
	}

	// Passing object of Config directly with out passing the UiDir doesn't work it throws error.
	if (c.EnableAPI == true) || (config.NoUi == false) || (config.EnableAPI == true) {
		EnableNeuronApi(config)
	}
	ConfResponse = config
	return nil
}

// API enablement will be done here by this function.
func EnableNeuronApi(config conf.ConfigResponse) {

	errCh := make(chan error, 1)
	//Initializing router to prepare neuron to serve endpoints
	rout := new(rou.MuxIn)
	if config.NoUi != true {
		rout.UiDir = config.UiDir
		rout.UiTemplatePath = config.UiTemplatePath
	}
	rout.Apilog = config.ApiLogPath
	router := rout.NewRouter()

	type r struct {
		router *mux.Router
		port   string
	}
	neurouter := r{router: router, port: config.Port}
	//starting the neuron on specified port
	var wg sync.WaitGroup
	wg.Add(1)
	go func(neurouter r) {
		starterr := http.ListenAndServe(":"+neurouter.port, neurouter.router)
		if starterr != nil {
			neulog.Error(err.StartNeuronError())
			errCh <- err.StartNeuronError()
		}
	}(neurouter)
	httperr := <-errCh
	if httperr != nil {
		log.Fatal(httperr)
	}
}

// All one has to do is to call this function to config neuron.
func NeuronCliMeta() (conf.CliMeta, error) {

	//Initializing log first before anyother thing
	logerr := neulog.Init()
	if logerr != nil {
		return conf.CliMeta{}, logerr
	}

	//configuring neuron to prepare it for operations
	config, conferr := conf.GetCliMeta()
	if conferr != nil {
		return conf.CliMeta{}, conferr
	}

	return config, nil
}
