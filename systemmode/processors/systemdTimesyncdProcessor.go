package processors

import (
	"github.com/gandrille/go-commons/ini"
	"github.com/gandrille/go-commons/result"
	"github.com/gandrille/go-commons/systemctl"
)

type systemdtimesyncd struct {
}

func (e systemdtimesyncd) Key() string {
	return "configure-systemd-timesyncd"
}

func (e systemdtimesyncd) Describe(args []string) string {
	return "NTP configuration using server " + args[0]
}

func (e systemdtimesyncd) Run(args []string) result.Result {
	server := args[0]
	service := "systemd-timesyncd.service"
	updated := false

	// Update config file
	fileupdated, err1 := ini.SetValue("/etc/systemd/timesyncd.conf", "Time", "NTP", server, false, true)
	if err1 != nil {
		return result.NewError(e.Describe(args) + " " + err1.Error())
	}
	msg := "NTP "
	if fileupdated {
		msg = "configuration updated, "
		updated = true
	} else {
		msg = "configuration already updated, "
	}

	// (re)start service
	if fileupdated {
		restarted, err2 := systemctl.Restart(service)
		if err2 != nil {
			return result.NewError(e.Describe(args) + " " + err2.Error())
		}
		if restarted {
			msg += "service restarted, "
		} else {
			msg += "service started, "
		}
		updated = true
	} else {
		activated, err3 := systemctl.Activate(service)
		if err3 != nil {
			return result.NewError(e.Describe(args) + " " + err3.Error())
		}
		if activated {
			msg += "service activated, "
			updated = true
		} else {
			msg += "service already started, "
		}
	}

	// Enable service
	enabled, err4 := systemctl.Enable(service)
	if err4 != nil {
		return result.NewError(e.Describe(args) + " " + err4.Error())
	}
	if enabled {
		msg += "service enabled"
		updated = true
	} else {
		msg += "service already enabled"
	}

	if updated {
		return result.NewUpdated(msg)
	}
	return result.NewUnchanged(msg)
}
