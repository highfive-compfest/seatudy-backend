package config

import (
	"github.com/highfive-compfest/seatudy-backend/internal/domain/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func NewPostgresql() *gorm.DB {
	db, err := gorm.Open(postgres.Open(Env.DbDsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	if err := migratePostgresqlTables(db); err != nil {
		log.Fatalln(err)
	}

	return db
}

func migratePostgresqlTables(db *gorm.DB) error {
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
		&user.User{},
	); err != nil {
		return err
	}

	return nil
}
