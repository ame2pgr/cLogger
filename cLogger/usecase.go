package cLogger

import (
	"Service/cTime"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gopkg.in/tucnak/telebot.v2"
)

func New(ctx context.Context, lC LogConfig) (l Logger, err error) {

	location, err := cTime.MskLocation()
	if err != nil {
		return
	}
	if location == nil {
		err = errors.New("location from cTime is not ok")
		return
	}

	bellowErrorSampler := &zerolog.BurstSampler{
		Burst:  bellowErrorLevelBurst,
		Period: bellowErrorLevelPeriod,
		NextSampler: &zerolog.BasicSampler{
			N: bellowErrorLevelNextNMessages,
		},
	}

	writer := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{
		Out:          os.Stdout,
		TimeFormat:   time.RFC3339,
		TimeLocation: location,
		FormatLevel: func(i any) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
		FormatMessage: func(i any) string {
			return fmt.Sprintf("| %s |", i)
		},
		FormatCaller: func(i any) string {
			return filepath.Base(fmt.Sprintf("%s", i))
		},
	})

	var tgBot *telebot.Bot = &telebot.Bot{}
	if lC.TgConfig.TgOn {

		if lC.TgConfig.TgChatId == 0 {
			err = errors.New("tgChatId is not ok")
			return
		}
		if lC.TgConfig.TgToken == "" {
			err = errors.New("tgToken is not ok")
			return
		}
		if !validateLogLeve(lC.TgConfig.TgLevel) {
			err = errors.New("tgLevel is not ok")
			return
		}

		tgBot, err = telebot.NewBot(telebot.Settings{
			Token: lC.TgConfig.TgToken,
			Poller: &telebot.LongPoller{
				Timeout: 10 * time.Second,
			},
		})
		if err != nil {
			return
		}
	}

	zLogger := zerolog.New(writer).
		With().
		Timestamp().
		CallerWithSkipFrameCount(callerWithSkipFrameCount).
		Int("pid", os.Getpid()).
		Logger().
		Sample(zerolog.LevelSampler{
			TraceSampler: bellowErrorSampler,
			DebugSampler: bellowErrorSampler,
			InfoSampler:  bellowErrorSampler,
			WarnSampler:  bellowErrorSampler,
		})

	l = &logger{
		z:         zLogger,
		goVersion: runtime.Version(),

		logChangeMutex: new(sync.RWMutex),

		tgOn:     lC.TgConfig.TgOn,
		tgLevel:  logLevelId(lC.TgConfig.TgLevel),
		tgChatId: lC.TgConfig.TgChatId,
		tgBot:    tgBot,

		serviceName: lC.SystemConfig.ServiceName,

		traceMutex: new(sync.RWMutex),
		traceMap:   make(map[string]Tracer),
	}

	l.Hook()
	return l, err
}

func (l *logger) Hook() {
	l.z = l.z.Hook(l)
}

func (l *logger) Trace(input any, addDescription ...string) {
	event := l.z.Trace().Str("go_version", l.goVersion)
	l.eventTypeHandler(event, input, addDescription)
}

func (l *logger) Debug(input any, addDescription ...string) {
	event := l.z.Debug()
	l.eventTypeHandler(event, input, addDescription)
}

func (l *logger) Info(input any, addDescription ...string) {
	event := l.z.Info()
	l.eventTypeHandler(event, input, addDescription)
}

func (l *logger) Warn(input any, addDescription ...string) {
	event := l.z.Warn()
	l.eventTypeHandler(event, input, addDescription)
}

func (l *logger) Error(err error, addDescription ...string) {
	event := l.z.Error().Stack()
	l.eventTypeHandler(event, err, addDescription)
}

func (l *logger) Fatal(err error, addDescription ...string) {
	event := l.z.Fatal().Stack()
	l.eventTypeHandler(event, err, addDescription)
}

func (l *logger) Panic(input any) {
	event := l.z.Panic().Stack()
	l.eventTypeHandler(event, input, []string{""})
}

func (l *logger) HttpInfo(p HttpInfo) {

	event := l.z.Info()

	if p.Method != "" {
		event.Str("method", p.Method)
	}
	if p.Path != "" {
		event.Str("url", p.Path)
	}
	if p.Body != "" {
		event.Str("body", p.Body)
	}

	if p.CorrelationId != "" {
		event.Str("correlation_id", p.CorrelationId)
	}

	if p.UserAgent != "" {
		event.Str("user_agent", p.UserAgent)
	}
	if p.Access != "" {
		p.Access = l.simpleSecret(p.Access)
		event.Str("access", p.Access)
	}

	if p.Public != "" {
		p.Public = l.simpleSecret(p.Public)
		event.Str("public", p.Public)
	}
	if p.Signature != "" {
		p.Signature = l.simpleSecret(p.Signature)
		event.Str("signature", p.Signature)
	}

	event.Dur("time (ms)", time.Duration(p.RequestTime.Milliseconds())).
		Msg(fmt.Sprintf("addDescription=[%s]", p.AddDescription))
}

