package jenkins

import (
	"fmt"
	"github.com/bndr/gojenkins"
	log "github.com/nikhilsbhat/neuron/logger"
	"os"
)

type jenkinsCred struct {
	JenkinsDomain   string `json:"jenkinsdomain,omitempty"`
	JenkinsUsername string `json:"jenkinsusername,omitempty"`
	JenkinsPassword string `json:"jenkinspassword,omitempty"`
}

var (
	jenkinsCredential jenkinsCred
)

func Init() {

	log.Info("Action of Configuring Jenkins Credentials inititated")
	log.Info("++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	log.Info("")
	var config map[string]interface{}
	if _, dir_err := os.Stat("/var/lib/neuron/neuron.json"); os.IsNotExist(dir_err) {

		log.Warn("I did not find any configuration file to read CI data")

		jenkinsCredentialptr := &jenkinsCredential
		*jenkinsCredentialptr = jenkinsCred{JenkinsDomain: "notfound"}

	} else {
		var cred_config map[string]interface{}
		log.Info("Found configuration file and reading CI data from there")
		configptr := &config
		*configptr = readCiConfig()

		var jenkinsdomain string

		if config["ci"] != nil {

			for _, domain_find := range config["ci"].([]interface{}) {

				jenkin := domain_find.(map[string]interface{})
				if jenkin["name"] == "jenkins" {

					jenkinsdomainptr := &jenkinsdomain
					*jenkinsdomainptr = jenkin["url"].(string)

					ci_dat := fmt.Sprintf("%s/ci_cred.json", config["config"].(map[string]interface{})["home"])
					if _, dir_err := os.Stat(ci_dat); os.IsNotExist(dir_err) {

						log.Info("couldn't find credentials of CI, guess you've not set that")
						jenkinsCredentialptr := &jenkinsCredential
						*jenkinsCredentialptr = jenkinsCred{JenkinsDomain: jenkinsdomain}

					} else {

						if cred_config["jenkins"] != nil {
							log.Info("Fetching credentials of the CI you've set")
							configptr := &cred_config
							*configptr = readCiCred(ci_dat)

							jenkinsCredentialptr := &jenkinsCredential
							*jenkinsCredentialptr = jenkinsCred{JenkinsDomain: jenkinsdomain, JenkinsUsername: cred_config["jenkins"].(map[string]interface{})["username"].(string), JenkinsPassword: cred_config["jenkins"].(map[string]interface{})["password"].(string)}

						} else {

							jenkinsCredentialptr := &jenkinsCredential
							*jenkinsCredentialptr = jenkinsCred{JenkinsDomain: jenkinsdomain, JenkinsUsername: "notfound", JenkinsPassword: "notfound"}

						}
					}

				} else {

					log.Warn("I did not find any configurations (URL) regarding jenkins in configuration file")

				}
			}

		} else {

			log.Info("I did not find any configurations for CI in configuration file")

		}
	}

	log.Info("Action of Configuring Jenkins Credentials Completed")
	log.Info("")
	log.Info("++++++++++++++++++++++++++++++++++++++++++++++++++++++")
}

func GetJobs() []*gojenkins.Job {

	var details []*gojenkins.Job
	if (jenkinsCredential.JenkinsDomain != "notfound") || (jenkinsCredential.JenkinsUsername != "notfound") || (jenkinsCredential.JenkinsPassword != "notfound") {

		/*jenkins := gojenkins.CreateJenkins(nil,jenkinsCredential.JenkinsDomain)
		err := jenkins.Init()

		if err != nil {
			panic("Something Went Wrong")
		}

		build, err := jenkins.GetAllJobs("Test")
		if err != nil {
			panic("Job Does Not Exist")
		}

		detailsptr := &details
		*detailsptr = build*/

	}
	return details
}
