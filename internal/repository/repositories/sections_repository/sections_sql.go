package sections_repository

import (
	"app/internal/logger"
	"app/internal/repository/sql_utils"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
	"strconv"
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
	sql_utils.Log("GetAllSections", logger.LogStatusInProgress, "Select all sections")

	secs, err := sql_utils.Query[models.Section](s.db, "SELECT `id`,`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id` FROM sections", nil)
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogError("GetAllSections", "Select all sections", err)
		return nil, err
	}

	sql_utils.Log("GetAllSections", logger.LogStatusSuccess, "Select all sections")
	return secs, nil
}

func (s *SectionsRepositorySQL) GetSectionById(id int) (models.Section, error) {
	sql_utils.Log("GetSectionById", logger.LogStatusInProgress, "Select section by id "+strconv.Itoa(id))

	args := make([]any, 1)
	args[0] = id
	sec, err := sql_utils.QueryRow[models.Section](s.db, "SELECT `id`,`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id` "+
		"FROM sections WHERE id = ?", args)

	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogError("GetSectionById", "Select section by id "+strconv.Itoa(id), err)
		return models.Section{}, err
	}

	sql_utils.Log("GetSectionById", logger.LogStatusSuccess, "Select section by id "+strconv.Itoa(id))
	return sec, nil
}

func (s *SectionsRepositorySQL) CreateSection(sc models.Section) (models.Section, error) {
	sql_utils.LogAudit("CreateSection", logger.LogStatusInProgress, "Insert section")

	args := []any{sc.SectionNumber, sc.CurrentTemperature, sc.MinimumTemperature, sc.CurrentCapacity, sc.MinimumCapacity, sc.MaximumCapacity, sc.WarehouseId, sc.ProductTypeId}
	lastId, err := sql_utils.Insert(s.db,
		"INSERT INTO `sections` (`section_number`,`current_temperature`,`minimum_temperature`,`current_capacity`,`minimum_capacity`,`maximum_capacity`,`warehouse_id`,`product_type_id`) "+
			"VALUES (?,?,?,?,?,?,?,?)", args,
	)
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("CreateSection", "Insert section", err)
		return models.Section{}, err
	}
	sc.ID = int(lastId)

	sql_utils.LogAudit("CreateSection", logger.LogStatusSuccess, "Insert section. Id: "+strconv.Itoa(sc.ID))
	return sc, nil
}

func (s *SectionsRepositorySQL) UpdateSectionById(section models.Section) (models.Section, error) {
	sql_utils.LogAudit("UpdateSectionById", logger.LogStatusInProgress, "Update section by id: "+strconv.Itoa(section.ID))

	args := []any{section.SectionNumber, section.CurrentTemperature, section.MinimumTemperature, section.CurrentCapacity, section.MinimumCapacity, section.MaximumCapacity, section.WarehouseId, section.ProductTypeId, section.ID}

	_, err := sql_utils.Update(s.db, `
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
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("UpdateSectionById", "Update section by id: "+strconv.Itoa(section.ID), err)
		return models.Section{}, err
	}

	sql_utils.LogAudit("UpdateSectionById", logger.LogStatusSuccess, "Update section by id: "+strconv.Itoa(section.ID))
	return section, nil
}

func (s *SectionsRepositorySQL) DeleteSectionById(id int) error {
	sql_utils.LogAudit("DeleteSectionById", logger.LogStatusInProgress, "Delete section by id: "+strconv.Itoa(id))

	args := make([]any, 1)
	args[0] = id
	rowsAffc, err := sql_utils.Delete(s.db, "DELETE FROM sections WHERE id = ?", args)
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("DeleteSectionById", "Delete section by id: "+strconv.Itoa(id), err)
		return err
	}
	if rowsAffc == 0 {
		err = custom_errors.ErrNotFound
		sql_utils.LogAuditError("DeleteSectionById", "Delete section by id: "+strconv.Itoa(id), err)
		return err
	}

	sql_utils.LogAudit("DeleteSectionById", logger.LogStatusSuccess, "Delete section by id: "+strconv.Itoa(id))
	return nil
}

func (s *SectionsRepositorySQL) GetProductBatchBySection(sectionId *int) (prod []models.ProductBatchResponse, err error) {
	if sectionId != nil {
		sql_utils.Log("GetProductBatchBySection", logger.LogStatusInProgress, "Select product batches by section id: "+strconv.Itoa(*sectionId))
	} else {
		sql_utils.Log("GetProductBatchBySection", logger.LogStatusInProgress, "Select product batches by section")
	}

	query := "SELECT `section_id`,`section_number`,SUM(`current_quantity`) `products_count` FROM `product_batches` INNER JOIN `sections` ON product_batches.section_id = sections.id"
	var args []any
	if sectionId != nil {
		query += " WHERE section_id = ?"
		args = append(args, *sectionId)
	}
	query += " GROUP BY section_id"
	prod, err = sql_utils.Query[models.ProductBatchResponse](s.db, query, args)
	err = sql_utils.HandleSqlError(err)
	if err != nil {
		sql_utils.LogError("GetProductBatchBySection", "Select product batches by section", err)
		return prod, err
	}

	if sectionId != nil {
		sql_utils.Log("GetProductBatchBySection", logger.LogStatusSuccess, "Select product batches by section id: "+strconv.Itoa(*sectionId))
	} else {
		sql_utils.Log("GetProductBatchBySection", logger.LogStatusSuccess, "Select product batches by section")
	}
	return prod, nil
}
