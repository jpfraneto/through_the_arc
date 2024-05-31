package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type FidResponse struct {
	Fid int `json:"fid"`
}

func FetchThisTokensAssociatedFid(tokenId int) (int, error) {
	err := godotenv.Load()
	if err != nil {
		return 0, fmt.Errorf("Error loading .env file")
	}

	apiKey := os.Getenv("NEYNAR_API_KEY")
	smartContractAddress := "0x9156d9f4459f92c3c7f7b898d22045b04a6363f6"

	ethereumAddress, err := queryTheSmartContractAndCallTheOwnerOfReadFunctionWithThisTokenIdOnBaseMainnet(smartContractAddress, tokenId)
	if err != nil {
		return 0, err
	}

	body, err := json.Marshal(map[string]string{"address": ethereumAddress})
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", "https://api.neynar.com/v2/farcaster/user/lookup", bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}

	req.Header.Set("api_key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch FID: status code %d", resp.StatusCode)
	}

	var fidResponse FidResponse
	if err := json.NewDecoder(resp.Body).Decode(&fidResponse); err != nil {
		return 0, err
	}

	return fidResponse.Fid, nil
}

func queryTheSmartContractAndCallTheOwnerOfReadFunctionWithThisTokenIdOnBaseMainnet(smartContractAddress string, tokenId int) (string, error) {
	// This function is a placeholder and should query the smart contract on Base Mainnet to get the Ethereum address for the token ID
	return "0xabcdefabcdefabcdefabcdefabcdefabcdef", nil
}
