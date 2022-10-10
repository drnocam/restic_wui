package fileop

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteSettings(data []byte) error {

	file, _ := json.MarshalIndent(string(data), "", " ")

	err := os.WriteFile("test.json", file, 0644)
	fmt.Println("Dosya yazildi")
	return err
}

func ReadSettings() error {

	if data, err := os.ReadFile("test.json"); err == nil {
		fmt.Print(string(data))
	}
	return nil
}
