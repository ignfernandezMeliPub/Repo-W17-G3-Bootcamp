package product_batch_repository

import (
	"app/internal/repository/sql_utils"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
)

type ProductBatchRepositorySQL struct {
	db *sql.DB
}

func NewProductBatchRepositorySQL(db *sql.DB) *ProductBatchRepositorySQL {
	return &ProductBatchRepositorySQL{db}
}

func (p *ProductBatchRepositorySQL) CreateProductBatch(pb models.ProductBatch) (prod models.ProductBatch, err error) {
	args := []any{pb.BatchNumber, pb.CurrentQuantity, pb.CurrentTemperature, pb.DueDate, pb.InitialQuantity, pb.ManufacturingDate, pb.ManufacturingHour, pb.MinimumTemperature, pb.ProductId, pb.SectionId}
	lastId, err := sql_utils.Insert(p.db,
		"INSERT INTO `product_batches` (`batch_number`,`current_quantity`,`current_temperature`,`due_date`,`initial_quantity`,`manufacturing_date`,`manufacturing_hour`,`minimum_temperature`,`product_id`,`section_id`) "+
			"VALUES (?,?,?,?,?,?,?,?)", args,
	)
	if err != nil {
		return
	}
	pb.ID = int(lastId)
	return pb, err
}

func (p *ProductBatchRepositorySQL) GetAllProductBatchesBySection() (prods []models.ProductBatchResponse, err error) {
	prods, err = sql_utils.Query[models.ProductBatchResponse](p.db, "SELECT `section_id`,`section_number`,SUM(`current_quantity`) `products_count` "+
		"FROM `product_batches` INNER JOIN `sections` ON product_batches.section_id = sections.id"+
		"GROUP BY `section_id`", nil)
	if err != nil {
		return
	}
	if len(prods) == 0 {
		err = &custom_errors.ResourceNotFoundError{}
		return
	}
	return
}

func (p *ProductBatchRepositorySQL) GetProductBatchBySectionId(sectionId int) (prod models.ProductBatchResponse, err error) {
	args := make([]any, 1)
	args[0] = sectionId
	prod, err = sql_utils.QueryRow[models.ProductBatchResponse](p.db, "SELECT `section_id`,`section_number`,SUM(`current_quantity`) `products_count` "+
		"FROM `product_batches` INNER JOIN `sections` ON product_batches.section_id = sections.id"+
		"WHERE section_id = ? GROUP BY `section_id`", args)
	return
}
