package mock

import "github.com/bacv/pow-wow/lib/protocol"

type ResponseWriter struct {
	Written protocol.Message
}

func (w *ResponseWriter) Write(m protocol.Message) error {
	w.Written = m
	return nil
}

func (w *ResponseWriter) Close() {}
