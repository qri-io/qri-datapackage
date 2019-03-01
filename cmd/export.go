package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/frictionlessdata/datapackage-go/datapackage"
	"github.com/frictionlessdata/datapackage-go/validator"
	"github.com/qri-io/dataset"
	"github.com/spf13/cobra"
)

// ExportCmd turns a Qri datasets into an open knowledge foundation datapackage
var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "write a qri dataset as a datapackage zip archive",
	Example: `$ qri-datapackage export me/dataset

$ qri-datapackage export me/dataset package.zip`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		stream := &bytes.Buffer{}
		qriCmd := exec.Command("qri", "export", "--format=json", "--output=temp.json", name)
		qriCmd.Stdout = stream
		if err := qriCmd.Run(); err != nil {
			fmt.Println(stream.String())
			fmt.Printf("qri export error: %s\n", err.Error())
			return
		}
		defer os.Remove("temp.json")
		data, err := ioutil.ReadFile("temp.json")
		if err != nil {
			fmt.Printf("error opening data: %s", err.Error())
			return
		}

		ds := &dataset.Dataset{}
		if err := ds.UnmarshalJSON(data); err != nil {
			fmt.Printf("error unmarshaling dataset: %s", err.Error())
			return
		}

		pkg, err := DatasetToDataPackage(ds)
		if err != nil {
			fmt.Printf("error creating datapackage: %s", err.Error())
			return
		}

		pkgName := ds.Name + "_datapackage.zip"
		if len(args) == 2 {
			pkgName = args[1]
		}
		if err = pkg.Zip(pkgName); err != nil {
			fmt.Printf("error writing zip: %s", err.Error())
			return
		}
		fmt.Printf("exported datapackage zip archive to: %s\n", pkgName)
	},
}

// DatasetToDataPackage converts a dataset
func DatasetToDataPackage(ds *dataset.Dataset) (pkg *datapackage.Package, err error) {
	des, err := packageDescriptor(ds.Name, ds.Meta)
	if err != nil {
		return nil, err
	}

	// madeWithImport := false
	// if ds.Meta != nil {
	// 	if str, ok := ds.Meta.Meta()["importer"]; ok && str == "qri-dataset" {
	// 		madeWithImport = true
	// 	}
	// }
	// if !madeWithImport {
	// 	log.Debug("couldn't find meta exporter tag, this might not work")
	// }

	if body, ok := ds.Body.(map[string]interface{}); ok {
		rsc := []interface{}{}
		for _, r := range body {
			if resource, ok := r.(map[string]interface{}); ok {
				delete(resource, "path")
			}
			rsc = append(rsc, r)
		}
		des["resources"] = rsc
	}

	return datapackage.New(des, "", validator.InMemoryLoader())
}

func packageDescriptor(name string, md *dataset.Meta) (map[string]interface{}, error) {
	des := map[string]interface{}{}

	des["name"] = name
	if md != nil {
		des["title"] = md.Title
		des["description"] = md.Description
		des["version"] = md.Version
		des["homepage"] = md.HomeURL
		des["license"] = map[string]interface{}{
			"type": md.License.Type,
			"url":  md.License.URL,
		}
	}

	return des, nil
}
