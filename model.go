package main

/*
	Models are stored here. Let's explain:
		Projects have multiple components (or packages).
			They have multiple jobs.
		Projects also have multiple repositories.
			Each repository has one type.
			Each repository has multiple architectures.
		Projects have multiple owners.
			Owners have a distinct username and a list of projects.
			Owners can be groups or users.
				Groups have multiple users.
				Users have an email address.
	tl;dr projects

	Do note, though, that the Authentication table is not in here for security reasons.
	See auth.go.

	We're using composite primary keys.
	http://weblogs.sqlteam.com/jeffs/archive/2007/08/23/composite_primary_keys.aspx
*/

import (
	"bytes"
	"database/sql"
	"encoding/json"
	_ "github.com/ziutek/mymysql/godrv"
	"strconv"
	"time"
)

type State byte

const (
	_                   = iota
	STATE_UNKNOWN State = iota
	STATE_PENDING
	STATE_ACTIVE
	STATE_FINISHING
	STATE_DEAD
	STATE_FAILED
	STATE_DONE
)

func Model_DB() (*sql.DB, error) {
	info := make(map[string]map[string]string)
	err := Config_GetConfig("database", &info)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	buffer.WriteString(info[Config_Environment]["proto"])
	buffer.WriteString(":")
	buffer.WriteString(info[Config_Environment]["addr"])
	if info[Config_Environment]["options"] != "" {
		buffer.WriteString(",")
		buffer.WriteString(info[Config_Environment]["options"])
	}
	buffer.WriteString("*")
	buffer.WriteString(info[Config_Environment]["dbname"])
	buffer.WriteString("/")
	buffer.WriteString(info[Config_Environment]["user"])
	buffer.WriteString("/")
	buffer.WriteString(info[Config_Environment]["pass"])

	return sql.Open("mymysql", buffer.String())
}

type Project struct { // projects
	Name string            // PRIMARY project_name
	Meta map[string]string // string JSONed (key => value) meta
}

// returns project list
func Model_GetProjectList() ([]string, error) {
	db, err := Model_DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT project_name FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	names := make([]string, 0)

	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	err = rows.Err()
	return names, err

}

// returns project, error
func Model_GetProject(projectname string) (*Project, error) {
	db, err := Model_DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	prep, err := db.Prepare("SELECT * FROM projects WHERE project_name='?' LIMIT 1")
	if err != nil {
		return nil, err
	}

	row := prep.QueryRow(projectname)
	var name string
	var meta string

	err = row.Scan(&name, &meta)
	if err != nil {
		return nil, err
	}

	var metamap map[string]string

	err = json.Unmarshal(meta, &metamap)
	return &Project{Name: name, Meta: metamap}, err
}

type Component struct { // components
	Project       string            // PRIMARY project_name
	Name          string            // PRIMARY component_name
	GitRepository string            // gitrepo
	Meta          map[string]string // string JSONed (key => value) meta
}

