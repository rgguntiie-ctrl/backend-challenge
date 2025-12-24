package meta

type MetaOK struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	// Detailed metadata response
	Meta *MetaPagination `json:"meta"`
}

type MetaError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	httpStatus int
	err        error
}

type MetaPagination struct {
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
}

// AppendError appends the provided error to the MetaError instance.
//
// Parameters:
// - err: The error to be appended.
func (e *MetaError) AppendError(err error) {
	e.err = err
}

// HttpStatus returns the HTTP status code associated with the MetaError.
// It provides a way to retrieve the specific HTTP status that should be
// used in response to the error.
//
// Returns:
//
//	int: The HTTP status code.
func (e *MetaError) HttpStatus() int {
	return e.httpStatus
}

// Error returns the error message contained in the MetaError.
// It implements the error interface.
func (e *MetaError) Error() string {
	return e.Message
}

// Unwrap returns the underlying error wrapped by MetaError.
// This method allows MetaError to be compatible with the standard library's error unwrapping mechanism.
func (e *MetaError) Unwrap() error {
	return e.err
}

// NewMetaError creates a new MetaError with the specified code and message.
// Additional options can be provided using optionFuncs to customize the MetaError.
func NewMetaError(code int, message string, optionFuncs ...func(*MetaErrorOptions)) *MetaError {
	defaultOptions := getDefaultMetaErrorOptions()
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &MetaError{
		Code:       code,
		Message:    message,
		httpStatus: options.httpStatus,
	}
}

func IsMetaError(err error) (*MetaError, bool) {
	m, ok := err.(*MetaError)
	return m, ok
}

// NewMetaOK creates a new instance of MetaSuccess with a success code, message, and data.
//
// Parameters:
//   - message: A string containing the success message.
//   - data: An interface{} containing the data to be included in the response.
//
// Returns:
//   - *MetaSuccess: A pointer to a MetaSuccess struct initialized with the provided code, message, and data.
func NewMetaOK(message string, data any, optionFuncs ...func(*MetaOKOptions)) *MetaOK {
	defaultOptions := getDefaultMetaOKOptions()
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	return &MetaOK{
		Code:    0,
		Message: message,
		Data:    data,
		Meta:    options.pagination,
	}
}
