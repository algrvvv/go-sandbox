package src

type ExecutionResult struct {
	Output string // вывод исполняемого кода
	Error  string // ошибка при выполнении кода
}

type Session struct {
	Code string // код сессии
}
