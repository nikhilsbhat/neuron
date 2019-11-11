package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	nwget "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/network/get"
	svget "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/server/get"
	"github.com/nikhilsbhat/neuron/ci/jenkins"
	"github.com/nikhilsbhat/neuron/database"
	"github.com/nikhilsbhat/neuron/database/common"
	/*      count "neuron/count"
	        build "neuron/neuronbuild"
	        buildim "neuron/neuronbuild/image"
	        imcreate "neuron/cloudoperations/image/create"
	        imdelete "neuron/cloudoperations/image/delete"
	        imget "neuron/cloudoperations/image/get"
	        lbcreate "neuron/cloudoperations/loadbalancer/create"
	        lbdelete "neuron/cloudoperations/loadbalancer/delete"
	        lbget "neuron/cloudoperations/loadbalancer/get"
	        nwcreate "neuron/cloudoperations/network/create"
	        nwdelete "neuron/cloudoperations/network/delete"
	        nwupdate "neuron/cloudoperations/network/update"
	        svcreate "neuron/cloudoperations/server/create"
	        svdelete "neuron/cloudoperations/server/delete"
	        svupdate "neuron/cloudoperations/server/update"*/)

// UiTemplatePath holds the path to folder consisting of UI gotemplates and other files to host UI.
var UiTemplatePath string

func neuron(w http.ResponseWriter, r *http.Request) {

	t := template.Must(template.ParseGlob(UiTemplatePath))
	err := t.ExecuteTemplate(w, "page_layout.tmpl", uiTemp{Title: "Dashboard", Cont: "dashboard"})
	if err != nil {
		log.Fatal("Cannot Get View ", err)
	}

}

func buildapp(w http.ResponseWriter, r *http.Request) {

	t := template.Must(template.ParseGlob(UiTemplatePath))
	err := t.ExecuteTemplate(w, "page_layout.tmpl", uiTemp{Title: "BuildView", Cont: "build"})
	if err != nil {
		log.Fatal("Cannot Get View ", err)
	}

}

func cloudview(w http.ResponseWriter, r *http.Request) {

	getserverdetails := svget.GetServersInput{}
	getserverdetails.Cloud.Name = "aws"
	getserverdetails.Cloud.Region = "ap-south-1"
	getserverdetails.Cloud.Profile = "niktest"
	getserverresponse, regerr := getserverdetails.GetAllServers()
	if regerr != nil {
		fmt.Fprintf(w, "%v\n", regerr)
	} else {
		jsonval, _ := json.Marshal(getserverresponse)

		value := []byte(string(jsonval))
		var data []map[string]interface{}
		err1 := json.Unmarshal(value, &data)
		if err1 != nil {
			fmt.Println(err1)
		}
		var mp4 []map[string]interface{}
		var mp []interface{}
		for _, mp1 := range data {
			for _, mp2 := range mp1 {
				mpt := &mp
				*mpt = mp2.([]interface{})
				for _, mp3 := range mp2.([]interface{}) {
					mp4 = append(mp4, mp3.(map[string]interface{}))
				}
			}
		}

		t := template.Must(template.ParseGlob(UiTemplatePath))
		err := t.ExecuteTemplate(w, "page_layout.tmpl", struct {
			Title     string
			Cont      string
			Pas       uiTemp
			AllServer []map[string]interface{}
		}{Title: "CloudView", Cont: "cloudview", Pas: uiTemp{Pass: mp}, AllServer: mp4})
		if err != nil {
			log.Fatal("Cannot Get View ", err)
		}
	}
}

func cloudsetting(w http.ResponseWriter, r *http.Request) {

	t := template.Must(template.ParseGlob(UiTemplatePath))
	err := t.ExecuteTemplate(w, "page_layout.tmpl", uiTemp{Title: "CloudSettings", Cont: "cloudsettings"})
	if err != nil {
		log.Fatal("Cannot Get View ", err)
	}

}

func ciview(w http.ResponseWriter, r *http.Request) {

	t := template.Must(template.ParseGlob(UiTemplatePath))
	err := t.ExecuteTemplate(w, "page_layout.tmpl", uiTemp{Title: "CIView", Cont: "ciview"})
	if err != nil {
		log.Fatal("Cannot Get View ", err)
	}

}

