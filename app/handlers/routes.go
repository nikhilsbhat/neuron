package handlers

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type UiRoute struct {
	Name        string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routs []Route
type UiRouts []UiRoute

var Routes = Routs{

	//routes for various other operations
	//Route{"BuildMachine", "POST", "/startbuildmachine", createbuild},
	Route{"CreateConsul", "CREATE", "/createconsul", createconsul},
	Route{"GetRegions", "GET", "/getregions", getregions},
	Route{"GetSubnets", "GET", "/getsubnets", getsubnets},
	//Route{"StartImageServer", "POST", "/startimagemachine", startimagemachine},
	Route{"GetCount", "GET", "/getcount", getcount},
	//routes which deals with servers
	Route{"CreateServer", "CREATE", "/createserver", createserver},
	Route{"CreateServerMock", "CREATE", "/createservermock", createservermock},
	Route{"GetServers", "GET", "/getservers", getservers},
	Route{"GetAllServers", "GET", "/getallservers", getallservers},
	Route{"DeleteServers", "DELETE", "/deleteservers", deleteservers},
	Route{"UpdateServers", "UPDATE", "/updateservers", updateservers},
	//routes which deals with network
	Route{"CreateNetwork", "CREATE", "/createnetwork", createnetwork},
	Route{"DeleteNetwork", "DELETE", "/deletenetworks", deletenetworks},
	Route{"GetAllNetworks", "GET", "/getallnetworks", getallnetworks},
	Route{"GetNetworks", "GET", "/getnetworks", getnetworks},
	Route{"UpdateNetworks", "UPDATE", "/updatenetwork", updatenetwork},
	//routes which deals with loadbalancers
	Route{"CreateLoadBalancer", "CREATE", "/createloadbalancer", createloadbalancer},
	Route{"DeleteLoadBalancer", "DELETE", "/deleteloadbalancer", deleteloadbalancer},
	Route{"GetLoadBalancer", "GET", "/getloadbalancers", getloadbalancers},
	Route{"GetAllLoadBalancer", "GET", "/getallloadbalancers", getallloadbalancers},
	//routes which deals with images
	Route{"CreateImage", "CREATE", "/createimage", createimage},
	Route{"DeleteImage", "DELETE", "/deleteimage", deleteimage},
	Route{"GetImage", "GET", "/getimages", getimages},
	Route{"GetAllImage", "GET", "/getallimages", getallimages},
}

var UiRoutes = UiRouts{
	UiRoute{"NeuRon", "/", neuron},
	UiRoute{"NeuRon", "/neuron", neuron},
	UiRoute{"BuildApp", "/buildapp", buildapp},
	UiRoute{"CloudView", "/cloudview", cloudview},
	UiRoute{"CloudSetting", "/cloudsettings", cloudsetting},
	UiRoute{"CIView", "/ciview", ciview},
	UiRoute{"CISetting", "/cisettings", cisetting},
	UiRoute{"SetCi", "/setci", setci},
	UiRoute{"JenkinsView", "/jenkinsview", jenkinsview},
	UiRoute{"CloudView2", "/cloudview2", cloudview2},
}
