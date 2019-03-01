package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/frictionlessdata/datapackage-go/datapackage"
	"github.com/qri-io/dataset"
	"github.com/spf13/cobra"
)

// ImportCmd brings an open knowledge datapackage into Qri
var ImportCmd = &cobra.Command{
	Use:   "import",
	Short: "import a datapackage into qri",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkgPath := args[0]
		pkg, err := datapackage.Load(pkgPath)
		if err != nil {
			fmt.Printf("error loading data package: %s\n", err.Error())
			return
		}
		log.Debugf("Data package '%s' loaded.\n", pkg.Descriptor()["name"])

		ds, err := dataPackageToDataset(pkg)
		if err != nil {
			fmt.Printf("error converting to dataset: %s\n", err.Error())
			return
		}

		body := ds.Body
		ds.Body = nil

		dsPath := filepath.Join(filepath.Dir(pkgPath), "dataset.json")
		f, err := os.Create(dsPath)
		if err != nil {
			fmt.Printf("error creating temp dataset file: %s\n", err.Error())
			return
		}
		if err := json.NewEncoder(f).Encode(ds); err != nil {
			fmt.Printf("error encoding dataset: %s\n", err.Error())
			return
		}
		f.Close()
		defer os.Remove(dsPath)

		bodyPath := filepath.Join(filepath.Dir(pkgPath), "body.json")
		bodyF, err := os.Create(bodyPath)
		if err != nil {
			fmt.Printf("error creating temp dataset file: %s\n", err.Error())
			return
		}
		if err := json.NewEncoder(bodyF).Encode(body); err != nil {
			fmt.Printf("error encoding dataset: %s\n", err.Error())
			return
		}
		bodyF.Close()
		defer os.Remove(bodyPath)

		name := filepath.Base(filepath.Dir(pkgPath))
		if str, ok := pkg.Descriptor()["name"].(string); ok {
			name = str
		}
		log.Debugf("using name: %s", name)

		qriCmd := exec.Command("qri", "save", "--file="+dsPath, "--body="+bodyPath, name)
		qriCmd.Stderr = os.Stderr
		qriCmd.Stdout = os.Stdout
		if err := qriCmd.Run(); err != nil {
			fmt.Printf("qri error: %s\n", err.Error())
			return
		}
	},
}

func dataPackageToDataset(pkg *datapackage.Package) (ds *dataset.Dataset, err error) {
	ds = &dataset.Dataset{}
	if ds.Meta, err = meta(pkg); err != nil {
		return
	}

	if ds.Body, err = combineResourcesBody(pkg); err != nil {
		return
	}

	return
}

func meta(pkg *datapackage.Package) (*dataset.Meta, error) {
	md := &dataset.Meta{}
	des := pkg.Descriptor()

	if err := md.Set("importer", "qri-datapackage"); err != nil {
		return nil, err
	}

	if str, ok := des["title"].(string); ok {
		md.Title = str
	}
	if str, ok := des["description"].(string); ok {
		md.Description = str
	}
	if str, ok := des["version"].(string); ok {
		md.Version = str
	}
	if str, ok := des["homepage"].(string); ok {
		md.HomeURL = str
	}
	if lisc, ok := des["license"].(map[string]interface{}); ok {
		md.License = &dataset.License{}
		// TODO (B5) - name field is required
		// if str, ok := lisc["name"].(string); ok {
		// 	md.License.Name = str
		// }
		if str, ok := lisc["url"].(string); ok {
			md.License.URL = str
		}
		if str, ok := lisc["type"].(string); ok {
			md.License.Type = str
		}
	}

	return md, nil
}

func combineResourcesBody(pkg *datapackage.Package) (body interface{}, err error) {
	bd := map[string]interface{}{}
	for _, res := range pkg.Resources() {
		r := res.Descriptor()
		r["data"], err = res.ReadAll()
		if err != nil {
			return
		}
		bd[res.Name()] = r
	}
	return bd, nil
}
