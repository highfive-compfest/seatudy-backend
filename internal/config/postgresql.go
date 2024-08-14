package config

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresql(migrations ...any) *gorm.DB {
	gormLogger := logger.Default
	if Env.ENV != "production" {
		gormLogger = gormLogger.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			Env.PostgresHost,
			Env.PostgresUser,
			Env.PostgresPassword,
			Env.PostgresDbName,
			Env.PostgresPort,
		),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: gormLogger,
	})
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

	if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE midtrans_status AS ENUM (
				'challenge',
				'success',
				'failure',
				'pending'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
        DO $$ BEGIN
            CREATE TYPE course_category AS ENUM (
                'Web Development',
                'Game Development',
                'Cloud Computing',
                'Data Science & Analytics',
                'Programming Languages',
                'Cybersecurity',
                'Mobile App Development',
                'Database Management',
                'Software Development',
                'DevOps & Automation',
                'Networking',
                'AI & Machine Learning',
                'Internet of Things (IoT)',
                'Blockchain & Cryptocurrency',
                'Augmented Reality (AR) & Virtual Reality (VR)'
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
