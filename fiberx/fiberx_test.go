package fiberx

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

type respCase struct {
	app        *fiber.App
	method     string
	target     string
	reqBody    io.Reader
	statusCode int
	respBody   string
}

func Test_Fiberx_ErrorHandler(t *testing.T) {
	t.Run("StatusUnprocessableEntity", func(t *testing.T) {
		t.Parallel()

		app := fiber.New(fiber.Config{
			ErrorHandler: ErrHandler,
		})
		app.Get("/422", func(c *fiber.Ctx) error {
			type User struct {
				Username string `validate:"required"`
				Field1   string `validate:"required,lt=10"`
				Field2   string `validate:"required,gt=1"`
			}

			user := User{
				Username: "kiyon",
				Field1:   "This field is always too long.",
				Field2:   "1",
			}

			return V.Struct(user)
		})

		assertRespCase(t, respCase{
			app:        app,
			method:     fiber.MethodGet,
			target:     "/422",
			statusCode: fiber.StatusUnprocessableEntity,
			respBody:   `{"code":422, "data":{"Field1":" must be less than 10 characters in length", "Field2":" must be greater than 1 character in length"}}`,
		})
	})

	t.Run("normal error", func(t *testing.T) {
		t.Parallel()

		app := fiber.New(fiber.Config{
			ErrorHandler: ErrHandler,
		})
		app.Get("/", func(c *fiber.Ctx) error {
			return errors.New("hi, i'm an error")
		})

		assertRespCase(t, respCase{
			app:        app,
			method:     fiber.MethodGet,
			target:     "/",
			statusCode: fiber.StatusInternalServerError,
			respBody:   `{"code":500, "message":"Internal Server Error"}`,
		})
	})

	t.Run("fiber error", func(t *testing.T) {
		t.Parallel()

		app := fiber.New(fiber.Config{
			ErrorHandler: ErrHandler,
		})
		app.Get("/400", func(c *fiber.Ctx) error {
			return fiber.ErrBadRequest
		})

		assertRespCase(t, respCase{
			app:        app,
			method:     fiber.MethodGet,
			target:     "/400",
			statusCode: fiber.StatusBadRequest,
			respBody:   `{"code":400, "message":"Bad Request"}`,
		})
	})
}

func Test_Fiberx_ValidateBody(t *testing.T) {
	at := assert.New(t)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		app := fiber.New()

		app.Post("/", func(c *fiber.Ctx) error {
			type User struct {
				Username string `validate:"required" json:"username"`
			}

			var u User
			if err := ValidateBody(c, &u); err != nil {
				return err
			}

			return c.SendString(u.Username)
		})

		assertRespCase(t, respCase{
			app:        app,
			method:     fiber.MethodPost,
			target:     "/?username=kiyon",
			reqBody:    bytes.NewReader([]byte("username=kiyon")),
			statusCode: fiber.StatusOK,
			respBody:   `kiyon`,
		})
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		c := fiber.New().AcquireCtx(&fasthttp.RequestCtx{})
		at.NotNil(ValidateBody(c, nil))
	})
}

func Test_Fiberx_ValidateQuery(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		app := fiber.New()

		app.Get("/", func(c *fiber.Ctx) error {
			type User struct {
				Username string `validate:"required" json:"username"`
			}

			var u User
			if err := ValidateQuery(c, &u); err != nil {
				return err
			}

			return c.SendString(u.Username)
		})

		assertRespCase(t, respCase{
			app:        app,
			method:     fiber.MethodGet,
			target:     "/?username=kiyon",
			statusCode: fiber.StatusOK,
			respBody:   `kiyon`,
		})
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		at := assert.New(t)

		fctx := &fasthttp.RequestCtx{}
		fctx.Request.URI().SetQueryString("a=b")
		c := fiber.New().AcquireCtx(fctx)
		at.NotNil(ValidateQuery(c, nil))
	})
}

