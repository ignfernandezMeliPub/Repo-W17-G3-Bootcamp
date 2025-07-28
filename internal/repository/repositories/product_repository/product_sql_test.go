package product_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func setupProductRepository(t *testing.T) (*ProductRepositoryMySQL, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewProductRepositoryMySQL(db)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestProductRepositoryMySQL_GetAllProducts(t *testing.T) {
	repo, mock, cleanup := setupProductRepository(t)
	defer cleanup()

	t.Run("should return all products successfully when database has products", func(t *testing.T) {
		// Arrange
		expectedProducts := []models.Product{
			{
				ID:                             1,
				ProductCode:                    "PROD001",
				Description:                    "Product 1",
				Width:                          10.5,
				Height:                         20.0,
				Length:                         30.0,
				NetWeight:                      5.5,
				ExpirationRate:                 30,
				RecommendedFreezingTemperature: -18.0,
				FreezingRate:                   5,
				ProductTypeId:                  1,
				SellerId:                       nil,
			},
			{
				ID:                             2,
				ProductCode:                    "PROD002",
				Description:                    "Product 2",
				Width:                          15.0,
				Height:                         25.0,
				Length:                         35.0,
				NetWeight:                      7.0,
				ExpirationRate:                 45,
				RecommendedFreezingTemperature: -20.0,
				FreezingRate:                   3,
				ProductTypeId:                  2,
				SellerId:                       nil,
			},
		}

		rows := sqlmock.NewRows([]string{
			"id", "product_code", "description", "width", "height", "length",
			"net_weight", "expiration_rate", "recommended_freezing_temperature",
			"freezing_rate", "product_type_id", "seller_id",
		}).AddRow(
			1, "PROD001", "Product 1", 10.5, 20.0, 30.0, 5.5, 30, -18.0, 5, 1, nil,
		).AddRow(
			2, "PROD002", "Product 2", 15.0, 25.0, 35.0, 7.0, 45, -20.0, 3, 2, nil,
		)

		mock.ExpectQuery(`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id FROM products`).
			WillReturnRows(rows)

		// Act
		products, err := repo.GetAllProducts()

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedProducts, products)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return not found error when database has no products", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{
			"id", "product_code", "description", "width", "height", "length",
			"net_weight", "expiration_rate", "recommended_freezing_temperature",
			"freezing_rate", "product_type_id", "seller_id",
		})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id FROM products`)).
			WillReturnRows(rows)

		// Act
		_, err := repo.GetAllProducts()

		// Assert
		require.Error(t, err)
		require.Equal(t, custom_errors.ErrNotFound, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return connection error when database connection fails", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id FROM products`)).
			WillReturnError(&mysql.MySQLError{
				Number:  2006,
				Message: "MySQL server has gone away",
			})

		// Act
		_, err := repo.GetAllProducts()

		// Assert
		require.Error(t, err)
		require.IsType(t, &mysql.MySQLError{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestProductRepositoryMySQL_GetProductById(t *testing.T) {
	repo, mock, cleanup := setupProductRepository(t)
	defer cleanup()

	t.Run("should return product successfully when product exists with given id", func(t *testing.T) {
		// Arrange
		expectedProduct := models.Product{
			ID:                             1,
			ProductCode:                    "PROD001",
			Description:                    "Product 1",
			Width:                          10.5,
			Height:                         20.0,
			Length:                         30.0,
			NetWeight:                      5.5,
			ExpirationRate:                 30,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   5,
			ProductTypeId:                  1,
			SellerId:                       nil,
		}

		rows := sqlmock.NewRows([]string{
			"id", "product_code", "description", "width", "height", "length",
			"net_weight", "expiration_rate", "recommended_freezing_temperature",
			"freezing_rate", "product_type_id", "seller_id",
		}).AddRow(
			1, "PROD001", "Product 1", 10.5, 20.0, 30.0, 5.5, 30, -18.0, 5, 1, nil,
		)

		mock.ExpectQuery(`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id FROM products WHERE id = ?`).
			WithArgs(1).
			WillReturnRows(rows)

		// Act
		product, err := repo.GetProductById(1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedProduct, product)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return not found error when product does not exist with given id", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{
			"id", "product_code", "description", "width", "height", "length",
			"net_weight", "expiration_rate", "recommended_freezing_temperature",
			"freezing_rate", "product_type_id", "seller_id",
		})

		mock.ExpectQuery(`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id FROM products WHERE id = ?`).
			WithArgs(999).
			WillReturnRows(rows)

		// Act
		_, err := repo.GetProductById(999)

		// Assert
		require.Error(t, err)
		require.Equal(t, custom_errors.ErrNotFound, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return connection error when database connection fails", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id FROM products WHERE id = ?`).
			WithArgs(1).
			WillReturnError(&mysql.MySQLError{
				Number:  2006,
				Message: "MySQL server has gone away",
			})

		// Act
		_, err := repo.GetProductById(1)

		// Assert
		require.Error(t, err)
		require.IsType(t, &mysql.MySQLError{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepositoryMySQL_GetProductByCode(t *testing.T) {
	repo, mock, cleanup := setupProductRepository(t)
	defer cleanup()

	t.Run("should return product successfully when product exists with given code", func(t *testing.T) {
		// Arrange
		expectedProduct := models.Product{
			ID:                             1,
			ProductCode:                    "PROD001",
			Description:                    "Product 1",
			Width:                          10.5,
			Height:                         20.0,
			Length:                         30.0,
			NetWeight:                      5.5,
			ExpirationRate:                 30,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   5,
			ProductTypeId:                  1,
			SellerId:                       nil,
		}

		rows := sqlmock.NewRows([]string{
			"id", "product_code", "description", "width", "height", "length",
			"net_weight", "expiration_rate", "recommended_freezing_temperature",
			"freezing_rate", "product_type_id", "seller_id",
		}).AddRow(
			1, "PROD001", "Product 1", 10.5, 20.0, 30.0, 5.5, 30, -18.0, 5, 1, nil,
		)

		mock.ExpectQuery(`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id FROM products WHERE product_code = ?`).
			WithArgs("PROD001").
			WillReturnRows(rows)

		// Act
		product, err := repo.GetProductByCode("PROD001")

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedProduct, product)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return not found error when product does not exist with given code", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{
			"id", "product_code", "description", "width", "height", "length",
			"net_weight", "expiration_rate", "recommended_freezing_temperature",
			"freezing_rate", "product_type_id", "seller_id",
		})

		mock.ExpectQuery(`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id FROM products WHERE product_code = ?`).
			WithArgs("NONEXISTENT").
			WillReturnRows(rows)

		// Act
		_, err := repo.GetProductByCode("NONEXISTENT")

		// Assert
		require.Error(t, err)
		require.Equal(t, custom_errors.ErrNotFound, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return connection error when database connection fails", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(`SELECT id, product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id FROM products WHERE product_code = ?`).
			WithArgs("PROD001").
			WillReturnError(&mysql.MySQLError{
				Number:  2006,
				Message: "MySQL server has gone away",
			})

		// Act
		_, err := repo.GetProductByCode("PROD001")

		// Assert
		require.Error(t, err)
		require.IsType(t, &mysql.MySQLError{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepositoryMySQL_CreateProduct(t *testing.T) {
	repo, mock, cleanup := setupProductRepository(t)
	defer cleanup()

	t.Run("should create product successfully when all required fields are valid", func(t *testing.T) {
		// Arrange
		product := models.Product{
			ProductCode:                    "PROD001",
			Description:                    "Product 1",
			Width:                          10.5,
			Height:                         20.0,
			Length:                         30.0,
			NetWeight:                      5.5,
			ExpirationRate:                 30,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   5,
			ProductTypeId:                  1,
			SellerId:                       nil,
		}

		expectedProduct := product
		expectedProduct.ID = 1

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO products (product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)).
			WithArgs("PROD001", "Product 1", 10.5, 20.0, 30.0, 5.5, 30, -18.0, 5, 1, nil).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Act
		result, err := repo.CreateProduct(product)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedProduct, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return unique violation error when product code already exists", func(t *testing.T) {
		// Arrange
		product := models.Product{
			ProductCode:                    "EXISTING",
			Description:                    "Product 1",
			Width:                          10.5,
			Height:                         20.0,
			Length:                         30.0,
			NetWeight:                      5.5,
			ExpirationRate:                 30,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   5,
			ProductTypeId:                  1,
			SellerId:                       nil,
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO products (product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)).
			WithArgs("EXISTING", "Product 1", 10.5, 20.0, 30.0, 5.5, 30, -18.0, 5, 1, nil).
			WillReturnError(&mysql.MySQLError{
				Number:  1062,
				Message: "Duplicate entry 'EXISTING' for key 'products.product_code'",
			})

		// Act
		_, err := repo.CreateProduct(product)
		// Assert
		require.Error(t, err)
		require.IsType(t, &custom_errors.UniqueAttributeViolationErr{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return foreign key violation error when product type id does not exist", func(t *testing.T) {
		// Arrange
		product := models.Product{
			ProductCode:                    "PROD001",
			Description:                    "Product 1",
			Width:                          10.5,
			Height:                         20.0,
			Length:                         30.0,
			NetWeight:                      5.5,
			ExpirationRate:                 30,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   5,
			ProductTypeId:                  999,
			SellerId:                       nil,
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO products (product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)).
			WithArgs("PROD001", "Product 1", 10.5, 20.0, 30.0, 5.5, 30, -18.0, 5, 999, nil).
			WillReturnError(&mysql.MySQLError{
				Number:  1452,
				Message: "Cannot add or update a child row: a foreign key constraint fails",
			})

		// Act
		_, err := repo.CreateProduct(product)

		// Assert
		require.Error(t, err)
		require.IsType(t, &custom_errors.ForeignKeyViolationError{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return foreign key violation error when seller id does not exist", func(t *testing.T) {
		// Arrange
		sellerId := 999
		product := models.Product{
			ProductCode:                    "PROD001",
			Description:                    "Product 1",
			Width:                          10.5,
			Height:                         20.0,
			Length:                         30.0,
			NetWeight:                      5.5,
			ExpirationRate:                 30,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   5,
			ProductTypeId:                  1,
			SellerId:                       &sellerId,
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO products (product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)).
			WithArgs("PROD001", "Product 1", 10.5, 20.0, 30.0, 5.5, 30, -18.0, 5, 1, 999).
			WillReturnError(&mysql.MySQLError{
				Number:  1452,
				Message: "Cannot add or update a child row: a foreign key constraint fails",
			})

		// Act
		_, err := repo.CreateProduct(product)

		// Assert
		require.Error(t, err)
		require.IsType(t, &custom_errors.ForeignKeyViolationError{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return connection error when database connection fails", func(t *testing.T) {
		// Arrange
		product := models.Product{
			ProductCode:                    "PROD001",
			Description:                    "Product 1",
			Width:                          10.5,
			Height:                         20.0,
			Length:                         30.0,
			NetWeight:                      5.5,
			ExpirationRate:                 30,
			RecommendedFreezingTemperature: -18.0,
			FreezingRate:                   5,
			ProductTypeId:                  1,
			SellerId:                       nil,
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO products (product_code, description, width, height, length, net_weight, expiration_rate, recommended_freezing_temperature, freezing_rate, product_type_id, seller_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)).
			WithArgs("PROD001", "Product 1", 10.5, 20.0, 30.0, 5.5, 30, -18.0, 5, 1, nil).
			WillReturnError(&mysql.MySQLError{
				Number:  2006,
				Message: "MySQL server has gone away",
			})

		// Act
		_, err := repo.CreateProduct(product)

		// Assert
		require.Error(t, err)
		require.IsType(t, &mysql.MySQLError{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepositoryMySQL_UpdateProductById(t *testing.T) {
	repo, mock, cleanup := setupProductRepository(t)
	defer cleanup()

	t.Run("should update product successfully when product exists and all fields are valid", func(t *testing.T) {
		// Arrange
		product := models.Product{
			ID:                             1,
			ProductCode:                    "PROD001",
			Description:                    "Updated Product",
			Width:                          15.0,
			Height:                         25.0,
			Length:                         35.0,
			NetWeight:                      7.0,
			ExpirationRate:                 45,
			RecommendedFreezingTemperature: -20.0,
			FreezingRate:                   3,
			ProductTypeId:                  2,
			SellerId:                       nil,
		}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET product_code = ?, description = ?, width = ?, height = ?, length = ?, net_weight = ?, expiration_rate = ?, recommended_freezing_temperature = ?, freezing_rate = ?, product_type_id = ?, seller_id = ? WHERE id = ?`)).
			WithArgs("PROD001", "Updated Product", 15.0, 25.0, 35.0, 7.0, 45, -20.0, 3, 2, nil, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Act
		result, err := repo.UpdateProductById(product)

		// Assert
		require.NoError(t, err)
		require.Equal(t, product, result)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return unique violation error when updated product code already exists", func(t *testing.T) {
		// Arrange
		product := models.Product{
			ID:                             1,
			ProductCode:                    "EXISTING_CODE",
			Description:                    "Updated Product",
			Width:                          15.0,
			Height:                         25.0,
			Length:                         35.0,
			NetWeight:                      7.0,
			ExpirationRate:                 45,
			RecommendedFreezingTemperature: -20.0,
			FreezingRate:                   3,
			ProductTypeId:                  2,
			SellerId:                       nil,
		}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET product_code = ?, description = ?, width = ?, height = ?, length = ?, net_weight = ?, expiration_rate = ?, recommended_freezing_temperature = ?, freezing_rate = ?, product_type_id = ?, seller_id = ? WHERE id = ?`)).
			WithArgs("EXISTING_CODE", "Updated Product", 15.0, 25.0, 35.0, 7.0, 45, -20.0, 3, 2, nil, 1).
			WillReturnError(&mysql.MySQLError{
				Number:  1062,
				Message: "Duplicate entry 'EXISTING_CODE' for key 'products.product_code'",
			})

		// Act
		_, err := repo.UpdateProductById(product)

		// Assert
		require.Error(t, err)
		require.IsType(t, &custom_errors.UniqueAttributeViolationErr{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return foreign key violation error when updated product type id does not exist", func(t *testing.T) {
		// Arrange
		product := models.Product{
			ID:                             1,
			ProductCode:                    "PROD001",
			Description:                    "Updated Product",
			Width:                          15.0,
			Height:                         25.0,
			Length:                         35.0,
			NetWeight:                      7.0,
			ExpirationRate:                 45,
			RecommendedFreezingTemperature: -20.0,
			FreezingRate:                   3,
			ProductTypeId:                  999,
			SellerId:                       nil,
		}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET product_code = ?, description = ?, width = ?, height = ?, length = ?, net_weight = ?, expiration_rate = ?, recommended_freezing_temperature = ?, freezing_rate = ?, product_type_id = ?, seller_id = ? WHERE id = ?`)).
			WithArgs("PROD001", "Updated Product", 15.0, 25.0, 35.0, 7.0, 45, -20.0, 3, 999, nil, 1).
			WillReturnError(&mysql.MySQLError{
				Number:  1452,
				Message: "Cannot add or update a child row: a foreign key constraint fails",
			})

		// Act
		_, err := repo.UpdateProductById(product)

		// Assert
		require.Error(t, err)
		require.IsType(t, &custom_errors.ForeignKeyViolationError{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return connection error when database connection fails", func(t *testing.T) {
		// Arrange
		product := models.Product{
			ID:                             1,
			ProductCode:                    "PROD001",
			Description:                    "Updated Product",
			Width:                          15.0,
			Height:                         25.0,
			Length:                         35.0,
			NetWeight:                      7.0,
			ExpirationRate:                 45,
			RecommendedFreezingTemperature: -20.0,
			FreezingRate:                   3,
			ProductTypeId:                  2,
			SellerId:                       nil,
		}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET product_code = ?, description = ?, width = ?, height = ?, length = ?, net_weight = ?, expiration_rate = ?, recommended_freezing_temperature = ?, freezing_rate = ?, product_type_id = ?, seller_id = ? WHERE id = ?`)).
			WithArgs("PROD001", "Updated Product", 15.0, 25.0, 35.0, 7.0, 45, -20.0, 3, 2, nil, 1).
			WillReturnError(&mysql.MySQLError{
				Number:  2006,
				Message: "MySQL server has gone away",
			})

		// Act
		_, err := repo.UpdateProductById(product)

		// Assert
		require.Error(t, err)
		require.IsType(t, &mysql.MySQLError{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepositoryMySQL_DeleteProductById(t *testing.T) {
	repo, mock, cleanup := setupProductRepository(t)
	defer cleanup()

	t.Run("should delete product successfully when product exists with given id", func(t *testing.T) {
		// Arrange
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM products WHERE id = ?`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Act
		err := repo.DeleteProductById(1)

		// Assert
		require.NoError(t, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return foreign key violation error when product has related records", func(t *testing.T) {
		// Arrange
		mock.ExpectExec(`DELETE FROM products WHERE id = ?`).
			WithArgs(1).
			WillReturnError(&mysql.MySQLError{
				Number:  1451,
				Message: "Cannot delete or update a parent row: a foreign key constraint fails",
			})

		// Act
		err := repo.DeleteProductById(1)

		// Assert
		require.Error(t, err)
		require.IsType(t, &custom_errors.ForeignKeyViolationError{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return not found error when product does not exist", func(t *testing.T) {
		// Arrange
		mock.ExpectExec(`DELETE FROM products WHERE id = ?`).
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))
		// Act
		err := repo.DeleteProductById(999)

		// Assert
		require.Error(t, err)
		require.IsType(t, &custom_errors.ResourceNotFoundError{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return connection error when database connection fails", func(t *testing.T) {
		// Arrange
		mock.ExpectExec(`DELETE FROM products WHERE id = ?`).
			WithArgs(1).
			WillReturnError(&mysql.MySQLError{
				Number:  2006,
				Message: "MySQL server has gone away",
			})

		// Act
		err := repo.DeleteProductById(1)

		// Assert
		require.Error(t, err)
		require.IsType(t, &mysql.MySQLError{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestProductRepositoryMySQL_GetReportRecords(t *testing.T) {
	repo, mock, cleanup := setupProductRepository(t)
	defer cleanup()

	t.Run("should return product records report successfully when filtering by product id", func(t *testing.T) {
		// Arrange
		productId := 1
		expectedReports := []models.ProductRecordReport{
			{
				ProductID:    1,
				Description:  "Product 1",
				RecordsCount: 5,
			},
		}

		rows := sqlmock.NewRows([]string{
			"product_id", "description", "records_count",
		}).AddRow(1, "Product 1", 5)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id product_id, p.description, COUNT(r.id) as records_count FROM products p LEFT JOIN product_records r ON r.product_id = p.id WHERE p.id = ? GROUP BY p.id, p.description`)).
			WithArgs(1).
			WillReturnRows(rows)

		// Act
		reports, err := repo.GetReportRecords(&productId)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedReports, reports)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return all product records reports successfully when no product id filter", func(t *testing.T) {
		// Arrange
		expectedReports := []models.ProductRecordReport{
			{
				ProductID:    1,
				Description:  "Product 1",
				RecordsCount: 5,
			},
			{
				ProductID:    2,
				Description:  "Product 2",
				RecordsCount: 3,
			},
		}

		rows := sqlmock.NewRows([]string{
			"product_id", "description", "records_count",
		}).
			AddRow(1, "Product 1", 5).
			AddRow(2, "Product 2", 3)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id product_id, p.description, COUNT(r.id) as records_count FROM products p LEFT JOIN product_records r ON r.product_id = p.id GROUP BY p.id, p.description`)).
			WillReturnRows(rows)

		// Act
		reports, err := repo.GetReportRecords(nil)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedReports, reports)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return not found error when no product records exist", func(t *testing.T) {
		// Arrange
		rows := sqlmock.NewRows([]string{
			"product_id", "description", "records_count",
		})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id product_id, p.description, COUNT(r.id) as records_count FROM products p LEFT JOIN product_records r ON r.product_id = p.id GROUP BY p.id, p.description`)).
			WillReturnRows(rows)

		// Act
		reports, err := repo.GetReportRecords(nil)

		// Assert
		require.Error(t, err)
		require.Equal(t, custom_errors.ErrNotFound, err)
		require.Empty(t, reports)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return connection error when database connection fails", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id product_id, p.description, COUNT(r.id) as records_count FROM products p LEFT JOIN product_records r ON r.product_id = p.id GROUP BY p.id, p.description`)).
			WillReturnError(&mysql.MySQLError{
				Number:  2006,
				Message: "MySQL server has gone away",
			})

		// Act
		_, err := repo.GetReportRecords(nil)

		// Assert
		require.Error(t, err)
		require.IsType(t, &mysql.MySQLError{}, err)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
