package loader

import "app/pkg/models"

func NewWarehouseJSONFile(path string) *WarehouseJSONFile {
	return &WarehouseJSONFile{
		path: path,
	}
}

type WarehouseJSONFile struct {
	path string
}

func (l *WarehouseJSONFile) Load() (map[int]models.Warehouse, error) {
	wh, err := LoadDataFromFile[models.Warehouse](l.path)
	if err != nil {
		return nil, err
	}

	w := make(map[int]models.Warehouse)
	for _, warehouse := range wh {
		w[warehouse.Id] = warehouse
	}

	return w, nil
}
