package src

type ExecutionResult struct {
	Code   string // код, который был отправлен на сборку
	Output string // вывод исполняемого кода
	Error  string // ошибка при выполнении кода
}

type Session struct {
	Code string // код сессии
}
