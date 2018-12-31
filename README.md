# Prerequisites:

  * Place your AWS credentails under .aws/credentails
  * Place your AZURE credentails under .azure/credentails
  * Now you can configure neuron to an extent. Place your "neuron.json" under /var/lib/neuron and set back neuron takes care rest of it.
  # sample neuron.json:-

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

# Steps:

  * Give the package execution permission.
  * Execute the package using command "./neuron".
  * To make UI available place a directory contaning UI templates in a userdefined pah and mention the same in configuration file.

# Usage:

  # Examples:-

	curl -H "Content-Type: application/json" -X CREATE -d '{"InstanceName":"shashank-machine","ImageId":"ami-46eea129","SubnetId":"subnet-d81893b1","KeyName":"chef-coe-ind","Flavor":"t2.micro","UserData":"echo Iam nikhil","Cloud":"aws","Region":"ap-south-1","Count":1,"AssignPubIp":true}' http://104.211.76.61/createserver

	curl -H "Content-Type: application/json" -X UPDATE -d '{"AppVersion":"2.8","UniqueId":"2.8"}' http://35.224.18.58:8080/startbuildmachine

	curl -H "Content-Type: application/json" -X GET -d '{"SubnetId":"subnet-d81893b1","Cloud": {"Cloud":"aws","Region":"ap-south-1"}}' http://104.211.76.61/getservers

	curl -H "Content-Type: application/json" -X DELETE -d '{"InstanceIds":["i-0a8fe9854a6e939f1"],"Cloud":"aws","Region":"ap-south-1"}' http://104.211.76.61/deleteservers

	curl -H "Content-Type: application/json" -X CREATE -d '{"Name":"Test_Network","VpcCidr":"192.168.0.0/16","SubCidr":["192.168.10.0/24","192.168.20.0/24"],"Type":"public","Ports":["8080","22","80","443"],"Cloud":"aws","Region":"ap-south-1"}' http://104.211.76.61/createnetwork

	curl -H "Content-Type: application/json" -X GET -d '{"Cloud":"aws","Region":"eu-west-3"}' http://104.211.76.61/getregions

	curl -H "Content-Type: application/json" -X GET -d '{"Cloud":"aws","Region":"eu-west-3"}' http://104.211.76.61/getsubnets

	curl -H "Content-Type: application/json" -X GET -d '{"Cloud":"aws","Region":"eu-west-3"}' http://104.211.76.61/getallservers

	curl -H "Content-Type: application/json" -X DELETE -d '{"VpcId":["vpc-0b61987128b98bee0","vpc-01f61d4f8108feaa4"],"Cloud":"aws","Region":"eu-west-3"}' http://104.211.76.61/deletenetworks

	curl -H "Content-Type: application/json" -X GET -d '{"Cloud":"aws","Region":"ap-south-1"}' http://104.211.76.61/getallnetworks

	curl -H "Content-Type: application/json" -X GET -d '{"VpcIds":["vpc-0b61987128b98bee0","vpc-01f61d4f8108feaa4"],"Cloud":"aws","Region":"ap-south-1"}' http://104.211.76.61/getnetworks

	curl -H "Content-Type: application/json" -X DELETE -d '{"VpcId":["vpc-fd8e1394"],"Cloud":"aws","Region":"eu-west-3"}' http://104.211.76.61/deletevpcservers

	curl -H "Content-Type: application/json" -X UPDATE -d '{"AppVersion":"1.0","AppName":"gameoflife-web","RepoEmail":"iaac.devops.coe@gmail.com","RepoUsername":"devopsiac","RepoPasswd":"devops123","ArtDomain":"13.126.216.231","ArtUsername":"admin","ArtPasswd":"password","InstanceName":"app-server","SubnetId":"subnet-d81893b1","KeyName":"chef-coe-ind","Flavor":"t2.micro","AssignPubIp":"true","Cloud":{"Cloud":"aws","Region":"ap-south-1"}}' http://35.224.18.58:8080/startimagemachine

	curl -H "Content-Type: application/json" -X CREATE -d '{"InstanceName":"consul-machine","ImageId":"ami-46eea129","SubnetId":"subnet-d81893b1","KeyName":"chef-coe-ind","Flavor":"t2.micro","UserData":"echo 'nothing'","AssignPubIp":true,"Cloud": {"Cloud":"aws","Region":"ap-south-1"}}' http://localhost:80/createserver

	curl -H "Content-Type: application/json" -X CREATE -d '{"Name":"testing","VpcId":"vpc-5675fc3f","Scheme":"external","Type":"classic","LbPort":80,"InstPort":80,"Lbproto":"HTTP","Instproto":"HTTP","Cloud":"azure","Region":"ap-south-1"}' http://104.211.76.61/createloadbalancer

	curl -H "Content-Type: application/json" -X CREATE -d '{"Name":"test123","VpcId":"vpc-5675fc3f","Scheme":"external","Type":"application","LbPort":80,"InstPort":80,"Lbproto":"HTTP","Instproto":"HTTP","HttpCode":"201","HealthPath": "/index.html", "Cloud":"azure","Region":"ap-south-1"}' http://104.211.76.61/createloadbalancer

	curl -H "Content-Type: application/json" -X DELETE -d '{"LbArn":["arn:aws:elasticloadbalancing:ap-south-1:983899670911:loadbalancer/app/test123/82dab8f14a8707f0"],"Cloud":"azure","Region":"ap-south-1"}' http://104.211.76.61/deleteloadbalancer

	curl -H "Content-Type: application/json" -X DELETE -d '{"LbName":["testing123","testing"],"Cloud":"aws","Region":"ap-south-1"}' http://104.211.76.61/deleteloadbalancer

	curl -H "Content-Type: application/json" -X GET -d '{"Cloud":"azure","Region":"ap-south-1"}' http://104.211.76.61/getallloadbalancer

	curl -H "Content-Type: application/json" -X CREATE -d '{"Cloud":"azure","Region":"ap-south-1","InstanceId":["i-0fc1a42f8f0a29e0f"]}' http://104.211.76.61/createimage

	curl -H "Content-Type: application/json" -X DELETE -d '{"Cloud":"azure","Region":"ap-south-1","ImageId":["ami-0e8249c2a42f3ed89"]}' http://104.211.76.61/deleteimage

	curl -H "Content-Type: application/json" -X GET -d '{"Cloud":"azure","Region":"ap-south-1","ImageId":["ami-0e8249c2a42f3ed89"]}' http://104.211.76.61/getimages

	curl -H "Content-Type: application/json" -X GET -d '{"Cloud":"aws","Region":"ap-south-1"}' http://104.211.76.61/getallimages

	curl -H "Content-Type: application/json" -X GET -d '{"Cloud":"azure","Region":"ap-south-1"}' http://104.211.76.61/getcount

	curl -H "Content-Type: application/json" -X UPDATE -d '{"InstanceIds":["i-038c90521dac45a67","i-0da0046be7c988931"],"StateAction":"stop","Cloud":"aws","Region":"ap-south-1"}' http://104.211.76.61/updateservers

	curl -H "Content-Type: application/json" -X UPDATE -d '{"Cloud":"aws","Region":"ap-south-1","Catageory": {"Resource":"subnets","Action":"create","Subnets":{"Name":"updatesubnet","Cidr":["192.168.30.0/24"],"Type":"private","VpcId":"vpc-0f3f1ddac4f0fb680","Zones":["ap-south-1b"]}}}' http://104.211.76.61/updatenetwork

 # Examples With Profile :-

	curl -H "Content-Type: application/json" -X CREATE -d '{"InstanceName":"consul-machine","ImageId":"ami-46eea129","SubnetId":"subnet-d81893b1","KeyName":"chef-coe-ind","Flavor":"t2.micro","UserData":"echo 'nothing'","AssignPubIp":true, "Cloud":"aws","Region":"ap-south-1", "Profile":"niktest"}' http://localhost:80/createserver
  
# For UI:

  * For UI of neuron, mention the path of UI folder in the configuration file.
  * This UI has very minimal feature when compared to neron application alone.
  * Once it is hosted one can access UI by hitting: http://<your domain/ip><:port>

# Coming Soon:

  * We are into Azure now, but it is in alfa and has very minimal feature.
  * Few more feature on updation endpoints has to be added.
  * Few more endpoints for the operations.
  
For doubts and more info please write to us @: [iaac.devops.coe]
[iaac.devops.coe]: iaac.devops.coe@gmail.com
