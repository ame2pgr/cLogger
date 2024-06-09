package cLogger

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/tucnak/telebot.v2"
)

type LogConfig struct {
	//настройка подключения телеграма
	TgConfig TgConfig
	//другие настройки
	SystemConfig SystemConfig
}

type TgConfig struct {
	//подключение телеграма
	TgOn bool
	//определяет минимальный (включительно) уровень лога в телеграм
	TgLevel string

	//токен телеграм бота
	TgToken string
	//идентификатор чата в телеграм
	TgChatId int64
}

type SystemConfig struct {
	//имя сервиса (используется в телеграме для отображения исходного сервиса отправки лога)
	ServiceName string
}

type logger struct {
	z         zerolog.Logger
	goVersion string

	logChangeMutex *sync.RWMutex

	tgOn     bool
	tgLevel  zerolog.Level
	tgBot    *telebot.Bot
	tgChatId int64

	serviceName string

	traceMutex *sync.RWMutex
	traceMap   map[string]Tracer
}

type HttpInfo struct {
	Method string
	Path   string
	Body   string

	CorrelationId string

	UserAgent string
	Access    string

	Public    string
	Signature string

	RequestTime time.Duration

	AddDescription string
}

type Tracer struct {
	//время начала отсчёта трейсера. Используется для вычисления времени
	//испольнения точек трейсера и общего времени
	Start time.Time
	//заголовок трейсера
	Title string
	//набор точек трейсера. Включает в себя описание и набор интерфейсов -
	//структур, данные которых будут парситься в лог трейсера
	Points []Point
}

type Point struct {
	Description   string
	Duration      time.Duration
	StructToParse []any
}
