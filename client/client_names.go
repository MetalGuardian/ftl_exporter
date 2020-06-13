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

// GetClientNames retrieves ordered list of client's names from
// response of `>client-names` command
func (client *FTLClient) GetClientNames() (*[]Client, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if err := sendCommand(conn, ">client-names"); err != nil {
		return nil, err
	}

	var clients []Client
	for {
		name, err := readString(conn)
		if err == EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		address, err := readString(conn)
		if err != nil {
			return nil, err
		}

		clients = append(clients, struct {
			Name    string
			Address string
		}{Name: name, Address: address})
	}

	return &clients, nil
}