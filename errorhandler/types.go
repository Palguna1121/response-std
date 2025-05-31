package errorhandler

type NotFoundError struct {
	Message string `json:"message"`
}

type BadRequestError struct {
	Message string `json:"message"`
}

type InternalServerError struct {
	Message string `json:"message"`
}

type UnauthorizedError struct {
	Message string `json:"message"`
}

type ForbiddenError struct {
	Message string `json:"message"`
}

type UnprocessableEntityError struct {
	Message string `json:"message"`
}

type ConflictError struct {
	Message string `json:"message"`
}

type GoneError struct {
	Message string `json:"message"`
}

type PreconditionFailedError struct {
	Message string `json:"message"`
}

type RequestTimeoutError struct {
	Message string `json:"message"`
}

type TooManyRequestsError struct {
	Message string `json:"message"`
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func (e *UnprocessableEntityError) Error() string {
	return e.Message
}

func (e *ConflictError) Error() string {
	return e.Message
}

func (e *GoneError) Error() string {
	return e.Message
}

func (e *PreconditionFailedError) Error() string {
	return e.Message
}

func (e *RequestTimeoutError) Error() string {
	return e.Message
}

func (e *TooManyRequestsError) Error() string {
	return e.Message
}

