package main_test

import (
	"testing"

	"."
)

var a main.App

func TestMain(m *testing.M) {
	a = main.App{}
	a.Init()
	a.Start()
}

func clearEmployeeTable() {
	a.DB.Exec("DELETE FROM employee")
	a.DB.Exec("ALTER SEQUENCE employee")
}
