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

package ftl_client

import (
	"encoding/binary"
	"io"
	"net"
)

func (client *Client) GetTopClients() (*DomainEntries, error) {
	return topClientsFor(">top-clients", client)
}

func (client *Client) GetTopBlockedClients() (*DomainEntries, error) {
	return topClientsFor(">top-clients blocked", client)
}

func topClientsFor(command string, client *Client) (*DomainEntries, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if _, err := conn.Write([]byte(command)); err != nil {
		return nil, err
	}

	var result DomainEntries

	if err := binary.Read(conn, binary.BigEndian, &result.Total); err != nil {
		return nil, err
	}

	for {
		var format uint8
		err := binary.Read(conn, binary.BigEndian, &format)

		if err == io.EOF || format == formatEOF {
			break
		}

		if err != nil {
			return nil, err
		}

		var emptyString struct {
			_ uint32
			_ [0]byte
		}

		if err := binary.Read(conn, binary.BigEndian, &emptyString); err != nil {
			return nil, err
		}

		err = binary.Read(conn, binary.BigEndian, &format)
		if err != nil {
			return nil, err
		}

		var length uint32

		err = binary.Read(conn, binary.BigEndian, &length)
		if err != nil {
			return nil, err
		}

		address := make([]byte, length)

		err = binary.Read(conn, binary.BigEndian, &address)
		if err != nil {
			return nil, err
		}

		var count UInt32Block

		err = binary.Read(conn, binary.BigEndian, &count)
		if err != nil {
			return nil, err
		}

		result.List = append(result.List, struct {
			Domain string
			Count  UInt32Block
		}{Domain: string(address), Count: count})
	}

	return &result, nil
}