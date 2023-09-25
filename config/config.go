package config

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"knaq-wallet/controller/erra"
	"knaq-wallet/tools/converter"
	"log"
	"os"
)

var Config config

func New() {
	Config = config{
		port: os.Getenv(envPort),
		database: []database{
			{
				endpoint:     os.Getenv(envDatabaseEndpoint),
				kind:         getKind(),
				logLevel:     getLogLevel(),
				maxIdleConns: converter.StringToInt(os.Getenv(envDatabaseMaxIdleConns)),
				maxOpenConns: converter.StringToInt(os.Getenv(envDatabaseMaxOpenConns)),
			},
		},
		stage:                Stage(os.Getenv(envStage)),
		jwtSecret:            os.Getenv(envJwtSecret),
		nearEnv:              os.Getenv(envNearEnv),
		nearID:               os.Getenv(envNearID),
		nearPrivateKey:       os.Getenv(envNearPrivateKey),
		nearFtContract:       os.Getenv(envNearFtContract),
		nearFtStorageDeposit: os.Getenv(envNearFtStorageDeposit),
	}

	if err := GenerateNearMasterAccountFile(); err != nil {
		log.Fatal(err)
	}
}
func getKind() int {
	switch os.Getenv(envDatabaseKind) {
	case "mysql":
		return 1
	case "postgresql":
		return 2
	default:
		return 0
	}
}
func getLogLevel() int {
	switch os.Getenv(envDatabaseLogLevel) {
	case "silent":
		return 1
	case "error":
		return 2
	case "warn":
		return 3
	case "info":
		return 4
	default:
		return 1
	}
}
func decodeHexString(encodedString string) []byte {
	decodedString, err := hex.DecodeString(encodedString)
	if err != nil {
		log.Fatal(err)
	}
	return decodedString
}
func GenerateNearMasterAccountFile() error {
	var (
		err                error
		homeDir            string
		filePath           string
		nearMasterJsonByte []byte
	)
	homeDir, err = os.UserHomeDir()
	if err != nil {
		return erra.Error(err)
	}
	filePath = fmt.Sprintf("%s/.near-credentials/%s", homeDir, Config.GetNearEnv())
	if err = os.MkdirAll(filePath, 0755); err != nil {
		return erra.Error(err)
	}
	nearMasterJsonByte, err = json.Marshal(map[string]interface{}{
		"private_key": Config.GetNearPrivateKey(),
	})
	if err != nil {
		return erra.Error(err)
	}
	if err = os.WriteFile(
		fmt.Sprintf("%s/%s.json", filePath, Config.GetNearID()),
		nearMasterJsonByte,
		0755,
	); err != nil {
		return erra.Error(err)
	}

	return nil
}
