package main

import (
	"log"
	"slices"
	"strings"
	"time"

	"github.com/ski7777/nuki-cron/internal/config"
	"github.com/ski7777/nuki-cron/internal/datetimeinterval"
	"github.com/ski7777/nuki-cron/internal/nuki"
	"github.com/ski7777/nuki-cron/internal/recurringinterval"
)

const (
	MANAGED_SUFFIX = " (managed by nuki-cron)"
)

func main() {
	//config parsing
	c, err := config.NewConfigFromFile("/config/config.json")
	if err != nil {
		log.Panicln(err)
	}
	permissionteams := map[string][]int{}
	for name := range c.Permissions {
		permissionteams[name] = []int{}
	}
	var teamschedules [][]recurringinterval.RecurringInterval
	var teamextraevents []datetimeinterval.DateTimeIntervals
	for tn, team := range c.Teams {
		var schedules []recurringinterval.RecurringInterval
		for _, schedule := range team.Schedules {
			ri, err := schedule.ToRecurringInterval()
			if err != nil {
				log.Panicln(err)
			}
			schedules = append(schedules, ri)
		}
		teamschedules = append(teamschedules, schedules)

		var extraevents datetimeinterval.DateTimeIntervals
		for _, extraevent := range team.ExtraEvents {
			extraevents.Add(extraevent.ToDateTimeInterval())
		}
		teamextraevents = append(teamextraevents, extraevents)

		for _, permission := range team.Permissions {
			if _, ok := permissionteams[permission]; !ok {
				log.Printf("permission %s for team %s not found in config", permission, team.Name)
			} else {
				permissionteams[permission] = append(permissionteams[permission], tn)
			}
		}
	}

	//evaluation
	//teams
	var teamintervals []datetimeinterval.DateTimeIntervals
	for tn, schedules := range teamschedules {
		intervals := teamextraevents[tn].Copy()
		for _, schedule := range schedules {
			interval, err := schedule.GetNextTimes()
			if err != nil {
				log.Panicln(err)
			}
			intervals.Add(interval)
		}
		teamintervals = append(teamintervals, intervals)
	}
	//permissions
	permissionintervals := map[string]*datetimeinterval.DateTimeInterval{}
	for permission, tns := range permissionteams {
		intervals := datetimeinterval.DateTimeIntervals{}
		for _, tn := range tns {
			intervals.AddAll(teamintervals[tn])
		}
		interval, ok := intervals.GetNextExtended()
		if !ok {
			log.Printf("no intervals for permission %s", permission)
			permissionintervals[permission] = nil
		} else {
			permissionintervals[permission] = &interval
		}
	}

	permissionhumannames := map[string]string{}
	n := nuki.NewNukiClient(c.ApiKey)
	auths, err := n.GetAuths()
	if err != nil {
		log.Panicln(err)
	}
	for pn, perm := range c.Permissions {
		if perm.Omit {
			continue
		}
		i := slices.IndexFunc(auths, func(auth nuki.GetAuthReponseAuth) bool {
			return auth.Id == perm.Id
		})
		if i == -1 {
			log.Printf("auth with id %s for permission %s not found in nuki api", perm.Id, pn)
			continue
		}
		name := auths[i].Name
		if !strings.HasSuffix(name, MANAGED_SUFFIX) {
			name = name + MANAGED_SUFFIX
		}
		permissionhumannames[pn] = name
	}

	updates := nuki.SmartLockAuthMultiUpdateRequest{}
	zeroTime := time.Time{}.UTC().Format(time.RFC3339)
	for pn, phn := range permissionhumannames {
		pid := c.Permissions[pn].Id
		pi := permissionintervals[pn]
		update := nuki.SmartLockAuthMultiUpdate{}
		update.Id = pid
		update.Name = phn
		if pi == nil {
			log.Printf("permission %s (id: %s, human name: %s) has no intervals", pn, pid, phn)
			update.Enabled = false
			update.AllowedFromDate = zeroTime
			update.AllowedUntilDate = zeroTime
		} else {
			log.Printf("permission %s (id: %s, human name: %s) has next interval %s - %s", pn, pid, phn, pi.Start.String(), pi.End.String())
			update.Enabled = true
			update.AllowedWeekDays = 127 //all
			update.AllowedFromDate = pi.Start.UTC().Format(time.RFC3339)
			update.AllowedUntilDate = pi.End.UTC().Format(time.RFC3339)
		}
		updates = append(updates, update)
	}
	err = n.UpdateAuths(updates)
	if err != nil {
		log.Fatalln(err)
	}
}
