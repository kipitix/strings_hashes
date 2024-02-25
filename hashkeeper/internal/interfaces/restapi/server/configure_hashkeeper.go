// This file is safe to edit. Once it exists it will not be overwritten

package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"strconv"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"hashkeeper/internal/application"
	"hashkeeper/internal/domain/calculate"
	"hashkeeper/internal/domain/datahash"
	"hashkeeper/internal/domain/find"
	"hashkeeper/internal/interfaces/restapi/models"
	"hashkeeper/internal/interfaces/restapi/server/operations"
	"hashkeeper/pkg/hashlog"

	"github.com/sirupsen/logrus"
)

//go:generate swagger generate server --target ../../restapi --name Hashkeeper --spec ../../../../pkg/interfaces/restapi/spec/hashkeeper.yaml --server-package server --principal interface{}

type RESTAPICfg struct {
	RESTAPISendRequestTimeout  string `arg:"--rest-api-send-request-timeout,env:REST_API_SEND_REQUEST_TIMEOUT" default:"60s" help:"Send request maximum execution time"`
	RESTAPICheckRequestTimeout string `arg:"--rest-api-check-request-timeout,env:REST_API_CHECK_REQUEST_TIMEOUT" default:"10s" help:"Check request maximum execution time"`
}

var _server struct {
	app                 application.App
	sendRequestTimeout  time.Duration
	checkRequestTimeout time.Duration
}

func InitServer(app application.App, cfg RESTAPICfg) error {
	_server.app = app

	if dur, err := time.ParseDuration(cfg.RESTAPISendRequestTimeout); err != nil {
		return hashlog.WithStackErrorf("wrong send timeout: %w", err)
	} else {
		_server.sendRequestTimeout = dur
	}

	if dur, err := time.ParseDuration(cfg.RESTAPICheckRequestTimeout); err != nil {
		return hashlog.WithStackErrorf("wrong check timeout: %w", err)
	} else {
		_server.checkRequestTimeout = dur
	}

	return nil
}

func configureFlags(api *operations.HashkeeperAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.HashkeeperAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.Logger = logrus.Tracef

	api.UseSwaggerUI()

	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Restore hashes
	api.GetCheckHandler = operations.GetCheckHandlerFunc(func(params operations.GetCheckParams) middleware.Responder {
		// Prepare request
		findReqIn := make([]datahash.HashID, len(params.Ids))
		for i, inStrID := range params.Ids {
			inID, err := strconv.Atoi(inStrID)
			if err != nil {
				hashlog.LogErrorWithStack(err).Error("can`t parse input data")
				return operations.NewGetCheckBadRequest()
			}

			findReqIn[i] = datahash.HashID(inID)
		}
		findReq := find.NewFindRequest(findReqIn)

		ctxReqID := hashlog.AppendReqID(context.Background(), findReq.RequestID())

		hashlog.LogReqID(ctxReqID).Debug("new request received - find hashes")
		hashlog.LogReqID(ctxReqID).WithField("ids", findReq.IDs()).Trace("request data")

		// Make request for find hashes
		ctxFind, cancelFind := context.WithTimeout(ctxReqID, _server.checkRequestTimeout)
		defer cancelFind()

		findResp, err := _server.app.FindHandler().Find(ctxFind, findReq)
		if err != nil {
			hashlog.LogWithReqID(hashlog.LogErrorWithStack(err), ctxFind).Error("find failed")
			return operations.NewPostSendInternalServerError()
		}

		// Compose responses
		hashes := findResp.Hashes()
		resPayload := make(models.ArrayOfHash, len(hashes))
		for i, hash := range hashes {
			id := int64(hash.ID())
			str := string(hash.Hash())
			resPayload[i] = &models.Hash{
				ID:   &id,
				Hash: &str,
			}
		}

		hashlog.LogReqID(ctxReqID).WithField("hashes", findResp.Hashes()).Trace("response data")
		hashlog.LogReqID(ctxReqID).Debug("find hashes - request finished successfully")

		return operations.NewPostSendOK().WithPayload(resPayload)
	})

	// Calc and store hashes
	api.PostSendHandler = operations.PostSendHandlerFunc(func(params operations.PostSendParams) middleware.Responder {
		// Prepare request
		calcReqIn := make([]datahash.StringContent, len(params.Params))
		for i, inStr := range params.Params {
			calcReqIn[i] = datahash.StringContent(inStr)
		}
		calcReq := calculate.NewCalculateRequest(calcReqIn)

		ctxReqID := hashlog.AppendReqID(context.Background(), calcReq.RequestID())

		hashlog.LogReqID(ctxReqID).Debug("new request received - calculate and store hashes")
		hashlog.LogReqID(ctxReqID).WithField("strings", calcReq.Strings()).Trace("request data")

		// Make request for calculate and store hashes
		ctxCalc, cancelCalc := context.WithTimeout(ctxReqID, _server.sendRequestTimeout)
		defer cancelCalc()

		calcResp, err := _server.app.CalculateHandler().Calculate(ctxCalc, calcReq)
		if err != nil {
			hashlog.LogWithReqID(hashlog.LogErrorWithStack(err), ctxCalc).Error("calculate and store failed")
			return operations.NewPostSendInternalServerError()
		}

		// Compose response
		hashes := calcResp.Hashes()
		resPayload := make(models.ArrayOfHash, len(hashes))
		for i, hash := range hashes {
			id := int64(hash.ID())
			str := string(hash.Hash())
			resPayload[i] = &models.Hash{
				ID:   &id,
				Hash: &str,
			}
		}

		hashlog.LogReqID(ctxReqID).WithField("hashes", calcResp.Hashes()).Trace("response data")
		hashlog.LogReqID(ctxReqID).Debug("calculate and store hashes - request finished successfully")

		return operations.NewPostSendOK().WithPayload(resPayload)
	})

	api.PreServerShutdown = func() {
		_server.app.Shutdown()
	}

	api.ServerShutdown = func() {
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
