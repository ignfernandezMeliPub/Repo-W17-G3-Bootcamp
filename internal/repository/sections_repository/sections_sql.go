package sections_repository

import (
	"app/internal/repository/sql_utils"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
)

type SectionsRepositorySQL struct {
	db *sql.DB
}

func NewSectionsRepositorySQL(db *sql.DB) *SectionsRepositorySQL {
	return &SectionsRepositorySQL{
		db,
	}
}

func (s *SectionsRepositorySQL) GetAllSections() ([]models.Section, error) {
	secs, err := sql_utils.Query[models.Section](s.db, "SELECT `id`,`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id` FROM sections", nil)
	if err != nil {
		return nil, err
	}
	if len(secs) == 0 {
		err = &custom_errors.ResourceNotFoundError{}
		return nil, err
	}

	return secs, nil
}

func (s *SectionsRepositorySQL) GetSectionById(id int) (models.Section, error) {
	args := make([]any, 1)
	args[0] = id
	sec, err := sql_utils.QueryRow[models.Section](s.db, "SELECT `id`,`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id` "+
		"FROM sections WHERE id = ?", args)

	if err != nil {
		return models.Section{}, err
	}
	return sec, nil
}

func (s *SectionsRepositorySQL) CreateSection(sc models.Section) (models.Section, error) {
	args := []any{sc.SectionNumber, sc.CurrentTemperature, sc.MinimumTemperature, sc.CurrentCapacity, sc.MinimumCapacity, sc.MaximumCapacity, sc.WarehouseId, sc.ProductTypeId}
	lastId, err := sql_utils.Insert(s.db,
		"INSERT INTO `sections` (`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id`) "+
			"VALUES (?,?,?,?,?,?,?,?)", args,
	)
	if err != nil {
		return models.Section{}, err
	}
	sc.ID = int(lastId)
	return sc, nil
}

func (s *SectionsRepositorySQL) UpdateSectionById(section models.Section) (models.Section, error) {
	args := []any{section.SectionNumber, section.CurrentTemperature, section.MinimumTemperature, section.CurrentCapacity, section.MinimumCapacity, section.MaximumCapacity, section.WarehouseId, section.ProductTypeId, section.ID}

	rowsAffc, err := sql_utils.Update(s.db, `
		UPDATE sections SET
			section_number = ?,
			current_temperature = ?,
			minimum_temperature = ?,
			current_capacity = ?,
			minimum_capacity = ?,
			maximum_capacity = ?,
			warehouse_id = ?,
			product_type_id = ?
		WHERE id = ?`,
		args,
	)
	if err != nil {

		return models.Section{}, err

	}

	if rowsAffc == 0 {

		return models.Section{}, &custom_errors.ResourceNotFoundError{}

	}

	return section, nil
}

func (s *SectionsRepositorySQL) DeleteSectionById(id int) error {
	args := make([]any, 1)
	args[0] = id
	rowsAffc, err := sql_utils.Delete(s.db, "DELETE FROM sections WHERE id = ?", args)
	if rowsAffc == 0 {
		return &custom_errors.ResourceNotFoundError{}
	}
	return err
}

func (s *SectionsRepositorySQL) GetAllProductBatchesBySection() (prods []models.ProductBatchResponse, err error) {
	prods, err = sql_utils.Query[models.ProductBatchResponse](s.db, "SELECT `section_id`,`section_number`,SUM(`current_quantity`) `products_count` FROM `product_batches` INNER JOIN `sections` ON product_batches.section_id = sections.id GROUP BY section_id", nil)
	if err != nil {
		return
	}
	if len(prods) == 0 {
		err = &custom_errors.ResourceNotFoundError{}
		return
	}
	return
}
func (s *SectionsRepositorySQL) GetProductBatchBySectionId(sectionId int) (prod models.ProductBatchResponse, err error) {
	args := make([]any, 1)
	args[0] = sectionId
	prod, err = sql_utils.QueryRow[models.ProductBatchResponse](s.db, "SELECT `section_id`,`section_number`,SUM(`current_quantity`) `products_count` FROM `product_batches` INNER JOIN `sections` ON product_batches.section_id = sections.id WHERE section_id = ? GROUP BY `section_id`", args)
	return
}
