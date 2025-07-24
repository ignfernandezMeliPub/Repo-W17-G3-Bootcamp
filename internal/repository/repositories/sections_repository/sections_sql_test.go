package sections_repository

import (
	"database/sql"
	"log"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	txdb.Register("txdb", "mysql", "root@/txdb_test")
}

func main() {
	db, err := sql.Open("txdb", "identifier")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Your code
}
