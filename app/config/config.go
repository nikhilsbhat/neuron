package config

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/nikhilsbhat/neuron/database"
	"github.com/nikhilsbhat/neuron/database/common"
	err "github.com/nikhilsbhat/neuron/error"
	log "github.com/nikhilsbhat/neuron/logger"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// AppConfig holds the entire confiration of the application neuron, while reading config file we decode the configuration to this struct below.
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

	// By enabling this you are asking neuron to serve UI.
	EnableUI bool `json:"enableui"`

	// It holds the details of the database that has to be connected with neuron.
	Database []*db `json:"database"`

	Cloud []*Cloud
}
type db struct {
	Name *string `json:"name"`
	Addr *string `json:"addr"`
}

// Cloud holds the details of the cloud, this will be set in the configuration.
type Cloud struct {
	Name    string `json:"name"`
	Region  string `json:"region"`
	Default bool   `json:"default"`
}

// ConfigResponse will be the response type of function Init.
type ConfigResponse struct {
	Port           string
	UiDir          string
	UiTemplatePath string
	ApiLogPath     io.Writer
	NoUi           bool
	EnableAPI      bool
}

// CliMeta holds data whether cli is enabled along with decoded config.
type CliMeta struct {
	CliSet bool
	*AppConfig
}

func (conf *AppConfig) createDirectories() error {

	log.Info("Creating required directories for Neuron to store data")
	for _, dir := range conf.Directories {
		dirpath := filepath.Join(conf.Home, dir)
		log.Info(fmt.Sprintf("The directory in creation is: %s", dirpath))
		if _, direrr := os.Stat(dirpath); os.IsNotExist(direrr) {
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

// ConfigNeuron is responsible for configuring neuron, this has to be called if api has to be exposed.
func ConfigNeuron(path ...string) error {

	conf, cferr := loadConfig(path)
	if cferr != nil {
		return cferr
	}

	//creatinig directories
	// just append to this array if in case any new directories has to be created for neuron in future
	conf.Directories = []string{"data"}
	direrr := conf.createDirectories()
	if direrr != nil {
		return direrr
	}
	log.Info(printSpace)
	log.Info(pritnDash)
	return nil
}

// Init will configure the application by reading the configuration file at '/var/lib/neuron'.
// Be sure what you pass as path to this, because only the first element is considered while setting path.
func Init(path ...string) (ConfigResponse, error) {

	conf, cferr := loadConfig(path)
	if cferr != nil {
		return ConfigResponse{}, cferr
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

	if (conf.EnableAPI == false) && (conf.EnableUI == true) {
		return ConfigResponse{}, fmt.Errorf("You cannot enable ui alone without api. Quitting installation")
	}

	cnfgerr := ConfigNeuron(path[0])
	if cnfgerr != nil {
		return ConfigResponse{}, nil
	}

	// configuring db
	conf.configDB()

	if conf.EnableAPI == true {
		api, apierr := conf.configApi()
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

func loadConfig(path []string) (AppConfig, error) {

	conf, conferr := findConfig(setCOnfigPath(path))
	if conferr != nil {
		return AppConfig{}, conferr
	}
	return conf, nil

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
			log.Info(fmt.Sprintf(" Provided configs are: %v ", conf.Database))
			if strings.ToLower(*dataBase.Name) == "mongodb" {
				log.Info("Found a compatible databse. Establishing connection....")
				dbsession, dberr := mgo.Dial(*dataBase.Addr)
				if dberr != nil {
					log.Error(fmt.Sprintf("Unable to reach %s which you provided", *dataBase.Name))
					log.Warn(switchToFs)
					_, dataerr := dbcommon.ConfigDb(database.Storage{Fs: fmt.Sprintf("%s/data/", conf.Home)})
					if dataerr != nil {
						log.Error(err.DbSessionError())
						log.Error(dataerr)
					}
				} else {
					_, dataerr := dbcommon.ConfigDb(database.Storage{Db: dbsession})
					if dataerr != nil {
						log.Error(err.DbSessionError())
						log.Error(dataerr)
					}
				}
			}
			if database.Db == nil {
				log.Warn("We do not support other database, only MongoDb is compatible for now")
				log.Warn(switchToFs)
				_, dataerr := dbcommon.ConfigDb(database.Storage{Fs: fmt.Sprintf("%s/data/", conf.Home)})
				if dataerr != nil {
					log.Error(err.DbSessionError())
					log.Error(dataerr)
				}
			}

		} else {
			log.Warn("Couldn't find database Config in neuron.json")
			log.Warn(switchToFs)
			_, dataerr := dbcommon.ConfigDb(database.Storage{Fs: fmt.Sprintf("%s/data/", conf.Home)})
			if dataerr != nil {
				log.Error(err.DbSessionError())
				log.Error(dataerr)
			}
		}
	}
}

// Configuring API happens here.
func (conf *AppConfig) configApi() (ConfigResponse, error) {

	var ui ConfigResponse
	// configuring ui
	if conf.EnableUI == true {
		uiresp, uierr := conf.configUI()
		if uierr != nil {
			return ConfigResponse{}, uierr
		}
		if reflect.DeepEqual(uiresp, ConfigResponse{}) {
			uiresp.NoUi = true
		}
		ui = uiresp
	}

	// configuring api log path
	apilogpath, uierr := conf.configapilogs()
	if uierr != nil {
		return ConfigResponse{}, uierr
	}
	ui.ApiLogPath = apilogpath
	ui.Port = conf.Port
	ui.EnableAPI = conf.EnableAPI
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
		if _, direrr := os.Stat(conf.UiDir); os.IsNotExist(direrr) {
			log.Error(err.UiNotFound())
			log.Warn(uiNotAvailable)
			log.Warn(printSpace)
			return response, nil
		}
		response.UiDir = conf.UiDir
		response.UiTemplatePath = fmt.Sprintf("%s/pages/*", conf.UiDir)

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

// GetCliMeta will be give back data of type CliMeta, to help decision making while we initialize cli.
// Be sure what you pass as path to this, because only the first element is considered while setting path.
func GetCliMeta(path ...string) (CliMeta, error) {

	conf, conferr := readConfig(setCOnfigPath(path))
	if conferr != nil {
		return CliMeta{}, conferr
	}

	if reflect.DeepEqual(conf, AppConfig{}) {
		log.Info(printSpace)
		log.Error(err.UninitializedCli())
		log.Info(printSpace)
		log.Error(err.ConfigNotfound())
		return CliMeta{}, err.CliFailure()
	}
	dberr := conf.prepareMinimalCli()
	if dberr != nil {
		return CliMeta{}, dberr
	}
	return CliMeta{true, &conf}, nil
}

// Database will be set here if it was mentioned in config file.
func (conf *AppConfig) prepareMinimalCli() error {

	for _, dataBase := range conf.Database {
		if val := reflect.DeepEqual(*dataBase, db{}); val != true {
			if strings.ToLower(*dataBase.Name) == "mongodb" {
				dberr := dataBase.switchtoDB(conf.Home)
				if dberr != nil {
					return dberr
				}
				return nil
			}
			if database.Db == nil {
				fserr := switchtoFS(conf.Home)
				if fserr != nil {
					return fserr
				}
				return nil
			}
		} else {
			fserr := switchtoFS(conf.Home)
			if fserr != nil {
				return fserr
			}
			return nil
		}
	}
	return fmt.Errorf("An Unknown error occurred while prepearing cli")
}

func switchtoFS(home string) error {
	_, dataerr := dbcommon.ConfigDb(database.Storage{Fs: fmt.Sprintf("%s/data/", home)})
	if dataerr != nil {
		return err.DbSessionError()
	}
	return nil
}

func (data *db) switchtoDB(home string) error {

	dbsession, dberr := mgo.Dial(*data.Addr)
	if dberr != nil {
		_, dataerr := dbcommon.ConfigDb(database.Storage{Fs: fmt.Sprintf("%s/data/", home)})
		if dataerr != nil {
			return err.DbSessionError()
		}
		return nil
	}
	_, dataerr := dbcommon.ConfigDb(database.Storage{Db: dbsession})
	if dataerr != nil {
		return err.DbSessionError()
	}
	return nil
}

func setCOnfigPath(path []string) string {

	if path != nil {
		return path[0]
	}
	return "/var/lib/neuron/neuron.json"
}