func cisetting(w http.ResponseWriter, r *http.Request) {

	/*      if database.Db == nil {
	                if _, dir_err := os.Stat("/var/lib/neuron/neuron.json"); os.IsNotExist(dir_err) {
	                        fmt.Println("I did not find any configuration file to read CI data")
	                        config_byt := []byte(`{"ci": [{"name": "Not Connected","domain": "Not Connected" }]}`)
	                        var dumy_config map[string]interface{}
	                        if err := json.Unmarshal(config_byt, &dumy_config); err != nil {
	                                panic(err)
	                        }

	                        t := template.Must(template.ParseGlob(UiTemplatePath))
	                        err := t.ExecuteTemplate(w, "page_layout.tmpl", struct {
	                                Title  string
	                                Cont   string
	                                Cred   string
	                                Pas    uiTemp
	                                CiCred interface{}
	                        }{Title: "CISettings", Cont: "cisettings", Cred: "no", Pas: uiTemp{Pass: dumy_config["ci"].([]interface{})}, CiCred: "dummy"})
	                        if err != nil {
	                                log.Fatal("Cannot Get View ", err)
	                        }

	                } else {
	                        fmt.Println("Found configuration file and reading CI data from there")
	                        if config == nil {
	                                fmt.Fprintf(w, "Encountered error while reading config file")
	                        } else {
	                                if (database.Db) != nil {
	                                        fmt.Println("Found Database connecting to fetch further data")
	                                        t := template.Must(template.ParseGlob(UiTemplatePath))
	                                        err := t.ExecuteTemplate(w, "page_layout.tmpl", struct {
	                                                Title  string
	                                                Cont   string
	                                                Cred   string
	                                                CiCred interface{}
	                                        }{Title: "CISettings", Cont: "cisettings", Cred: "no", CiCred: "dummy"})
	                                        if err != nil {
	                                                log.Fatal("Cannot Get View ", err)
	                                        }
	                                } else {
	                                        ci_dat_file := fmt.Sprintf("%s/data/ci_cred.json", config["home"])
	                                        if _, dir_err := os.Stat(ci_dat_file); os.IsNotExist(dir_err) {
	                                                fmt.Println("couldn't find credentials of CI, guess you've not set that")

	                                                // redering template with no CI credentials
	                                                t := template.Must(template.ParseGlob(UiTemplatePath))
	                                                err := t.ExecuteTemplate(w, "page_layout.tmpl", struct {
	                                                        Title  string
	                                                        Cont   string
	                                                        Cred   string
	                                                        Pas    uiTemp
	                                                        CiCred interface{}
	                                                }{Title: "CISettings", Cont: "cisettings", Cred: "no", CiCred: "dummy"})
	                                                if err != nil {
	                                                        log.Fatal("Cannot Get View ", err)
	                                                }

	                                        } else {
	                                                fmt.Println(ci_dat_file)
	                                                fetch_data := readCiCred(ci_dat_file)
	                                                fmt.Println(fetch_data)
	                                                fmt.Println("Fetching credentials of the CI you've set")
	                                                // redering template with CI credentials
	                                                t := template.Must(template.ParseGlob(UiTemplatePath))
	                                                err := t.ExecuteTemplate(w, "page_layout.tmpl", struct {
	                                                        Title  string
	                                                        Cont   string
	                                                        Cred   string
	                                                        Pas    uiTemp
	                                                        CiCred interface{}
	                                                }{Title: "CISettings", Cont: "cisettings", Cred: "yes", CiCred: fetch_data["ci"].([]interface{})})
	                                                if err != nil {
	                                                        log.Fatal("Cannot Get View ", err)
	                                                }
	                                        }
	                                }
	                        }
	                }
	        } else {

	                records, err := dbcommon.FetchCiData(database.DataDetail{"nikhil", "ci"})
	                if err != nil {
	                        fmt.Fprintf(w, fmt.Sprintf("%s", err))
	                } else {
	                        t := template.Must(template.ParseGlob(UiTemplatePath))
	                        err := t.ExecuteTemplate(w, "page_layout.tmpl", struct {
	                                Title  string
	                                Cont   string
	                                Cred   string
	                                Pas    uiTemp
	                                CiCred interface{}
	                        }{Title: "CISettings", Cont: "cisettings", Cred: "yes", CiCred: records})
	                        if err != nil {
	                                log.Fatal("Cannot Get View ", err)
	                        }
	                }
	        } */
}

