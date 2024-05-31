package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"through_the_arc/utils"
)

func main_weekOne() {
	var wg sync.WaitGroup
	var mu sync.Mutex

	subscribersOfThroughTheArcFids := make([]int, 0)

	for i := 1; i < 162; i++ { // Adjust the range as needed
		wg.Add(1)
		go func(tokenId int) {
			defer wg.Done()
			thisSubscribersFid, err := utils.FetchThisTokensAssociatedFid(tokenId)
			if err != nil {
				log.Printf("Error fetching FID for token %d: %v", tokenId, err)
				return
			}
			mu.Lock()
			subscribersOfThroughTheArcFids = append(subscribersOfThroughTheArcFids, thisSubscribersFid)
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	for _, thisSubscribersFid := range subscribersOfThroughTheArcFids {
		wg.Add(1)
		go func(fid int) {
			defer wg.Done()

			last2222casts, err := utils.GetRecentCasts(fid, 2222)
			if err != nil {
				log.Printf("Error fetching casts for FID %d: %v", fid, err)
				return
			}

			frenInformation, err := utils.FetchThisUsersInformationFromFarcaster(fid)
			if err != nil {
				log.Printf("Error fetching user information for FID %d: %v", fid, err)
				return
			}

			fren := utils.Fren{
				Fid:            fid,
				Week:           "one",
				Last2222Casts:  last2222casts,
				Username:       frenInformation.Username,
				PfpURL:         frenInformation.PfpURL,
				Profile:        frenInformation.Profile.Bio.Text,
				FollowerCount:  frenInformation.FollowerCount,
				FollowingCount: frenInformation.FollowingCount,
				EthAddress:     frenInformation.VerifiedAddresses.EthAddresses[0],
			}

			imagePath, err := utils.GenerateImageOfThisWeekForFren(fren)
			if err != nil {
				log.Printf("Error generating image for FID %d: %v", fid, err)
				return
			}

			arweaveLink, err := utils.UploadImageToArweave(imagePath)
			if err != nil {
				log.Printf("Error uploading image to Arweave for FID %d: %v", fid, err)
				return
			}
			fren.ArweaveHashForThisWeek = arweaveLink

			openEditionContractAddress, err := utils.CreateCollectionFromArweaveLink(arweaveLink)
			if err != nil {
				log.Printf("Error creating open edition NFT for FID %d: %v", fid, err)
				return
			}

			message := fmt.Sprintf("Congratulations! Your edition has been minted and is available at %s", openEditionContractAddress)
			err = utils.SendDirectCast(fid, message)
			if err != nil {
				log.Printf("Error sending direct cast to FID %d: %v", fid, err)
				return
			}

			dmSent := true

			fren.OpenEditionContractAddress = openEditionContractAddress
			fren.DmSent = dmSent
			fren.ProcessComplete = true

			filename := filepath.Join("subscriber_data", "week_one", fmt.Sprintf("%d.json", fid))
			err = utils.StoreData(filename, fren)
			if err != nil {
				log.Printf("Error storing data for FID %d: %v", fid, err)
				return
			}

			fmt.Printf("Data stored for FID %d in %s\n", fid, filename)
		}(thisSubscribersFid)
	}

	wg.Wait()
}
