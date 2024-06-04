package submit

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

/**
 * @Author: ygzhang
 * @Date: 2024/6/3 21:19
 * @Func:
 **/

func SubmitTransactions() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			fmt.Printf("Failed to populate wallet contents: %s\n", err)
			os.Exit(1)
		}
	}

	ccpPath := filepath.Join(
		"/home",
		"zhangyiguang",
		"fabric",
		"fabric-samples-2.2",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	contract := network.GetContract("fabcar")

	// result, err := contract.EvaluateTransaction("queryAllCars")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))

	// result, err = contract.SubmitTransaction("createCar", "CAR10", "VW", "Polo", "Grey", "Mary")
	// if err != nil {
	// 	fmt.Printf("Failed to submit transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))

	// result, err = contract.EvaluateTransaction("queryCar", "CAR10")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))

	// _, err = contract.SubmitTransaction("changeCarOwner", "CAR10", "Archie")
	// if err != nil {
	// 	fmt.Printf("Failed to submit transaction: %s\n", err)
	// 	os.Exit(1)
	// }

	// result, err = contract.EvaluateTransaction("queryCar", "CAR10")
	// if err != nil {
	// 	fmt.Printf("Failed to evaluate transaction: %s\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(string(result))
	submitTransactions(contract)
}

func populateWallet(wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		"/home",
		"zhangyiguang",
		"fabric",
		"fabric-samples-2.2",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return errors.New("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	err = wallet.Put("appUser", identity)
	if err != nil {
		return err
	}
	return nil
}

func submitTransactions(contract *gateway.Contract) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			car := "CAR" + strconv.Itoa(i)
			owner := "owner" + strconv.Itoa(i)
			_, err := contract.SubmitTransaction("changeCarOwner", car, owner)
			if err != nil {
				fmt.Printf("Failed to submit transaction: %s\n", err)
				os.Exit(1)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
