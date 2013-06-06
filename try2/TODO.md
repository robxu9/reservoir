TODO List
=========

* Need methods for manipulating structs in databases
	* This could be easily accomplished with the completion of `database_old.go`, but until then.
	* Create methods for handling owner/project/user/team
* Create structs for Job, which represents a task.
	* Such as as creation job, build job, a signing job, a publishing job which the worker can do.
* Create a scheduler - a channel of Jobs.
	* Scheduler should monitor and occasionally check up on how things are doing.
	* I'm thinking just have a method open for cron to add default tasks to scheduler.
* Create controllers for jobs, owners, projects, users, teams to handle API requests.
	* Most methods require apikey and permission access.
	* Additionally, applications calling for users and teams should go through Owners to get an
		abstract view of it (is (not) user/team) before proceeding to call users and/or teams.