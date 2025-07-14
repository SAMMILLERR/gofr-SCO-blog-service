package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

func create_authors_table() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			// Create the authors table
			_, err := d.SQL.Exec(`
				CREATE TABLE IF NOT EXISTS authors (
					id SERIAL PRIMARY KEY,
					username VARCHAR(50) UNIQUE NOT NULL,
					email VARCHAR(100) UNIQUE NOT NULL,
					password_hash VARCHAR(255) NOT NULL,
					first_name VARCHAR(50) NOT NULL,
					last_name VARCHAR(50) NOT NULL,
					bio TEXT,
					avatar_url VARCHAR(255),
					is_active BOOLEAN DEFAULT true,
					is_verified BOOLEAN DEFAULT false,
					created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
				);
			`)
			if err != nil {
				return err
			}

			// Create indexes for better performance
			_, err = d.SQL.Exec(`
				CREATE INDEX IF NOT EXISTS idx_authors_username ON authors(username);
				CREATE INDEX IF NOT EXISTS idx_authors_email ON authors(email);
				CREATE INDEX IF NOT EXISTS idx_authors_is_active ON authors(is_active);
			`)
			if err != nil {
				return err
			}

			// Create updated_at trigger
			_, err = d.SQL.Exec(`
				CREATE OR REPLACE FUNCTION update_authors_updated_at()
				RETURNS TRIGGER AS $$
				BEGIN
					NEW.updated_at = CURRENT_TIMESTAMP;
					RETURN NEW;
				END;
				$$ language 'plpgsql';

				CREATE TRIGGER update_authors_updated_at
					BEFORE UPDATE ON authors
					FOR EACH ROW
					EXECUTE FUNCTION update_authors_updated_at();
			`)
			return err
		},
	}
}
