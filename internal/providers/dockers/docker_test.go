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
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/deweppro/core/pkg/filesystem/shell"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDockersModule_HUB(t *testing.T) {
	temp := fmt.Sprintf("/tmp/%s", uuid.New().String())
	data, err := shell.Command("mkdir "+temp, "/tmp", []string{})
	require.NoError(t, err, data)

	mod := NewDockersModule(&ConfigDockers{Docker: ConfigDockerData{
		Images: temp,
		Store:  "docker.io/library",
	}})

	require.NoError(t, mod.Up())

	c := make(chan []byte, 1)
	go func() {
		for d := range c {
			fmt.Println(string(d))
		}
	}()

	cl, er := mod.NewClient("alpine", c)
	require.NoError(t, er)

	l := []string{
		`echo "hello world"`,
		`nslookup google.com`,
		`ping -c 4 google.com`,
	}

	for _, li := range l {
		require.NoError(t, cl.Exec(li))
	}

	require.Error(t, cl.Exec(`p0ng -c 10 google.com`))

	require.NoError(t, mod.Down())
}

func TestDockersModule_FILE(t *testing.T) {
	fl := `
FROM golang:1.13-alpine
ENV GO111MODULE=on
RUN apk update && \
    apk add --virtual build-dependencies build-base \
    bash git
WORKDIR /app
`
	temp := fmt.Sprintf("/tmp/%s", uuid.New().String())
	data, err := shell.Command("mkdir "+temp, "/tmp", []string{})
	require.NoError(t, err, data)

	err = ioutil.WriteFile(temp+"/go13.dockerfile", []byte(fl), 0777)
	require.NoError(t, err)

	mod := NewDockersModule(&ConfigDockers{Docker: ConfigDockerData{
		Images: temp,
		Store:  "docker.io/library",
	}})

	require.NoError(t, mod.Up())

	c := make(chan []byte, 1)
	go func() {
		for d := range c {
			fmt.Println(string(d))
		}
	}()

	cl, er := mod.NewClient("go13", c)
	require.NoError(t, er)

	l := []string{
		`echo "hello world"`,
		`nslookup google.com`,
		`ping -c 4 google.com`,
	}

	for _, li := range l {
		require.NoError(t, cl.Exec(li))
	}

	require.Error(t, cl.Exec(`p0ng -c 10 google.com`))

	require.NoError(t, mod.Down())
}
