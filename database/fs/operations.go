package fs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/nikhilsbhat/neuron/database"
	err "github.com/nikhilsbhat/neuron/error"
	log "github.com/nikhilsbhat/neuron/logger"
	"os"
	"sort"
	"strings"
	//"time"
)

/*const (
	usernotfound =
    cloudcredentialsnotfound =
)*/

func readUserData() ([]database.UserData, error) {

	if _, searcherr := os.Stat(fmt.Sprintf("%s/users.json", (database.Db).(string))); os.IsNotExist(searcherr) {
		return nil, err.UsersNotFound()
	}
	usrdata, loadusrerr := ioutil.ReadFile(fmt.Sprintf("%s/users.json", (database.Db).(string)))
	if loadusrerr != nil {
		return nil, err.ReadFileError()
	}

	// decoding userdata to json
	decoder := json.NewDecoder(bytes.NewReader([]byte(usrdata)))
	usr := make([]database.UserData, 0)
	decoderr := decoder.Decode(&usr)
	if decoderr != nil {
		log.Error(err.JsonDecodeError())
		return nil, err.InvalidUsersFile()
	}
	return usr, nil
}

func readCiData() ([]database.CiData, error) {

	if _, searcherr := os.Stat(fmt.Sprintf("%s/cidata.json", (database.Db).(string))); os.IsNotExist(searcherr) {
		return nil, err.UsersNotFound()
	}
	ciraw, loadusrerr := ioutil.ReadFile(fmt.Sprintf("%s/cidata.json", (database.Db).(string)))
	if loadusrerr != nil {
		return nil, err.ReadFileError()
	}

	// decoding userdata to json
	decoder := json.NewDecoder(bytes.NewReader([]byte(ciraw)))
	cidata := make([]database.CiData, 0)
	decoderr := decoder.Decode(&cidata)
	if decoderr != nil {
		log.Error(err.JsonDecodeError())
		return nil, err.InvalidCiDataFile()
	}
	return cidata, nil
}

func StoreCIdata(d database.DataDetail, data database.CiData) (interface{}, error) {

	return nil, fmt.Errorf("The details you enetered matches with existing records")
}

func GetCiData(c string) (database.CiData, error) {

	cidata, cierr := readCiData()
	if cierr != nil {
		return database.CiData{}, cierr
	}
	//sorting the struct
	id := func(p1, p2 *database.CiData) bool {
		return p1.Id < p2.Id
	}
	Byci(id).Sort(cidata)

	for _, value := range cidata {
		if strings.ToLower(value.CiName) == c {
			return value, nil
		}
	}

	return database.CiData{}, fmt.Errorf("There are no CI data stored so far with the name you provided.")
}

func CreateUser(d database.DataDetail, data database.UserData) (interface{}, error) {

	return "We cannot take the data you entered, Because we found the data matchces your entries", nil
}

func UpdateUser(d database.DataDetail, data database.UserData) (interface{}, error) {

	return "User profile updated successfully", nil
}

func GetUsersDetails(d database.DataDetail, data database.UserData) ([]database.UserData, error) {

	return nil, fmt.Errorf("Something went wrong while fetching user details")
}

func GetUserDetails(d database.DataDetail, data database.UserData) (database.UserData, error) {

	/*session := ((database.Db).(*mgo.Session)).Copy()
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
	}*/
	return database.UserData{}, fmt.Errorf("Something went wrong while fetching user details")
}

func GetCloudCredentails(data database.UserData, cred database.GetCloudAccess) (database.CloudProfiles, error) {

	usr, usrerr := readUserData()
	if usrerr != nil {
		return database.CloudProfiles{}, usrerr
	}
	//sorting the struct
	id := func(p1, p2 *database.UserData) bool {
		return p1.Id < p2.Id
	}
	Byusr(id).Sort(usr)

	for _, value := range usr {
		if (value.UserName == data.UserName) && (value.Password == data.Password) {
			for _, v := range value.CloudProfiles {
				if (strings.ToLower(v.Cloud) == (strings.ToLower(cred.Cloud))) && (v.Name == cred.ProfileName) {
					return database.CloudProfiles{Name: v.Name, Cloud: v.Cloud, KeyId: v.KeyId, SecretAccess: v.SecretAccess}, nil
				}
			}
		}
	}

	return database.CloudProfiles{}, fmt.Errorf("Unable to find cloud credentials for the profile enetered")
}

//miscellaneous functions to support above functions
//functions used to support usr management if database is of type filesystem
type usrSorter struct {
	usrs []database.UserData
	by   func(p1, p2 *database.UserData) bool
}
type Byusr func(p1, p2 *database.UserData) bool

func (by Byusr) Sort(users []database.UserData) {
	ps := &usrSorter{
		usrs: users,
		by:   by,
	}
	sort.Sort(ps)
}

func (s *usrSorter) Len() int {
	return len(s.usrs)
}

func (s *usrSorter) Less(i, j int) bool {
	return s.by(&s.usrs[i], &s.usrs[j])
}

func (s *usrSorter) Swap(i, j int) {
	s.usrs[i], s.usrs[j] = s.usrs[j], s.usrs[i]
}

//functions used to manage ci credentials if database is of type filesystem
type ciSorter struct {
	cicrds []database.CiData
	by     func(p1, p2 *database.CiData) bool
}
type Byci func(p1, p2 *database.CiData) bool

func (by Byci) Sort(cis []database.CiData) {
	ps := &ciSorter{
		cicrds: cis,
		by:     by,
	}
	sort.Sort(ps)
}

func (s *ciSorter) Len() int {
	return len(s.cicrds)
}

func (s *ciSorter) Less(i, j int) bool {
	return s.by(&s.cicrds[i], &s.cicrds[j])
}

func (s *ciSorter) Swap(i, j int) {
	s.cicrds[i], s.cicrds[j] = s.cicrds[j], s.cicrds[i]
}
