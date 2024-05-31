package utils

import (
	"fmt"
	"sync"
)

// ProcessChannelData processes data for a specific channel
func ProcessChannelData(channelID string, interactions []Cast) {
	// Placeholder function to simulate processing
	// In a real implementation, this would generate part of the art piece
	fmt.Printf("Processing data for channel %s with %d interactions\n", channelID, len(interactions))
}

// ProcessUserInteractions concurrently processes user interactions
func ProcessUserInteractions(interactions []Cast) map[string]int {
	channelMap := make(map[string]int)
	for _, interaction := range interactions {
		channelID := interaction.ThreadHash // or another field that represents the channel
		channelMap[channelID]++
	}

	var wg sync.WaitGroup
	for channelID, count := range channelMap {
		wg.Add(1)
		go func(chID string, cnt int) {
			defer wg.Done()
			ProcessChannelData(chID, nil) // update to pass actual data if needed
		}(channelID, count)
	}
	wg.Wait()
	return channelMap
}
