package middleware

import (
	"bufio"
	"bytes"
	"cro_test/pkg/logger"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

type (
	bodyDumpResponseWriter struct {
		io.Writer
		http.ResponseWriter
		status int
	}
)

func Logging(logger logger.Logger, reqDump, respDump bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			reqBody, _ := ioutil.ReadAll(req.Body)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
			respBody := new(bytes.Buffer)
			mw := io.MultiWriter(respBody, resp)
			w := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: resp}
			resp = w
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = xid.New().String()
			}
			resp.Header().Set(echo.HeaderXRequestID, rid)

			if strings.Contains(req.URL.Path, "swagger") {
				next.ServeHTTP(resp, req)
				return
			}

			log := logger.Field("request_id", rid).
				Field("method", req.Method).Field("path", req.URL.Path).Field("bytes_in", req.Header.Get("Content-Length"))
			ctx := log.Attach(req.Context())
			req = req.WithContext(ctx)

			if reqDump && len(reqBody) > 0 && len(reqBody) <= 300 {
				log.Infof("request dump\n%s\nrequest dump", string(reqBody))
			}

			start := time.Now()
			defer func() {
				end := time.Now()
				log = log.Field("status", http.StatusText(w.status)).
					Field("latency", end.Sub(start).String()).
					Field("bytes_out", respBody.Len())
				if respDump && len(respBody.String()) <= 300 && len(respBody.String()) > 0 {
					log.Infof("response dump\n%s\nresponse dump", respBody.String())
				}
				if w.status < 400 {
					log.Info("request done")
				}
			}()
			next.ServeHTTP(resp, req)
		})
	}
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
