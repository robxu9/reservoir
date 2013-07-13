package main

import (
	"database/sql"
)

// Up is executed when this migration is applied
func Up_20130702202327(txn *sql.Tx) {
	_, err := txn.Exec("CREATE TABLE projects (project_name TEXT, meta TEXT, PRIMARY KEY (project_name))")
	if err != nil {
		panic(err)
	}
	_, err = txn.Exec("CREATE TABLE components (project_name TEXT, component_name TEXT, gitrepo TEXT, meta TEXT, PRIMARY KEY (project_name, component_name))")
	if err != nil {
		panic(err)
	}
	_, err = txn.Exec("CREATE TABLE repositories (project_name TEXT, repository_name TEXT, type TEXT, arch TEXT, state TINYINT UNSIGNED, PRIMARY KEY (project_name, repository_name))")
	if err != nil {
		panic(err)
	}
	_, err = txn.Exec("CREATE TABLE types (type_name TEXT, display_name TEXT, description TEXT, needsuper BOOL, applsh TEXT, buldsh TEXT, publsh TEXT, PRIMARY KEY (type_name))")
	if err != nil {
		panic(err)
	}
	_, err = txn.Exec("CREATE TABLE jobs (project_name TEXT, repository_name TEXT, repository_arch TEXT, component_name TEXT, job_id BIGINT AUTO_INCREMENT, time TEXT, state TINYINT UNSIGNED, results TEXT, PRIMARY KEY (project_name, repository_name, repository_arch, component_name, job_id))")
	if err != nil {
		panic(err)
	}
	_, err = txn.Exec("CREATE TABLE owners (username TEXT, projects TEXT, email TEXT, isgroup BOOL, groupusers TEXT, PRIMARY KEY (username))")
	if err != nil {
		panic(err)
	}
	_, err = txn.Exec("CREATE TABLE authentication (username TEXT, password TEXT, token TEXT, admin TINYINT, PRIMARY KEY (username))")
	if err != nil {
		panic(err)
	}
	err = txn.Commit()
	if err != nil {
		panic(err)
	}
}

// Down is executed when this migration is rolled back
func Down_20130702202327(txn *sql.Tx) {
	_, err := txn.Exec("DROP TABLE IF EXISTS projects, components, repositories, types, jobs, owners, authentication;")
	if err != nil {
		panic(err)
	}
	err = txn.Commit()
	if err != nil {
		panic(err)
	}
}
