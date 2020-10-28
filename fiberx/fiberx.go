package fiberx

import (
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

// ErrHandler is Dawn's error handler
var ErrHandler = func(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	res := Response{
		Code:    code,
		Message: utils.StatusMessage(code),
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		res.Code = fiber.StatusUnprocessableEntity
		res.Message = ""
		res.Data = removeTopStruct(errs.Translate(trans))
	} else if e, ok := err.(*fiber.Error); ok {
		res.Code = e.Code
		res.Message = e.Message
	}

	return Resp(c, res.Code, res)
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, msg := range fields {
		stripStruct := field[strings.Index(field, ".")+1:]
		res[stripStruct] = strings.TrimLeft(msg, stripStruct)
	}
	return res
}

// ValidateBody accepts a obj holds results from BodyParser
// and then do the validation by a validator
func ValidateBody(c *fiber.Ctx, obj interface{}) (err error) {
	if err := c.BodyParser(obj); err != nil {
		return err
	}

	return V.Struct(obj)
}

// ValidateQuery accepts a obj holds results from QueryParser
// and then do the validation by a validator
func ValidateQuery(c *fiber.Ctx, obj interface{}) (err error) {
	if err := c.QueryParser(obj); err != nil {
		return err
	}

	return V.Struct(obj)
}

// Response is a unified format for api results
type Response struct {
	// Code is the status code by default, but also can be
	// a custom code
	Code int `json:"code,omitempty"`
	// Message shows detail thing back to caller
	Message string `json:"message,omitempty"`
	// RequestID needs to be used with middleware
	RequestID string `json:"request_id,omitempty"`
	// Data accepts any thing as the response data
	Data interface{} `json:"data,omitempty"`
}

// Resp returns the custom response
func Resp(c *fiber.Ctx, statusCode int, res Response) error {
	if res.Code == 0 {
		res.Code = statusCode
	}

	if id := c.Response().Header.Peek(fiber.HeaderXRequestID); len(id) > 0 && res.RequestID == "" {
		res.RequestID = utils.GetString(id)
	}

	return c.Status(statusCode).JSON(res)
}

// Data returns data with status code OK by default
func Data(c *fiber.Ctx, data interface{}) error {
	return Resp(c, fiber.StatusOK, Response{Data: data})
}

// Message wraps for RespOK with required message
func Message(c *fiber.Ctx, msg string) error {
	return RespOK(c, msg)
}

// respCommon
func respCommon(c *fiber.Ctx, code int, msg ...string) error {
	res := Response{
		Message: utils.StatusMessage(code),
	}

	if len(msg) > 0 {
		res.Message = msg[0]
	}
	return Resp(c, code, res)
}

// RespOK responses with status code 200 RFC 7231, 6.3.1
func RespOK(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusOK, msg...)
}

// RespCreated responses with status code 201 RFC 7231, 6.3.2
func RespCreated(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusCreated, msg...)
}

// RespAccepted responses with status code 202 RFC 7231, 6.3.3
func RespAccepted(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusAccepted, msg...)
}

// RespNonAuthoritativeInformation responses with status code 203 RFC 7231, 6.3.4
func RespNonAuthoritativeInformation(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusNonAuthoritativeInformation, msg...)
}

// RespNoContent responses with status code 204 RFC 7231, 6.3.5
func RespNoContent(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusNoContent, msg...)
}

// RespResetContent responses with status code 205 RFC 7231, 6.3.6
func RespResetContent(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusResetContent, msg...)
}

// RespPartialContent responses with status code 206 RFC 7233, 4.1
func RespPartialContent(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusPartialContent, msg...)
}

// RespMultiStatus responses with status code 207 RFC 4918, 11.1
func RespMultiStatus(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusMultiStatus, msg...)
}

// RespAlreadyReported responses with status code 208 RFC 5842, 7.1
func RespAlreadyReported(c *fiber.Ctx, msg ...string) error {
	return respCommon(c, fiber.StatusAlreadyReported, msg...)
}

var pid = os.Getpid()

// Logger logs request and response info to os.Stdout
// or os.Stderr. The format is:
// time #pid[-request-id]: latency status clientIP method protocol://host_path[ error]
func Logger() fiber.Handler {
	return func(ctx *fiber.Ctx) (err error) {
		start := time.Now()

		err = ctx.Next()

		end := time.Now()
		latency := end.Sub(start).Truncate(time.Microsecond)

		bb := bytebufferpool.Get()
		defer bytebufferpool.Put(bb)

		// append time
		bb.B = end.AppendFormat(bb.B, "2006/01/02 15:04:05.000")

		// append pid
		_, _ = bb.WriteString(" #")
		bb.B = fasthttp.AppendUint(bb.B, pid)

		// append request id
		if requestId := ctx.Response().Header.Peek(fiber.HeaderXRequestID); len(requestId) > 0 {
			_ = bb.WriteByte('-')
			_, _ = bb.Write(requestId)
		}
		_, _ = bb.WriteString(": ")

		// append latency
		_, _ = bb.WriteString(latency.String())
		_ = bb.WriteByte(' ')

		// append status code
		statusCode := ctx.Response().StatusCode()
		bb.B = fasthttp.AppendUint(bb.B, statusCode)
		_ = bb.WriteByte(' ')

		// append client ip
		_, _ = bb.WriteString(ctx.IP())
		_ = bb.WriteByte(' ')

		// append http method
		_, _ = bb.WriteString(ctx.Method())
		_ = bb.WriteByte(' ')

		// append http protocol://host/uri
		_, _ = bb.WriteString(ctx.Protocol())
		_, _ = bb.WriteString("://")
		_, _ = bb.Write(ctx.Request().URI().Host())
		_, _ = bb.Write(ctx.Request().RequestURI())

		w := os.Stdout
		// append error
		if err != nil {
			_ = bb.WriteByte(' ')
			_, _ = bb.WriteString(err.Error())
			w = os.Stderr
		}

		// append newline
		_ = bb.WriteByte('\n')

		// ignore error on purpose
		_, _ = bb.WriteTo(w)

		return
	}
}
