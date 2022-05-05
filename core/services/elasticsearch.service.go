package services

import (
	"github.com/andrwnv/event-aggregator/misc"
	es "github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/estransport"
	"os"
)

type EsService struct {
	Client *es.Client

	locationMap string
}

func NewEsService() *EsService {
	client, err := es.NewClient(es.Config{
		Addresses: []string{
			"http://elasticsearch:9200",
		},
		Logger: &estransport.ColorLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		},
	})

	if err != nil {
		misc.ReportCritical("Cant connect to elastic search")
	}

	return &EsService{
		Client: client,
		locationMap: `
			{
				"mappings":{
					"records":{
						"properties":{
							"store_name":{
								"type":"text"
							},
							"location":{
								"type":"geo_point"
							}
						}
					}
				}
			}`,
	}
}