func Test_Fiberx_2xx(t *testing.T) {
	tt := []struct {
		code int
		fn   func(c *fiber.Ctx, msg ...string) error
	}{
		{fiber.StatusOK, RespOK},
		{fiber.StatusCreated, RespCreated},
		{fiber.StatusAccepted, RespAccepted},
		{fiber.StatusNonAuthoritativeInformation, RespNonAuthoritativeInformation},
		{fiber.StatusNoContent, RespNoContent},
		{fiber.StatusResetContent, RespResetContent},
		{fiber.StatusPartialContent, RespPartialContent},
		{fiber.StatusMultiStatus, RespMultiStatus},
		{fiber.StatusAlreadyReported, RespAlreadyReported},
	}

	for _, tc := range tt {
		t.Run(strconv.Itoa(tc.code), func(t *testing.T) {
			t.Parallel()

			fn := tc.fn
			app := fiber.New()
			app.Get("/", func(c *fiber.Ctx) error {
				if tc.code == fiber.StatusNoContent {
					return fn(c, "I will be removed")
				}
				return fn(c)
			})

			c := respCase{
				app:        app,
				method:     fiber.MethodGet,
				target:     "/",
				statusCode: tc.code,
				respBody:   fmt.Sprintf("{\"code\":%d,\"message\":\"%s\"}", tc.code, utils.StatusMessage(tc.code)),
			}

			if tc.code == fiber.StatusNoContent {
				c.respBody = ""
			}

			assertRespCase(t, c)
		})
	}
}

func Test_Fiberx_RequestID(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderXRequestID, "id")
		return RespOK(c)
	})

	assertRespCase(t, respCase{
		app:        app,
		method:     fiber.MethodGet,
		target:     "/",
		statusCode: fiber.StatusOK,
		respBody:   `{"code":200,"message":"OK","request_id":"id"}`,
	})
}

func Test_Fiberx_Logger(t *testing.T) {
	t.Parallel()

	app := fiber.New(fiber.Config{ErrorHandler: ErrHandler})
	app.Use(Logger())

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderXRequestID, "id")
		return RespOK(c)
	})

	assertRespCase(t, respCase{
		app:        app,
		method:     fiber.MethodGet,
		target:     "/",
		statusCode: fiber.StatusOK,
		respBody:   `{"code":200,"message":"OK","request_id":"id"}`,
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.ErrForbidden
	})

	assertRespCase(t, respCase{
		app:        app,
		method:     fiber.MethodGet,
		target:     "/error",
		statusCode: fiber.StatusForbidden,
		respBody:   `{"code":403,"message":"Forbidden"}`,
	})
}

func Test_Fiberx_Response_Message(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return Message(c, "message")
	})

	assertRespCase(t, respCase{
		app:        app,
		method:     fiber.MethodGet,
		target:     "/",
		statusCode: fiber.StatusOK,
		respBody:   `{"code":200,"message":"message"}`,
	})
}

func Test_Fiberx_Response_Data(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return Data(c, []string{"data1", "data2"})
	})

	assertRespCase(t, respCase{
		app:        app,
		method:     fiber.MethodGet,
		target:     "/",
		statusCode: fiber.StatusOK,
		respBody:   `{"code":200,"data":["data1","data2"]}`,
	})
}

func assertRespCase(t *testing.T, c respCase) {
	t.Helper()

	at := assert.New(t)

	isJson := strings.HasPrefix(c.respBody, "{") && strings.HasSuffix(c.respBody, "}")

	res := httptest.NewRequest(c.method, c.target, c.reqBody)
	if c.method == fiber.MethodPost {
		res.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationForm)
	}

	resp, err := c.app.Test(res)
	at.Nil(err)
	at.Equal(c.statusCode, resp.StatusCode)
	if isJson {
		at.Equal(fiber.MIMEApplicationJSON, resp.Header.Get(fiber.HeaderContentType))
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	at.Nil(err)
	if isJson {
		at.JSONEq(c.respBody, string(body))
	} else {
		at.Equal(c.respBody, string(body))
	}
}
