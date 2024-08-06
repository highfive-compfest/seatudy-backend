package config

import (
	"log"


	"github.com/highfive-compfest/seatudy-backend/internal/schema"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresql() *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
        DSN:                  Env.DbDsn,
        PreferSimpleProtocol: true, // disables implicit prepared statement usage
    }), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

	if err := migratePostgresqlTables(db, migrations...); err != nil {
		log.Fatalln(err)
	}

	return db
}

func migratePostgresqlTables(db *gorm.DB, migrations ...any) error {
	if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE user_role AS ENUM (
				'student',
				'instructor'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE course_difficulty AS ENUM (
				'beginner',
				'intermediate',
				'advanced',
				'expert'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE post_type AS ENUM (
				'material',
				'assignment'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`).Error; err != nil {
		return err
	}

	if err := db.AutoMigrate(
		migrations..., // BREAKING: entities should be passed from cmd/api/main.go due to circular dependency issue
	); err != nil {
		return err
	}

	return nil
}
