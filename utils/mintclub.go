package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type EditionParameters struct {
	MaxSupply  int    `json:"maxSupply"`
	Price      int    `json:"price"`
	Fid        int    `json:"fid"`
	Currency   string `json:"currency"`
	ArweaveLink string `json:"arweaveLink"`
}

type MintClubResponse struct {
	ContractAddress string `json:"contractAddress"`
}

func createCollectionFromArweaveLink(arweaveLink string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("Error loading .env file")
	}

	apiKey := os.Getenv("MINT_CLUB_API_KEY")

	editionParams := EditionParameters{
		MaxSupply:  8,
		Price:      0,
		Currency:   "ETH",
		ArweaveLink: "https://arweave.net/rkLlX6b5Pf9syi7FDK6cMOglclCPQdjqF_Ft_OhQnFc",
	}

	jsonData, err := json.Marshal(editionParams)
	if err != nil {
		return "", fmt.Errorf("Error marshalling edition parameters: %v", err)
	}

	url := "https://api.mint.club/v2/nft/create"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("Error creating HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Mint Club API request failed with status: %s", resp.Status)
	}

	var mintClubResponse MintClubResponse
	err = json.NewDecoder(resp.Body).Decode(&mintClubResponse)
	if err != nil {
		return "", fmt.Errorf("Error decoding response: %v", err)
	}

	return mintClubResponse.ContractAddress, nil
}
