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

package internal

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func WrapError(err1 error, prefix string, err2 error) error {
	if err2 == nil {
		return err1
	}

	if err1 == nil {
		return fmt.Errorf("[%s] %s", prefix, err2.Error())
	}

	return errors.WithMessagef(err1, "[%s] %s", prefix, err2.Error())
}

func LogReader(c chan []byte, r io.ReadCloser, cb func([]byte) ([]byte, error)) error {
	defer r.Close()
	rd := bufio.NewReader(r)
	for {
		l, _, err := rd.ReadLine()
		switch err {
		case nil:
			if cb == nil {
				c <- l
			} else {
				b, e := cb(l)
				if e != nil {
					return e
				} else {
					c <- b
				}
			}
		case io.EOF:
			return nil
		default:
			return err
		}
	}
}

func ToTar(in, out string) error {
	f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	tw := tar.NewWriter(f)

	return filepath.Walk(in, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsDir() || !info.Mode().IsRegular() {
			return nil
		}

		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		name := info.Name()
		if in != path {
			name = strings.TrimLeft(path, in)
		}

		h := &tar.Header{Name: name, Mode: 0755, Size: int64(len(buf))}
		if err := tw.WriteHeader(h); err != nil {
			return err
		}
		if _, err := tw.Write(buf); err != nil {
			return err
		}

		return nil
	})
}
