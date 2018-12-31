package neuron

import (
	"net/http"
	conf "neuron/app/config"
	rou "neuron/app/routers"
	err "neuron/error"
	log "neuron/logger"
)

func StartNeuron() error {

	//Initializing log first before anyother thing
	logerr := log.Init()
	if logerr != nil {
		return logerr
	}

	//configuring neuron to prepare it for operations
	config, conferr := conf.ConfigNeuron()
	if conferr != nil {
		return conferr
	}

	//Initializing router to prepare neuron to serve endpoints
	rout := new(rou.MuxIn)
	if config.NoUi != true {
		rout.UiDir = config.UiDir
		rout.UiTemplatePath = config.UiTemplatePath
	}
	rout.Apilog = config.ApiLogPath
	router := rout.NewRouter()

	//starting the neuron on specified port
	starterr := http.ListenAndServe(":"+config.Port, router)
	if starterr != nil {
		log.Error(err.StartNeuronError())
		return err.StartNeuronError()
	}

	return nil
}
