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

	"github.com/sirupsen/logrus"

	"github.com/deweppro/core/pkg/server/http"
)

type InjectableInterface interface {
	Method() string
	Model() json.Unmarshaler
	CallBack(interface{}) (json.Marshaler, error)
}

type JRPCModule struct {
	allowedMethods map[string]InjectableInterface
}

func NewJRPCModule() *JRPCModule {
	return &JRPCModule{
		allowedMethods: make(map[string]InjectableInterface),
	}
}

func (m *JRPCModule) Inject(i InjectableInterface) {
	logrus.Info("add method = ", i.Method())
	m.allowedMethods[i.Method()] = i
}

func (m *JRPCModule) Up() error {
	return nil
}

func (m *JRPCModule) Down() error {
	return nil
}

func (m *JRPCModule) CallBack(message *http.Message) error {
	message.Encode(func() (int, map[string]string, interface{}) {
		var req http.JsonRPCRequest

		if err := message.Decode(func(data []byte) error { return json.Unmarshal(data, &req) }); err != nil {
			return 999, nil, err
		}

		module, ok := m.allowedMethods[req.Method]
		if !ok {
			return 1, nil, &http.JsonRPCResponseError{
				ID: req.ID,
				Error: http.JsonRPCErrorBody{
					Code:    1,
					Message: errNotAllowedMethod.Error(),
				},
			}
		}

		model := module.Model()
		if err := json.Unmarshal(req.Params, model); err != nil {
			return 2, nil, &http.JsonRPCResponseError{
				ID: req.ID,
				Error: http.JsonRPCErrorBody{
					Code:    2,
					Message: err.Error(),
				},
			}
		}

		resp, err := module.CallBack(model)
		if err != nil {
			return 3, nil, &http.JsonRPCResponseError{
				ID: req.ID,
				Error: http.JsonRPCErrorBody{
					Code:    3,
					Message: err.Error(),
				},
			}
		}

		raw, err := resp.MarshalJSON()
		if err != nil {
			return 4, nil, &http.JsonRPCResponseError{
				ID: req.ID,
				Error: http.JsonRPCErrorBody{
					Code:    4,
					Message: err.Error(),
				},
			}
		}

		return 0, nil, &http.JsonRPCResponse{
			ID:     req.ID,
			Result: raw,
		}
	})
	return nil
}
