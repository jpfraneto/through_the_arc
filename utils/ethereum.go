package utils

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

const contractABI = `[{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"}]`

func queryTheSmartContractAndCallTheOwnerOfReadFunctionWithThisTokenIdOnBaseMainnet(smartContractAddress string, tokenId int) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("Error loading .env file")
	}

	alchemyRPC := os.Getenv("ALCHEMY_HTTPS_RPC")
	client, err := ethclient.Dial(alchemyRPC)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	contractAddress := common.HexToAddress(smartContractAddress)
	instance, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: instance.Methods["ownerOf"].ID,
	}

	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to call contract: %v", err)
	}

	var ownerAddress common.Address
	err = instance.Unpack(&ownerAddress, "ownerOf", result)
	if err != nil {
		return "", fmt.Errorf("Failed to unpack result: %v", err)
	}

	return ownerAddress.Hex(), nil
}
