package application

const (
	ExitSuccess = 0
	ExitError   = 1
)

func GetExitCode(err error) int {
	switch {
	case err == nil:
		return ExitSuccess
	default:
		return ExitError
	}
}
