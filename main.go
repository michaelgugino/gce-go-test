package main

import "golang.org/x/net/context"
import "golang.org/x/oauth2"
import "golang.org/x/oauth2/google"
import compute "google.golang.org/api/compute/v1"

//import "fmt"
import "log"
import "io/ioutil"

func doCall(service *compute.Service) {
	// Handy example: https://github.com/hashicorp/go-discover/blob/master/provider/gce/gce_discover.go
	f := func(page *compute.InstanceList) error {
		for i, inst := range page.Items {

			log.Printf("Got compute.Instance, i: %#v, %v", inst.Name, i)
		}
		return nil
	}
    // If we have less than MaxResults, we wouldn't need Pages, but can't tell that
    // without calling, can we?
	//insts, err := service.Instances.List("openshift-gce-devel", "us-east1-c").Do()
	call := service.Instances.List("openshift-gce-devel", "us-east1-c")
    // MaxResults Default is 500
	call.MaxResults(10)
	call.Pages(oauth2.NoContext, f)
}

func main() {
	// Your credentials should be obtained from the Google
	// Developer Console (https://console.developers.google.com).
	data, err := ioutil.ReadFile("/home/mgugino/clouds/aos-serviceaccount.json")

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
	doCall(service)

	// Reads credentials filepath from env GOOGLE_APPLICATION_CREDENTIALS
	// https://cloud.google.com/docs/authentication/production#obtaining_and_providing_service_account_credentials_manually
	client2, err := google.DefaultClient(context.TODO(), compute.ComputeScope)
	if err != nil {
		log.Fatalf("Unable to create client2: %v", err)
	}
	service2, err := compute.New(client2)
	if err != nil {
		log.Fatalf("Unable to create Compute service: %v", err)
	}
	log.Print("setting up second call")
	doCall(service2)

}
