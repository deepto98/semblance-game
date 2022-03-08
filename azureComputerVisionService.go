package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.0/computervision"
	"github.com/Azure/go-autorest/autorest"
	"github.com/joho/godotenv"
)

var computerVisionContext context.Context

func CreateComputerVisionClient() computervision.BaseClient {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	computerVisionKey := os.Getenv("COMPUTER_VISION_KEY")

	endpointURL := os.Getenv("COMPUTER_VISION_ENDPOINT")

	computerVisionClient := computervision.New(endpointURL)
	computerVisionClient.Authorizer = autorest.NewCognitiveServicesAuthorizer(computerVisionKey)

	computerVisionContext = context.Background()
	return computerVisionClient

}

func TagRemoteImage(client computervision.BaseClient, remoteImageURL string) map[string]bool {

	var remoteImage computervision.ImageURL
	remoteImage.URL = &remoteImageURL

	remoteImageTags, err := client.TagImage(
		computerVisionContext,
		remoteImage,
		"")
	if err != nil {
		log.Fatal(err)
	}

	var tags []string
	mymap := make(map[string]bool)

	if len(*remoteImageTags.Tags) == 0 {
		// fmt.Println("No tags detected.")
	} else {
		for _, tag := range *remoteImageTags.Tags {
			tags = append(tags, *tag.Name)
			mymap[*tag.Name] = true

			fmt.Printf("'%v' with confidence %.2f%%\n", *tag.Name, *tag.Confidence*100)
		}
	}
	return mymap
}
