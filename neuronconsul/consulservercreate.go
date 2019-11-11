package DengineConsul

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"github.com/nikhilsbhat/neuron-cloudy/cloudoperations/server/create"
)

var ()

type ConsulCreateInput struct {
	InstanceName string
	ImageId      string
	SubnetId     string
	ConsulConfig interface{}
	KeyName      string
	Flavor       string
	Cloud        string
	Region       string
	AssignPubIp  bool
}

type ConsulServerConfig struct {
	Bind_addr               string
	Datacenter              string
	Ui                      bool
	Client_Addr             string
	Advertise_addr          string
	Data_dir                string
	Log_level               string
	Enable_syslog           bool
	Enable_debug            bool
	Enable_script_checks    bool
	Node_name               string
	Server                  bool
	Leave_on_terminate      bool
	Skip_leave_on_interrupt bool
	Rejoin_after_leave      bool
	Bootstrap_expect        int64
	Retry_join              []string
	Check                   Check
}

type Check struct {
	Id       string
	Name     string
	Tcp      string
	Interval string
	Timeout  string
}

type ConsulCreateResponse struct {
	ConsulResponse ServerCreate.ServerCreateResponse
}

func convdatatojson(data interface{}) string {

	json_val, _ := json.MarshalIndent(data, "", " ")
	return strings.ToLower(string(json_val))

}

func encodeusrdata(consul_data string) string {

	userdata := fmt.Sprintf("#!/bin/sh \n mkdir -p /etc/consul.d/agent \n echo '%s' > /etc/consul.d/agent/agent.json \n ip=$(/sbin/ip -o -4 addr list eth0 | awk '{print $4}' | cut -d/ -f1) \n sed -i s/127.0.0.1/$ip/g /etc/consul.d/agent/agent.json \n sv start consul", consul_data)
	enc_data := b64.StdEncoding.EncodeToString([]byte(userdata))
	return enc_data

}

func (consul ConsulCreateInput) ConsulServerCreate() ServerCreate.ServerCreateResponse {

	// I will be preparing the userdata in the required format to pass
	json_userdata := convdatatojson(consul.ConsulConfig)
	encoded_userdata := encodeusrdata(json_userdata)

	// I will be calling ServerCreate function of Dengine by orchestrating things which It reqires
	server_create_input := ServerCreate.ServerCreateInput{consul.InstanceName, consul.ImageId, consul.SubnetId, consul.KeyName, consul.Flavor, encoded_userdata, consul.Cloud, consul.Region, consul.AssignPubIp}
	server_response := server_create_input.CreateServer()
	return server_response
}
