package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Cast struct {
	Hash              string   `json:"hash"`
	ParentHash        string   `json:"parentHash"`
	ThreadHash        string   `json:"threadHash"`
	Text              string   `json:"text"`
	Timestamp         string   `json:"timestamp"`
	Author            User     `json:"author"`
	Embeds            []EmbedUrl `json:"embeds"`
	MentionedProfiles []User    `json:"mentionedProfiles"`
}

type User struct {
	Fid              int             `json:"fid"`
	Username         string          `json:"username"`
	DisplayName      string          `json:"display_name"`
	CustodyAddress   string          `json:"custody_address"`
	PfpURL           string          `json:"pfp_url"`
	Profile          Profile         `json:"profile"`
	FollowerCount    int             `json:"follower_count"`
	FollowingCount   int             `json:"following_count"`
	Verifications    []string        `json:"verifications"`
	VerifiedAddresses VerifiedAddresses `json:"verified_addresses"`
	ActiveStatus     string          `json:"active_status"`
	PowerBadge       bool            `json:"power_badge"`
	ViewerContext    ViewerContext   `json:"viewer_context"`
}

type Profile struct {
	Bio Bio `json:"bio"`
}

type Bio struct {
	Text string `json:"text"`
}

type VerifiedAddresses struct {
	EthAddresses []string `json:"eth_addresses"`
	SolAddresses []string `json:"sol_addresses"`
}

type ViewerContext struct {
	Following  bool `json:"following"`
	FollowedBy bool `json:"followed_by"`
}

type EmbedUrl struct {
	URL string `json:"url"`
}

type RecentCastsResponse struct {
	Result struct {
		Casts []Cast `json:"casts"`
		Next  struct {
			Cursor string `json:"cursor"`
		} `json:"next"`
	} `json:"result"`
}

type Fren struct {
	Fid                    int        `json:"fid"`
	Week                   string     `json:"week"`
	Last2222Casts          []Cast     `json:"last2222Casts"`
	Username               string     `json:"username"`
	PfpURL                 string     `json:"pfp_url"`
	Profile                string     `json:"profile"`
	FollowerCount          int        `json:"follower_count"`
	FollowingCount         int        `json:"following_count"`
	EthAddress             string     `json:"eth_address"`
	ArweaveHashForThisWeek string     `json:"arweaveHashForThisWeek"`
	OpenEditionContractAddress string `json:"openEditionContractAddress"`
	DmSent                 bool       `json:"dmSent"`
	ProcessComplete        bool       `json:"processComplete"`
}

func GetRecentCasts(fid int, limit int) ([]Cast, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file")
	}

	apiKey := os.Getenv("NEYNAR_API_KEY")
	url := fmt.Sprintf("%s/farcaster/recent-casts?viewerFid=%d&limit=%d", baseURL, fid, limit)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("api_key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var recentCastsResponse RecentCastsResponse
	err = json.Unmarshal(body, &recentCastsResponse)
	if err != nil {
		return nil, err
	}

	return recentCastsResponse.Result.Casts, nil
}

type FrenInformationResponse struct {
	Username         string          `json:"username"`
	PfpURL           string          `json:"pfp_url"`
	Profile          Profile         `json:"profile"`
	FollowerCount    int             `json:"follower_count"`
	FollowingCount   int             `json:"following_count"`
	VerifiedAddresses VerifiedAddresses `json:"verified_addresses"`
}

func FetchThisUsersInformationFromFarcaster(fid int) (FrenInformationResponse, error) {
	err := godotenv.Load()
	if err != nil {
		return FrenInformationResponse{}, fmt.Errorf("Error loading .env file")
	}

	apiKey := os.Getenv("NEYNAR_API_KEY")
	url := fmt.Sprintf("https://api.neynar.com/v2/farcaster/user/bulk?fids=%d&viewer_fid=16098", fid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return FrenInformationResponse{}, err
	}

	req.Header.Set("api_key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return FrenInformationResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return FrenInformationResponse{}, err
	}

	var frenInformationResponse FrenInformationResponse
	err = json.Unmarshal(body, &frenInformationResponse)
	if err != nil {
		return FrenInformationResponse{}, err
	}

	return frenInformationResponse, nil
}

type DirectCastRequest struct {
	RecipientFid   int    `json:"recipientFid"`
	Message        string `json:"message"`
	IdempotencyKey string `json:"idempotencyKey"`
}

func SendDirectCast(recipientFid int, message string) error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file")
	}

	apiKey := os.Getenv("WARPCAST_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("Warpcast API key is not set in the environment variables")
	}

	idempotencyKey := uuid.New().String()
	apiUrl := "https://api.warpcast.com/v2/ext-send-direct-cast"

	directCastRequest := DirectCastRequest{
		RecipientFid:   recipientFid,
		Message:        message,
		IdempotencyKey: idempotencyKey,
	}

	jsonData, err := json.Marshal(directCastRequest)
	if err != nil {
		return fmt.Errorf("Error marshalling direct cast request: %v", err)
	}

	req, err := http.NewRequest("PUT", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Error creating HTTP request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Warpcast API request failed with status: %s", resp.Status)
	}

	var responseMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		return fmt.Errorf("Error decoding response: %v", err)
	}

	return nil
}