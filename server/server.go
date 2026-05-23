package server

import (
	"context"
	"encoding/json"
	"io"

	"github.com/sourcegraph/jsonrpc2"
)

func Serve(ctx context.Context, inout io.ReadWriteCloser) error {
	stream := jsonrpc2.NewBufferedStream(inout, jsonrpc2.VSCodeObjectCodec{})
	conn := jsonrpc2.NewConn(ctx, stream, jsonrpc2.HandlerWithError(handle))
	defer conn.Close()

	select {
	case <-conn.DisconnectNotify():
		return nil
	case <-ctx.Done():
		return inout.Close()
	}
}

func handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (any, error) {
	switch req.Method {
	case "evaluate":
		var p EvaluateParams
		if req.Params != nil {
			if err := json.Unmarshal(*req.Params, &p); err != nil {
				return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()}
			}
		}
		return evaluate(ctx, p)
	default:
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: req.Method}
	}
}