func Model_GetComponentList(projectname string) ([]string, error) {
	db, err := Model_DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	prep, err := db.Prepare("SELECT component_name FROM components WHERE project_name='?'")
	if err != nil {
		return nil, err
	}

	rows, err := prep.Query(projectname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	names := make([]string, 0)

	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	err = rows.Err()
	return names, err

}

func Model_GetComponent(projectname string, componentname string) (*Component, error) {
	db, err := Model_DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	prep, err := db.Prepare("SELECT * FROM components WHERE project_name='?' AND component_name='?' LIMIT 1")
	if err != nil {
		return nil, err
	}

	row := prep.QueryRow(projectname, componentname)

	var project string
	var name string
	var gitrepo string
	var meta string

	err = row.Scan(&project, &name, &gitrepo, &meta)

	if err != nil {
		return nil, err
	}

	var metamap map[string]string

	err = json.Unmarshal(meta, &metamap)
	return &Component{Project: project, Name: name, GitRepository: gitrepo, Meta: metamap}, err

}

type Repository struct { // repositories
	Project       string   // PRIMARY project_name
	Name          string   // PRIMARY repository_name
	Type          string   // *Type type
	Architectures []string // pass into type's BuildSh and PublishSh as ENV var arch
	State         State    // state
}

func Model_GetRepositoryList(projectname string) ([]string, error) {
	db, err := Model_DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	prep, err := db.Prepare("SELECT repository_name FROM repositories WHERE project_name='?'")
	if err != nil {
		return nil, err
	}

	rows, err := prep.Query(projectname)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	names := make([]string, 0)

	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	err = rows.Err()
	return names, err

}

func Model_GetRepository(projectname string, repositoryname string) (*Repository, error) {
	db, err := Model_DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	prep, err := db.Prepare("SELECT * FROM repositories WHERE project_name='?' AND repository_name='?' LIMIT 1")
	if err != nil {
		return nil, err
	}

	row := prep.QueryRow(projectname, repositoryname)

	var project string
	var name string
	var repotype string
	var repoarches string
	var state State

	err = row.Scan(&project, &name, &repotype, &repoarches, &state)

	if err != nil {
		return nil, err
	}

	var repoarchesslice []string

	err = json.Unmarshal(meta, &repoarchesslice)
	return &Repository{Project: project, Name: name, Type: repotype, Architectures: repoarchesslice, State: state}, err

}

type Type struct { // types
	Name          string // PRIMARY type_name
	DisplayName   string // display_name
	Description   string // description
	NeedSuperUser bool   // some types need superuser permissions, so this will run sudo with the respective scripts. needsuper
	ApplicableSh  string // determine if component is applicable for this repo. applsh
	BuildSh       string // build a component (some components may need chroot - this happens here). When done, move new components to structure emulating repository folder for merging. buldsh
	PublishSh     string // given a repository folder, publish/update the repository. publsh
}

func Model_GetTypeList(projectname string) ([]string, error) {
	db, err := Model_DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT type_name FROM types")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	names := make([]string, 0)

	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	err = rows.Err()
	return names, err

}

func Model_GetType(typename string) (*Type, error) {
	db, err := Model_DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	prep, err := db.Prepare("SELECT * FROM types WHERE type_name='?' LIMIT 1")
	if err != nil {
		return nil, err
	}

	row := prep.QueryRow(typename)

	var name string
	var displayname string
	var description string
	var superuser bool
	var applicablesh string
	var buildsh string
	var publishsh string

	err = row.Scan(&name, &displayname, &description, &superuser, &applicablesh, &buildsh, &publishsh)

	if err != nil {
		return nil, err
	}

	return &Type{Name: name, DisplayName: displayname, Description: description, NeedSuperUser: superuser, ApplicableSh: applicablesh, BuildSh: buildsh, PublishSh: publishsh}, nil
}

type Job struct { // jobs
	Project      string            // PRIMARY project_name
	Repository   string            // PRIMARY repository_name
	Architecture string            // PRIMARY repository_arch
	Component    string            // PRIMARY component_name
	Id           uint64            // PRIMARY AUTO_INCREMENT job_id
	Time         *time.Time        // string UNIXTIME time
	State        State             // state
	Results      map[string]string // string JSON (filename => sha256 for filestore) results
}

// This works differently for a variety of reasons, most notably because it'd be hell to try every permutation.
// Available attributes so you can filter jobs:
// 		project => filter with project
//			repository => filter with repository in project (requires project)
//				arch => filter with arch in repository (requires repository)
//			component => filter with component in project (requires project)
// 		fromtime => filter with after time period unix time UTC
//		totime => filter with before time period unix time UTC
//		state => filter with state
//
func Model_GetJobList(attrs map[string]string) ([]*Job, error) {
	db, err := Model_DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var prepStr bytes.Buffer
	prepStr.WriteString("SELECT project_name, repository_name, repository_arch, component_name, job_id FROM jobs")

	stmtSlice := make([]interface{}, 0)

	if len(attrs) != 0 {
		prepStr.WriteString(" WHERE ")

		useAnd := false

		if prjname, ok := attrs["project"]; ok {
			prepStr.WriteString("project_name='?' ")
			stmtSlice = append(stmtSlice, prjname)
			useAnd = true

			if reponame, ok := attrs["repository"]; ok {
				prepStr.WriteString("AND repository_name='?' ")
				stmtSlice = append(stmtSlice, reponame)

				if archname, ok := attrs["arch"]; ok {
					prepStr.WriteString("AND repository_arch='?' ")
					stmtSlice = append(stmtSlice, archname)
				}

			}
			if compname, ok := attrs["component"]; ok {
				prepStr.WriteString("AND component_name='?' ")
				stmtSlice = append(stmtSlice, compname)
			}
		}
		if fromtime, ok := attrs["fromtime"]; ok {
			if useAnd {
				prepStr.WriteString("AND ")
			} else {
				useAnd = true
			}
			prepStr.WriteString("time >='?' ")
			stmtSlice = append(stmtSlice, fromtime)
		}

		if totime, ok := attrs["totime"]; ok {
			if useAnd {
				prepStr.WriteString("AND ")
			} else {
				useAnd = true
			}
			prepStr.WriteString("time <='?' ")
			stmtSlice = append(stmtSlice, totime)
		}

		if state, ok := attrs["state"]; ok {
			if useAnd {
				prepStr.WriteString("AND ")
			} else {
				useAnd = true
			}
			prepStr.WriteString("state='?' ")
			stmtSlice = append(stmtSlice, state)
		}
	}

	prep, err := db.Prepare(prepStr.String())
	if err != nil {
		return nil, err
	}

	rows, err := prep.Query(stmtSlice...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*Job, 0)

	for rows.Next() {
		var prjname string
		var reponame string
		var archname string
		var compname string
		var jobid uint64
		rows.Scan(&prjname, &reponame, &archname, &compname, &jobid)
		job, err := Model_GetJob(prjname, reponame, archname, compname, jobid)
		if err != nil {
			return nil, err
		}
		results = append(results, job)
	}

	err = rows.Err()
	return results, err

}

func Model_GetJob(projectname string, repositoryname string, architecture string, component string, id uint64) (*Job, error) {
	db, err := Model_DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	prep, err := db.Prepare("SELECT * FROM jobs WHERE project_name='?' AND repository_name='?' AND repository_arch='?' AND component_name='?' AND job_id='?' LIMIT 1")
	if err != nil {
		return nil, err
	}

	row := prep.QueryRow(projectname, repositoryname, architecture, component, id)

	var project string
	var repo string
	var arch string
	var comp string
	var jobid uint64
	var unixtime string
	var state State
	var rawresults string

	err = row.Scan(&project, &repo, &arch, &comp, &jobid, &unixtime, &state, &rawresults)

	if err != nil {
		return nil, err
	}

	timeS := time.Unix(strconv.ParseInt(unixtime, 10, 0), -1)

	var resultmap map[string]string

	err = json.Unmarshal(rawresults, &resultmap)

	return &Job{Project: project, Repository: repo, Architecture: arch, Component: comp, Id: jobid, Time: timeS, State: state, Results: resultmap}, err

}

type Owner struct { // owners
	Username string   `db:"username"` // PRIMARY username
	Projects []string `db:"projects"` // []*Project (now list of project names) projects
	Email    string   `db:"email"`    // email
}

type Group struct { // same table as above
	Owner          // isgroup BOOL
	Users []string `db:"groupusers"` // []*User (now list of usernames) groupusers TEXT
}

// returns username, isGroup
func Model_GetOwnerList() (map[string]bool, error) {
	db, err := Model_DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT username, isgroup FROM owners")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	names := make(map[string]bool, 0)

	for rows.Next() {
		var name string
		var isgroup bool
		rows.Scan(&name, &isgroup)
		names[name] = isgroup
	}

	err = rows.Err()
	return names, err
}

// returns Owner, isGroup, error
func Model_GetOwner(username string) (*Owner, bool, error) {
	db, err := Model_DB()
	if err != nil {
		return nil, false, err
	}
	defer db.Close()

	prep, err := db.Prepare("SELECT * FROM owners WHERE username='?' LIMIT 1")
	if err != nil {
		return nil, false, err
	}

	row := prep.QueryRow(username)

	//username TEXT, projects TEXT, email TEXT, isgroup BOOL, groupusers TEXT

	var user string
	var prjJSON string
	var email string
	var isgroup bool
	var groupUsersJSON sql.NullString

	err = row.Scan(&user, &prjJSON, &email, &isgroup, &groupUsersJSON)

	if err != nil {
		return nil, false, err
	}

	prjArray := make([]string, 0)
	err = json.Unmarshal(prjJSON, prjArray)

	if err != nil {
		return nil, false, err
	}

	if isgroup {
		groupUsersArray := make([]string, 0)
		if groupUsersJSON.Valid {
			err = json.Unmarshal(groupUsersJSON.String, groupUsersArray)
			if err != nil {
				return nil, true, err
			}
		}
		return &Group{Owner{Username: user, Projects: prjArray, Email: email}, groupUsersArray}, true, nil
	} else {
		return &Owner{Username: user, Projects: prjArray, Email: email}, false, nil
	}
}
