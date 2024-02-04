package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
)

func main() {
	allProcess := flag.Bool("A", false, "Показать все процессы")
	exceptLeaders := flag.Bool("d", false, "Не показывать процессы сессии")

	flag.Parse() // обработка всех объявленных флагов

	args := []string{"ps"} // предварительно объявляем аргументы для команды ps

	if *allProcess {
		args = append(args, "-A") // добавляем флаг -A, если указан флаг -A
	}
	if *exceptLeaders {
		args = append(args, "-d") // добавляем флаг -d, если указан флаг -d
	}

	// создаем структуру ProcAttr для передачи потоков ввода/вывода команды ps
	process := syscall.ProcAttr{
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
	}

	//запускаем команду ps с указанными аргументами
	pid, err := syscall.ForkExec("/bin/ps", args, &process)
	if err != nil {
		fmt.Println("Ошибка выполнения команды ps:", err)
		return
	}

	// ожидаем завершения процесса
	syscall.Wait4(pid, nil, 0, nil)
}
