package msg

type Logger interface{
	Info(msg string)
	Success(msg string)
	Warn(msg string)
	Error(msg string)
	Fail(msg string)
}
