// Package neuron will help one to config neuron as per the entries made in the configuration file.
// This will help in enabling API and UI as well.
package neuron

import (
        "github.com/gorilla/mux"
        conf "github.com/nikhilsbhat/neuron/app/config"
        rou "github.com/nikhilsbhat/neuron/app/routers"
        err "github.com/nikhilsbhat/neuron/error"
        neulog "github.com/nikhilsbhat/neuron/logger"
        "log"
        "net/http"
        "sync"
)

// Config will help one to set few things before they call the methods to initialize/start neuron.
type Config struct {
        EnableAPI  bool
        ConfigPath string
}

// This holds configuration response from ConfigNeuron, this will let other functions/methods take decision further.
var (
        ConfResponse = conf.ConfigResponse{}
)

// Init will help in configuring neuron, by default other methods takes care of it.
func (c *Config) Init() error {

        //Initializing log first before anyother thing
        logerr := neulog.Init()
        if logerr != nil {
                return logerr
        }

        //configuring neuron to run it's service
        config, conferr := conf.Init(c.ConfigPath)
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

// ConfigureNeuron will call method ConfigNeuron of config to configure neron, so that it will be operable.
func (c *Config) ConfigureNeuron() error {
        //Initializing log first before anyother thing
        logerr := neulog.Init()
        if logerr != nil {
                return logerr
        }

        //configuring neuron to prepare it for operations
        conferr := conf.ConfigNeuron(c.ConfigPath)
        if conferr != nil {
                return conferr
        }
        return nil
}

// EnableNeuronApi will help in enablement of API.
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

// CliMeta will help in configuring neuron amd make it callable from cli, all one has to do is to call this function before using cli.
func CliMeta() (conf.CliMeta, error) {

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
