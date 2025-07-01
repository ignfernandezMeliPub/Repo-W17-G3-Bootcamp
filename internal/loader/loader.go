package loader

import (
	"encoding/json"
	"os"
)

func LoadDataFromFile[T any](filename string) ([]T, error) {
	// ! Obtenemos data del archivo
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// ! Instanciamos T a partir de la data
	var instances []T

	err = json.Unmarshal(data, &instances)
	if err != nil {
		return nil, err
	}

	return instances, nil
}

func SaveDataToFile[T any](filename string, data []T) error {
	// ! Convertimos la data a JSON
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// ! Escribimos el JSON al archivo
	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
