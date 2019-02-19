package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
)

const noUIHTML = `<html>
<head>
<title>NeuRon</title>

<style>
p.b {
    font-family: "Lucida Console", Monaco, monospace;
    font-weight: bold;
}
body {
  background-color: #C0CFD1;
}
</style>
</head>

<body>
<h1 class="b" align="center"><font color="#23545B">Hey buddy you have not enabled UI for NeuRon</h1>
<p  class="b" align="center"><font color="#64717B">Please mention the UI directory path in the neuron.json (refer the README.md for more info)</p>
<p  class="b" align="center"><font color="#64717B">Even after setting the UI path you are seeing this page then have a look at the app's logfile for more info &#9786;</p>
</body>
</html>`

// FillMyStruct helps converting map[string]interface to strut type you need.
func FillMyStruct(s FillStructs) error {
	val1 := s.Data
	mystructjson, _ := json.Marshal(s.Data)
	json.Unmarshal([]byte(string(mystructjson)), &val1)
	setfield := FillStructs{Type: s.Type, Data: s.Data}
	for k, v := range val1.(map[string]interface{}) {
		err := setfield.SetField(s.Type, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetField is a method where actual converstion of map[string]interface to strut happens.
func (s *FillStructs) SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(s.Type).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		switch newval := value.(type) {
		case []string:
			structFieldValue.Set(reflect.ValueOf(newval))
		case []interface{}:
			s := make([]string, len(newval))
			for i, v := range newval {
				s[i] = fmt.Sprint(v)
			}
			structFieldValue.Set(reflect.ValueOf(s))
		case float64:
			structFieldValue.Set(reflect.ValueOf(int(newval)))
		default:
			invalidTypeError := errors.New("Provided value type didn't match obj field type")
			return invalidTypeError
		}
	} else if structFieldType == val.Type() {
		structFieldValue.Set(val)
	} else {
		invalidTypeError := errors.New("Provided value type didn't match obj field type")
		return invalidTypeError
	}
	return nil
}

func readCiCred(credfile string) map[string]interface{} {

	plan, _ := ioutil.ReadFile(credfile)
	var data map[string]interface{}
	err := json.Unmarshal(plan, &data)
	if err != nil {
		fmt.Println("ERROR: Configuration file provided is not valid, I cannot proceed further")
	}
	//config_value := data["ci"].([]interface{})

	return data
}
