package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type Employee struct {
	ID           int     `json:"id,omitempty"`
	Name         string  `json:"name"`
	Department   string  `json:"department"`
	Title        string  `json:"title"`
	Remuneration float64 `json:"remuneration"`
	Expenses     float64 `json:"expenses"`
	Year         int     `json:"year"`
}

func (em *Employee) getEmployee(db *sql.DB) error {
	getEmpQuery := "SELECT name, department, title, remuneration, expenses, year FROM remuneration WHERE id = $1"
	return db.QueryRow(getEmpQuery, em.ID).Scan(&em.Name, &em.Department, &em.Title, &em.Remuneration, &em.Expenses, &em.Year)
}

func (em *Employee) createEmployee(db *sql.DB) error {
	createEmpQuery := "INSERT INTO employeee(name, department, title, remuneration, expenses, year) VALUES($1, $2, $3, $4, $5, $6) RETURNING id"
	err := db.QueryRow(createEmpQuery, em.Name, em.Department, em.Title, em.Remuneration, em.Expenses, em.Year).Scan(&em.ID)
	if err != nil {
		return err
	}
	return nil
}

func (em *Employee) updateEmployee(db *sql.DB) error {
	updateEmpQuery := "UPDATE employee SET name=$1, department=$1, title=$3, remuneration=$4, expenses=$5, year=$6 WHERE id=$7"
	_, err := db.Exec(updateEmpQuery, em.Name, em.Department, em.Title, em.Remuneration, em.Expenses, em.Year, em.ID)
	return err
}

func (em *Employee) deleteEmployee(db *sql.DB) error {
	deleteEmpQuery := "DELETE from employee WHERE id=$1"
	_, err := db.Exec(deleteEmpQuery, em.ID)
	return err
}

func getAllEmployees(db *sql.DB) ([]Employee, error) {
	getAllEmpQuery := "SELECT name, department, title, remuneration, expenses, year FROM employee"
	rows, err := db.Query(getAllEmpQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	allEmployees := []Employee{}

	for rows.Next() {
		var currEmp Employee
		if err := rows.Scan(&currEmp.Name, &currEmp.Department, &currEmp.Title, &currEmp.Remuneration, &currEmp.Expenses, &currEmp.Year); err != nil {
			fmt.Println(err)
			return nil, err
		}
		allEmployees = append(allEmployees, currEmp)
	}
	return allEmployees, nil
}

func isAnEmployeeDepartment(filter string) bool {
	s := []string{
		"Engineering Services",
		"Dev Svcs, Bldg & Licensing",
		"Mayor & City Council",
		"Office of the City Manager",
		"City Clerk's Office",
		"Planning, Urban Des & Sustain",
		"IT, Digital Strategy & 311",
		"Community Services",
		"Real Estate & Facilities Mgmt",
		"Board of Parks & Recreation",
		"Vancouver Public Library Board",
		"Finance, Risk&Supply Chain Mgt",
		"VFRS & OEM",
		"Human Resources",
		"Law Department",
	}

	for _, dept := range s {
		if dept == filter {
			fmt.Println("Matches: " + filter)
			return true
		}
		fmt.Println("the other cases.")
	}
	fmt.Println("Called")

	return false
}

func getSomeEmployees(db *sql.DB, filter string) ([]Employee, error) {
	getSomeEmpQuery := "SELECT name, department, title, remuneration, expenses, year FROM employee WHERE department = $1"

	if !isAnEmployeeDepartment(filter) {
		fmt.Println("it hits.")
		return []Employee{}, errors.New("The entered subdomain is not a department")
	}

	rows, err := db.Query(getSomeEmpQuery, filter)
	if err != nil {
		fmt.Println("query error" + err.Error())
		return nil, err
	}

	defer rows.Close()
	allEmployees := []Employee{}
	for rows.Next() {
		var currEmp Employee
		if err := rows.Scan(&currEmp.Name, &currEmp.Department, &currEmp.Title, &currEmp.Remuneration, &currEmp.Expenses, &currEmp.Year); err != nil {
			fmt.Println(err)
			return nil, err
		}
		allEmployees = append(allEmployees, currEmp)
	}
	return allEmployees, nil

}
