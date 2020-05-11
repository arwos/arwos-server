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

package main

import (
	"flag"

	"arwos-server/internal/dockers"

	"github.com/deweppro/core/pkg/provider/db/sqlite"

	"arwos-server/internal/jsonrpc"

	"github.com/deweppro/core/pkg/app"
	"github.com/deweppro/core/pkg/provider/debug"
)

var cfile = flag.String("config", "./configs/config.yaml", "path to config file")
var pidfile = flag.String("pid", "/tmp/arwos-server.pid", "path to pid file")

func main() {
	flag.Parse()

	app.New(
		*cfile,
		app.NewInterfaces().Add(
			&debug.ConfigDebug{},
			&sqlite.ConfigSQLite{},
			&jsonrpc.ConfigHttp{},
			&dockers.ConfigDockers{},
		),
		app.NewInterfaces().Add(
			// providers
			debug.New,
			sqlite.MustNew,
			// modules
			jsonrpc.NewHTTPModule,
			jsonrpc.NewJRPCModule,
			dockers.NewDockersModule,
		),
	).PidFile(*pidfile).Run()
}
