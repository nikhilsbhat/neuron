package config

import (
	"encoding/json"
	neuerr "neuron/error"
	log "neuron/logger"
	"os"
	"bytes"
	"io/ioutil"
	//"path/filepath"
)

func findConfig() (config, error) {

	if _, dir_neuerr := os.Stat("/var/lib/neuron/neuron.json"); os.IsNotExist(dir_neuerr) {

		log.Info("+++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		log.Info("")
		log.Info("You did not provide the configuration file hence setting configurations to default")

		decoder := json.NewDecoder(bytes.NewReader([]byte(`{"port": "80","logfile": "neuron.log","logfile_location": "/var/log/neuron/", "home": "/var/lib/neuron"}`)))
		var config_data map[string]interface{}
		if decodneuerr := decoder.Decode(&config_data); decodneuerr != nil {
			log.Error(neuerr.JsonDecodeError())
			log.Error("Please provide us valid file")
			log.Error("Hence quitting installation...")
			return config{}, neuerr.JsonDecodeError()
		}
		conf := new(config)
		conf.appconfig = config_data
		return *conf, nil
	} else {

		log.Info("+++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		log.Info("")
		log.Info("Found configuration file, configuring application as per the entries")
		config_data, confneuerr := readConfig()
		if confneuerr != nil {
			return config{}, confneuerr
		}
		conf := new(config)
		conf.appconfig = config_data
		return *conf, nil
	}
}

func readConfig() (map[string]interface{}, error) {
	conf_file, confneuerr := ioutil.ReadFile("/var/lib/neuron/neuron.json")
	if confneuerr != nil {
		log.Error(neuerr.InvalidConfig())
		return nil, neuerr.InvalidConfig()
	}

	decoder := json.NewDecoder(bytes.NewReader([]byte(conf_file)))
	var confdata map[string]interface{}
	if decodneuerr := decoder.Decode(&confdata); decodneuerr != nil {
		log.Error(neuerr.JsonDecodeError())
		log.Error("Hence quitting installation...")
		return nil, neuerr.JsonDecodeError()
	}
	return confdata["config"].(map[string]interface{}), nil
}
