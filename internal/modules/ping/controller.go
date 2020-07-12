/*
 *  Copyright (c) 2020.  Mikhail Knyazhev <markus621@gmail.com>
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/gpl-3.0.html>.
 */

package ping

import (
	"encoding/json"

	"arwos-server/internal/providers/jsonrpc"
)

type Ping struct {
}

func NewController(j *jsonrpc.JRPCModule) *Ping {
	ob := &Ping{}
	j.Inject(ob)
	return ob
}

func (c *Ping) Up() error {
	return nil
}

func (c *Ping) Down() error {
	return nil
}

func (c *Ping) Method() string {
	return "ping"
}

func (c *Ping) Model() json.Unmarshaler {
	return &Input{}
}

func (c *Ping) CallBack(in interface{}) (json.Marshaler, error) {
	return &Output{Text: "pong"}, nil
}
