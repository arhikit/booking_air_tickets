# Cheatsheet по синтаксису GoLang

> Дополни меня, если нужно =)

## Hello World
Наша первая программа напечатает классическое сообщение “hello world”
```go
package main

import "fmt"

func main() {
    fmt.Println("hello world")
}
```

```go
// Для запуска программы, добавьте код в файл hello-world.go и выполните go run.
$ go run hello-world.go
> hello world

// Иногда необходимо собрать программу в бинарный файл. Мы можем сделать это 
// с помощью команды go build.
$ go build hello-world.go
$ ls
> hello-world    hello-world.go

//Мы можем выполнить бинарный файл напрямую.
$ ./hello-world
> hello world
```

## Пакеты
Объявление нового пакета используется ключевое слово package и имя пакета
```go
package main
```

## Импорты
Импорты нужны для использования существующих пакетов
```go
import "fmt"
```

Пример многострочного ипорта
```go
import (
    "fmt"
    "time"
)
```

Использование импортированного пакета в коде
```go
package main
import (
    "fmt"
    "time"
)

func main() {

    fmt.Println(time.Now())

}
```

## Область видимости
Для определения области видимости константы/функции или метода используется нижний или верхний регистр. Ограничение работает в рамках пакетов
```go
package helloworld

const Hello = "Hello"
const world = "world"
```
## Переменные
В Go, переменные объявляются явно и используются компилятором, например, для проверки корректного вызова функции (типы аргументов)
```go
package main

import "fmt"

func main() {

	// var объявляет 1 или более переменных
    var a = "initial"
    fmt.Println(a)

	// Можно объявить несколько переменных за раз
    var b, c int = 1, 2
    fmt.Println(b, c)

	// Тип определяется по инициализированной переменной
    var d = true
    fmt.Println(d)

	// Переменные без инициализации имеют нулевое значение соответствующего типа
    var e int
    fmt.Println(e)

	// Короткое объявление и инициализация переменной
    f := "apple"
    fmt.Println(f)
}
```

## Константы
Константы могут определяться в любом месте кода, как и переменные. Им могут быть присвоены строки, числа и логические значения
```go
package main

import (
    "fmt"
    "math"
)

const s string = "constant"

func main() {
    fmt.Println(s)

    const n = 500000000

    const d = 3e20 / n
    fmt.Println(d)

    fmt.Println(int64(d))

    fmt.Println(math.Sin(n))
}
```

## If
```go
    if num := 9; num < 0 {
    fmt.Println(num, "is negative")
    } else if num < 10 {
    fmt.Println(num, "has 1 digit")
    } else {
    fmt.Println(num, "has multiple digits")
    }
```

## For (classic, short, while, range)
```go

    // Короткое объявление
	i := 1
	for i <= 3 {
		fmt.Println(i)
		i = i + 1
	}

	// Классическое объявление цикла
	for j := 7; j <= 9; j++ {
		fmt.Println(j)
	}

	// Бесконечное выполнение цикла
	for {
		fmt.Println("loop")
	}
	
	nums := []int{2, 3, 4}
	sum := 0
	// Пример итерирования
	for _, num := range nums {
	    sum += num
	}
```
## Switch
```go
    // Стандартное использование switch
    switch time.Now().Weekday() {
    case time.Saturday, time.Sunday: // Так можно объединять условия
        fmt.Println("It's the weekend")
    default:
        fmt.Println("It's a weekday")
    }
	
	// Вариант использования switch без условия
	t := time.Now()
	switch {
    case t.Hour() < 12:
        fmt.Println("It's before noon")
    default:
        fmt.Println("It's after noon")
    }
	
```

## Функции
```go
package main

import "fmt"

/* 
При объявлении функции используется ключевое слово func. 
Так же, необходимо помимо имен передаваемых значений объявлять и их типы.
Перед телом функции нужно указать возвращаемый тип
*/
func plus(a int, b int) int {

    return a + b
}

// Пример объединения типа нескольких переменных
func plusPlus(a, b, c int) int {
    return a + b + c
}

func main() {
    
	// Примеры вызова функций
    res := plus(1, 2)
    fmt.Println("1+2 =", res)

    res = plusPlus(1, 2, 3)
    fmt.Println("1+2+3 =", res)
}
```

## Multiple return values
Go имеет встроенную поддержку нескольких возвращаемых значений. 
Эта особенность часто применяется в Go, например, для возврата результата функции и ошибки.
```go
package main

import "fmt"

func vals() (int, int) {
    return 3, 7
}

func main() {

    a, b := vals()
    fmt.Println(a)
    fmt.Println(b)

    _, c := vals()
    fmt.Println(c)
}
```

