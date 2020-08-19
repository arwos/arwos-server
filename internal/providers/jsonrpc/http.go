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
	http2 "net/http"

	"github.com/deweppro/core/pkg/server/http"
)

type HTTPModule struct {
	jrpc   *JRPCModule
	cfg    *ConfigHttp
	server *http.Server
}

func NewHTTPModule(cfg *ConfigHttp, j *JRPCModule) *HTTPModule {
	return &HTTPModule{
		jrpc:   j,
		cfg:    cfg,
		server: http.NewServer(),
	}
}

func (h *HTTPModule) Up() error {
	h.server.SetAddr(h.cfg.Http.Addr)
	h.server.AddRoute(h)
	return h.server.Up()
}

func (h *HTTPModule) Down() error {
	return h.server.Down()
}

func (h *HTTPModule) Handlers() []http.Handler {
	return []http.Handler{
		{Method: http2.MethodPost, Path: "/rpc", Call: h.jrpc.CallBack},
	}
}

func (h *HTTPModule) Formatter() http.FormatterFunc {
	return http.JsonRPCFormatter
}

func (h *HTTPModule) Middleware() http.CallFunc {
	return func(message *http.Message) error {
		return nil
	}
}
