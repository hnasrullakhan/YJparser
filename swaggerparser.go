package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"encoding/json"
	"github.com/buger/jsonparser"
	"strings"
	"YJparser/yamlparser"
)


type Swagger struct{
	Swagger string   `json:"swagger"`
	Paths json.RawMessage
	//Paths  map[string]interface{} `json:"paths"`
	//Definitions  map[string]interface{}  `json:"definitions"`
	Definitions json.RawMessage
}

var Defs []Definitionsprops
type Definitionsprops struct {
	name string
	Properties map[string]interface{} `json:"properties"`
	//Properties map[string]json.RawMessage  `json:"properties"`
	indprop []property
	//ppty []*property

}

type property struct {
	Name string
	Type string	`json:"type"`
	Format string	`json:"format"`
	Items interface{}	`json:"items"`
	Enum interface{}	`json:"enums"`
	Refs string 		`json:"$refs"`
	Default bool		`json:"default"`

}
func main() {
	filePath := "./swagger_webapp.json";
	fmt.Printf( "// reading file %s\n", filePath )
	file, err1 := ioutil.ReadFile( filePath )

	// parsing Yaml to populate structures
	if err1 != nil {
		fmt.Printf( "// error while reading file %s\n", filePath )
		fmt.Printf("File error: %v\n", err1)
		os.Exit(1)
	}
	var swag *Swagger
	err2 := json.Unmarshal(file,&swag)
	fmt.Printf("this is swag value : %s \n",swag.Swagger)
	fmt.Println("================================")
	meta := yamlparser.Model{}
	yamlparser.ParseYaml(&meta)
	fmt.Println(meta)
	for i := range meta.Relationships {
		fmt.Println(string(meta.Relationships[i].Name))

		v, _, _, _ := jsonparser.Get(swag.Paths, "/" + meta.Relationships[i].Name, "get", "responses", "200", "schema", "$ref")
		fmt.Printf("%s\n", string(v))
		fmt.Println("================================")
		defintion := strings.SplitAfter(string(v), "#/definitions/")
		fmt.Println(defintion[1])
		//fmt.Printf( "The defintion %s \n", (swag.Definitions));
		def, _, _, _ := jsonparser.Get(swag.Definitions, defintion[1])
		//fmt.Println(string(def))
		var vardef Definitionsprops
		_ = json.Unmarshal(def, &vardef)
		vardef.name = defintion[1]
		fmt.Println(vardef)
		//v,_,_,_ = jsonparser.Get(swag.Definitions,defintion[1],"properties")
		fmt.Println("================================")

		for key, val := range vardef.Properties {
			lname := key
			/*var tmpProp property
			_ = json.Unmarshal(val, &tmpProp )
			tmpProp.Name = lname
			fmt.Println("================================")

			vardef.indprop = append(vardef.indprop, tmpProp)
			fmt.Println("Name:",tmpProp.Name)
			fmt.Println("Type:",tmpProp.Type)
			fmt.Println("Format:",tmpProp.Format)
			fmt.Println("Enum:",tmpProp.Enum)
			fmt.Println("Items:",tmpProp.Items)
			fmt.Println("Default:",tmpProp.Default)

			fmt.Println("================================")*/


			ltype := val.(map[string]interface{})["type"]
			if ltype == nil {
				ltype = ""
			}

			lFormat := val.(map[string]interface{})["format"]
			if lFormat == nil {
				lFormat = ""
			}
			lItems := val.(map[string]interface{})["items"]
			if lItems == nil {
				lItems = ""
			}
			lEnum := val.(map[string]interface{})["enum"]
			if lEnum == nil {
				lEnum = ""
			}
			lDefault := val.(map[string]interface{})["default"]
			if lDefault == nil {
				lDefault = false
			}
			lRefs := val.(map[string]interface{})["$ref"]
			if lRefs == nil {
				lRefs = ""
			}
			fmt.Println("================================")
			fmt.Println("new property")
			fmt.Println("Name:", lname)
			fmt.Println("Type:", ltype)
			fmt.Println("Format:", lFormat)
			fmt.Println("Items:", lItems)
			fmt.Println("Enum:", lEnum)
			fmt.Println("Refs:", lRefs)
			fmt.Println("Default:", lDefault)
			fmt.Println("================================")

			tmpProperty := property{Name:string(lname), Type:ltype.(string), Format:lFormat.(string), Items:lItems, Enum:lEnum,Refs:lRefs.(string), Default:lDefault.(bool)}

			vardef.indprop = append(vardef.indprop, tmpProperty)

		}
		Defs = append(Defs,vardef)

	}
	fmt.Println("**************************************************************")
	fmt.Println("**************************************************************")
	fmt.Println("**************************************************************")
	fmt.Println("**************************************************************")

	//var propertymap map[string]json.RawMessage
	for iter := range Defs{
		fmt.Println(Defs[iter].name)
		fmt.Println(Defs[iter].Properties)
		fmt.Println(Defs[iter].indprop)


		fmt.Println("**************************************************************")

	}
	fmt.Print(err2)



}

