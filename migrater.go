package migrate

import (
	"errors"
	"reflect"

	"github.com/gookit/color"

	"github.com/gogf/gf/util/gconv"

	"gorm.io/gorm"
)

type Migration struct {
	Id        uint   `gorm:"primaryKey;column:ID"`
	Migration string `gorm:"size:255;not null;uniqueIndex;column:MIGRATION"`
	Batch     uint   `gorm:"not null;column:BATCH"`
}

func initMigration(db *gorm.DB) error {

	if !db.Migrator().HasTable(&Migration{}) {
		return db.Migrator().CreateTable(&Migration{})
	}

	return nil
}

func hasMigrated(db *gorm.DB) []string {
	ms := make([]*Migration, 0)
	db.Select("MIGRATION").Find(&ms)
	mString := make([]string, 0)
	for _, v := range ms {
		mString = append(mString, v.Migration)
	}
	return mString
}

func shouldMigrate(hasMigrated []string, migrations []interface{}) []interface{} {
	should := make([]interface{}, 0)
	for _, v := range migrations {
		name := gconv.String(reflect.ValueOf(v).Elem().FieldByName("Name"))
		in := false
		for _, h := range hasMigrated {
			if h == name {
				in = true
				break
			}
		}
		if in == false {
			should = append(should, v)
		}
	}
	return should
}

func getBatch(db *gorm.DB) uint {
	m := Migration{}
	db.Order("batch desc").First(&m)
	return m.Batch + 1
}

func Up(db *gorm.DB, migrations []interface{}) error {
	if err := initMigration(db); err != nil {
		return err
	}

	sm := shouldMigrate(hasMigrated(db), migrations)

	if len(sm) == 0 {
		return errors.New("nothing to migrate")
	}

	batch := getBatch(db)
	for _, v := range sm {
		name := gconv.String(reflect.ValueOf(v).Elem().FieldByName("Name"))
		color.Green.Println("migrating " + name)
		rv := reflect.ValueOf(v)
		up := rv.MethodByName("Up")
		result := up.Call([]reflect.Value{})
		wrong := false
		if len(result) == 0 {
			wrong = true
		} else {
			wrong = !result[0].IsNil()
		}
		if wrong {
			color.Red.Println("please check migration " + name)
			return errors.New("migration error")
		}
		db.Create(&Migration{
			Migration: name,
			Batch:     batch,
		})
		color.Green.Println("migrated " + name + " success")
	}
	return nil
}

func Down(db *gorm.DB, migrations []interface{}) error {
	if err := initMigration(db); err != nil {
		return err
	}

	m := Migration{}
	db.Order("batch desc").Last(&m)
	if m.Batch == 0 {
		return errors.New("nothing to rollback")
	}

	ms := make([]*Migration, 0)
	db.Where(&Migration{Batch: m.Batch}).Order("id desc").Find(&ms)
	type RollBack struct {
		MigrationTable *Migration
		Migration      interface{}
	}
	should := make([]RollBack, 0)
	for _, v := range ms {
		for _, vv := range migrations {
			if gconv.String(reflect.ValueOf(vv).Elem().FieldByName("Name")) == v.Migration {
				should = append(should, RollBack{
					MigrationTable: v,
					Migration:      vv,
				})
				break
			}
		}
	}
	if len(should) == 0 {
		return errors.New("migrate file not exist")
	}

	for _, v := range should {
		name := gconv.String(reflect.ValueOf(v.Migration).Elem().FieldByName("Name"))
		color.Green.Println("starting rollback " + name)
		rv := reflect.ValueOf(v.Migration)
		up := rv.MethodByName("Down")
		result := up.Call([]reflect.Value{})
		wrong := false
		if len(result) == 0 {
			wrong = true
		} else {
			wrong = !result[0].IsNil()
		}
		if wrong {
			color.Red.Println("please check migration " + name)
			return errors.New("rollback error")
		}
		db.Delete(&Migration{
			Id: v.MigrationTable.Id,
		})
		color.Green.Println("rollback " + name + " success")
	}

	return nil
}
