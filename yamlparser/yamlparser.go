package yamlparser

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)
type Model struct{
	Source_dir string
	Relationships	[]Relations
}

type Relations struct {
	Name string
}

func ParseYaml(mod *Model) {
	filePath := "./hx.yaml";
	fmt.Printf( "// reading file %s\n", filePath )
	file, err1 := ioutil.ReadFile( filePath )
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
	for k := range mod.Relationships {
		fmt.Printf( "The name %s \n", mod.Relationships[k].Name);
	}
}