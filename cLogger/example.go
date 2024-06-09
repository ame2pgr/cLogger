package cLogger

/*
package main

import (
	"Service/cLogger"
	"context"
	"errors"
	"time"
)

func main() {

	type b struct {
		a int
	}
	type c struct {
		a int
		b *int
	}

	type a struct {
		l int
		h int
		b b
		c *c
	}
	four := 4

	l, err := cLogger.New(context.TODO(), cLogger.LogConfig{
		TgConfig: cLogger.TgConfig{
			TgOn:     true,
			TgLevel:  "trace",
			TgToken:  "7454253821:AAFhO9NWM7seEiKltTf2ZEPceJMaPuxlzMk",
			TgChatId: -1002247246509,
		},
		SystemConfig: cLogger.SystemConfig{
			ServiceName: "test",
		},
	})

	if err != nil {
		panic(err)
	}
	l.Trace("trace 1")
	l.Debug("debug 1")
	l.Info("info 1")
	l.Warn("warn 1")
	l.Trace(a{l: 333, h: 333})
	l.Error(errors.New("error 1"), "1212", "121212")
	l.Error(errors.New("error 2"))

	lHash := l.FireTrace("main func test")
	l.Point(lHash, "1", "test")
	time.Sleep(time.Second)
	l.Point(lHash, "2", 223)
	l.Point(lHash, "3", map[string]string{"a": "b", "c": "d"}, errors.New("2323"), a{l: 1, h: 2, b: b{a: 32}, c: &c{a: 3, b: &four}})
	l.Point(lHash, "4", errors.New("ewwewewe"))
	time.Sleep(time.Second)
	l.Point(lHash, "5")
	l.Show(lHash)

	time.Sleep(time.Second * 10)
}
*/

/*
2024-06-09T02:51:16+03:00 [TRACE] usecase.go:122 | trace 1 | go_version=go1.21.5 pid=47891 type=string
2024-06-09T02:51:16+03:00 [DEBUG] usecase.go:127 | debug 1 | pid=47891 type=string
2024-06-09T02:51:16+03:00 [INFO] usecase.go:132 | info 1 | pid=47891 type=string
2024-06-09T02:51:16+03:00 [WARN] usecase.go:137 | warn 1 | pid=47891 type=string
2024-06-09T02:51:16+03:00 [TRACE] usecase.go:122 | {l:333 h:333 b:{a:0} c:<nil>} | go_version=go1.21.5 pid=47891 type=main.a
2024-06-09T02:51:16+03:00 [ERROR] usecase.go:142 | error 1; addDescription = [1212; 121212] | pid=47891 type=*errors.errorString
2024-06-09T02:51:16+03:00 [ERROR] usecase.go:142 | error 2 | pid=47891 type=*errors.errorString
2024-06-09T02:51:18+03:00 [TRACE] main.go:59 | main func test [2024-06-09 02:51:16.109344 +0300 MSK]
        - 1  {(string) (len=4) "test"
}  [0 ms]
        - 2  {(int) 223
}  [1001 ms]
        - 3  {(map[string]string) (len=2) {
 (string) (len=1) "a": (string) (len=1) "b",
 (string) (len=1) "c": (string) (len=1) "d"
}
}  {(*errors.errorString)(0x140003100c0)(2323)
}  {(main.a) {
 l: (int) 1,
 h: (int) 2,
 b: (main.b) {
  a: (int) 32
 },
 c: (*main.c)(0x140003100d0)({
  a: (int) 3,
  b: (*int)(0x140000b1688)(4)
 })
}
}  [1001 ms]
        - 4  {(*errors.errorString)(0x140003100e0)(ewwewewe)
}  [1001 ms]
        - 5  [2003 ms]
 | go_version=go1.21.5 pid=47891
*/
