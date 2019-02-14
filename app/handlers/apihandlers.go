package handlers

import (
	"encoding/json"
	"fmt"
	count "github.com/nikhilsbhat/neuron/count"
	"io/ioutil"
	"log"
	"net/http"
	/*build "neuron/neuronbuild"
	  buildim "neuron/neuronbuild/image"*/
	imcreate "github.com/nikhilsbhat/neuron/cloudoperations/image/create"
	imdelete "github.com/nikhilsbhat/neuron/cloudoperations/image/delete"
	imget "github.com/nikhilsbhat/neuron/cloudoperations/image/get"
	lbcreate "github.com/nikhilsbhat/neuron/cloudoperations/loadbalancer/create"
	lbdelete "github.com/nikhilsbhat/neuron/cloudoperations/loadbalancer/delete"
	lbget "github.com/nikhilsbhat/neuron/cloudoperations/loadbalancer/get"
	misc "github.com/nikhilsbhat/neuron/cloudoperations/miscellaneous"
	nwcreate "github.com/nikhilsbhat/neuron/cloudoperations/network/create"
	nwdelete "github.com/nikhilsbhat/neuron/cloudoperations/network/delete"
	nwget "github.com/nikhilsbhat/neuron/cloudoperations/network/get"
	nwupdate "github.com/nikhilsbhat/neuron/cloudoperations/network/update"
	svcreate "github.com/nikhilsbhat/neuron/cloudoperations/server/create"
	svdelete "github.com/nikhilsbhat/neuron/cloudoperations/server/delete"
	svget "github.com/nikhilsbhat/neuron/cloudoperations/server/get"
	svupdate "github.com/nikhilsbhat/neuron/cloudoperations/server/update"
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
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData svcreate.ServerCreateInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			create_server_response, serv_err := myData.CreateServer()
			if serv_err != nil {
				fmt.Fprintf(rw, "%v\n", serv_err)
			} else {
				json_val, _ := json.MarshalIndent(create_server_response, "", " ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

//The below two function is suspended temporarily, and this will be back soon
/*func createbuild(rw http.ResponseWriter, req *http.Request) {
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
                json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
                fmt.Fprintf(rw, "%v\n", string(json_val))
        } else {
                var t createbuildmachine

                err = json.Unmarshal(body, &t)
                if err != nil {
                        json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
                        fmt.Fprintf(rw, "%v\n", string(json_val))
                } else {
                        server_create_input := build.BuildServerCreateInput{t.AppVersion, t.UniqueId, "build-machine", "subnet-d81893b1", "chef-coe-ind", "t2.micro", "aws", "ap-south-1", true}
                        server_response := server_create_input.BuildServerCreate()
                        json_val, _ := json.MarshalIndent(server_response, "", " ")
                        fmt.Fprintf(rw, "%v\n", string(json_val))
                }
        }
}

func startimagemachine(rw http.ResponseWriter, req *http.Request) {
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
                json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
                fmt.Fprintf(rw, "%v\n", string(json_val))
        } else {
                var t startimageinput

                err = json.Unmarshal(body, &t)
                if err != nil {
                        json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
                        fmt.Fprintf(rw, "%v\n", string(json_val))
                } else {
                        server_create_input := buildim.PackerServerCreateInput{t.AppVersion, t.AppName, t.RepoEmail, t.RepoUsername, t.RepoPasswd, t.ArtDomain, t.ArtUsername, t.ArtPasswd, t.InstanceName, t.SubnetId, t.KeyName, t.Flavor, t.Cloud.Cloud, t.Cloud.Region, t.AssignPubIp}
                        server_response, serv_err := server_create_input.BuildServerCreate()
                        if serv_err != nil {
                                fmt.Fprintf(rw, "%v\n", serv_err)
                        } else {
                                json_val, _ := json.MarshalIndent(server_response, "", " ")
                                fmt.Fprintf(rw, "%v\n", string(json_val))
                        }
                }
        }
}*/

func deleteservers(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData svdelete.DeleteServersInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			delete_server_response, del_err := myData.DeleteServer()
			if del_err != nil {
				fmt.Fprintf(rw, "%v\n", del_err)
			} else {
				json_val, _ := json.MarshalIndent(delete_server_response, "", " ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func createnetwork(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData nwcreate.NetworkCreateInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			server_response, ser_resp_err := myData.CreateNetwork()
			if ser_resp_err != nil {
				fmt.Fprintf(rw, "%v\n", ser_resp_err)
			} else {
				json_val, _ := json.MarshalIndent(server_response, "", " ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func getregions(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData misc.GetRegionInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			get_region_response, reg_err := myData.GetRegions()
			if reg_err != nil {
				fmt.Fprintf(rw, "%v\n", reg_err)
			} else {
				json_val, _ := json.Marshal(get_region_response)
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func getsubnets(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData nwget.GetNetworksInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			get_subnets_response, sub_err := myData.GetSubnets()
			if sub_err != nil {
				fmt.Fprintf(rw, "%v\n", sub_err)
			} else {
				json_val, _ := json.MarshalIndent(get_subnets_response, "", " ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func getservers(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData svget.GetServersInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			get_server_response, get_err := myData.GetServersDetails()
			if get_err != nil {
				fmt.Fprintf(rw, "%v\n", get_err)
			} else {
				json_val, _ := json.MarshalIndent(get_server_response, "", " ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func getallservers(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData svget.GetServersInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			getserver_response, serv_err := myData.GetAllServers()
			if serv_err != nil {
				fmt.Fprintf(rw, "%v\n", serv_err)
			} else {
				json_val, _ := json.MarshalIndent(getserver_response, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}

}

func updateservers(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData svupdate.UpdateServersInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			update_server_response, update_err := myData.UpdateServers()
			if update_err != nil {
				fmt.Fprintf(rw, "%v\n", update_err)
			} else {
				json_val, _ := json.MarshalIndent(update_server_response, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}

}

func deletenetworks(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData nwdelete.DeleteNetworkInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			delete_network_response, net_err := myData.DeleteNetwork()
			if net_err != nil {
				fmt.Fprintf(rw, "%v\n", net_err)
			} else {
				json_val, _ := json.MarshalIndent(delete_network_response, "", " ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func createloadbalancer(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData lbcreate.LbCreateInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			create_lb_response, create_lb_err := myData.CreateLoadBalancer()
			if create_lb_err != nil {
				fmt.Fprintf(rw, "%v\n", create_lb_err)
			} else {
				json_val, _ := json.MarshalIndent(create_lb_response, "", " ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func deleteloadbalancer(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData lbdelete.LbDeleteInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			delete_lb_response, del_lb_response := myData.DeleteLoadBalancer()
			if del_lb_response != nil {
				fmt.Fprintf(rw, "%v\n", del_lb_response)
			} else {
				json_val, _ := json.MarshalIndent(delete_lb_response, "", " ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func getloadbalancers(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData lbget.GetLoadbalancerInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			loadbalancers_get_response, get_lb_response := myData.GetLoadbalancers()
			if get_lb_response != nil {
				fmt.Fprintf(rw, "%v\n", get_lb_response)
			} else {
				json_val, _ := json.MarshalIndent(loadbalancers_get_response, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func getallloadbalancers(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData lbget.GetLoadbalancerInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			loadbalancers_get_response, get_lb_response := myData.GetAllLoadbalancer()
			if get_lb_response != nil {
				fmt.Fprintf(rw, "%v\n", get_lb_response)
			} else {
				json_val, _ := json.MarshalIndent(loadbalancers_get_response, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func createimage(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData imcreate.CreateImageInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			image_create_response, img_err := myData.CreateImage()
			if img_err != nil {
				fmt.Fprintf(rw, "%v\n", img_err)
			} else {
				json_val, _ := json.MarshalIndent(image_create_response, "", " ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func deleteimage(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData imdelete.DeleteImageInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			image_delete_response, del_img_err := myData.DeleteImage()
			if del_img_err != nil {
				fmt.Fprintf(rw, "%v\n", del_img_err)
			} else {
				json_val, _ := json.MarshalIndent(image_delete_response, "", " ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func getimages(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData imget.GetImagesInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			image_get_response, get_img_err := myData.GetImage()
			if get_img_err != nil {
				fmt.Fprintf(rw, "%v\n", get_img_err)
			} else {
				json_val, _ := json.MarshalIndent(image_get_response, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func getallimages(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData imget.GetImagesInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			image_get_response, getall_img_err := myData.GetAllImage()
			if getall_img_err != nil {
				fmt.Fprintf(rw, "%v\n", getall_img_err)
			} else {
				json_val, _ := json.MarshalIndent(image_get_response, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func getallnetworks(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData nwget.GetNetworksInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			get_all_networks_response, net_err := myData.GetAllNetworks()
			if net_err != nil {
				fmt.Fprintf(rw, "%v\n", net_err)
			} else {
				json_val, _ := json.MarshalIndent(get_all_networks_response, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func getnetworks(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData nwget.GetNetworksInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			get_networks_response, net_err := myData.GetNetworks()
			if net_err != nil {
				fmt.Fprintf(rw, "%v\n", net_err)
			} else {
				json_val, _ := json.MarshalIndent(get_networks_response, "", "  ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func getcount(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData count.GetCountInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			get_count_response := myData.GetCount()
			json_val, _ := json.MarshalIndent(get_count_response, "", "  ")
			fmt.Fprintf(rw, "%v\n", string(json_val))
		}
	}
}

func updatenetwork(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		json_val, _ := json.Marshal(Error{"Not recieved input in a valid format"})
		fmt.Fprintf(rw, "%v\n", string(json_val))
	} else {
		var myData nwupdate.NetworkUpdateInput

		err = json.Unmarshal(body, &myData)
		if err != nil {
			json_val, _ := json.Marshal(Error{"Unable to unmarshal the entered input. Provide input in valid format"})
			fmt.Fprintf(rw, "%v\n", string(json_val))
		} else {
			net_update_response, net_up_err := myData.UpdateNetwork()
			if net_up_err != nil {
				fmt.Fprintf(rw, "%v\n", net_up_err)
			} else {
				json_val, _ := json.MarshalIndent(net_update_response, "", " ")
				fmt.Fprintf(rw, "%v\n", string(json_val))
			}
		}
	}
}

func createservermock(rw http.ResponseWriter, req *http.Request) {
	create_server_response, serv_err := svcreate.CreateServerMock()
	if serv_err != nil {
		fmt.Fprintf(rw, "%v\n", serv_err)
	} else {
		fmt.Fprintf(rw, "%v\n", create_server_response.DefaultResponse)
	}
}