func (l *logger) FireTrace(title string) string {

	hash := uuid.NewString()
	tracer := Tracer{
		Start:  cTime.Msk(),
		Title:  title,
		Points: []Point{},
	}

	l.setToTracerMap(hash, tracer)
	return hash
}

func (l *logger) Point(hash, description string, values ...any) (err error) {

	tracer, err := l.fromTracerMap(hash)
	if err != nil {
		return err
	}

	tracer.Points = append(tracer.Points, Point{
		Description:   description,
		Duration:      cTime.Since(tracer.Start),
		StructToParse: values,
	})

	l.setToTracerMap(hash, tracer)
	return nil
}

func (l *logger) Show(hash string) error {

	tracer, err := l.fromTracerMap(hash)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("%s [%s]\n", tracer.Title, tracer.Start.String())
	for _, point := range tracer.Points {
		message = fmt.Sprintf("%s\t- %s ", message, point.Description)
		for _, str := range point.StructToParse {
			message = fmt.Sprintf("%s {%s} ", message, spew.Sdump(str))
		}
		message = fmt.Sprintf("%s [%d ms]\n", message, point.Duration.Milliseconds())
	}

	l.z.Trace().Caller().
		Str("go_version", l.goVersion).
		Msg(message)
	l.deleteFromTracerMap(hash)
	return nil
}

func (l *logger) fromTracerMap(hash string) (tracer Tracer, err error) {

	var ok bool

	if hash == "" {
		err = errors.New("hash is empty")
		return tracer, err
	}

	l.traceMutex.RLock()
	tracer, ok = l.traceMap[hash]
	l.traceMutex.RUnlock()

	if !ok {
		err = fmt.Errorf("hash [%s] not found", hash)
		return
	}

	return tracer, nil
}

func (l *logger) deleteFromTracerMap(hash string) {
	l.traceMutex.Lock()
	delete(l.traceMap, hash)
	l.traceMutex.Unlock()
}

func (l *logger) setToTracerMap(hash string, tracer Tracer) {
	l.traceMutex.Lock()
	l.traceMap[hash] = tracer
	l.traceMutex.Unlock()
}

func (l *logger) eventTypeHandler(event *zerolog.Event, input any, addDescription []string) {
	var (
		addDescriptionString string = strings.Join(addDescription, "; ")
		message              string
	)

	switch iType := input.(type) {
	case string:
		message = l.messageFormater(input.(string), addDescriptionString)
		event.
			Type("type", iType).
			Msg(message)
	case error:
		message = l.messageFormater(input.(error).Error(), addDescriptionString)
		event.
			Type("type", iType).
			Msg(message)
	default:
		message = l.messageFormater(fmt.Sprintf("%+v", input), addDescriptionString)
		event.
			Type("type", iType).
			Msg(message)
	}
}

func (l *logger) messageFormater(message, addDescription string) string {

	if len(addDescription) == 0 {
		return message
	}

	return fmt.Sprintf("%s; addDescription = [%s]", message, addDescription)
}

func (l *logger) simpleSecret(value string) string {

	valueLen := len(value)
	if valueLen > 8 {
		return fmt.Sprintf("%s**%s", value[:2], value[valueLen-2:])
	}
	return "**"
}

func validateLogLeve(level string) bool {
	_, ok := logLevelMap[level]
	return ok
}

func logLevelId(level string) zerolog.Level {

	value, ok := logLevelMap[level]
	if !ok {
		return zerolog.Disabled
	}

	return value
}

func (l *logger) tgSend(message string, level zerolog.Level) (err error) {

	if l.tgBot == nil {
		err = errors.New("tgBot is nil")
		return
	}

	l.logChangeMutex.RLock()
	chatId := l.tgChatId
	l.logChangeMutex.RUnlock()

	emoji := levelEmoji[level]
	message = fmt.Sprintf("%[1]s %[2]s [%[3]s] %[2]s\n\n%[4]s", l.serviceName, emoji, strings.ToUpper(level.String()), message)
	_, err = l.tgBot.Send(telebot.ChatID(chatId), message)
	return
}

func (l *logger) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	go func() {
		l.logChangeMutex.RLock()
		tgLevel := l.tgLevel
		l.logChangeMutex.RUnlock()

		if level >= tgLevel {
			if err := l.tgSend(msg, level); err != nil {
				l.z.Error().CallerSkipFrame(callerWithSkipFrameCount - 4).Err(err).Send()
			}
		}
	}()
}

func (l *logger) UpdateTgLevel(level zerolog.Level) {
	l.logChangeMutex.Lock()
	l.tgLevel = level
	l.logChangeMutex.Unlock()
}

func (l *logger) UpdateTgChatId(chatId int64) {
	l.logChangeMutex.Lock()
	l.tgChatId = chatId
	l.logChangeMutex.Unlock()
}
