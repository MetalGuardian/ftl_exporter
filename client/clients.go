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

// GetTopClients retrieves the list of clients together with amount of queries
// made by each client from response of `>top-clients` command
func (client *FTLClient) GetTopClients() (*TopEntries, error) {
	return topClientsFor(">top-clients", client)
}

// GetTopBlockedClients retrieves the list of clients together with amount of blocked
// queries made by each client from response of `>top-clients` command
func (client *FTLClient) GetTopBlockedClients() (*TopEntries, error) {
	return topClientsFor(">top-clients blocked", client)
}

func topClientsFor(command string, client *FTLClient) (*TopEntries, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if err := sendCommand(conn, command); err != nil {
		return nil, err
	}

	total, err := readInt32(conn)
	if err != nil {
		return nil, err
	}

	result := TopEntries{
		Total: total,
	}

	for {
		_, err := readString(conn)
		if err == errEndOfInput {
			break
		}
		if err != nil {
			return nil, err
		}

		address, err := readString(conn)
		if err != nil {
			return nil, err
		}

		count, err := readInt32(conn)
		if err != nil {
			return nil, err
		}

		result.Entries = append(result.Entries, entry{Label: address, Count: count})
	}

	return &result, nil
}
