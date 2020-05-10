/*
 *  Copyright (c) 2020.  Mikhail Knyazhev <markus621@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/gpl-3.0.html>.
 *
 */

package jsonrpc

import (
	"encoding/json"

	"github.com/deweppro/core/pkg/provider/server/http"
)

type JRPCModule struct {
	allowedMethods map[string]struct{}
}

func NewJRPCModule() *JRPCModule {
	return &JRPCModule{
		allowedMethods: make(map[string]struct{}),
	}
}

func (m *JRPCModule) Up() error {
	return nil
}

func (m *JRPCModule) Down() error {
	return nil
}

func (m *JRPCModule) Route(message *http.Message) {
	message.Encode(func() (int, map[string]string, interface{}) {
		// ---
		var req http.JsonRPCRequest
		er := message.Decode(func(data []byte) error {
			return json.Unmarshal(data, &req)
		})
		if er != nil {
			return 500, nil, er
		}

		// ---
		if _, ok := m.allowedMethods[req.Method]; !ok {
			return 500, nil, errNotAllowedMethod
		}

		// ---
		var resp json.RawMessage
		//if err := m.Nats.Request(req.Method, &req.Params, &resp, 30); err != nil {
		//	return 500, nil, err
		//}
		return 0, nil, &http.JsonRPCResponse{
			ID:     req.ID,
			Result: resp,
		}
	})
}
