package mock

import "github.com/bacv/pow-wow/lib"

type ResponseWriter struct {
	Written lib.Message
}

func (w *ResponseWriter) Write(m lib.Message) error {
	w.Written = m
	return nil
}

func (w *ResponseWriter) Close() {}
