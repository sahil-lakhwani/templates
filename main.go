package main

import (
	"fmt"
	"os"
	"log"
	"path/filepath"

	"text/template"
	"github.com/Masterminds/sprig"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	// variable values to be substituted
	values := map[string]interface{}{
		"replicaCount": 3,	
	}


	// example K8s manifest
	// checks if a storage class "gp2" is available, else uses "standard"
	manifest := `apiVersion: zookeeper.pravega.io/v1beta1
        kind: ZookeeperCluster
        metadata:
            name: zookeeper
            namespace: zookeeper
        spec:
        	{{- $storageClass := "gp2" -}}
            replicas: {{ .replicaCount}}
            storageClass: {{ if (lookup "storage.k8s.io/v1" "StorageClass" "" $storageClass) }}"gp2"{{ else }}"standard"{{ end }}
            `


    // fmap can be extended with any custom functions.
	fmap := funcMap()
    t := template.Must(template.New("test").Funcs(fmap).Parse(manifest))

	err := t.Execute(os.Stdout, values)
	if err != nil {
		fmt.Printf("Error during template execution: %s", err)
		return
	}
}

// Functions from sprig + lookup function from helm
func funcMap() template.FuncMap {
	fmap := sprig.TxtFuncMap()

	config := kubeConfig()

	fmap["lookup"] = NewLookupFunction(config)

	return fmap
}

func kubeConfig() *rest.Config {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	fmt.Println("Using kubeconfig file: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	return config
}