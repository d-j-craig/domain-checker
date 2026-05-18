package db

import "domain-checker/types"

func LookUpDomain(dl []string) (types.StatusDB, error) {	

	// insert data into db domains table
	rows, err := Conn.Query(ctx,
		`SELECT 
		domains.url, 
		status.status_code,
		status.response_time,
		status.response_size,
		status.error
		FROM domains
		JOIN status ON domains.id = status.domain_id
		WHERE domains.url = ANY($1)	
		`, dl)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// create results statusdb
	var results types.StatusDB
	// save the returned rows in a statusDB struct
	for rows.Next() {
		var status types.StatusList
		err := rows.Scan(
			&status.Name,
			&status.Status,
			&status.ResponseTime,
			&status.ResponseSize,
			&status.Err,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, status)
	}

	return results, nil
	

}