package commands

import (
	"fmt"
	//"github.com/spf13/cobra"
	//err "neuron/error"
	//"github.com/fatih/color"
	neuron "neuron/app"
	config "neuron/app/config"
	"os"
)

type cliMeta struct {
	*config.CliMeta
}

var (
	cc = &cliMeta{}
)

func init() {

	meta, clierr := neuron.NeuronCliMeta()
	if clierr != nil {
		fmt.Println(clierr)
		os.Exit(3)
	}
	cc = &cliMeta{&meta}
}
