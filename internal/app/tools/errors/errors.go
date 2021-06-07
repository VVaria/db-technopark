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

type Success struct {
	Description string `json:"description"`
	Data    string `json:"data"`
}

func JSONError(error *Error) []byte {
	jsonError, err := json.Marshal(error)
	if err != nil {
		return []byte("")
	}
	return jsonError
}

func JSONMessage(message ...interface{}) []byte {
	if len(message) > 1 {
		jsonSucc, err := json.Marshal(message[1])
		if err != nil {
			return []byte("")
		}
		return jsonSucc
	}
	jsonSucc, err := json.Marshal(Success{Description: message[0].(string)})
	if err != nil {
		return []byte("")
	}
	return jsonSucc
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
		Message:   "Пользователь отсутсвует в системе.",
	},
	ForumNotExist: {
		ErrorCode: ForumNotExist,
		HttpError: http.StatusNotFound,
		Message:   "Форум отсутсвует в системе.",
	},
	ThreadNotExist: {
		ErrorCode: ThreadNotExist,
		HttpError: http.StatusNotFound,
		Message:   "Ветка обсуждения отсутсвует в форуме.",
	},
	PostNotExist: {
		ErrorCode: PostNotExist,
		HttpError: http.StatusNotFound,
		Message:   "Сообщение отсутсвует в форуме.",
	},
	UserCreateConflict: {
		ErrorCode: UserCreateConflict,
		HttpError: http.StatusConflict,
		Message:   "Пользователь уже присутсвует в базе данных.",
	},
	UserProfileConflict: {
		ErrorCode: UserProfileConflict,
		HttpError: http.StatusConflict,
		Message:   "Новые данные профиля пользователя конфликтуют с имеющимися пользователями.",
	},
	ForumCreateConflict: {
		ErrorCode: ForumCreateConflict,
		HttpError: http.StatusConflict,
		Message:   "Форум уже присутсвует в базе данных.",
	},
	ForumCreateThreadConflict: {
		ErrorCode: ForumCreateThreadConflict,
		HttpError: http.StatusConflict,
		Message:   "Ветка обсуждения уже присутсвует в базе данных.",
	},
	PostWrongThread: {
		ErrorCode: PostWrongThread,
		HttpError: http.StatusConflict,
		Message:   "Хотя бы один родительский пост отсутсвует в текущей ветке обсуждения",
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
