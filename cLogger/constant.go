package cLogger

import (
	"time"

	"github.com/rs/zerolog"
)

//используется для вывода в консоль места вызова логгера во внутренней 
//бизнес логике сервиса.
const callerWithSkipFrameCount int = 3

//данные константы конфигурируют правило поведения лога, пример:
//если за 1 секунду будет запрос на 4 или более логов одного и того же уровня, то
//в течение следующей секунды в логах будет показываться каждое второе сообщение
const (
	bellowErrorLevelPeriod        time.Duration = 1 * time.Second
	bellowErrorLevelBurst         uint32        = 4
	bellowErrorLevelNextNMessages uint32        = 2
)

var logLevelMap map[string]zerolog.Level = map[string]zerolog.Level{
	zerolog.TraceLevel.String(): zerolog.TraceLevel,
	zerolog.DebugLevel.String(): zerolog.DebugLevel,
	zerolog.InfoLevel.String():  zerolog.InfoLevel,
	zerolog.WarnLevel.String():  zerolog.WarnLevel,
	zerolog.ErrorLevel.String(): zerolog.ErrorLevel,
	zerolog.FatalLevel.String(): zerolog.FatalLevel,
	zerolog.PanicLevel.String(): zerolog.PanicLevel,
	zerolog.NoLevel.String():    zerolog.NoLevel,
	zerolog.Disabled.String():   zerolog.Disabled,
}

var levelEmoji map[zerolog.Level]string = map[zerolog.Level]string{
	zerolog.InfoLevel:  "🟢",
	zerolog.WarnLevel:  "🟠",
	zerolog.ErrorLevel: "🔴",
	zerolog.FatalLevel: "🚨",
	zerolog.PanicLevel: "🚨",
}
