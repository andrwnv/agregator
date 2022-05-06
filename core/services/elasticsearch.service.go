package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/estransport"
	"github.com/google/uuid"
	"net/http"
	"os"
	"strings"
)

type EsService struct {
	client    *elasticsearch.Client
	indexName string
}

func NewEsService() *EsService {
	client, _ := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://127.0.0.1:9200",
		},
		Logger: &estransport.ColorLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		},
	})

	_, err := client.Ping()
	if err != nil {
		misc.ReportCritical("Cant connect to elastic search")
	}

	indexName := "aggregator_records_v1"
	status, err := client.Indices.Exists([]string{indexName})
	if err != nil {
		misc.ReportCritical("Broken connect to elastic search")
	}

	if status.StatusCode != http.StatusOK {
		locationMap := fmt.Sprintf(`
		{
			"mappings":{
				"%s":{
					"properties":{
						"location_name":{
							"type":"text"
						},
						"location":{
							"type":"geo_point"
						},
						"location_type":{
							"type":"text"
						}
					}
				}
			}
		}`, indexName)

		_, mappingErr := client.Indices.Create(
			indexName,
			client.Indices.Create.WithBody(strings.NewReader(locationMap)),
		)
		if mappingErr != nil {
			misc.ReportCritical("Cant map elastic search index")
		}
	}

	return &EsService{
		client:    client,
		indexName: indexName,
	}
}

func (es *EsService) Create(createDto dto.CreateAggregatorRecordDto) (err error) {
	if serialized, err := json.Marshal(createDto); err == nil {
		_, createErr := es.client.Create(es.indexName, uuid.New().String(), strings.NewReader(string(serialized)))
		return createErr
	}
	return err
}

func (es *EsService) Delete(id uuid.UUID) error {
	_, err := es.client.Delete(es.indexName, id.String())
	return err
}

func (es *EsService) Update(id uuid.UUID, updateDto dto.UpdateAggregatorRecordDto) (err error) {
	if serialized, err := json.Marshal(updateDto); err == nil {
		_, updateErr := es.client.Update(es.indexName, id.String(), strings.NewReader(string(serialized)))
		return updateErr
	}
	return err
}

func (es *EsService) SearchNearby(userLocation dto.LocationDto, limitSize int) ([]dto.AggregatorRecordDto, error) {
	exp := 0.1

	topLeftLocation := userLocation
	bottomRightLocation := userLocation

	// Extend west/east
	if userLocation.Lat > 0 {
		topLeftLocation.Lat += exp
		bottomRightLocation.Lat -= exp
	} else if userLocation.Lat <= 0 {
		topLeftLocation.Lat -= exp
		bottomRightLocation.Lat += exp
	}

	// Extend North/South
	topLeftLocation.Lon -= exp
	bottomRightLocation.Lon += exp

	geoBound := dto.GeoBounding{
		TopLeft:     topLeftLocation,
		BottomRight: bottomRightLocation,
	}

	query := map[string]interface{}{
		"query": dto.GeoSearch{
			GeoBoundingBox: geoBound,
		},
	}
	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(query); err != nil {
		return []dto.AggregatorRecordDto{}, err
	}

	result, err := es.client.Search(
		es.client.Search.WithSize(limitSize),
		es.client.Search.WithIndex(es.indexName),
		es.client.Search.WithBody(buffer),
		es.client.Search.WithPretty(),
	)
	defer result.Body.Close()
	if err == nil && result != nil {
		fmt.Println(*result)
	}

	return []dto.AggregatorRecordDto{}, nil
}
