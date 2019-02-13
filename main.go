package main

import (
	//"neuron/app"
	cli "github.com/nikhilsbhat/neuron/cli"
	/*err "neuron/error"
	log "neuron/logger"*/)

//This function is responsible for starting the application.
func main() {

	/*neuerr := neuron.StartNeuron()
		if neuerr != nil {
	                log.Error(neuerr)
			log.Error(err.FailStartError())
		}*/
	cli.CliMain()
}
