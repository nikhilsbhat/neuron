package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	count "github.com/nikhilsbhat/neuron/count"
	/*build "neuron/neuronbuild"
	  buildim "neuron/neuronbuild/image"*/
	imcreate "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/image/create"
	imdelete "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/image/delete"
	imget "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/image/get"
	lbcreate "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/loadbalancer/create"
	lbdelete "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/loadbalancer/delete"
	lbget "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/loadbalancer/get"
	misc "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/miscellaneous"
	nwcreate "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/network/create"
	nwdelete "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/network/delete"
	nwget "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/network/get"
	nwupdate "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/network/update"
	svcreate "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/server/create"
	svdelete "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/server/delete"
	svget "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/server/get"
	svupdate "github.com/nikhilsbhat/neuron-cloudy/cloudoperations/server/update"
)

func createconsul(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(body))
	var t servercreateinput

	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	//    consul := DengineConsul.ConsulConfig.
}

func createserver(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData svcreate.ServerCreateInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			createserverresponse, serverr := myData.CreateServer()
			if serverr != nil {
				fmt.Fprintf(rw, "%v\n", serverr)
			} else {
				jsonval, _ := json.MarshalIndent(createserverresponse, "", " ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

//The below two function is suspended temporarily, and this will be back soon
/*func createbuild(rw http.ResponseWriter, req *http.Request) {
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
                jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
                fmt.Fprintf(rw, "%v\n", string(jsonval))
        } else {
                var t createbuildmachine

                err = json.Unmarshal(body, &t)
                if err != nil {
                        jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
                        fmt.Fprintf(rw, "%v\n", string(jsonval))
                } else {
                        server_create_input := build.BuildServerCreateInput{t.AppVersion, t.UniqueId, "build-machine", "subnet-d81893b1", "chef-coe-ind", "t2.micro", "aws", "ap-south-1", true}
                        serverresponse := server_create_input.BuildServerCreate()
                        jsonval, _ := json.MarshalIndent(serverresponse, "", " ")
                        fmt.Fprintf(rw, "%v\n", string(jsonval))
                }
        }
}

func startimagemachine(rw http.ResponseWriter, req *http.Request) {
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
                jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
                fmt.Fprintf(rw, "%v\n", string(jsonval))
        } else {
                var t startimageinput

                err = json.Unmarshal(body, &t)
                if err != nil {
                        jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
                        fmt.Fprintf(rw, "%v\n", string(jsonval))
                } else {
                        server_create_input := buildim.PackerServerCreateInput{t.AppVersion, t.AppName, t.RepoEmail, t.RepoUsername, t.RepoPasswd, t.ArtDomain, t.ArtUsername, t.ArtPasswd, t.InstanceName, t.SubnetId, t.KeyName, t.Flavor, t.Cloud.Cloud, t.Cloud.Region, t.AssignPubIp}
                        serverresponse, serverr := server_create_input.BuildServerCreate()
                        if serverr != nil {
                                fmt.Fprintf(rw, "%v\n", serverr)
                        } else {
                                jsonval, _ := json.MarshalIndent(serverresponse, "", " ")
                                fmt.Fprintf(rw, "%v\n", string(jsonval))
                        }
                }
        }
}*/

func deleteservers(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData svdelete.DeleteServersInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			deleteserverresponse, delerr := myData.DeleteServer()
			if delerr != nil {
				fmt.Fprintf(rw, "%v\n", delerr)
			} else {
				jsonval, _ := json.MarshalIndent(deleteserverresponse, "", " ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func createnetwork(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData nwcreate.NetworkCreateInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			serverresponse, servresperr := myData.CreateNetwork()
			if servresperr != nil {
				fmt.Fprintf(rw, "%v\n", servresperr)
			} else {
				jsonval, _ := json.MarshalIndent(serverresponse, "", " ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func getregions(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData misc.GetRegionInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			getregionresponse, regerr := myData.GetRegions()
			if regerr != nil {
				fmt.Fprintf(rw, "%v\n", regerr)
			} else {
				jsonval, _ := json.Marshal(getregionresponse)
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func getsubnets(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData nwget.GetNetworksInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			getsubnetsresponse, suberr := myData.GetSubnets()
			if suberr != nil {
				fmt.Fprintf(rw, "%v\n", suberr)
			} else {
				jsonval, _ := json.MarshalIndent(getsubnetsresponse, "", " ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func getservers(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData svget.GetServersInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			getservesponse, geterr := myData.GetServersDetails()
			if geterr != nil {
				fmt.Fprintf(rw, "%v\n", geterr)
			} else {
				jsonval, _ := json.MarshalIndent(getservesponse, "", " ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func getallservers(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData svget.GetServersInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			getserverresponse, serverr := myData.GetAllServers()
			if serverr != nil {
				fmt.Fprintf(rw, "%v\n", serverr)
			} else {
				jsonval, _ := json.MarshalIndent(getserverresponse, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}

}

func updateservers(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData svupdate.UpdateServersInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			updateservresponse, updaterr := myData.UpdateServers()
			if updaterr != nil {
				fmt.Fprintf(rw, "%v\n", updaterr)
			} else {
				jsonval, _ := json.MarshalIndent(updateservresponse, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}

}

func deletenetworks(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData nwdelete.DeleteNetworkInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			deletenetworkresponse, neterr := myData.DeleteNetwork()
			if neterr != nil {
				fmt.Fprintf(rw, "%v\n", neterr)
			} else {
				jsonval, _ := json.MarshalIndent(deletenetworkresponse, "", " ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func createloadbalancer(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData lbcreate.LbCreateInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			createlbresponse, createlberr := myData.CreateLoadBalancer()
			if createlberr != nil {
				fmt.Fprintf(rw, "%v\n", createlberr)
			} else {
				jsonval, _ := json.MarshalIndent(createlbresponse, "", " ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func deleteloadbalancer(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData lbdelete.LbDeleteInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			deletelbresponse, delbresponse := myData.DeleteLoadBalancer()
			if delbresponse != nil {
				fmt.Fprintf(rw, "%v\n", delbresponse)
			} else {
				jsonval, _ := json.MarshalIndent(deletelbresponse, "", " ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func getloadbalancers(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData lbget.GetLoadbalancerInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			loadbalancersgetresponse, getlbresponse := myData.GetLoadbalancers()
			if getlbresponse != nil {
				fmt.Fprintf(rw, "%v\n", getlbresponse)
			} else {
				jsonval, _ := json.MarshalIndent(loadbalancersgetresponse, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func getallloadbalancers(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData lbget.GetLoadbalancerInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			loadbalancersgetresponse, getlbresponse := myData.GetAllLoadbalancer()
			if getlbresponse != nil {
				fmt.Fprintf(rw, "%v\n", getlbresponse)
			} else {
				jsonval, _ := json.MarshalIndent(loadbalancersgetresponse, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func createimage(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData imcreate.CreateImageInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			imagecreateresponse, imgerr := myData.CreateImage()
			if imgerr != nil {
				fmt.Fprintf(rw, "%v\n", imgerr)
			} else {
				jsonval, _ := json.MarshalIndent(imagecreateresponse, "", " ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func deleteimage(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData imdelete.DeleteImageInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			imagedeleteresponse, delimgerr := myData.DeleteImage()
			if delimgerr != nil {
				fmt.Fprintf(rw, "%v\n", delimgerr)
			} else {
				jsonval, _ := json.MarshalIndent(imagedeleteresponse, "", " ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func getimages(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData imget.GetImagesInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			imagegetresponse, getimgerr := myData.GetImage()
			if getimgerr != nil {
				fmt.Fprintf(rw, "%v\n", getimgerr)
			} else {
				jsonval, _ := json.MarshalIndent(imagegetresponse, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func getallimages(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData imget.GetImagesInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			imagegetresponse, getallimgerr := myData.GetAllImage()
			if getallimgerr != nil {
				fmt.Fprintf(rw, "%v\n", getallimgerr)
			} else {
				jsonval, _ := json.MarshalIndent(imagegetresponse, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func getallnetworks(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData nwget.GetNetworksInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			getallnetworksresponse, neterr := myData.GetAllNetworks()
			if neterr != nil {
				fmt.Fprintf(rw, "%v\n", neterr)
			} else {
				jsonval, _ := json.MarshalIndent(getallnetworksresponse, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func getnetworks(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData nwget.GetNetworksInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			getnetworksresponse, neterr := myData.GetNetworks()
			if neterr != nil {
				fmt.Fprintf(rw, "%v\n", neterr)
			} else {
				jsonval, _ := json.MarshalIndent(getnetworksresponse, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func getcount(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData count.GetCountInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			getcountresponse := myData.GetCount()
			jsonval, _ := json.MarshalIndent(getcountresponse, "", "  ")
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		}
	}
}

func updatenetwork(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.Marshal(Error{"Not received input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var myData nwupdate.NetworkUpdateInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			jsonval, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			netupdateresponse, netuperr := myData.UpdateNetwork()
			if netuperr != nil {
				fmt.Fprintf(rw, "%v\n", netuperr)
			} else {
				jsonval, _ := json.MarshalIndent(netupdateresponse, "", " ")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func createservermock(rw http.ResponseWriter, req *http.Request) {
	createserverresponse, serverr := svcreate.CreateServerMock()
	if serverr != nil {
		fmt.Fprintf(rw, "%v\n", serverr)
	} else {
		fmt.Fprintf(rw, "%v\n", createserverresponse.DefaultResponse)
	}
}
