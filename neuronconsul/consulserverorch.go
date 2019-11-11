package DengineConsul

import (
	"fmt"
	"time"
)

type ConsulConfig struct {
	Datacenter              string
	Data_dir                string
	Log_level               string
	Enable_syslog           bool
	Enable_debug            bool
	Enable_script_checks    bool
	Leave_on_terminate      bool
	Skip_leave_on_interrupt bool
	Rejoin_after_leave      bool
	Bootstrap_expect        int64
	Check                   Check
}

type consulresponce struct {
	Bind_addr      string
	Ui             bool
	Client         string
	Advertise_addr string
	Retry_join     []string
	Server         bool
}

func (concreate *ConsulCreateInput) ConsulOrch(conconf ConsulConfig) ServerCreate.ServerCreateResponse {

	// I will establish session so that we can carry out the process in cloud
	session_input := DengineAwsInterface.EstablishConnectionInput{concreate.Region, "ec2"}
	session_input.EstablishConnection()

	// Yo I will be setting a unique ID for each server that I will be creating
	now_time := time.Now().Local().Format("20060102150405")
	unique_name := fmt.Sprintf(concreate.InstanceName + "-" + now_time)

	// making decision for consul (deciding whether the instances should be consul server or agent)
	consul_responce := DengineAwsInterface.DecideConsulType(concreate.SubnetId, concreate.Region)
	consul_config := ConsulServerConfig{
		consul_responce.Bind_addr,
		conconf.Datacenter,
		consul_responce.Ui,
		consul_responce.Client,
		consul_responce.Advertise_addr,
		conconf.Data_dir,
		conconf.Log_level,
		conconf.Enable_syslog,
		conconf.Enable_debug,
		conconf.Enable_script_checks,
		unique_name,
		consul_responce.Server,
		conconf.Leave_on_terminate,
		conconf.Skip_leave_on_interrupt,
		conconf.Rejoin_after_leave,
		conconf.Bootstrap_expect,
		consul_responce.Retry_join,
		conconf.Check,
	}

	server_create_input := ConsulCreateInput{
		unique_name,
		concreate.ImageId,
		concreate.SubnetId,
		consul_config,
		concreate.KeyName,
		concreate.Flavor,
		concreate.Cloud,
		concreate.Region,
		concreate.AssignPubIp,
	}

	response := server_create_input.ConsulServerCreate()
	return response
}
