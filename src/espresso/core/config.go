/**
* The purpose of this file
* is to parse the configuration
*/

package core

import "encoding/json"
import "os"
import "log"
import "io/ioutil"

const ConfigFile = "config.json"

type Config struct {
	Params	ConfigParams
}

type ConfigParams struct {
	IP		string
	Port		int
	Name		string
	Webroot		string
	
	DBHost		string
	DBUsername	string
	DBPassword	string
	DBTable		string
	DBPrefix	string
	
	//Databases	[]DatabaseObject
}

// In the future, we could allow multiple database connections
// This structure serves as a placeholder if we decide to.
type DatabaseObject struct {
	Host		string
	Username	string
	Password	string
	Table		string
	Prefix		string
}

func ParseConfig() ConfigParams {
	if _, err := os.Stat(ConfigFile); err != nil {
		log.Fatal(ConfigFile + " cannot be loaded! ")
	}
	configFile, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		log.Fatal(ConfigFile + "cannot be read! ")
	}
	var tmpStruct ConfigParams
	err = json.Unmarshal(configFile, &tmpStruct)
	if err != nil {
		log.Fatal(err)
	}
	return tmpStruct
}

func InitConfig() *Config {
	C := &Config{}
	C.Params = ParseConfig()
	return C
	
}
