package models

import (
	"fmt"
	"github.com/luciferCN22/go-gin-example/pkg/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

var db *gorm.DB

type Model struct {
	ID         int            `gorm:"primaryKey" json:"id"`
	CreatedOn  time.Time      `json:"created_on"`
	ModifiedOn time.Time      `json:"modified_on"`
	DeletedOn  gorm.DeletedAt `json:"deleted_on"`
}

func Setup() {
	var (
		err error
		dsn string
	)

	// Construct DSN based on the database type
	switch setting.DatabaseSetting.Type {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			setting.DatabaseSetting.User,
			setting.DatabaseSetting.Password,
			setting.DatabaseSetting.Host,
			setting.DatabaseSetting.Name)
	default:
		log.Fatalf("Database type not supported: %s", setting.DatabaseSetting.Type)
	}

	// Open the database connection
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set table prefix
	db.NamingStrategy = schema.NamingStrategy{
		TablePrefix:   setting.DatabaseSetting.TablePrefix, // table name prefix, table for `User` would be `t_users`
		SingularTable: true,                                // use singular table name, table for `User` would be `user` with this option enabled
	}

	// Additional database settings can be set here
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	db.Callback().Create().Replace("gorm:before_create", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:before_update", updateTimeStampForUpdateCallback)
	//db.Callback().Delete().Replace("gorm:delete", deleteCallback)
}

func CloseDB() {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
	}
	defer sqlDB.Close()
}

func updateTimeStampForCreateCallback(db *gorm.DB) {
	if db.Error == nil {
		nowTime := time.Now()
		// Get the scope of the current operation
		if _, ok := db.Statement.Schema.FieldsByName["CreatedOn"]; ok {
			// If the field exists, set the value of the field to the current time
			db.Statement.SetColumn("CreatedOn", nowTime)
		}
		if _, ok := db.Statement.Schema.FieldsByName["ModifiedOn"]; ok {
			// If the field exists, set the value of the field to the current time
			db.Statement.SetColumn("ModifiedOn", nowTime)
		}
	}
}

func updateTimeStampForUpdateCallback(db *gorm.DB) {
	if _, ok := db.Statement.Schema.FieldsByName["ModifiedOn"]; ok {
		db.Statement.SetColumn("ModifiedOn", time.Now())
	}
}

//func deleteCallback(db *gorm.DB) {
//	if db.Error == nil {
//		var extraOption string
//		if str, ok := db.Get("gorm:delete_option"); ok {
//			extraOption = fmt.Sprint(str)
//		}
//
//		// Check if DeletedOn field exists
//		deletedOnField, hasDeletedOnField := db.Statement.Schema.FieldsByName["DeletedOn"]
//		if hasDeletedOnField {
//			// Update DeletedOn field with the current timestamp
//			logging.Info(db.Statement.SQL.String())
//			db.Exec(fmt.Sprintf(
//				"UPDATE %v SET %v=%v%v%v",
//				db.Statement.Table,
//				deletedOnField.DBName,
//				time.Now().Unix(),
//				addExtraSpaceIfExist(db.Statement.SQL.String()),
//				addExtraSpaceIfExist(extraOption),
//			))
//		} else {
//			// Perform a regular DELETE operation
//			db.Exec(fmt.Sprintf(
//				"DELETE FROM %v%v%v",
//				db.Statement.Table,
//				addExtraSpaceIfExist(db.Statement.SQL.String()),
//				addExtraSpaceIfExist(extraOption),
//			))
//		}
//	}
//}
//
//// addExtraSpaceIfExist adds a separator
//func addExtraSpaceIfExist(str string) string {
//	if str != "" {
//		return " " + str
//	}
//	return ""
//}
