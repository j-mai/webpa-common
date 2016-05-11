package context

import (
	"encoding/base64"
	"fmt"
	"github.com/Comcast/webpa-common/canonical"
	"github.com/Comcast/webpa-common/convey"
	"github.com/Comcast/webpa-common/logging"
	"net/http"
)

// ConveyPayload is a map type which stores the optional, deserialized convey value
type ConveyPayload map[string]interface{}

// Context is the core type of this package.
type Context interface {
	// logging.Logger returns the contextual logging.Logger.  It is never nil.
	Logger() logging.Logger

	// DeviceId returns the canonical device id associated with the request.
	DeviceId() canonical.Id

	// ConveyPayload returns the optional payload of a convey header
	ConveyPayload() convey.Payload
}

// defaultContext is the default implementation of Context
type defaultContext struct {
	logger        logging.Logger
	deviceId      canonical.Id
	conveyPayload convey.Payload
}

func (c *defaultContext) Logger() logging.Logger {
	return c.logger
}

func (c *defaultContext) DeviceId() canonical.Id {
	return c.deviceId
}

func (c *defaultContext) ConveyPayload() convey.Payload {
	return c.conveyPayload
}

// NewContext creates a new Context object from an HTTP request
func NewContext(logger logging.Logger, request *http.Request) (Context, error) {
	deviceName := request.Header.Get(DeviceNameHeader)
	if len(deviceName) == 0 {
		return nil, missingDeviceNameError
	}

	deviceId, err := canonical.ParseId(deviceName)
	if err != nil {
		return nil, NewHttpError(
			http.StatusBadRequest,
			fmt.Sprintf(InvalidDeviceNameHeaderPattern, deviceName),
		)
	}

	var conveyPayload convey.Payload
	rawConveyValue := request.Header.Get(ConveyHeader)
	if len(rawConveyValue) > 0 {
		// BUG: https://www.teamccp.com/jira/browse/WEBPA-787
		const notAvailable string = "not-available"
		if rawConveyValue == notAvailable {
			logger.Error("Invalid convey header: %s.  FIX ME: https://www.teamccp.com/jira/browse/WEBPA-787", notAvailable)
		} else if conveyPayload, err = convey.ParsePayload(base64.StdEncoding, rawConveyValue); err != nil {
			return nil, &httpError{
				http.StatusBadRequest,
				fmt.Sprintf(InvalidConveyPattern, rawConveyValue, err),
			}
		}
	}

	return &defaultContext{
		logger:        logger,
		deviceId:      deviceId,
		conveyPayload: conveyPayload,
	}, nil
}
