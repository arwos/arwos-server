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

	"github.com/docker/docker/client"
	"github.com/pkg/errors"

	"arwos-server/internal"
)

const (
	imagesExt = ".dockerfile"
	tarExt    = ".tar"
)

type (
	dockersRegistry map[string]DockerClientInterface
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

	//c := make(chan []byte, 100)
	//go func() {
	//	for d := range c {
	//		fmt.Println(string(d))
	//	}
	//}()
	//
	//cl, er := dm.NewClient("alpine", c)
	//if er != nil {
	//	return errors.Wrap(er, "[New Client]")
	//}
	//
	//l := []string{
	//	`echo "hello world"`,
	//	`nslookup google.com`,
	//	`ping -c 10 google.com`,
	//	`p0ng -c 10 google.com`,
	//}
	//var errcmd error
	//for _, li := range l {
	//	errcmd = internal.WrapError(errcmd, li, cl.Exec(li))
	//}
	//
	//return errcmd

	return err
}

func (dm *DockersModule) Down() error {
	var err error
	for uniq, i := range dm.list {
		if er := i.Close(); er != nil {
			err = internal.WrapError(err, "Client Down: "+uniq, er)
		}
	}
	return internal.WrapError(err, "Client Close", dm.cli.Close())
}

func (dm *DockersModule) NewClient(image string, clog chan []byte) (DockerClientInterface, error) {
	cli := newClient(image, dm.cli, clog)

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
		return nil, internal.WrapError(err, "on close", cli.Close())
	}

	cli.OnRegistering(func(s string, c DockerClientInterface) {
		dm.list[s] = c
	})

	cli.OnDeregistering(func(s string) {
		delete(dm.list, s)
	})

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
