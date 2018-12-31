package jenkins

import (
    "fmt"
    "encoding/json"
    "io/ioutil"
)

func readCiConfig() map[string]interface{} {

    plan, _ := ioutil.ReadFile("/var/lib/neuron/neuron.json")
    var data map[string]interface{}
    err := json.Unmarshal(plan, &data)
    if err !=nil {
        fmt.Println("ERROR: Configuration file provided is not valid, I cannot proceed further")
    }
    //config_value := data["ci"].([]interface{})

    return  data
}

func readCiCred(cred_file string) map[string]interface{} {

    plan, _ := ioutil.ReadFile(cred_file)
    var data map[string]interface{}
    err := json.Unmarshal(plan, &data)
    if err !=nil {
        fmt.Println("ERROR: Configuration file provided is not valid, I cannot proceed further")
    }
    //config_value := data["ci"].([]interface{})

    return  data
}
