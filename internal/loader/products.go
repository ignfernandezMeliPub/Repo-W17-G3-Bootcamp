package loader

import (
	"app/pkg/models"
)

func NewProductLoaderJSONFile(path string) *ProductJSONFile {
	return &ProductJSONFile{
		path: path,
	}
}

type ProductJSONFile struct {
	path string
}

// Load is a method that loads the Loaders
func (p *ProductJSONFile) Load() (b map[int]models.Product, err error) {

	LoadersJSON, err := LoadDataFromFile[models.Product](p.path)

	if err != nil {
		return
	}
	// serialize Loaders
	b = make(map[int]models.Product)
	for _, product := range LoadersJSON {
		b[product.ID] = product
	}

	return
}
