package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"encoding/json"
	"github.com/buger/jsonparser"
	"strings"
	"YJparser/yamlparser"
	"text/template"

)



type Swagger struct{
	SwagVersion string   `json:"swagger"`
	Paths json.RawMessage
	//Paths  map[string]interface{} `json:"paths"`
	//Definitions  map[string]interface{}  `json:"definitions"`
	Definitions json.RawMessage
	Defs []*Definitionsprops
}
var Swag Swagger
type Definitionsprops struct {
	Name string
	Type string
	Properties map[string]interface{} `json:"properties"`
	//Properties map[string]json.RawMessage  `json:"properties"`
	Indprop []Property
	//ppty []*property

}

type Property struct {
	Name string
	Type string	`json:"type"`
	Format string	`json:"format"`
	Items interface{}	`json:"items"`
	Enum interface{}	`json:"enums"`
	Refs string 		`json:"$refs"`
	AdditionalProperties  AdditionalProps	`json:additionalProperties`
	Default bool		`json:"default"`

}

type AdditionalProps struct{
	Type string
	Refs string
	Items interface{}
}

func ParseDefintions( aInDefintionName string, jsonRawDef []byte){

	var vardef Definitionsprops
	_ = json.Unmarshal(jsonRawDef, &vardef)
	vardef.Name = aInDefintionName
	vardef.Type = "object"
	fmt.Println(vardef)
	//v,_,_,_ = jsonparser.Get(swag.Definitions,defintion[1],"properties")
	fmt.Println("================================")

	for key, val := range vardef.Properties {
		lname := key

		ltype := val.(map[string]interface{})["type"]
		if ltype == nil {
			ltype = ""
		}

		lFormat := val.(map[string]interface{})["format"]
		if lFormat == nil {
			lFormat = ""
		}
		lItems := val.(map[string]interface{})["items"]

		if lItems != nil {
			fmt.Println("lItems:",lItems)
			for keyItem,valItem := range (lItems).(map[string]interface{}) {

				if keyItem == "$ref"{
					defintion := strings.SplitAfter((valItem.(string)), "#/definitions/")
					fmt.Println(defintion[1])
					//fmt.Printf( "The defintion %s \n", (swag.Definitions));
					def, _, _, _ := jsonparser.Get(Swag.Definitions, defintion[1])
					fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxTHISISCHILDxxxxxxxxxxxxxxxxxxxxxxxxxxx")

					ParseDefintions(defintion[1],def)
					fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxCHILDENDSHERExxxxxxxxxxxxxxxxxxxxxxxxxxx")

				}
			}
		}else{
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
		if lRefs != nil {
			defintion := strings.SplitAfter((lRefs).(string), "#/definitions/")
			fmt.Println(defintion[1])
			//fmt.Printf( "The defintion %s \n", (swag.Definitions));
			def, _, _, _ := jsonparser.Get(Swag.Definitions, defintion[1])
			fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxTHISISCHILDxxxxxxxxxxxxxxxxxxxxxxxxxxx")

			ParseDefintions(defintion[1],def)
			ltype = defintion[1]
			fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxCHILDENDSHERExxxxxxxxxxxxxxxxxxxxxxxxxxx")

		}else{
			lRefs = ""
		}
		var lAddProps AdditionalProps
		tmpAddProps := val.(map[string]interface{})["additionalProperties"]
		if tmpAddProps != nil {
			tmpProp := tmpAddProps.(map[string]interface{})["type"]
			if (tmpProp.(string) == "object" || tmpProp.(string) == "array" ){

				lAddProps.Type = tmpProp.(string)

			}

			tmpProp = tmpAddProps.(map[string]interface{})["$refs"]
			if tmpProp != nil {
				lAddProps.Refs = tmpProp.(string)

			}else{
				lAddProps.Refs =""
			}
			tmpProp = tmpAddProps.(map[string]interface{})["items"]
			if tmpProp != nil {
				fmt.Println("lItems:",tmpProp)
				for keyItem,valItem := range (tmpProp).(map[string]interface{}) {

					if keyItem == "$ref"{
						defintion := strings.SplitAfter((valItem.(string)), "#/definitions/")
						fmt.Println(defintion[1])
						//lAddProps.Type
						//fmt.Printf( "The defintion %s \n", (swag.Definitions));
						def, _, _, _ := jsonparser.Get(Swag.Definitions, defintion[1])
						fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxTHISISCHILDxxxxxxxxxxxxxxxxxxxxxxxxxxx")

						ParseDefintions(defintion[1],def)
						fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxCHILDENDSHERExxxxxxxxxxxxxxxxxxxxxxxxxxx")

					}
				}
			}else{
				lAddProps.Items = ""
			}
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
		fmt.Println("lAddProps:", tmpAddProps)

		fmt.Println("================================")

		tmpProperty := Property{Name:string(lname), Type:ltype.(string), Format:lFormat.(string), Items:lItems, Enum:lEnum,Refs:lRefs.(string), Default:lDefault.(bool)}

		vardef.Indprop = append(vardef.Indprop, tmpProperty)

	}
	Swag.Defs = append(Swag.Defs,&vardef)
	fmt.Println("Function ends here")
}

const (
	dirFilePerm os.FileMode = 0777

	modelDirectory string = "model"
	metaDirectory  string = "meta/"
	projectKey     string = "bitbucket-eng-sjc1.cisco.com/an"
)


func loadTemplate(name string) *template.Template {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		fmt.Print("Environment variable 'GOPATH' must be set.")
	}

	fname := name
	fmt.Printf("Parsing template '%s'", fname)
	t, err := template.ParseFiles(fname)
	if err != nil {
		fmt.Print("Invalid adgen template : ", err)
	}
	return t
}

func createDir(name string) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		err := os.MkdirAll(name, dirFilePerm)
		if err != nil {
			fmt.Printf("cannot make dir '%s' : %s", name, err)
		}
	}
}