func setci(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	datavalue := database.CiData{CiName: r.FormValue("ciname"), CiURL: r.FormValue("ciurl"), CiUsername: r.FormValue("ciusername"), CiPassword: r.FormValue("cipassword"), Timestamp: time.Now()}

	fmt.Fprintf(w, "%v\n", datavalue)
	records, err := dbcommon.StoreCIdata(database.DataDetail{"nikhil", "ci"}, datavalue)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("%s", err))
	} else {
		fmt.Fprintf(w, fmt.Sprintf("%s", records))
	}
}

// Nouifound will display appropriate page if UI is not configured.
func Nouifound(w http.ResponseWriter, r *http.Request) {

	t := template.Must(template.New("").Parse(noUIHTML))
	err := t.Execute(w, nil)
	if err != nil {
		log.Fatal("Cannot Get View ", err)
	}

}

func jenkinsview(w http.ResponseWriter, r *http.Request) {

	build := jenkins.GetJobs()
	for _, job := range build {
		fmt.Fprintf(w, "Job Separater")
		jsonval, _ := json.MarshalIndent(job.Raw, "", " ")
		fmt.Fprintf(w, "%v\n", string(jsonval))
	}

}

func cloudview2(w http.ResponseWriter, r *http.Request) {

	// Fetching the details of all servers present in the CLOUD AWS
	getserverdetails := svget.GetServersInput{}
	getserverdetails.Cloud.Name = "aws"
	getserverdetails.Cloud.Region = "ap-south-1"
	getserverdetails.Cloud.Profile = "niktest"
	getserverresponse, serverr := getserverdetails.GetAllServers()
	var allserver []map[string]interface{}
	if serverr != nil {
		fmt.Fprintf(w, "%v\n", serverr)
	} else {
		serverjson, _ := json.Marshal(getserverresponse)

		servervalue := []byte(string(serverjson))
		var allserverdata []map[string]interface{}
		getserverr := json.Unmarshal(servervalue, &allserverdata)
		if getserverr != nil {
			fmt.Println(getserverr)
		}

		for _, mp1 := range allserverdata {
			for _, mp2 := range mp1 {
				for _, mp3 := range mp2.([]interface{}) {
					allserver = append(allserver, mp3.(map[string]interface{}))
				}
			}
		}
	}

	// Fetching the details of all networks present in CLOUD AWS
	getnetworksdetails := nwget.GetNetworksInput{}
	getnetworksdetails.Cloud.Name = "aws"
	getnetworksdetails.Cloud.Region = "ap-south-1"
	getnetworksdetails.Cloud.Profile = "niktest"
	getnetworksresponse, neterr := getnetworksdetails.GetAllNetworks()

	type Allnetworks struct {
		Region string
		AllNet []map[string]interface{}
		Subets []map[string]interface{}
	}
	var allnetwork []Allnetworks

	if neterr != nil {
		fmt.Fprintf(w, "%v\n", neterr)
	} else {
		networkjson, _ := json.Marshal(getnetworksresponse)

		networksvalue := []byte(string(networkjson))
		var allnetworksdata []map[string]interface{}
		getnetwerr := json.Unmarshal(networksvalue, &allnetworksdata)
		if getnetwerr != nil {
			fmt.Println(getnetwerr)
		}

		for _, first := range allnetworksdata {
			for _, second := range first {
				var val []map[string]interface{}
				var sub []map[string]interface{}
				for _, third := range second.([]interface{}) {
					val = append(val, third.(map[string]interface{}))
					//sub_val := third.(map[string]interface{})["Subnets"]
					for _, subr := range third.(map[string]interface{})["Subnets"].([]interface{}) {
						sub = append(sub, subr.(map[string]interface{}))
					}
				}
				allnetwork = append(allnetwork, Allnetworks{(first["AwsResponse"].([]interface{})[0]).(map[string]interface{})["Region"].(string), val, sub})
			}
		}
	}

	// Rendering template with data
	t := template.Must(template.ParseGlob(UiTemplatePath))
	err := t.ExecuteTemplate(w, "page_layout.tmpl", struct {
		Title      string
		Cont       string
		AllServer  []map[string]interface{}
		AllNetwork []Allnetworks
	}{Title: "CloudView", Cont: "cloudview2", AllServer: allserver, AllNetwork: allnetwork})
	if err != nil {
		log.Fatal("Cannot Get View ", err)
	}
}
