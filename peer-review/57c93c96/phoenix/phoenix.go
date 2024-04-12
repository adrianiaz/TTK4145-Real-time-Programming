package phoenix

import (
	"syscall"
)

func RespawnAfterPanic() {
	if r := recover(); r != nil {
		// we die gracefully so deathBedProcedure in main() can resurrect us
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}
}
