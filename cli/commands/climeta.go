package commands

import (
	"fmt"
	"os"

	neuron "github.com/nikhilsbhat/neuron/app"
	config "github.com/nikhilsbhat/neuron/app/config"
	"github.com/nikhilsbhat/neuron/cli/ui"
)

type cliMeta struct {
	*config.CliMeta
	*ui.NeuronUi
}

var (
	cc = &cliMeta{}
)

func init() {

	meta, clierr := neuron.CliMeta()
	if clierr != nil {
		fmt.Println(clierr)
		os.Exit(3)
	}
	nui := ui.NeuronUi{&ui.UiWriter{os.Stdout}}
	cc = &cliMeta{&meta, &nui}

}
