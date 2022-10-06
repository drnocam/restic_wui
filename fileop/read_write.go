package fileop

import (
	"encoding/json"
	"io/ioutil"
)

func WriteSettings(){

	file, _ := json.MarshalIndent(data, "", " ")
 
	_ = ioutil.WriteFile("test.json", file, 0644)
	return "dfas";
}

func ReadSettings() {

	return "dfas";
}