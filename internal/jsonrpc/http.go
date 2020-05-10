/*
 *  Copyright (c) 2020.  Mikhail Knyazhev <markus621@gmail.com>
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
 *
 */

package jsonrpc

import (
	"github.com/deweppro/core/pkg/provider/server/http"
)

type HTTPModule struct {
	jrpc    *JRPCModule
	cfg     *ConfigHttp
	httpsrv *http.Server
}

func NewHTTPModule(cfg *ConfigHttp, j *JRPCModule) *HTTPModule {
	return &HTTPModule{
		jrpc:    j,
		cfg:     cfg,
		httpsrv: http.New(),
	}
}

func (h *HTTPModule) Up() error {
	h.httpsrv.SetAddr(h.cfg.Http.Addr)

	h.httpsrv.Route("POST", "/rpc", &http.HttpHandlerItem{
		Call:       h.jrpc.Route,
		Middelware: h.middelware,
		Formatter:  http.JsonRPCFormatter,
	})

	return h.httpsrv.Up()
}

func (h *HTTPModule) Down() error {
	return h.httpsrv.Down()
}

func (h *HTTPModule) middelware(message *http.Message) http.Success {
	return true
}
