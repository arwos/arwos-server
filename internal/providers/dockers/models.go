/*
 * Copyright (c) 2020.  Mikhail Knyazhev <markus621@gmail.com>
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
 */

package dockers

import "github.com/pkg/errors"

const (
	cLoopCMD  = `i=0; while [ $i -le 5 ]; do sleep 10s; done`
	imagesExt = ".dockerfile"
	tarExt    = ".tar"
)

var (
	errorBadExec = errors.New(`completed with error`)
)

type ConfigDockers struct {
	Docker ConfigDockerData `yaml:"docker"`
}

type ConfigDockerData struct {
	Images string `yaml:"images"`
	Store  string `yaml:"store"`
}

type DockerMessage struct {
	Stream string `json:"stream,omitempty"`
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (m DockerMessage) ByteValue() []byte {
	return []byte(m.Value())
}

func (m DockerMessage) Value() string {
	return m.Stream + m.Status + m.Error
}

func (m DockerMessage) IsError() bool {
	return len(m.Error) > 0
}
