package handler

// error number
const (
	ERRNO_OK               = 0
	ERRNO_INVALID_ARGUMENT = 4000
	ERRNO_INTERNAL_ERROR   = 5000
)

// error str
const (
	ERRSTR_OK               = "ok"
	ERRSTR_INVALID_ARGUMENT = "invalid parameter"
	ERRSTR_INTERNAL_ERROR    = "internal error"
)

var (
	errno2str = map[int]string{
		ERRNO_OK:               ERRSTR_OK,
		ERRNO_INVALID_ARGUMENT: ERRSTR_INVALID_ARGUMENT,
		ERRNO_INTERNAL_ERROR:   ERRSTR_INTERNAL_ERROR,
	}
)
