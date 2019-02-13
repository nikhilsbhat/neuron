package mongo

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/nikhilsbhat/neuron/database"
	"strings"
)

func StoreCIdata(d database.DataDetail, data database.CiData) (interface{}, error) {

	session := ((database.Db).(*mgo.Session)).Copy()
	defer session.Close()
	c := session.DB(d.Database).C(d.Collection)

	query := c.Find(bson.M{"ciname": data.CiName, "ciurl": data.CiURL})
	resp := []bson.M{}
	qry_err := query.All(&resp)
	if qry_err != nil {
		return nil, qry_err
	}
	if len(resp) == 0 {
		ins_err := c.Insert(data)
		if ins_err != nil {
			return nil, ins_err
		}
		return "Records Created Successfully, check previous page for the details", nil
	}
	return nil, fmt.Errorf("The details you enetered matches with existing records")
}

func GetCiData(n string, d database.DataDetail) (database.CiData, error) {

	session := ((database.Db).(*mgo.Session)).Copy()
	defer session.Close()
	c := session.DB(d.Database).C(d.Collection)

	query := c.Find(bson.M{})
	resp := []bson.M{}
	qry_err := query.All(&resp)
	if qry_err != nil {
		return database.CiData{}, qry_err
	}
	_ = resp
	return database.CiData{}, nil
}

func CreateUser(d database.DataDetail, data database.UserData) (interface{}, error) {

	session := ((database.Db).(*mgo.Session)).Copy()
	defer session.Close()
	c := session.DB(d.Database).C(d.Collection)

	query := c.Find(bson.M{"username": data.UserName, "password": data.Password})
	resp := []bson.M{}
	qry_err := query.All(&resp)
	if qry_err != nil {
		return nil, qry_err
	}
	if len(resp) == 0 {
		// making indexing dynamic
		index_query := c.Find(bson.M{})
		index_resp := []bson.M{}
		index_qry_err := index_query.Sort("_id").All(&index_resp)
		if index_qry_err != nil {
			return nil, index_qry_err
		}
		var insert_value database.UserData
		if data.CloudProfiles != nil {
			cld_prf := make([]database.CloudProfiles, 0)
			//cld_prf = append(cld_prf, data.CloudProfiles)
			insert_value = database.UserData{Id: (index_resp[len(index_resp)-1]["_id"].(int)) + 1, UserName: data.UserName, Password: data.Password, CloudProfiles: cld_prf}
		} else {
			insert_value = database.UserData{Id: (index_resp[len(index_resp)-1]["_id"].(int)) + 1, UserName: data.UserName, Password: data.Password}
		}
		ins_err := c.Insert(insert_value)
		if ins_err != nil {
			return nil, ins_err
		}
		return "User details saved successfully", nil

	}
	if len(resp) > 1 {
	}
	return "We cannot take the data you entered, Because we found the data matchces your entries", nil
}

func UpdateUser(d database.DataDetail, data database.UserData) (interface{}, error) {

	session := ((database.Db).(*mgo.Session)).Copy()
	defer session.Close()
	c := session.DB(d.Database).C(d.Collection)

	query := c.Find(bson.M{"username": data.UserName, "password": data.Password})
	resp := []bson.M{}
	qry_err := query.All(&resp)
	if qry_err != nil {
		return nil, qry_err
	}
	existing := bson.M{"username": data.UserName, "password": data.Password}

	var val_appnd []interface{}
	for _, value := range resp {
		for k, v := range value {
			if k == "cloudprofiles" {
				for _, v1 := range v.([]interface{}) {
					val_appnd = append(val_appnd, v1)
				}
			}
		}
	}

	val_appnd = append(val_appnd, data.CloudProfiles)
	change := bson.M{"$set": bson.M{"cloudprofiles": val_appnd}}
	_, up_err := c.Upsert(existing, change)
	if up_err != nil {
		return nil, up_err
	}
	return "User profile updated successfully", nil
}

func GetUserDetails(d database.DataDetail, data database.UserData) (interface{}, error) {

	session := ((database.Db).(*mgo.Session)).Copy()
	defer session.Close()
	c := session.DB(d.Database).C(d.Collection)

	query := c.Find(bson.M{"username": data.UserName, "password": data.Password})
	resp := []bson.M{}
	qry_err := query.Sort("_id").All(&resp)
	if qry_err != nil {
		return nil, qry_err
	}
	for _, value := range resp {
		return database.UserData{UserName: value["username"].(string), Password: value["password"].(string)}, nil
	}
	return nil, fmt.Errorf("Something went wrong while fetching user details")
}

func GetCloudCredentails(d database.DataDetail, data database.UserData, cred database.GetCloudAccess) (database.CloudProfiles, error) {

	session := ((database.Db).(*mgo.Session)).Copy()
	defer session.Close()
	c := session.DB(d.Database).C(d.Collection)

	query := c.Find(bson.M{"username": data.UserName, "password": data.Password})
	resp := []bson.M{}
	qry_err := query.Sort("_id").All(&resp)

	if qry_err != nil {
		fmt.Println(qry_err)
	}

	for _, value := range resp {
		for k, v := range value {
			if k == "cloudprofiles" {
				for _, v1 := range v.([]interface{}) {
					if (v1.(bson.M)["name"].(string) == cred.ProfileName) && (strings.ToLower(v1.(bson.M)["cloud"].(string)) == cred.Cloud) {
						return database.CloudProfiles{Name: v1.(bson.M)["name"].(string), Cloud: v1.(bson.M)["cloud"].(string), KeyId: v1.(bson.M)["keyid"].(string), SecretAccess: v1.(bson.M)["secretaccess"].(string)}, nil
					}
				}
			}
		}
	}
	return database.CloudProfiles{}, fmt.Errorf("Unable to find cloud credentials for the profile enetered")
}
