# Prerequisites:

 * It can be configured using configuration file. Place your "neuron.json" under /var/lib/neuron and set back neuron takes care rest of it.
 * If configuration is not placed in the default location. One can make use of cli to configure the application.
 * Using cli to configure the application is prefered. One can find the steps to do the same below.

### sample neuron.json:-

```json
   {
     "port": "80",
     "home": "/var/lib/neuron",
     "uidir": "/root/go/src/neuron/web",
     "logfile": "neuron.log",
     "ui_logfile": "neuron_ui.log",
     "logfile_location": "/var/log/neuron/",
     "database": {
                    "mongoDB": "XX.XX.XX.XX",
                    "<DB name>": "<ADDRESS>"
                 }
   }
```

 * Neuron will log its output now, for the API log refer "/var/log/neuron/neuron.log" (default file) else to a user defined file and for application/startup log refer "/var/log/neuron/neuronapp.log"
 * Now we can connect the application to Database, for now it supports MongoDb (It is in alfa). To connect database mentioned the details of it in configuration file. If not filesystem will be considerd for the actions.
 * If the UI is required it has to be enabled by mentioning in the config file, by default it will be disabled.
 * The application will work on the port mentioned in the configuration file, If not it will work on port 80 by default.
 * This supports profile based cloud interation, we can store various cloud profiles/accounts, and interact with the cloud by making use of stored profiles.
 * Neuron's API's are now capable of returning RAW outputs from cloud, this will be the unfiltered output from cloud (One has to enable GetRaw flag in the API while requesting).
 * NeuRon supports updation of resources, such as updating server/network and etc. With this one can start/stop server in future can attach additional volume, create additional subnets in the existing VM.

## Configure

 * Use cli to configure your neuron.
 * Place the package built under `/usr/local/neuron`, and you will be able to call it like any other cli.
 * Once placed, call work like this: `neuron <commands>`. To see the complete commands and subcommands,
   use `neuron -h` which list all the available commands along with its respective flags.
 * Use the command `neuron init --config 'path/to/conf.json'` and this will setup the neuron and make it usable.
 * One have to use the flag `--config` only if the configuration file is not placed at the default location `/var/lib/neuron`.
 * To make UI available place a directory contaning UI templates in a userdefined pah and mention the same in configuration file.
 * Configuring appliction by passing config values via cli is not available yet.

# Usage:

### Examples:-

Below are the examples of how the api can be used.

```bash
    curl -H "Content-Type: application/json" -X GET -d '{"SubnetId":"subnet-d81893b1","Cloud": {"Cloud":"aws","Region":"ap-south-1"}}' http://104.211.76.61/getservers

    curl -H "Content-Type: application/json" -X DELETE -d '{"InstanceIds":["i-0a8fe9854a6e939f1"],"Cloud":"aws","Region":"ap-south-1"}' http://104.211.76.61/deleteservers

    curl -H "Content-Type: application/json" -X CREATE -d '{"Name":"Test_Network","VpcCidr":"192.168.0.0/16","SubCidr":["192.168.10.0/24","192.168.20.0/24"],"Type":"public","Ports":["8080","22","80","443"],"Cloud":"aws","Region":"ap-south-1"}' http://104.211.76.61/createnetwork

    curl -H "Content-Type: application/json" -X GET -d '{"Cloud":"aws","Region":"eu-west-3"}' http://104.211.76.61/getregions

    curl -H "Content-Type: application/json" -X GET -d '{"Cloud":"aws","Region":"eu-west-3"}' http://104.211.76.61/getsubnets

    curl -H "Content-Type: application/json" -X GET -d '{"Cloud":"aws","Region":"eu-west-3"}' http://104.211.76.61/getallservers

    curl -H "Content-Type: application/json" -X DELETE -d '{"VpcId":["vpc-0b61987128b98bee0","vpc-01f61d4f8108feaa4"],"Cloud":"aws","Region":"eu-west-3"}' http://104.211.76.61/deletenetworks

```

### Examples With Profile :-

Example which let you know how to make use of profile while making call to api.

```bash
    curl -H "Content-Type: application/json" -X CREATE -d '{"InstanceName":"consul-machine","ImageId":"ami-46eea129","SubnetId":"subnet-d81893b1","KeyName":"chef-coe-ind","Flavor":"t2.micro","UserData":"echo 'nothing'","AssignPubIp":true, "Cloud":"aws","Region":"ap-south-1", "Profile":"niktest"}' http://localhost:80/createserver
```

# UI:

 * For UI of neuron, mention the path of UI folder in the configuration file.
 * This UI has very minimal feature when compared to neron application alone.
 * Once it is hosted one can access UI by hitting: `http://<your domain/ip><:port>`

# Coming Soon:

 * We are into Azure now, but it is in alfa and has very minimal feature.
 * Few more feature on updation endpoints has to be added.
 * Few more endpoints for the operations.

**For doubts and more info please write to us @:** [iaac.devops.coe]
[iaac.devops.coe]: iaac.devops.coe@gmail.com
