package imageService

import (
	"fmt"
    "io/ioutil"
    "net/http"
	"encoding/json"
	"errors"
)

type ImageService interface {
	GetTag(imageUrl string) (string, error)
}

func NewImageService() ImageService {
	return &ImaggaService{}
}

type ImaggaService struct {
}

func (s *ImaggaService) GetTag(imageUrl string) (string, error) {
    client := &http.Client{}
    api_key := "acc_ccc4acbbccafbb8"
    api_secret := "732210487e462a32ccd08167157d8951"

    req, _ := http.NewRequest("GET", "https://api.imagga.com/v2/tags?image_url="+imageUrl, nil)
    req.SetBasicAuth(api_key, api_secret)

    resp, err := client.Do(req)

    if err != nil {
        fmt.Println("Error when sending request to the server")
        return "", err
    }

    defer resp.Body.Close()
    resp_body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", errors.New("Error getting image tag: " + string(resp_body))
	}
	type Result struct {
		ResultTag map[string]interface{} `json:"result"`
	}

	var res Result
	json.Unmarshal(resp_body, &res)
	tags := res.ResultTag["tags"].([]interface{})
	for _, tag := range tags {
		tagMap := tag.(map[string]interface{})
		fmt.Println(tagMap)
		confidence := tagMap["confidence"].(float64)
		if confidence > 0.5 {
			en := tagMap["tag"].(map[string]interface{})["en"].(string)
			if en == "vehicle"{
				return "vehicle", nil
			}
		}
	}
	return "", errors.New("Image is not clear enough")
}
