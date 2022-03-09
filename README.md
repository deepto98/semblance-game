[![Semblance](https://github.com/deepto98/semblance-game/blob/main/resources/semblance-logo.png?raw=true)](https://github.com/deepto98/semblance-game)

Semblance is a web based game, where you see an image, and have to describe it in a word.

![Screenshot](https://github.com/deepto98/semblance-game/blob/main/resources/semblance-screenshot.jpg?raw=true)

## Stack 

- Golang to create the server, run the game, use external APIs and for HTML templating
- [Random Image API](https://random.imagecdn.app/) to fetch a random image via a http request 
- Azure Cognitive Services' [Computer Vision API](https://docs.microsoft.com/en-us/azure/cognitive-services/computer-vision/) to generate tags for each image
- For the UI, I used Webflow, much improvement is needed there though 😂 
## Local Installation
1. To run Semblance locally, you must have access to an Azure Computer Vision instance, you can create one from [here](https://portal.azure.com/#create/Microsoft.CognitiveServicesComputerVision)
Once created, get the API Key and Endpoint from **Resource Management > Keys and Endpoint**
2. If Go isn't installed, set it up following [this](https://go.dev/doc/install)
3. Clone this repository, via HTTPS or SSH 
    ``` 
    git clone https://github.com/deepto98/semblance-game.git
4. Open ```.env``` and add your api key and endpoint
    ```
    # Add the Computer Vision key here
    COMPUTER_VISION_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
    # Add the Computer Vision endpoint here
    COMPUTER_VISION_ENDPOINT=https://sampleapp.cognitiveservices.azure.com
    ```
5. Build and run the app
    ```
    go build .
    go run .
    ```
6. If all has gone well, you should be able to play the game at ```localhost```
    
## Acknowledgement
 - Cognitive Services' [Quickstart code](https://github.com/Azure-Samples/cognitive-services-quickstart-code/blob/master/go/ComputerVision/ImageAnalysisQuickstart.go)
## License
- [MIT](https://github.com/deepto98/semblance-game/blob/main/LICENSE)


