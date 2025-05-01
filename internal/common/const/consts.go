package consts

const (
	ErrnoSuccess      = 0
	ErrnoUnknownError = 1

	ErrnoBindRequestError          = 1000
	ErrnoRequestValidateError      = 1001
	ErrnoInternalServerError       = 500
	ErrnoFailedConnectionError     = 504
	ErrnoParameterInputError       = 424
	ErrnoResourceNotFoundException = 404
	ErrnoDataParseException        = 405

	ErrnoUnmarshalError = 422
	ErrnoCastError      = 423

	ErrnoCacheSetError = 510
	ErrnoCacheGetError = 511
	ErrnoCacheDelError = 512

	ErrnoUserTokenInvalid      = 602
	ErrnoUserNotFoundError     = 604
	ErrnoUserModifyFailedError = 605
)

var ErrMsg = map[int]string{
	ErrnoSuccess:      "success",
	ErrnoUnknownError: "unknown error",

	ErrnoBindRequestError:      "bind request error",
	ErrnoRequestValidateError:  "validate request error",
	ErrnoUnmarshalError:        "unmarshal error",
	ErrnoCastError:             "cast error",
	ErrnoFailedConnectionError: "failed connection error",

	ErrnoUserNotFoundError:         "user not found",
	ErrnoParameterInputError:       "parameter input error",
	ErrnoResourceNotFoundException: "resource not found",
	ErrnoInternalServerError:       "internal server error",
	ErrnoUserTokenInvalid:          "token invalid",
	ErrnoUserModifyFailedError:     "user modify failed",
	ErrnoDataParseException:        "data parse exception",
	ErrnoCacheSetError:             "cache set error",
	ErrnoCacheGetError:             "cache get error",
	ErrnoCacheDelError:             "cache del error",
}
