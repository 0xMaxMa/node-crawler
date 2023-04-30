package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MariusVanDerWijden/node-crawler-backend/input"
	"github.com/MariusVanDerWijden/node-crawler-backend/parser"
)

func createDB(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE nodes (
		ID text not null,
		ip text,
		conn_type text,
		name text,
		version_major number,
		version_minor number,
		version_patch number,
		version_tag text,
		version_build text,
		version_date text,
		os_name text,
		os_architecture text,
		language_name text,
		language_version text,
		total_difficulty text,
		last_crawled datetime,
		country_name text,
		validator boolean,
		PRIMARY KEY (ID)
	);
	delete from nodes;
	`
	_, err := db.Exec(sqlStmt)
	return err
}

func InsertCrawledNodes(db *sql.DB, crawledNodes []input.CrawledNode) error {
	fmt.Printf("Writing nodes to db: %v\n", len(crawledNodes))

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		`insert into nodes(
			ID,
			ip,
			conn_type,
			name, 
			version_major, version_minor, version_patch, version_tag, version_build, version_date, 
			os_name, os_architecture, 
			language_name, language_version, total_difficulty, last_crawled, country_name, validator)
			values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(ID) DO UPDATE SET 
			name=excluded.name,
			ip=excluded.ip,
			conn_type=excluded.conn_type,
			version_major=excluded.version_major,
			version_minor=excluded.version_minor,
			version_patch=excluded.version_patch,
			version_tag=excluded.version_tag,
			version_build=excluded.version_build,
			version_date=excluded.version_date,
			os_name=excluded.os_name,
			os_architecture=excluded.os_architecture,
			language_name=excluded.language_name,
			language_version=excluded.language_version,
			total_difficulty=excluded.total_difficulty,
			last_crawled=excluded.last_crawled,
			country_name=excluded.country_name,
			validator=excluded.validator
			WHERE name=excluded.name OR excluded.name != "unknown"`)
	if err != nil {
		return err
	}

	for _, node := range crawledNodes {
		parsed, e := parser.ParseVersionString(node.ClientType)
		if parsed != nil {
			_, err = stmt.Exec(
				node.ID,
				node.IP,
				node.ConnType,
				parsed.Name,
				parsed.Version.Major,
				parsed.Version.Minor,
				parsed.Version.Patch,
				parsed.Version.Tag,
				parsed.Version.Build,
				parsed.Version.Date,
				parsed.Os.Os,
				parsed.Os.Architecture,
				parsed.Language.Name,
				parsed.Language.Version,
				node.TotalDifficulty,
				time.Now(),
				node.Country,
				node.Validator,
			)
			if err != nil {
				panic(err)
			}
		} else {
			if e == nil {
				_, err = stmt.Exec(
					node.ID,
					node.IP,
					node.ConnType,
					"",
					0,
					0,
					0,
					"",
					"",
					"",
					"",
					"",
					"",
					"",
					node.TotalDifficulty,
					time.Now(),
					node.Country,
					node.Validator,
				)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	return tx.Commit()
}

func dropOldNodes(db *sql.DB, minTimePassed time.Duration) error {
	fmt.Printf("Dropping all nodes older than: %v\n", minTimePassed)
	oldest := time.Now().Add(-minTimePassed)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`DELETE FROM nodes WHERE last_crawled < ?`)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(oldest)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	fmt.Printf("Dropped %v nodes\n", affected)
	return tx.Commit()
}
