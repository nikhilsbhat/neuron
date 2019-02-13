package NeuronLogger

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	err "github.com/nikhilsbhat/neuron/error"
	"os"
)

const (
	pritnDash       = "+++++++++++++++++++++++++++++++++++++++++++++++++++++++"
	loginitializing = "Logging is atmost important than anything other step hence preparing for it"
	printemptyline  = ""
)

func getlog() (string, error) {

	if _, dir_err := os.Stat("/var/lib/neuron/neuron.json"); os.IsNotExist(dir_err) {

		Info(pritnDash)
		Info(loginitializing)
		Info(printemptyline)
		return "/var/log/neuron", nil
	} else {

		config_data, conferr := getloglocation()
		if conferr != nil {
			switch conferr.(type) {
			case err.NoLogFound:
				return "/var/log/neuron", nil
			default:
				return "", conferr
			}
		}
		return config_data, nil
	}
}

func getloglocation() (string, error) {
	conf_file, conferr := ioutil.ReadFile("/var/lib/neuron/neuron.json")
	if conferr != nil {
		return "", err.ReadFileError()
	}
	decoder := json.NewDecoder(bytes.NewReader([]byte(conf_file)))

	var confdata map[string]interface{}
	if decoderr := decoder.Decode(&confdata); decoderr != nil {
		Error(err.JsonDecodeError())
		return "", err.InvalidConfig()
	}

	if confdata["loglocation"] != nil {
		return (confdata["loglocation"]).(string), nil
	}
	return "", err.LogNotFound()
}
