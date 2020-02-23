package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mvilche/cron"
)

type Runner struct {
	cron *cron.Cron
}

func NewRunner() *Runner {
	r := &Runner{
		cron: cron.New(),
	}
	return r
}

// Add crontab entry
func (r *Runner) Add(cronjob CrontabEntry) error {
	cronSpec := cronjob.Spec
	if !strings.HasPrefix(cronjob.Spec, "@") {
		cronSpec = fmt.Sprintf("0 %s", cronjob.Spec)
	}

	err := r.cron.AddFunc(cronSpec, r.cmdFunc(cronjob, func(execCmd *exec.Cmd) bool {
		// before exec callback
		LoggerInfo.CronjobExec(cronjob)
		return true
	}))

	if err != nil {
		LoggerError.Printf("Failed add cron job spec:%v cmd:%v err:%v", cronjob.Spec, cronjob.Command, err)
	} else {
		LoggerInfo.CronjobAdd(cronjob)
	}

	return err
}

// Add crontab entry with user
func (r *Runner) AddWithUser(cronjob CrontabEntry) error {

	cronSpec := cronjob.Spec
	if !strings.HasPrefix(cronjob.Spec, "@") {
		cronSpec = fmt.Sprintf("0 %s", cronjob.Spec)
	}

	err := r.cron.AddFunc(cronSpec, r.cmdFunc(cronjob, func(execCmd *exec.Cmd) bool {
		// before exec callback
		LoggerInfo.CronjobExec(cronjob)

		// lookup username
		u, err := user.Lookup(cronjob.User)
		if err != nil {
			LoggerError.Printf("user lookup failed: %v", err)
			return false
		}

		// convert userid to int
		userId, err := strconv.ParseUint(u.Uid, 10, 32)
		if err != nil {
			LoggerError.Printf("Cannot convert user to id:%v", err)
			return false
		}

		// convert groupid to int
		groupId, err := strconv.ParseUint(u.Gid, 10, 32)
		if err != nil {
			LoggerError.Printf("Cannot convert group to id:%v", err)
			return false
		}

		// add process credentials
		execCmd.SysProcAttr = &syscall.SysProcAttr{}
		execCmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(userId), Gid: uint32(groupId)}
		return true
	}))

	if err != nil {
		LoggerError.Printf("Failed add cron job %v; Error:%v", LoggerError.CronjobToString(cronjob), err)
	} else {
		LoggerInfo.Printf("Add cron job %v", LoggerError.CronjobToString(cronjob))
	}

	return err
}

// Return number of jobs
func (r *Runner) Len() int {
	return len(r.cron.Entries())
}

// Start runner
func (r *Runner) Start() {
	LoggerInfo.Printf("Iniciando servicio con %d cron activos\n", r.Len())
	r.cron.Start()
}

// Stop runner
func (r *Runner) Stop() {
	r.cron.Stop()
	LoggerInfo.Println("Stop runner")
}

// ENVIAR CORREO ELECTRONICO LUEGO DE EJECUTAR UNA TAREA

// Execute crontab command
func (r *Runner) cmdFunc(cronjob CrontabEntry, cmdCallback func(*exec.Cmd) bool) func() {

	var email_notification = "false"

	cmdFunc := func() {
		// fall back to normal shell if not specified
		taskShell := cronjob.Shell
		if taskShell == "" {
			taskShell = DEFAULT_SHELL
		}

		fmt.Println(cronjob.Email)

		start := time.Now()

		// Init command
		execCmd := exec.Command(taskShell, "-c", cronjob.Command)

		// add custom env to cronjob
		if len(cronjob.Env) >= 1 {
			execCmd.Env = append(os.Environ(), cronjob.Env...)
		}

		for i := 0; i < len(cronjob.Env); i++ {

			if cronjob.Env[i] == "EMAIL_NOTIFICACION=true" {
				email_notification = "true"
				LoggerInfo.Println("Notificaciones de envio por email activadas")
			} else {
				LoggerInfo.Println("Notificaciones de envio por email desactivadas")
			}
		}

		// exec custom callback
		if cmdCallback(execCmd) {

			// exec job
			out, err := execCmd.CombinedOutput()

			elapsed := time.Since(start)

			if err != nil {
				LoggerError.CronjobExecFailed(cronjob, string(out), err, elapsed)

				if email_notification == "true" && len(opts.ConfigFile) >= 1 {
					SendEmail(configfile, "ERROR", "<h3>Detalles:</h3><br>Error: "+string(out)+"<br><h4>Cron:</h4> "+cronjob.Command+"<br><h4>Hora ejecución:</h4>"+NowDate()+" "+NowHora()+"<br><h4>Duracion ejecución: </h4>"+elapsed.String()+"<br>")

				}

			} else {

				LoggerInfo.CronjobExecSuccess(cronjob, elapsed)

				if email_notification == "true" && len(opts.ConfigFile) >= 1 {

					SendEmail(configfile, "OK", "<h3>Detalles:</h3> Cron ejecutado de forma correcta<br><h4>Cron:</h4> "+cronjob.Command+"<br><h4>Hora ejecución:</h4>"+NowDate()+" "+NowHora()+"<br><h4>Duracion ejecución: </h4>"+elapsed.String()+"<br>")
				}

			}
		}
	}
	return cmdFunc
}
