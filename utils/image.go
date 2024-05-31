package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func GenerateImageOfThisWeekForFren(fren Fren) (string, error) {
	const width = 800
	const height = 600
	dc := gg.NewContext(width, height)
	dc.SetColor(color.White)
	dc.Clear()
	dc.SetColor(color.Black)
	dc.DrawStringAnchored(fmt.Sprintf("FID: %d", fren.Fid), width/2, height/2, 0.5, 0.5)
	imagePath := fmt.Sprintf("images/fren_%d.png", fren.Fid)

	err := os.MkdirAll("images", os.ModePerm)
	if err != nil {
		return "", err
	}

	err = dc.SavePNG(imagePath)
	if err != nil {
		return "", err
	}

	return imagePath, nil
}

func UploadImageToArweave(imagePath string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("Error loading .env file")
	}

	apiKey := os.Getenv("AKORD_API_KEY")

	tags := []Tag{
		{Name: "file-category", Value: "photo"},
		{Name: "file-description", Value: "FID image"},
	}

	jsonTags, err := json.Marshal(tags)
	if err != nil {
		return "", err
	}

	base64EncodedTags := base64.URLEncoding.EncodeToString(jsonTags)

	fileData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", err
	}

	url := "https://api.akord.com"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(fileData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Tags", base64EncodedTags)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to upload image to Arweave: status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		return "", err
	}

	arweaveHash, ok := responseMap["id"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get Arweave hash from response")
	}

	return fmt.Sprintf("https://arweave.net/%s", arweaveHash), nil
}
