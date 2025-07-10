package database

import (
	"gofr.dev/pkg/gofr/migration"
)

// GetMigrations returns all database migrations for the application
func GetMigrations() map[int64]migration.Migrate {
	return map[int64]migration.Migrate{
		1: createPostsTable(),
	}
}

// createPostsTable defines the migration for creating posts table
func createPostsTable() migration.Migrate {
	return migration.Migrate{
		UP: func(datasource migration.Datasource) error {
			_, err := datasource.SQL.Exec(`
				CREATE TABLE IF NOT EXISTS posts (
					id SERIAL PRIMARY KEY,
					title VARCHAR(200) NOT NULL,
					content TEXT NOT NULL,
					slug VARCHAR(200) NOT NULL UNIQUE,
					author_id INTEGER NOT NULL,
					status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
					created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
				);

				-- Create indexes for better performance
				CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id);
				CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
				CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);
				CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
			`)
			return err
		},
	}
}
