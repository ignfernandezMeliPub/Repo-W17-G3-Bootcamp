package loader

import (
	"app/pkg/models"
)

func NewBuyerLoaderJSONFile(path string) *LoaderJSONFile {
	return &LoaderJSONFile{
		path: path,
	}
}

type LoaderJSONFile struct {
	path string
}

// Load is a method that loads the Loaders
func (l *LoaderJSONFile) Load() (b map[int]models.Buyer, err error) {

	LoadersJSON, err := LoadDataFromFile[models.Buyer](l.path)

	if err != nil {
		return
	}
	// serialize Loaders
	b = make(map[int]models.Buyer)
	for _, buyer := range LoadersJSON {
		b[buyer.Id] = buyer
	}

	return
}
