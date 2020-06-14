// Copyright 2020 Ivan Pushkin
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"net"
)

// GetQueriesOverTime retrieves amount of allowed and blocked queries
// for the last 24 hours aggregated over 10 minute intervals
// from response of `>overTime` command
func (client *FTLClient) GetQueriesOverTime() (*QueriesOverTime, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if err := sendCommand(conn, ">overTime"); err != nil {
		return nil, err
	}

	var result QueriesOverTime

	lines, err := readMapCount(conn)
	if err != nil {
		return nil, err
	}

	for i := 0; i < lines; i++ {
		timestamp, err := readInt32(conn)
		if err != nil {
			return nil, err
		}

		count, err := readInt32(conn)
		if err != nil {
			return nil, err
		}

		result.Forwarded = append(result.Forwarded, timestampCount{
			Timestamp: timestamp,
			Count:     count,
		})
	}

	lines, err = readMapCount(conn)
	if err != nil {
		return nil, err
	}

	for i := 0; i < lines; i++ {
		timestamp, err := readInt32(conn)
		if err != nil {
			return nil, err
		}

		count, err := readInt32(conn)
		if err != nil {
			return nil, err
		}

		result.Blocked = append(result.Blocked, timestampCount{
			Timestamp: timestamp,
			Count:     count,
		})
	}

	return &result, nil
}
