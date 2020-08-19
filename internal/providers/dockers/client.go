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
	"os"
	"time"

	"arwos-server/internal/providers"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type (
	dockerClient struct {
		cli   *client.Client
		cid   string
		name  string
		uniq  string
		clog  chan []byte
		dereg func(s string)
	}

	DockerClientInterface interface {
		Close() error
		Exec(cmd string) error
	}
)

func newClient(name string, cli *client.Client, clog chan []byte) *dockerClient {
	return &dockerClient{
		cli:  cli,
		name: name,
		clog: clog,
		uniq: uuid.New().String(),
	}
}

func (dc *dockerClient) Create() error {
	cfg := &container.Config{
		Image: dc.name,
		Cmd:   []string{"sh", "-ce", cLoopCMD},
	}
	resp, err := dc.cli.ContainerCreate(context.TODO(), cfg, nil, nil, dc.uniq)
	if err != nil {
		return errors.Wrap(err, "ContainerCreate: "+dc.name)
	}
	dc.cid = resp.ID
	return dc.cli.ContainerStart(context.TODO(), dc.cid, types.ContainerStartOptions{})
}

func (dc *dockerClient) Close() error {
	var err error
	if len(dc.cid) > 0 {
		timeout := 0 * time.Second
		if er := dc.cli.ContainerStop(context.TODO(), dc.cid, &timeout); er != nil {
			err = providers.WrapError(err, "ContainerStop: "+dc.name, er)
		}
		if er := dc.cli.ContainerRemove(context.TODO(), dc.cid, types.ContainerRemoveOptions{}); er != nil {
			err = providers.WrapError(err, "ContainerRemove: "+dc.name, er)
		}
	}
	dc.dereg(dc.uniq)
	return err
}

func (dc *dockerClient) Exec(cmd string) error {
	ec := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          []string{"sh", "-ce", cmd},
	}
	resp, err1 := dc.cli.ContainerExecCreate(context.TODO(), dc.cid, ec)
	if err1 != nil {
		return errors.Wrap(err1, "ContainerExecCreate: "+dc.name)
	}
	hr, err2 := dc.cli.ContainerExecAttach(context.TODO(), resp.ID, types.ExecStartCheck{})
	if err2 != nil {
		return errors.Wrap(err2, "ContainerExecAttach: "+dc.name)
	}
	defer hr.Close()
	if err := providers.LogReader(dc.clog, hr.Reader, nil); err != nil {
		return errors.Wrap(err, "read data from ContainerExecAttach")
	}
	eci, err3 := dc.cli.ContainerExecInspect(context.TODO(), resp.ID)
	if err3 != nil {
		return errors.Wrap(err2, "ContainerExecInspect: "+dc.name)
	}
	if eci.ExitCode > 0 {
		return errors.Wrap(errorBadExec, cmd)
	}
	return nil
}

func (dc *dockerClient) ImageBuild(img imagesItem) error {
	f, err := os.Open(img.Tar)
	if err != nil {
		return err
	}
	defer f.Close()
	resp, err := dc.cli.ImageBuild(context.TODO(), f, types.ImageBuildOptions{
		Dockerfile: img.Name,
		Tags:       []string{img.Tag},
	})
	if err != nil {
		return errors.Wrap(err, "error on ImageBuild")
	}
	if err := providers.LogReader(dc.clog, resp.Body, decodeMessage); err != nil {
		return errors.Wrap(err, "read data from ImageBuild")
	}
	return resp.Body.Close()
}

func (dc *dockerClient) ImagePull(img string) error {
	resp, err := dc.cli.ImagePull(context.TODO(), img, types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "error on ImagePull")
	}
	if err := providers.LogReader(dc.clog, resp, decodeMessage); err != nil {
		return errors.Wrap(err, "read data from ImagePull")
	}
	return resp.Close()
}

func (dc *dockerClient) OnRegistering(call func(s string, c DockerClientInterface)) {
	call(dc.uniq, dc)
}

func (dc *dockerClient) OnUnregistering(call func(s string)) {
	dc.dereg = call
}
