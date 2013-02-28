package reservoir

import (
	"database/sql"
	"github.com/astaxie/beedb"
	_ "github.com/mattn/go-sqlite3"
)

var DBTYPE string = "sqlite3"
var DBLOCATION string = "database.db"

func init() {
	beedb.OnDebug = true
}

// Retrieve an instance of the DB.
// REMEMBER TO CLOSE THE DB AFTER USING IT.
func RetrieveDB() (*beedb.Model, error) {
	db, err := sql.Open(DBTYPE, DBLOCATION)
	if err != nil {
		return nil, err
	}
	orm := beedb.New(db)
	return &orm, nil
}
