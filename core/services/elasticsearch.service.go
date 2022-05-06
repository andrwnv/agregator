package services

import (
	"encoding/json"
	"fmt"
	"github.com/andrwnv/event-aggregator/core/dto"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/estransport"
	"github.com/google/uuid"
	"io"
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

func (es *EsService) SearchNearby(userLocation dto.LocationDto, from int, limitSize int) ([]dto.AggregatorRecordDto, error) {
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

	query := map[string]interface{}{
		"from": from * limitSize,
		"size": limitSize,
		"query": map[string]interface{}{
			"geo_bounding_box": map[string]interface{}{
				"location": dto.GeoBounding{
					TopLeft:     topLeftLocation,
					BottomRight: bottomRightLocation,
				},
			},
		},
	}

	buffer, err := json.Marshal(query)
	if err != nil {
		return []dto.AggregatorRecordDto{}, err
	}

	response, err := es.client.Search(
		es.client.Search.WithIndex(es.indexName),
		es.client.Search.WithBody(strings.NewReader(string(buffer))),
		es.client.Search.WithPretty(),
	)
	defer response.Body.Close()

	if err != nil || response.IsError() {
		return []dto.AggregatorRecordDto{}, err
	}

	// extract [hits][hits]
	jsonResult, _ := readerToMapInterface(response.Body)
	hits := jsonResult["hits"].(map[string]interface{})["hits"]

	// convert hits to dto
	var result []dto.AggregatorRecordDto
	if hitsList, ok := hits.([]interface{}); ok {
		for _, hitInfo := range hitsList {
			if j, marshalErr := json.Marshal(hitInfo.(map[string]interface{})); marshalErr != nil {
				var item dto.AggregatorRecordElasticDto
				unmarshalErr := json.Unmarshal(j, &item)
				if unmarshalErr == nil {
					result = append(result, dto.AggregatorRecordDto{
						ID:           item.ID,
						LocationName: item.AggregatorRecordDto.LocationName,
						Location:     item.AggregatorRecordDto.Location,
						LocationType: item.AggregatorRecordDto.LocationType,
					})
				}
			}
		}
	}

	return result, nil
}

func readerToMapInterface(reader io.Reader) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	err := json.NewDecoder(reader).Decode(&jsonMap)
	return jsonMap, err
}