func createFile(name string) *os.File {

	f, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, dirFilePerm)
	if err != nil {
		fmt.Printf("Failed to create file '%s' : %s", name, err.Error())
	}

	return f
}

func generateModel(mg Swagger){
	var t *template.Template
	var f *os.File
	fmt.Print("Generating Model File")
	t = loadTemplate("Model.gotmpl")

	createDir("GenModel")

	f = createFile("GenModel/ModelGen.yaml")
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			fmt.Print("Failed to close file: ", err)
		}
	}(f)

	err := t.Execute(f, mg)
	if err != nil {
		fmt.Print("Error processing template")
	}


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
	//var swag *Swagger
	err2 := json.Unmarshal(file,&Swag)
	fmt.Printf("this is swag value : %s \n",Swag.SwagVersion)
	fmt.Println("================================")
	meta := yamlparser.Model{}
	yamlparser.ParseYaml(&meta)
	fmt.Println(meta)
	for i := range meta.Relationships {
		fmt.Println(string(meta.Relationships[i].Name))

		v, _, _, _ := jsonparser.Get(Swag.Paths, "/" + meta.Relationships[i].Name, "get", "responses", "200", "schema", "$ref")
		if string(v) == "" {
			it, _, _, _ := jsonparser.Get(Swag.Paths, "/" + meta.Relationships[i].Name, "get", "responses", "200", "schema", "type")
			if string(it) == "array"{

				v, _, _, _ = jsonparser.Get(Swag.Paths, "/" + meta.Relationships[i].Name, "get", "responses", "200", "schema", "items","$ref")
			}
		}
		fmt.Printf("%s\n", string(v))
		fmt.Println("================================")
		defintion := strings.SplitAfter(string(v), "#/definitions/")
		fmt.Println(defintion[1])
		//fmt.Printf( "The defintion %s \n", (swag.Definitions));
		def, _, _, _ := jsonparser.Get(Swag.Definitions, defintion[1])

		ParseDefintions(defintion[1],def)
		//fmt.Println(string(def))


	}
	/*
	fmt.Println("**************************************************************")
	fmt.Println("**************************************************************")
	fmt.Println("**************************************************************")
	fmt.Println("**************************************************************")

	//var propertymap map[string]json.RawMessage
	for iter := range Swag.Defs{
		fmt.Println(Swag.Defs[iter].Name)
		fmt.Println(Swag.Defs[iter].Properties)
		fmt.Println(Swag.Defs[iter].Indprop)


		fmt.Println("**************************************************************")

	}*/
	generateModel(Swag)
	fmt.Print(err2)



}

