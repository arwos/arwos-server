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

package dockers

import (
	"context"
	"time"

	"arwos-server/internal"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type dockerClient struct {
	cli  *client.Client
	cid  string
	name string
}

func newClient(name string, cli *client.Client) *dockerClient {
	return &dockerClient{
		cli:  cli,
		name: name,
	}
}

func (dc *dockerClient) Down() error {
	var err error

	if len(dc.cid) > 0 {
		timeout := 10 * time.Second
		if er := dc.cli.ContainerStop(context.TODO(), dc.cid, &timeout); er != nil {
			err = internal.WrapError(err, "ContainerStop: "+dc.name, er)
		}

		if er := dc.cli.ContainerRemove(context.TODO(), dc.cid, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   true,
			Force:         true,
		}); er != nil {
			err = internal.WrapError(err, "ContainerRemove: "+dc.name, er)
		}
	}

	return err
}
