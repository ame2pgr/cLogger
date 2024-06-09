package cLogger

import "github.com/rs/zerolog"

type Logger interface {
	Trace(input any, addDescription ...string) // -1
	Debug(input any, addDescription ...string) // 0
	Info(input any, addDescription ...string)  // 1
	Warn(input any, addDescription ...string)  // 2
	Error(err error, addDescription ...string) // 3
	Fatal(err error, addDescription ...string) // 4
	Panic(input any)                           // 5

	//специальная ручка, парсящая http запрос и выводящая Info лог
	HttpInfo(p HttpInfo)

	//имплементация трейсера функции в логгер. Может парсить вложенные указатели как
	//значение переменной. При использовании парсинга переменных не рекомендуется
	//использовать в продакшн версии где скорость выполнения функции играет ключевую роль.
	//Результатом функции будет строка (хэш-ключ). Этот ключ нужно будет передавать в точки
	//трейсера и для получения финального трейса в консоли
	FireTrace(title string) string

	//точка трейсера. Передаётся хэш, полученный при инициализации трейсера, описание
	//точки. Дополнительным полем можно передать набор интересующих переменных. Стейс этих
	//переменных будет записан в точке и в момент формирования общего лога трейса в консоль
	//каждый элемент вне зависимости от вложенности и на каком уровне указатель. Ошибкой может
	//быть только отсутствие хэш-ключа в реестре трейсеров*/
	Point(hash, description string, values ...any) (err error)

	//функция, генерирующая лог трейсера. Принимает хэш трейсера. Ошибкой может быть только
	//отсутствие хэш-ключа в реестре трейсеров
	Show(hash string) (err error)

	//хук функция нужна для добавления телеграм оповещений в zerolog
	Hook()

	UpdateTgLevel(level zerolog.Level)
	UpdateTgChatId(chatId int64)
}
