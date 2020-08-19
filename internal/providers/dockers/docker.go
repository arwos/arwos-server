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
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"arwos-server/internal/providers"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

type (
	dockersRegistry map[string]DockerClientInterface
	imagesItem      struct {
		Tar    string
		Origin string
		Name   string
		Tag    string
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
		list:   make(dockersRegistry, 0),
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
		if err = providers.ToTar(path, tarfile); err != nil {
			return err
		}
		tag := strings.TrimRight(info.Name(), imagesExt)
		dm.images[tag] = imagesItem{
			Tar: tarfile, Origin: path, Name: info.Name(), Tag: tag,
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
	for uniq, i := range dm.list {
		if er := i.Close(); er != nil {
			err = providers.WrapError(err, "Client Down: "+uniq, er)
		}
	}
	return providers.WrapError(err, "Client Close", dm.cli.Close())
}

func (dm *DockersModule) NewClient(image string, clog chan []byte) (DockerClientInterface, error) {
	cli := newClient(image, dm.cli, clog)
	cli.OnRegistering(func(s string, c DockerClientInterface) {
		dm.list[s] = c
	})
	cli.OnUnregistering(func(s string) {
		delete(dm.list, s)
	})
	if path, exist := dm.existImage(image); exist {
		if err := cli.ImageBuild(path); err != nil {
			return nil, err
		}
	} else {
		if err := cli.ImagePull(image); err != nil {
			return nil, err
		}
	}
	if err := cli.Create(); err != nil {
		return nil, providers.WrapError(err, "on close", cli.Close())
	}
	return cli, nil
}

func decodeMessage(b []byte) ([]byte, error) {
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
