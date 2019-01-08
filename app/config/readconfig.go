package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	neuerr "neuron/error"
	log "neuron/logger"
	"os"
	//"path/filepath"
)

func findConfig(pathtofile string) (AppConfig, error) {

	if pathtofile != "" {
		if _, dir_neuerr := os.Stat(pathtofile); os.IsNotExist(dir_neuerr) {
			return AppConfig{}, neuerr.NoFileFoundError()
		} else {
			config_data, confneuerr := readConfig(pathtofile)
			if confneuerr != nil {
				return AppConfig{}, confneuerr
			}
			return config_data, nil
		}
	}

	if _, dir_neuerr := os.Stat("/var/lib/neuron/neuron.json"); os.IsNotExist(dir_neuerr) {

		log.Info("+++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		log.Info("")
		log.Info("You did not provide the configuration file hence setting configurations to default")

		decoder := json.NewDecoder(bytes.NewReader([]byte(`{"port": "80","logfile": "neuron.log","loglocation": "/var/log/neuron/", "home": "/var/lib/neuron"}`)))
		var config_data AppConfig
		if decodneuerr := decoder.Decode(&config_data); decodneuerr != nil {
			log.Error(neuerr.JsonDecodeError())
			log.Error("Please provide us valid file")
			log.Error("Hence quitting installation...")
			return AppConfig{}, neuerr.JsonDecodeError()
		}
		return config_data, nil
	} else {

		log.Info("+++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		log.Info("")
		log.Info("Found configuration file, configuring application as per the entries")
		config_data, confneuerr := readConfig("/var/lib/neuron/neuron.json")
		if confneuerr != nil {
			return AppConfig{}, confneuerr
		}
		return config_data, nil
	}
}

func readConfig(pathtofile string) (AppConfig, error) {
	conf_file, confneuerr := ioutil.ReadFile(pathtofile)
	if confneuerr != nil {
		log.Error(neuerr.InvalidConfig())
		return AppConfig{}, neuerr.InvalidConfig()
	}

	decoder := json.NewDecoder(bytes.NewReader([]byte(conf_file)))
	var confdata AppConfig
	if decodneuerr := decoder.Decode(&confdata); decodneuerr != nil {
		log.Error(neuerr.JsonDecodeError())
		log.Error("Hence quitting installation...")
		return AppConfig{}, decodneuerr
	}
	return confdata, nil
}
