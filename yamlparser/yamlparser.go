package yamlparser

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)
type Model struct{
	Source_dir string
	Package string
	Objects	[]*OBJECTS
	//Apitypes []OBJECTS
}

type OBJECTS struct {
	SourceName string	`yaml:"sourceName"`
	TargetName string	`yaml:"targetName"`
	Usage string		`yaml:"usage`
}


func ParseYaml(mod *Model) bool {
	filePath := "./hx.yaml";
	fmt.Printf( "// reading file %s\n", filePath )
	file, err1 := ioutil.ReadFile( filePath )
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Meta Model file doesnt exist:",filePath)
		return false
	}
	if err1 != nil {
		fmt.Printf( "// error while reading file %s\n", filePath )
		fmt.Printf("File error: %v\n", err1)
		os.Exit(1)
	}

	//var mod *Model

	err2 := yaml.Unmarshal(file, &mod)
	if err2 != nil {
		fmt.Println("error:", err2)
		os.Exit(1)
	}

	fmt.Println( "// loop over array of structs of Relations" )
	fmt.Printf( "The Model '%s\n", mod.Source_dir  );
	fmt.Printf( "The Package '%s\n", mod.Package  );

	for k := range mod.Objects {
		fmt.Printf( "The Source name %s \n", mod.Objects[k].SourceName);
		fmt.Printf( "The Target name %s \n", mod.Objects[k].TargetName);
		fmt.Printf( "The Target name %s \n", mod.Objects[k].Usage);
	}
	return true

}