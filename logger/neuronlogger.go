package NeuronLogger

import (
	//"errors"
	"fmt"
	"io"
	"log"
	ui "github.com/nikhilsbhat/neuron/cli/ui"
	err "github.com/nikhilsbhat/neuron/error"
	"os"
	"runtime"
	"time"
)

//AppName is set to "neuron" by default, one can change it by passing the required name to neuronlogInitializer()
var (
	AppName = "neuron"
	Logpath io.Writer
)

type NeuronLogInput struct {
	Logpath string
	AppName string `json:"AppName,omitempty"`
}

type loggerOptions struct {
	appname string
	logpath io.Writer
	msg     string
	level   string
	caller  string `json:"caller,omitempty"`
}

func Init() error {

	conf, conferr := getlog()
	if conferr != nil {
		return conferr
	}
	login := new(NeuronLogInput)
	/*ui := NeuronUi{&UiWriter{os.Stdout}}
	  login.Ui = &ui*/
	login.Logpath = conf
	logerr := login.neuronlogInitializer()
	if logerr != nil {
		return err.SetupLogError()
	}
	return nil
}

// To use the custom logging functionality from this package once has to invoke below function
func (loger *NeuronLogInput) neuronlogInitializer() error {

	if _, dir_err := os.Stat(loger.Logpath); os.IsNotExist(dir_err) {

		logdirerr := os.Mkdir(loger.Logpath, 0644)
		if logdirerr != nil {
			return err.LogDirError()
		}
	}

	if loger.Logpath != "" {
		if _, err1 := os.Stat(loger.Logpath + "/neuronapp.log"); os.IsNotExist(err1) {
			newfile, err2 := os.Create(loger.Logpath + "/neuronapp.log")
			if err2 != nil {
				return err.LogCreationError()
			}
			newfile.Close()

			if loger.AppName != "" {
				AppName = loger.AppName
			}
			logpath, logfilerr := os.OpenFile(loger.Logpath+"/neuronapp.log", os.O_APPEND|os.O_WRONLY, 0644)
			if logfilerr != nil {
				return err.LogOpenError()
			}
			Logpath = logpath
			return nil
		} else {
			if loger.AppName != "" {
				AppName = loger.AppName
			}
			logpath, logfilerr := os.OpenFile(loger.Logpath+"/neuronapp.log", os.O_APPEND|os.O_WRONLY, 0644)
			if logfilerr != nil {
				return err.LogOpenError()
			}
			Logpath = logpath
			return nil
		}
	} else {
		if _, err1 := os.Stat("/var/log/neuron/neuronapp.log"); os.IsNotExist(err1) {
			newfile, err2 := os.Create("/var/log/neuron/neuronapp.log")
			if err2 != nil {
				return err.LogCreationError()
			}
			newfile.Close()

			if loger.AppName != "" {
				AppName = loger.AppName
			}
			logpath, logfilerr := os.OpenFile("/var/log/neuron/neuronapp.log", os.O_APPEND|os.O_WRONLY, 0644)
			if logfilerr != nil {
				return err.LogOpenError()
			}
			Logpath = logpath
			return nil
		} else {
			if loger.AppName != "" {
				AppName = loger.AppName
			}
			logpath, logfilerr := os.OpenFile("/var/log/neuron/neuronapp.log", os.O_APPEND|os.O_WRONLY, 0644)
			if logfilerr != nil {
				return err.LogOpenError()
			}
			Logpath = logpath
			return nil
		}
	}
}

func Info(data interface{}) {
	login := loggerOptions{
		level:   " [INFO] ",
		appname: AppName,
		msg:     getStringOfMessage(data),
		logpath: Logpath,
	}
	if _, file, no, ok := runtime.Caller(3); ok {
		login.caller = fmt.Sprintf("%s:%d ", file, no)
	}
	login.appLog()
}

func Error(data interface{}) {
	login := loggerOptions{
		level:   " [ERROR] ",
		appname: AppName,
		msg:     getStringOfMessage(data),
		logpath: Logpath,
	}
	if _, file, no, ok := runtime.Caller(1); ok {
		login.caller = fmt.Sprintf("%s:%d ", file, no)
	}
	login.appLog()
}

func Warn(data interface{}) {
	login := loggerOptions{
		level:   " [WARN] ",
		appname: AppName,
		msg:     getStringOfMessage(data),
		logpath: Logpath,
	}
	if _, file, no, ok := runtime.Caller(3); ok {
		login.caller = fmt.Sprintf("%s:%d ", file, no)
	}
	login.appLog()
}

func Debug(data interface{}) {
	login := loggerOptions{
		level:   " [DEBUG] ",
		appname: AppName,
		msg:     getStringOfMessage(data),
		logpath: Logpath,
	}
	if _, file, no, ok := runtime.Caller(3); ok {
		login.caller = fmt.Sprintf("%s:%d ", file, no)
	}
	login.appLog()
}

func (loger *loggerOptions) appLog() {
	if loger.logpath != nil {
		newlog := log.New(loger.logpath, " [INFO ]", 0)
		newlog.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + loger.level)
		newlog.Println(loger.caller + loger.appname + ": " + loger.msg)
	}
	switch loger.level {
	case " [WARN] ":
		fmt.Println(ui.Warn(loger.level + loger.caller + loger.appname + ": " + loger.msg))
	case " [DEBUG] ":
		fmt.Println(ui.Debug(loger.level + loger.caller + loger.appname + ": " + loger.msg))
	case " [ERROR] ":
		fmt.Println(ui.Error(loger.level + loger.caller + loger.appname + ": " + loger.msg))
	case " [INFO] ":
		fmt.Println(ui.Info(loger.level + loger.caller + loger.appname + ": " + loger.msg))
	}
}

func getStringOfMessage(g interface{}) string {
	switch g.(type) {
	case string:
		return g.(string)
	case error:
		return g.(error).Error()
	default:
		return "unknown messagetype"
	}
	return ""
}
