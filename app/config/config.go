package config

import (
	"fmt"
	"github.com/globalsign/mgo"
	"io"
	"neuron/database"
	"neuron/database/common"
	err "neuron/error"
	log "neuron/logger"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	directories []string
	appconfig   map[string]interface{}
}

type ConfigResponse struct {
	Port           string
	UiDir          string
	UiTemplatePath string
	ApiLogPath     io.Writer
	NoUi           bool
}

func (c *config) createDirectories() error {

	log.Info("Creating required directories for Neuron to store data")
	for _, dir := range c.directories {
		dirpath := filepath.Join(c.appconfig["home"].(string), dir)
		log.Info(fmt.Sprintf("The directory in creation is: %s", dirpath))
		if _, dir_err := os.Stat(dirpath); os.IsNotExist(dir_err) {
			direrr := os.Mkdir(dirpath, 0644)
			if direrr != nil {
				log.Error(fmt.Sprintf("%s : %s", err.DirCreateError(), dirpath))
				return direrr
			}
		} else {
			log.Info("Skipping directories creation as directories exists")
		}
	}
	return nil
}

func ConfigNeuron() (ConfigResponse, error) {

	conf, conferr := findConfig()
	if conferr != nil {
		return ConfigResponse{}, conferr
	}

	if conf.appconfig == nil {
		log.Info(printSpace)
		log.Error(err.ConfigNotfound())
		log.Error(quitInstallation)
		log.Error(".....Quitting the installation process.....")
		log.Error(endOfLog)
		log.Error(pritnDash)
	}

	log.Info("Found configuration, configuring application as per the config file.....")

	//creatinig directories
	// just append to this array if in case any new directories has to be created in neuron in future
	conf.directories = []string{"data"}
	direrr := conf.createDirectories()
	if direrr != nil {
		return ConfigResponse{}, direrr
	}

	// configuring ui
	ui := conf.configUI()
	if ui == (ConfigResponse{}) {
		ui.NoUi = true
	}

	// configuring db
	conf.configDB()

	// configuring ui log path
	uilogpath, uierr := conf.configapilogs()
	if uierr != nil {
		return ConfigResponse{}, uierr
	}
	ui.ApiLogPath = uilogpath
	ui.Port = conf.appconfig["port"].(string)
	log.Info(printSpace)
	log.Info(endOfLog)
	log.Info(printSpace)
	log.Info(pritnDash)
	return ui, nil
}

func convertKeysToSlice(m map[string]interface{}) []string {

	ret := make([]string, 0)
	for key, _ := range m {
		ret = append(ret, key)
	}
	return ret
}

func (conf *config) configDB() {

	if conf.appconfig["database"] != nil {
		log.Info("Found Config for database")
		log.Info(fmt.Sprintf(" Provided configs are: %s ", convertKeysToSlice(conf.appconfig["database"].(map[string]interface{}))))
		if len(conf.appconfig["database"].(map[string]interface{})) > 1 {
			log.Warn(multipleDb)
			log.Warn("Anyways we will check for the compatible databse. If we find one, will establish connection with it")
		}
		for name, ip := range conf.appconfig["database"].(map[string]interface{}) {
			if strings.ToLower(name) == "mongodb" {
				log.Info("Found a compatible databse. Establishing connection....")
				db_session, dberr := mgo.Dial(ip.(string))
				if dberr != nil {
					data_path := fmt.Sprintf("%s/data/", conf.appconfig["home"].(string))
					log.Error(fmt.Sprintf("Unable to reach %s which you provided", name))
					log.Warn(switchToFs)
					_, data_err := dbcommon.ConfigDb(database.Storage{Fs: data_path})
					if data_err != nil {
						log.Error(err.DbSessionError())
						log.Error(data_err)
					}
				} else {
					_, data_err := dbcommon.ConfigDb(database.Storage{Db: db_session})
					if data_err != nil {
						log.Error(err.DbSessionError())
						log.Error(data_err)
					}
				}
			}
		}
		if database.Db == nil {
			data_path := fmt.Sprintf("%s/data/", conf.appconfig["home"].(string))
			log.Warn("We do not support other database, only MongoDb is compatible for now")
			log.Warn(switchToFs)
			_, data_err := dbcommon.ConfigDb(database.Storage{Fs: data_path})
			if data_err != nil {
				log.Error(err.DbSessionError())
				log.Error(data_err)
			}
		}

	} else {
		data_path := fmt.Sprintf("%s/data/", conf.appconfig["home"].(string))
		log.Warn("Couldn't find database config in neuron.json")
		log.Warn(switchToFs)
		_, data_err := dbcommon.ConfigDb(database.Storage{Fs: data_path})
		if data_err != nil {
			log.Error(err.DbSessionError())
			log.Error(data_err)
		}
	}
}

func (conf *config) configUI() ConfigResponse {

	var response ConfigResponse
	if conf.appconfig["uidir"] != "" {

		log.Info(printSpace)
		log.Info("Found configuration file, reading it to gather information regarding UI")
		if _, dir_err := os.Stat(conf.appconfig["uidir"].(string)); os.IsNotExist(dir_err) {
			log.Error(err.UiNotFound())
			log.Warn(uiNotAvailable)
			log.Warn(printSpace)
			return response
		}
		response.UiDir = conf.appconfig["uidir"].(string)
		response.UiTemplatePath = fmt.Sprintf("%s/pages/*", conf.appconfig["uidir"].(string))
		log.Info("...Awesome UI configured successfully...")
		log.Info("")
		return response
	}

	log.Warn("I could not find any ui directory path in the configuration")
	log.Warn(uiNotAvailable)
	log.Warn(printSpace)
	return response
}

func (conf *config) configapilogs() (io.Writer, error) {

	var (
		logfile string
		logpath string
	)
	if conf.appconfig["loglocation"] != nil {
		logpath = conf.appconfig["loglocation"].(string)
	} else {
		logpath = "/var/log/neuron"
	}

	if conf.appconfig["logfile"] != nil {
		logfile = conf.appconfig["logfile"].(string) + "/"
	} else {
		logfile = "/neuron.log"
	}

	loglocation := filepath.Join(logpath, logfile)

	if _, err1 := os.Stat(loglocation); os.IsNotExist(err1) {
		newfile, err2 := os.Create(loglocation)
		if err2 != nil {
			log.Error(err.UiLogCreationError())
			return nil, err.UiLogCreationError()
		}
		newfile.Close()

		path, logfilerr := os.OpenFile(loglocation, os.O_APPEND|os.O_WRONLY, 0644)
		if logfilerr != nil {
			log.Error(err.UiLogOpenError())
			return nil, err.UiLogOpenError()
		}
		return path, nil
	} else {

		path, logfilerr := os.OpenFile(loglocation, os.O_APPEND|os.O_WRONLY, 0644)
		if logfilerr != nil {
			log.Error(err.UiLogOpenError())
			return nil, err.UiLogOpenError()
		}
		return path, nil
	}
}
