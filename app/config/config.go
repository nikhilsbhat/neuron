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
	"reflect"
	"strings"
)

// This holds the entire confiration of the application neuron
type AppConfig struct {
	// Directories that has to be created to store application related data.
	// This will be used if filesystem is used as storage mode.
	Directories []string `json:"directories"`

	// Pass raw application configuration in a json format
	Rawconfig map[string]interface{} `json:"rawconfig"`

	// Port on which neuron has to operate.
	Port string `json:"port"`

	// Home directory or neuron.
	Home string `json:"home"`

	// Path of folder containing all the files which powers up UI.
	UiDir string `json:"uidir"`

	// Name of the log file, specify if you want to override the default name.
	LogFile string `json:"logfile"`

	// Path where logs has to be written.
	LogLocation string `json:"loglocation"`

	// By enabling this you are telling application to expose its API.
	EnableAPI bool `json:"enableapi"`

	// It holds the details of the database that has to be connected with neuron.
	Database []*db `json:"database"`
}
type db struct {
	name *string `json:"name"`
	addr *string `json:"addr"`
}

type ConfigResponse struct {
	Port           string
	UiDir          string
	UiTemplatePath string
	ApiLogPath     io.Writer
	NoUi           bool
}

func (c *AppConfig) createDirectories() error {

	log.Info("Creating required directories for Neuron to store data")
	for _, dir := range c.Directories {
		dirpath := filepath.Join(c.Home, dir)
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

// This will configure the application by reading the configuration file at '/var/lib/neuron'.
func ConfigNeuron(path string) (ConfigResponse, error) {

	var conf AppConfig
	var pathtofile string

	if path != "" {
		pathtofile = path
	}
	conf, conferr := findConfig(pathtofile)
	if conferr != nil {
		return ConfigResponse{}, conferr
	}

	if reflect.DeepEqual(conf, AppConfig{}) {
		log.Info(printSpace)
		log.Error(err.ConfigNotfound())
		log.Error(quitInstallation)
		log.Error(".....Quitting the installation process.....")
		log.Error(endOfLog)
		log.Error(pritnDash)
		return ConfigResponse{NoUi: true}, nil
	}

	log.Info("Found configuration, configuring application as per the config file.....")

	//creatinig directories
	// just append to this array if in case any new directories has to be created for neuron in future
	conf.Directories = []string{"data"}
	direrr := conf.createDirectories()
	if direrr != nil {
		return ConfigResponse{}, direrr
	}

	// configuring db
	conf.configDB()

	if conf.EnableAPI == true {
		api, apierr := conf.ConfigApi()
		if apierr != nil {
			return ConfigResponse{}, apierr
		}
		return api, nil
	}
	return ConfigResponse{}, nil
}

func convertKeysToSlice(m map[string]interface{}) []string {

	ret := make([]string, 0)
	for key, _ := range m {
		ret = append(ret, key)
	}
	return ret
}

// Database will be set here if it was mentioned in config file.
func (conf *AppConfig) configDB() {

	if len(conf.Database) > 1 {
		log.Warn(multipleDb)
		log.Warn("Anyways we will check for the compatible databse. If we find one, will establish connection with it")
	}
	for _, dataBase := range conf.Database {
		if val := reflect.DeepEqual(*dataBase, db{}); val != true {
			log.Info("Found Config for database")
			log.Info(fmt.Sprintf(" Provided configs are: %s ", conf.Database))
			if strings.ToLower(*dataBase.name) == "mongodb" {
				log.Info("Found a compatible databse. Establishing connection....")
				db_session, dberr := mgo.Dial(*dataBase.addr)
				if dberr != nil {
					log.Error(fmt.Sprintf("Unable to reach %s which you provided", *dataBase.name))
					log.Warn(switchToFs)
					_, data_err := dbcommon.ConfigDb(database.Storage{Fs: fmt.Sprintf("%s/data/", conf.Home)})
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
			if database.Db == nil {
				log.Warn("We do not support other database, only MongoDb is compatible for now")
				log.Warn(switchToFs)
				_, data_err := dbcommon.ConfigDb(database.Storage{Fs: fmt.Sprintf("%s/data/", conf.Home)})
				if data_err != nil {
					log.Error(err.DbSessionError())
					log.Error(data_err)
				}
			}

		} else {
			log.Warn("Couldn't find database Config in neuron.json")
			log.Warn(switchToFs)
			_, data_err := dbcommon.ConfigDb(database.Storage{Fs: fmt.Sprintf("%s/data/", conf.Home)})
			if data_err != nil {
				log.Error(err.DbSessionError())
				log.Error(data_err)
			}
		}
	}
}

// Configuring API happens here.
func (conf *AppConfig) ConfigApi() (ConfigResponse, error) {

	// configuring ui
	ui, uierr := conf.configUI()
	if uierr != nil {
		return ConfigResponse{}, uierr
	}
	if ui == (ConfigResponse{}) {
		ui.NoUi = true
	}

	ui.Port = conf.Port
	log.Info(printSpace)
	log.Info(endOfLog)
	log.Info(printSpace)
	log.Info(pritnDash)
	return ui, nil
}

// UI enablement happens here.
func (conf *AppConfig) configUI() (ConfigResponse, error) {

	var response ConfigResponse
	if conf.UiDir != "" {

		log.Info(printSpace)
		log.Info("Found configuration file, reading it to gather information regarding UI")
		if _, dir_err := os.Stat(conf.UiDir); os.IsNotExist(dir_err) {
			log.Error(err.UiNotFound())
			log.Warn(uiNotAvailable)
			log.Warn(printSpace)
			return response, nil
		}
		response.UiDir = conf.UiDir
		response.UiTemplatePath = fmt.Sprintf("%s/pages/*", conf.UiDir)
		// configuring ui log path
		uilogpath, uierr := conf.configapilogs()
		if uierr != nil {
			return ConfigResponse{}, uierr
		}
		response.ApiLogPath = uilogpath
		log.Info("...Awesome UI configured successfully...")
		log.Info("")
		return response, nil
	}

	log.Warn("I could not find any ui directory path in the configuration")
	log.Warn(uiNotAvailable)
	log.Warn(printSpace)
	return response, nil
}

// Name of the method specifies the work of it.
func (conf *AppConfig) configapilogs() (io.Writer, error) {

	var (
		logfile string
		logpath string
	)
	if conf.LogLocation != "" {
		logpath = conf.LogLocation
	} else {
		logpath = "/var/log/neuron"
	}

	if conf.LogFile != "" {
		logfile = conf.LogFile + "/"
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
