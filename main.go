package main

import "golang.org/x/oauth2"
import "golang.org/x/oauth2/google"
import compute "google.golang.org/api/compute/v1"

//import "fmt"
import "log"
import "io/ioutil"

func main() {
	// Your credentials should be obtained from the Google
	// Developer Console (https://console.developers.google.com).
	data, err := ioutil.ReadFile("serviceaccount.json")

	conf, err := google.JWTConfigFromJSON(data, "https://www.googleapis.com/auth/compute")
	if err != nil {
		log.Fatal(err)
	}
    // save token? https://github.com/googleapis/google-api-go-client/blob/master/examples/main.go#L125
	client := conf.Client(oauth2.NoContext)

	service, err := compute.New(client)
	if err != nil {
		log.Fatalf("Unable to create Compute service: %v", err)
	}
    // Handy example: https://github.com/hashicorp/go-discover/blob/master/provider/gce/gce_discover.go
	f := func(page *compute.InstanceList) error {
		for _, inst := range page.Items {

			log.Printf("Got compute.Instance, err: %#v, %v", inst.Name, err)
		}
		return nil
	}
    // If we have less than MaxResults, we wouldn't need Pages, but can't tell that
    // without calling, can we?
	//insts, err := service.Instances.List("openshift-gce-devel", "us-east1-c").Do()
	call := service.Instances.List("openshift-gce-devel", "us-east1-c")
    // MaxResults Default is 500
	call.MaxResults(1)
	call.Pages(oauth2.NoContext, f)
}
