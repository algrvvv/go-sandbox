package src

type ExecutionResult struct {
	Code   string `json:"code"`   // код, который был отправлен на сборку
	Output string `json:"output"` // вывод исполняемого кода
	Error  string `json:"error"`  // ошибка при выполнении кода
}

type SessionID string

type Session struct {
	Code string // код сессии
	Uid  string // временный идентификатор пользователя
}
