// go:build ignore

package main

import (
	"log"
	"strings"

	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/rawsql"
)

func main() {
	db, err := gorm.Open(rawsql.New(rawsql.Config{FilePath: []string{"../po/sql/table.sql"}}))
	if err != nil {
		log.Fatal(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           "../dao",
		OutFile:           "",
		ModelPkgPath:      "po",
		WithUnitTest:      false,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.WithQueryInterface,
	})
	g.UseDB(db)

	g.WithDataTypeMap(map[string]func(gorm.ColumnType) (dataType string){
		// bool mapping
		"tinyint": func(columnType gorm.ColumnType) (dataType string) {
			ct, _ := columnType.ColumnType()
			if strings.HasPrefix(ct, "tinyint(1)") {
				return "bool"
			}
			return "byte"
		},
	})

	tables := g.GenerateAllTable()
	g.ApplyBasic(tables...)
	g.ApplyInterface(func(IBase) {}, tables...)

	g.Execute()
}

type IBase interface {
	// SELECT * FROM @@table WHERE id=@id
	FindByID(id int64) (*gen.T, error)

	// DELETE FROM @@table WHERE id=@id
	DeleteByID(id int64) error
}
