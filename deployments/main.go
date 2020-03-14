package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiclient "traines.eu/time-space-train-planner/client"
	"traines.eu/time-space-train-planner/client/operations"
)

func main() {
	r := httptransport.New(apiclient.DefaultHost, apiclient.DefaultBasePath, apiclient.DefaultSchemes)
	r.DefaultAuthentication = httptransport.BearerToken(os.Getenv("DB_API_ACCESS_TOKEN"))
	r.DefaultMediaType = runtime.XMLMime
	r.Consumers = map[string]runtime.Consumer{
		runtime.XMLMime: runtime.XMLConsumer(),
	}
	r.Producers = map[string]runtime.Producer{
		"application/xhtml+xml": runtime.XMLProducer(),
	}
	client := apiclient.New(r, strfmt.Default)
	var params = operations.NewGetTimetableFchgEvaNoParams()
	params.EvaNo = "8000240"
	res, err := client.Operations.GetTimetableFchgEvaNo(params)

	log.Print(err)
	fmt.Print(*res.Payload.S[0].Tl.C)
}
