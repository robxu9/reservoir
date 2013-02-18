package reservoir

// This was an attempt to build a beans-like serializer,
// which would use reflection to map out a struct to a field.
// Sadly, that's not happening right now, so this is old.

import (
	"bytes"
	"fmt"
	mysql "github.com/ziutek/mymysql/thrsafe"
	"reflect"
	"strings"
)

type Database struct {
	connection *mysql.Conn
	prefix     string
	open       bool
}

func NewDatabase(address string, username string, password string, dbname string, prefix string) (*Database, error) {
	db := new(Database)
	db.connection = mysql.New("tcp", "", address, username, password, dbname)

	err := db.connection.Connect()
	if err != nil {
		return nil, err
	}

	db.connection.Register("set names utf8")

	db.open = true

	if prefix == nil {
		prefix = ""
	}

	db.prefix = prefix

	return db, nil
}

func (d *Database) Disconnect() {
	if db.open {
		d.connection.Close()
	}
}

func (d *Database) Insert(value interface{}) error {
	if value == nil {
		return fmt.Errorf("Cannot pass nil value to Insert!")
	}

	// sending an empty struct is enough for creation
	reflectedValue := reflect.ValueOf(&value).Elem()
	typeOf := reflectedValue.Type()

	// database name is represented by struct name
	tablename := typeOf.Name()

	var buffer bytes.Buffer

	buffer.WriteString("CREATE TABLE IF NOT EXISTS ")
	buffer.WriteString(tablename)
	buffer.WriteString(" (")

	for i := 0; i < reflectedValue.NumField(); i++ {
		field := reflectedValue.Field(i)
		if typeOf.Field(i).Tag == "db" {
			colname := typeOf.Field(i).Name()
			coltype := d.get_mysql_type(field.Type().Kind())

			if coltype == "" {
				return fmt.Errorf("We can't handle exported field %s.", colname)
			}

			buffer.WriteString(colname)
			buffer.WriteString(" ")
			buffer.WriteString(coltype)
			if i != reflectedValue.NumField()-1 {
				buffer.WriteString(", ")
			}
		}
	}

	buffer.WriteString(");")

	query := buffer.String()
	stmt, err := d.connection.Prepare(query)

	if err != nil {
		return err
	}

	_, result, err := stmt.Exec()

	return err
}

func (d *Database) Register(value interface{}) {
	if value == nil {
		panic("Cannot pass nil value to Register!")
	}

	// sending an empty struct is enough for creation
	reflectedValue := reflect.ValueOf(&value).Elem()
	typeOf := reflectedValue.Type()

	// database name is represented by struct name
	tablename := typeOf.Name()

	var buffer bytes.Buffer

	buffer.WriteString("CREATE TABLE IF NOT EXISTS ")
	buffer.WriteString(tablename)
	buffer.WriteString(" (")

	hasID := false

	for i := 0; i < reflectedValue.NumField(); i++ {
		field := reflectedValue.Field(i)
		colname := typeOf.Field(i).Name()

		if typeOf.Field(i).Tag == "field" {
			coltype := d.get_mysql_type(field.Type().Kind())

			if coltype == "" {
				panic(fmt.Sprintf("We can't handle exported field %s.", colname))
			}

			buffer.WriteString(colname)
			buffer.WriteString(" ")
			buffer.WriteString(coltype)
			if i != reflectedValue.NumField()-1 {
				buffer.WriteString(", ")
			}
		} else if strings.HasPrefix(typeOf.Field(i).Tag, "id:") {
			split := strings.Split(typeOf.Field(i).Tag, ":")
			if len(split) != 2 {
				panic("field tagged id not formatted properly to indicate table!")
			}
			table := strings.TrimSpace(split[1])
			if field.Type().Kind() != reflect.Int64 {
				panic("field tagged id must be of int64!")
			}
			if hasID {
				panic(fmt.Sprintf("Can't have two ID fields! Duplicate: %s", colname))
			}
			coltype := "BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY"

			buffer.WriteString(colname)
			buffer.WriteString(" ")
			buffer.WriteString(coltype)
			if i != reflectedValue.NumField()-1 {
				buffer.WriteString(", ")
			}
			hasID = true
		}
	}

	buffer.WriteString(");")

	query := buffer.String()
	stmt, err := d.connection.Prepare(query)

	if err != nil {
		return err
	}

	_, result, err := stmt.Exec()

	return err
}

func (d *Database) get_mysql_type(typeKind reflect.Kind) string {
	switch typeKind {
	// boolean
	case reflect.Bool:
		return "BOOL"
	// numeric types
	case reflect.Int8:
		return "TINYINT"
	case reflect.Int16:
		return "SMALLINT"
	case reflect.Int32:
		return "INT"
	case reflect.Int:
	case reflect.Int64:
		return "BIGINT"
	case reflect.Uint8:
		return "TINYINT UNSIGNED"
	case reflect.Uint16:
		return "SMALLINT UNSIGNED"
	case reflect.Uint32:
		return "INT UNSIGNED"
	case reflect.Uint:
	case reflect.Uint64:
	case reflect.Uintptr:
		return "BIGINT UNSIGNED"
	case reflect.Float32:
	case reflect.Float64:
		return "DOUBLE" // we do NOT use FLOAT because MySQL does everything with double precision
	case reflect.Complex64:
	case reflect.Complex128:
		return "BLOB" // oh god why
	// string types
	case reflect.String:
		return "TEXT"
	default: // can't handle Array, Chan, Func, Interface, Map, Ptr, Slice, Struct, UnsafePointer
		return ""
	}
}
