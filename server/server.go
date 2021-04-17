package server

import (
	"bytes"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

//requestHandler is used to handle test request
func requestHandler(ctx *fasthttp.RequestCtx) {
	var uriPath = ctx.Path()
	if len(uriPath) > 1 {

		var payload = bytes.NewBuffer([]byte{})
		size := ctx.QueryArgs().PeekBytes([]byte("size"))
		length, err := strconv.ParseInt(string(size), 10, 64)
		if err != nil {
			payload.Write([]byte("ERROR"))
		} else {
			payload.Write(make([]byte, int(length)))
		}

		ctx.Response.Header.SetStatusCode(200)
		//ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetBytesKV([]byte("Content-Type"), []byte("text/plain; charset=utf8"))
		ctx.Response.SetBody(payload.Bytes())
	}

}

//requestHandler2 is used to handle test request
func requestHandler2(ctx *fasthttp.RequestCtx) {
	var uriPath = ctx.Path()
	if len(uriPath) > 1 {

		var payload = bytes.NewBuffer([]byte{})
		size := ctx.QueryArgs().PeekBytes([]byte("size"))
		length, err := strconv.ParseInt(string(size), 10, 64)
		if err != nil {
			payload.Write([]byte("ERROR"))
		} else {
			payload.Write(make([]byte, int(length)))
		}

		ctx.Response.Header.SetStatusCode(200)
		ctx.Response.Header.SetConnectionClose()
		ctx.Response.Header.SetBytesKV([]byte("Content-Type"), []byte("text/plain; charset=utf8"))
		ctx.Response.SetBody(payload.Bytes())
	}

}

func Run(address string, ConnClose bool) {
	server := &fasthttp.Server{
		Handler: requestHandler,
		Name:    "sdwan speed test",
	}

	if ConnClose {
		server.Handler = requestHandler2
	}

	if err := server.ListenAndServe(address); err != nil {
		logrus.Fatal("error in ListenAndServe", err)
	}

}
