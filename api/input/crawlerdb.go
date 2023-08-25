package input

import (
	"database/sql"
	"time"
)

type CrawledNode struct {
	ID              string
	IP              string
	ConnType        string
	Now             string
	ClientType      string
	SoftwareVersion uint64
	Capabilities    string
	NetworkID       uint64
	Country         string
	ForkID          string
	TotalDifficulty string
	Validator       bool
}

func ReadRecentNodes(db *sql.DB, lastCheck time.Time, networkID string) ([]CrawledNode, error) {
	queryStmt := "SELECT ID, IP, Now, ClientType, SoftwareVersion, Capabilities, TotalDifficulty, NetworkID, Country, " +
		"ForkID, ConnType, Validator FROM nodes WHERE Now > ?"
	if networkID != "" {
		queryStmt = queryStmt + " AND NetworkID = " + networkID
	}

	// TODO do a proper check here ^
	rows, err := db.Query(queryStmt, lastCheck.String())

	if err != nil {
		return nil, err
	}

	var nodes []CrawledNode
	for rows.Next() {
		var node CrawledNode
		err = rows.Scan(&node.ID, &node.IP, &node.Now, &node.ClientType, &node.SoftwareVersion, &node.Capabilities, &node.TotalDifficulty, &node.NetworkID, &node.Country, &node.ForkID, &node.ConnType, &node.Validator)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
