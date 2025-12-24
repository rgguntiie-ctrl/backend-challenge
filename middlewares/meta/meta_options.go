package meta

type MetaErrorOptions struct {
	httpStatus int
}

func getDefaultMetaErrorOptions() MetaErrorOptions {
	return MetaErrorOptions{
		httpStatus: 500,
	}
}

// WithMetaErrorOptionsHttpStatus sets the HTTP status code for MetaErrorOptions.
//
// Parameters:
// - status: An integer representing the HTTP status code to be set.
//
// Returns:
// - A function that takes a pointer to MetaErrorOptions and sets its httpStatus field to the provided status.
func WithMetaErrorOptionsHttpStatus(status int) func(*MetaErrorOptions) {
	return func(options *MetaErrorOptions) {
		options.httpStatus = status
	}
}

type MetaErrorHandlerOptions struct {
	isLog bool
}

func getDefaultMetaErrorHandlerOptions() MetaErrorHandlerOptions {
	return MetaErrorHandlerOptions{
		isLog: false,
	}
}

// WithMetaErrorHandlerOptionsLogging enables logging for the given MetaErrorHandlerOptions.
// It sets the isLog field of the options to true.
//
// Parameters:
//
//	options (*MetaErrorHandlerOptions): The options for which logging should be enabled.
func WithMetaErrorHandlerOptionsLogging(options *MetaErrorHandlerOptions) {
	options.isLog = true
}

type MetaOKOptions struct {
	pagination *MetaPagination
}

func getDefaultMetaOKOptions() MetaOKOptions {
	return MetaOKOptions{
		pagination: nil,
	}
}

// WithMetaOKOptionsPagination sets the pagination options for MetaOKOptions.
// It takes a MetaPagination object and returns a function that modifies
// the pagination field of a MetaOKOptions instance.
//
// Parameters:
//
//	m (MetaPagination): The pagination settings to be applied.
//
// Returns:
//
//	func(*MetaOKOptions): A function that sets the pagination field of a MetaOKOptions instance.
func WithMetaOKOptionsPagination(m MetaPagination) func(*MetaOKOptions) {
	return func(options *MetaOKOptions) {
		options.pagination = &m
	}
}
