package loader

import (
	"app/pkg/models"
)

func NewProductTypeLoaderJSONFile(path string) *ProductTypeJSONFile {
	return &ProductTypeJSONFile{
		path: path,
	}
}

type ProductTypeJSONFile struct {
	path string
}

// Load is a method that loads the Loaders
func (p *ProductTypeJSONFile) Load() (b map[int]models.ProductType, err error) {

	LoadersJSON, err := LoadDataFromFile[models.ProductType](p.path)

	if err != nil {
		return
	}
	// serialize Loaders
	b = make(map[int]models.ProductType)
	for _, productType := range LoadersJSON {
		b[productType.Id] = productType
	}

	return
}
