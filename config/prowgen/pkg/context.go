// +skip_license_check
/*
Copyright 2022 The cert-manager Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pkg

import (
	"fmt"
	"math"
)

// ProwContext holds jobs and information required to configure jobs for a given release channel.
type ProwContext struct {
	// Branch is the name of the branch corresponding to the release channel modelled by this ProwContext.
	// While it's possible to define a presubmit for multiple branches, this often doesn't correctly model
	// how cert-manager uses prow in practice - usually, we want a different set of supported kubernetes
	// versions for each major cert-manager release (and therefore branch), and in any case want a different
	// dashboard for each supported release channel.
	Branch string

	// Image is the common test image used for running prow jobs.
	Image string

	// PresubmitDashboard, if set, will generate a presubmit TestGrid dashboard name based on the branch name
	// for each presubmit job. If false, no presubmits will be added to a TestGrid dashboard.
	PresubmitDashboard bool

	// PeriodicDashboard, if set, will generate a periodic TestGrid dashboard name based on the branch name
	// for each periodic job. If false, no periodics will be added to a TestGrid dashboard.
	PeriodicDashboard bool

	// Org is the GitHub organisation of the repository under test.
	Org string

	// Repo is the GitHub repository name of the repository under test.
	Repo string

	presubmits []*PresubmitJob
	periodics  []*PeriodicJob
}

// RequiredPresubmit adds a presubmit which is run by default and required to pass before a PR can be merged
func (pc *ProwContext) RequiredPresubmit(job *Job) {
	pc.addPresubmit(job, true, false, "")
}

// RequiredPresubmits adds a list of jobs to the context
func (pc *ProwContext) RequiredPresubmits(jobs []*Job) {
	for _, job := range jobs {
		pc.addPresubmit(job, true, false, "")
	}
}

// OptionalPresubmit adds a presubmit which is not run by default and is optional
func (pc *ProwContext) OptionalPresubmit(job *Job) {
	pc.addPresubmit(job, false, true, "")
}

// OptionalPresubmitIfChanged adds a presubmit which is not run by default and is optional unless a file has been
// changed which matches changedFileRegex. In that situation, the job is always run.
// See https://docs.prow.k8s.io/docs/jobs/#triggering-jobs-based-on-changes
func (pc *ProwContext) OptionalPresubmitIfChanged(job *Job, changedFileRegex string) {
	pc.addPresubmit(job, false, true, changedFileRegex)
}

func (pc *ProwContext) addPresubmit(job *Job, alwaysRun bool, optional bool, changedFileRegex string) {
	job.Name = pc.presubmitJobName(job.Name)

	if pc.PresubmitDashboard {
		addTestGridAnnotations(pc.presubmitDashboardName())(job)
	}

	pc.presubmits = append(pc.presubmits, &PresubmitJob{
		Job: *job,
		// see the comment on ProwContext.Branch for why we only support a single branch here
		Branches:     []string{pc.Branch},
		AlwaysRun:    alwaysRun,
		Optional:     optional,
		RunIfChanged: changedFileRegex,
	})
}

// Periodic adds periodic jobs which will run every `periodicityHours` hours, at some minute
// within the hour, one job for each configured branch
func (pc *ProwContext) Periodics(job *Job, periodicityHours int) {
	originalName := job.Name

	job.Name = pc.periodicJobName(originalName)

	if pc.PeriodicDashboard {
		addTestGridAnnotations(pc.periodicDashboardName())(job)
	}

	pc.periodics = append(pc.periodics, &PeriodicJob{
		Job: *job,
		ExtraRefs: []ExtraRef{
			{
				Org:     pc.Org,
				Repo:    pc.Repo,
				BaseRef: pc.Branch,
			},
		},
		PeriodicityHours: periodicityHours,
		// Use a minute and startHour of 0 when adding the period here, but in JobFile
		// when we actually generate the JobFile struct we'll recalculate values
		// to spread the periodics evenly across every hour / day
		Cron: cronSchedule(0, 0, periodicityHours),
	})
}

func (pc *ProwContext) JobFile() *JobFile {
	presubmitKey := fmt.Sprintf("%s/%s", pc.Org, pc.Repo)

	// By dividing 60 by the number of periodics we get the maximum number of minutes
	// apart which we can schedule jobs, in an effort to maximise the amount of time
	// that each job could theoretically run without anything else running in parallel.

	// We aim to maximise the amount of CPU available to each job at startup, which is
	// often when most jobs need a lot of CPU (e.g. to set up tests / clusters)
	// We use ceil because more spreading isn't a bad thing
	minuteSpread := int(math.Ceil(60.0 / float64(len(pc.periodics))))

	// Count the number of jobs with each periodicity to make it easier to spread them later
	periodicityCounts := map[int]int{}

	for _, p := range pc.periodics {
		periodicityCounts[p.PeriodicityHours] += 1
	}

	// periodicitySeen is used to track how many times we've seen a job with a given
	// periodicity, which allows us to spread jobs with the same periodicity across the day
	periodicitySeen := map[int]int{}

	for i, p := range pc.periodics {
		minute := (i * minuteSpread) % 60

		// hourCounter is how many jobs with the same periodicity we've already seen
		hourCounter := periodicitySeen[p.PeriodicityHours]

		periodicitySeen[p.PeriodicityHours] = hourCounter + 1

		// Hour spread is how far apart we can place jobs in starting hours
		// E.g. Say we have 5 jobs with periodicity 4.
		// (Bear in mind these jobs can only start in hours 0,1,2,3 [1])
		// ceil(4/5) = 1, meaning we have to place a job every hour and can't space them out any more.

		// If we instead had 3 jobs with periodicity 8, ceil(8/3) = 3 and so we can place jobs at
		// 0,3,6 (or 1,4,7) to spread them out more.

		// [1] Effectively we operate modulo the periodicity; if (start hour >= periodicity) it
		// reduces the number of invocations possible in a calendar day, e.g.:
		// "0 7-23/8 * * *" runs 3 times in a day (at 07:00, 15:00 and 23:00), but
		// "0 8-23/8 * * *" only runs twice (at 08:00 and 16:00)
		hourSpread := int(math.Floor(float64(p.PeriodicityHours) / float64(periodicityCounts[p.PeriodicityHours])))
		if hourSpread == 0 {
			hourSpread = 1
		}

		if p.PeriodicityHours == 24 {
			// 24h periodicity is different because it can be started at any hour and will always be run
			// once per day
			// In this case, just set the spread to 7 which gives a pretty good spread modulo 24:
			// [0, 7, 14, 21, 4, 11, 18, 1, 8, 15, 22, 5, 12, 19, 2, 9, 16, 23, 6, 13, 20, 3, 10, 17]
			hourSpread = 7
		}

		startHour := (hourSpread * hourCounter) % p.PeriodicityHours

		p.Cron = cronSchedule(minute, startHour, p.PeriodicityHours)
	}

	return &JobFile{
		Presubmits: map[string][]*PresubmitJob{
			presubmitKey: pc.presubmits,
		},
		Periodics: pc.periodics,
	}
}

// presubmitJobName returns a prow name for the given presubmit job. For example,
// for the branch "release-1.0" and the test "foo", this would return "pull-cert-manager-release-1.0-foo"
func (pc *ProwContext) presubmitJobName(name string) string {
	return fmt.Sprintf("pull-%s-%s-%s", pc.Repo, pc.Branch, name)
}

// periodicJobName returns a prow name for the given periodic job. For example,
// for the branch "release-1.0" and the test "foo", this would return "ci-cert-manager-release-1.0-foo"
func (pc *ProwContext) periodicJobName(name string) string {
	return fmt.Sprintf("ci-%s-%s-%s", pc.Repo, pc.Branch, name)
}

func (pc *ProwContext) presubmitDashboardName() string {
	return fmt.Sprintf("%s-presubmits-%s", pc.Repo, pc.Branch)
}

func (pc *ProwContext) periodicDashboardName() string {
	return fmt.Sprintf("%s-periodics-%s", pc.Repo, pc.Branch)
}

func cronSchedule(minute int, startHour int, periodicityHours int) string {
	return fmt.Sprintf("%02d %02d-23/%02d * * *", minute, startHour, periodicityHours)
}
