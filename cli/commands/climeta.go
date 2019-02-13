package commands

import (
	"fmt"
	neuron "github.com/nikhilsbhat/neuron/app"
	config "github.com/nikhilsbhat/neuron/app/config"
	"github.com/nikhilsbhat/neuron/cli/ui"
	"os"
)

type cliMeta struct {
	*config.CliMeta
	*ui.NeuronUi
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
	nui := ui.NeuronUi{&ui.UiWriter{os.Stdout}}
	cc = &cliMeta{&meta, &nui}

}
