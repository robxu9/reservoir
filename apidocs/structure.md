Reservoir API structure
=======================

Currently consists of:

==Workers==

* Workers execute Jobs. Supported includes Linux 2.6/3, Darwin with MacPorts, and Cygwin.
* These take the jobs and basically run the scripts attached to them.

==Jobs==

* Define the tasks for workers to do. Jobs are pretty much javascript files.
* Jobs can queue other jobs for the scheduler.

==Scheduler==

* The job scheduler schedules jobs that are queued, either from an API or other jobs.
* This is basically a channel that sends out jobs to available workers.

(so, it would be like <-workerChan, followed by <-jobChan, and do it all over again).


==FUTURE==

===Projects===

* Projects hold projects inside of them.
* Different types of projects can have different things.
* Building will trigger jobs.

===Owners===

* Owners own projects.
* Owners can either be users or teams.