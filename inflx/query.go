package inflx

import (
	"fmt"

	"github.com/influxdata/influxdb1-client/v2"
)

// queryDB convenience function to query the database
func Query(db, cmd string) (response *client.Response, res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: db,
	}
	if response, err = Client.Query(q); err == nil {
		if response.Error() != nil {
			return response, res, response.Error()
		}
		res = response.Results
	} else {
		fmt.Println("err!", err)
		return response, res, err
	}
	return response, res, nil
}
