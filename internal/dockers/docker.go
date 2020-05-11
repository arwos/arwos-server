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
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/client"

	"arwos-server/internal"
)

const (
	imagesExt = ".dockerfile"
	tarExt    = ".tar"
)

type (
	dockersRegistry map[string]*dockerClient
	imagesItem      struct {
		Tar    string
		Origin string
		Name   string
	}
	imagesList map[string]imagesItem

	DockersModule struct {
		cli    *client.Client
		cfg    *ConfigDockers
		list   dockersRegistry
		images imagesList
	}
)

func NewDockersModule(cfg *ConfigDockers) *DockersModule {
	return &DockersModule{
		cfg:    cfg,
		list:   make(dockersRegistry),
		images: make(imagesList),
	}
}

func (dm *DockersModule) Up() error {
	err := filepath.Walk(dm.getPath(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsDir() ||
			!info.Mode().IsRegular() ||
			filepath.Ext(info.Name()) != imagesExt {
			return nil
		}

		tarfile := path + tarExt
		if err = internal.ToTar(path, tarfile); err != nil {
			return err
		}

		dm.images[strings.TrimRight(info.Name(), imagesExt)] = imagesItem{
			Tar: tarfile, Origin: path, Name: info.Name(),
		}
		return nil
	})
	if err != nil {
		return err
	}

	dm.cli, err = client.NewClientWithOpts()
	return err
}

func (dm *DockersModule) Down() error {
	var err error
	for n, i := range dm.list {
		if er := i.Down(); er != nil {
			err = internal.WrapError(err, "Client Down: "+n, er)
		}
	}
	return internal.WrapError(err, "Client Close", dm.cli.Close())
}

func (dm *DockersModule) NewClient(image string, clog chan []byte) (*dockerClient, error) {
	path, exist := dm.existImage(image)
	if exist {
		f, err := os.Open(path.Tar)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		resp, err := dm.cli.ImageBuild(context.TODO(), f, types.ImageBuildOptions{
			Dockerfile: path.Name,
			//NetworkMode: "host",
		})
		if err != nil {
			return nil, err
		}

		if err := internal.LogReader(clog, resp.Body, dm.decode); err != nil {
			return nil, err
		}
	} else {
		//image = dm.cfg.Docker.Store + "/" + image

		resp, err := dm.cli.ImagePull(context.TODO(), image, types.ImagePullOptions{})
		if err != nil {
			return nil, err
		}

		if err := internal.LogReader(clog, resp, dm.decode); err != nil {
			return nil, err
		}
	}

	return newClient(image, dm.cli), nil
}

func (dm *DockersModule) decode(b []byte) ([]byte, error) {
	var msg DockerMessage
	if err := json.Unmarshal(b, &msg); err != nil {
		return nil, err
	}
	if msg.IsError() {
		return nil, errors.New(msg.Error)
	}
	return msg.ByteValue(), nil
}

func (dm *DockersModule) getPath() string {
	return dm.cfg.Docker.Images
}

func (dm *DockersModule) existImage(name string) (imagesItem, bool) {
	path, ok := dm.images[name]
	return path, ok
}
