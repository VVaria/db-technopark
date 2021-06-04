package errors

import (
	"encoding/json"
	"net/http"
)

type ErrorType uint8

const (
	WrongErrorCode ErrorType = iota + 1
	InternalError
	UserNotExist
	ForumNotExist
	ThreadNotExist
	PostNotExist
	UserCreateConflict
	UserProfileConflict
	ForumCreateConflict
	ForumCreateThreadConflict
	PostWrongThread
)

type Error struct {
	ErrorCode ErrorType `json:"-"`
	HttpError int       `json:"-"`
	Message   string    `json:"message"`
}

func JSONError(error *Error, w http.ResponseWriter) {
	jsonError, err := json.Marshal(error)
	if err != nil {
		return
	}
	w.WriteHeader(error.HttpError)
	w.Write(jsonError)
}

func JSONSuccess(status int, data interface{}, w http.ResponseWriter) {
	jsonSuccess, err := json.Marshal(data)
	if err != nil {
		return
	}

	w.WriteHeader(status)
	w.Write(jsonSuccess)
}

var CustomErrors = map[ErrorType]*Error{
	WrongErrorCode: {
		ErrorCode: WrongErrorCode,
		HttpError: http.StatusInternalServerError,
	},
	InternalError: {
		ErrorCode: InternalError,
		HttpError: http.StatusInternalServerError,
		Message:   "something wrong",
	},
	UserNotExist: {
		ErrorCode: UserNotExist,
		HttpError: http.StatusNotFound,
		Message:   "user doesn't exist\n",
	},
	ForumNotExist: {
		ErrorCode: ForumNotExist,
		HttpError: http.StatusNotFound,
		Message:   "forum doesn't exist\n",
	},
	ThreadNotExist: {
		ErrorCode: ThreadNotExist,
		HttpError: http.StatusNotFound,
		Message:   "thread doesn't exist\n",
	},
	PostNotExist: {
		ErrorCode: PostNotExist,
		HttpError: http.StatusNotFound,
		Message:   "post doesn't exist\n",
	},
	UserCreateConflict: {
		ErrorCode: UserCreateConflict,
		HttpError: http.StatusConflict,
		Message:   "conflicts with user creation\n",
	},
	UserProfileConflict: {
		ErrorCode: UserProfileConflict,
		HttpError: http.StatusConflict,
		Message:   "conflicts with user profile\n",
	},
	ForumCreateConflict: {
		ErrorCode: ForumCreateConflict,
		HttpError: http.StatusConflict,
		Message:   "conflicts with forum creation\n",
	},
	ForumCreateThreadConflict: {
		ErrorCode: ForumCreateThreadConflict,
		HttpError: http.StatusConflict,
		Message:   "conflicts with thread creation\n",
	},
	PostWrongThread: {
		ErrorCode: PostWrongThread,
		HttpError: http.StatusConflict,
		Message:   "Parent post was created in another thread\n",
	},
}

func Cause(code ErrorType) *Error {
	err, ok := CustomErrors[code]
	if !ok {
		return CustomErrors[WrongErrorCode]
	}

	return err
}

func UnexpectedInternal(err error) *Error {
	unexpErr := CustomErrors[InternalError]
	unexpErr.Message = err.Error()

	return unexpErr
}
