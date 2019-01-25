package dbcommon

import (
	//"github.com/globalsign/mgo/bson"
	"fmt"
	"github.com/globalsign/mgo"
	"neuron/database"
	"neuron/database/fs"
	"neuron/database/mongodb"
	err "neuron/error"
)

func ConfigDb(d database.Storage) (interface{}, error) {

	if d.Db != nil {
		switch (d.Db).(type) {
		case *mgo.Session:
			database.Db = d.Db
			return nil, nil
		default:
			return nil, err.UnknownDbType()
		}
	}
	if d.Fs != "" {
		database.Db = d.Fs
		return nil, nil
	}
	return nil, fmt.Errorf("Oops..!! an error occured. We did not receive valid input to configure DB")
}

func StoreCIdata(d database.DataDetail, data database.CiData) (interface{}, error) {

	if database.Db != nil {
		switch database.Db.(type) {
		case *mgo.Session:
			status, stat_err := mongo.StoreCIdata(d, data)
			if stat_err != nil {
				return nil, stat_err
			}
			return status, nil
		default:
			return nil, err.UnknownDbType()
		}
	}
	return "Database is not configured, we are not supporting filesystem now", nil
}

func GetCiData(ci string, d ...database.DataDetail) (database.CiData, error) {

	if database.Db != nil {
		switch (database.Db).(type) {
		case *mgo.Session:
			status, stat_err := mongo.GetCiData(ci, d[0])
			if stat_err != nil {
				return database.CiData{}, stat_err
			}
			return status, nil
		case string:
			status, stat_err := fs.GetCiData(ci)
			if stat_err != nil {
				return database.CiData{}, stat_err
			}
			return status, nil
		default:
			return database.CiData{}, err.UnknownDbType()
		}
	}
	return database.CiData{}, fmt.Errorf("Database is not configured, we are not supporting filesystem now")
}

func CreateUser(d database.DataDetail, data database.UserData) (interface{}, error) {

	if database.Db != nil {
		switch (database.Db).(type) {
		case *mgo.Session:
			status, stat_err := mongo.CreateUser(d, data)
			if stat_err != nil {
				return nil, stat_err
			}
			return status, nil
		default:
			return nil, err.UnknownDbType()
		}
	}
	return "Database is not configured, we are not supporting filesystem now", nil
}

func UpdateUser(session interface{}, d database.DataDetail, data database.UserData) (interface{}, error) {

	if database.Db != nil {
		switch (database.Db).(type) {
		case *mgo.Session:
			status, stat_err := mongo.UpdateUser(d, data)
			if stat_err != nil {
				return nil, stat_err
			}
			return status, nil
		default:
			return nil, err.UnknownDbType()
		}
	}
	return "Database is not configured, we are not supporting filesystem now", nil
}

func GetUserDetails(d database.DataDetail, data database.UserData) (interface{}, error) {

	if database.Db != nil {
		switch (database.Db).(type) {
		case *mgo.Session:
			status, stat_err := mongo.GetUserDetails(d, data)
			if stat_err != nil {
				return nil, stat_err
			}
			return status, nil
		default:
			return nil, err.UnknownDbType()
		}
	}
	return "Database is not configured, we are not supporting filesystem now", nil
}

/*func GetUsersDetails(d database.DataDetail, data database.UserData) (interface{}, error) {

	if database.Db != nil {
		switch (database.Db).(type) {
		case *mgo.Session:
			status, stat_err := mongo.GetUsersDetails(d, data)
			if stat_err != nil {
				return nil, stat_err
			}
			return status, nil
		case string:
			status, stat_err := fs.GetUsersDetails(data, cred)
			if stat_err != nil {
				return database.CloudProfiles{}, stat_err
			}
			return status, nil
		default:
			return nil, err.UnknownDbType()
		}
	}
	return "Database is not configured, we are not supporting filesystem now", nil
}*/

func GetCloudCredentails(data database.UserData, cred database.GetCloudAccess, d ...database.DataDetail) (database.CloudProfiles, error) {

	if database.Db != nil {
		switch (database.Db).(type) {
		case *mgo.Session:
			status, stat_err := mongo.GetCloudCredentails(d[0], data, cred)
			if stat_err != nil {
				return database.CloudProfiles{}, stat_err
			}
			return status, nil
		case string:
			status, stat_err := fs.GetCloudCredentails(data, cred)
			if stat_err != nil {
				return database.CloudProfiles{}, stat_err
			}
			return status, nil
		default:
			return database.CloudProfiles{}, err.UnknownDbType()
		}
	}

	return database.CloudProfiles{}, err.DbNotConfiguredError()
}
