package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
)

var url = "http://127.0.0.1:9200"

func main() {
	fmt.Println("url", url)
	client, err := elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		log.Fatal("err1", err)
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("twitter").Do(context.Background())
	if err != nil {
		// Handle error
		log.Fatal("err2", err)
	}
	if !exists {
		// Create a new index.
		mapping := `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":1
	},
	"mappings":{
		"doc":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
                "retweets":{
                    "type":"long"
                },
				"tags":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}
`
		createIndex, err := client.CreateIndex("twitter").Body(mapping).IncludeTypeName(true).Do(context.Background())
		if err != nil {
			// Handle error
			log.Fatal("err3", err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	_, err = client.Index().Index("twitter").BodyJson(map[string]interface{}{"user": "lin", "message": "hello"}).Do(context.Background())
	fmt.Println("err update", err)

	result, err := client.Search().Index("twitter").Query(elastic.NewTermQuery("user", "lin")).Do(context.Background())
	fmt.Println("err query", err, result.Hits.TotalHits)
	var hitID string
	bytes, _ := json.Marshal(result)
	fmt.Println(string(bytes))
	for _, hit := range result.Hits.Hits {
		fmt.Println("hit", hit.Id)
		hitID = hit.Id
	}

	_, err = client.Index().Index("twitter").Id(hitID).BodyJson(map[string]interface{}{"message": "world"}).Do(context.Background())
	fmt.Println("err update", err)

	ch := make(chan int, 1)
	e := <-ch
	fmt.Printf("%+v\n", e)
}
