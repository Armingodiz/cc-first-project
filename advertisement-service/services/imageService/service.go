package imageService

import (
	"fmt"
    "io/ioutil"
    "net/http"
	"encoding/json"
)

type ImageService interface {
	GetTag(imageUrl string) (string, error)
}

func NewImageService() MailService {
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
	type Tag struct {
		result struct{
			tags []struct{
				confidence float64
				tag struct{
					en string
				}
			}
		}
	}
	var tag Tag
	json.Unmarshal(resp_body, &tag)
	for _, t := range tag.result.tags {
		if t.tag.en == "Vehicle" && t.confidence > 0.5 {
			return t.tag.en, nil
		}
	}
	return "", errors.New("Image is not clear enough")
}
