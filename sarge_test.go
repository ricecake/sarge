package main_test

import (
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/ricecake/sarge/models"
	mocket "github.com/selvatico/go-mocket"
	"time"
)

var _ = Describe("Sarge", func() {
	Describe("DB Model", func() {
		mocket.Catcher.Register()
		db, err := gorm.Open(mocket.DRIVER_NAME, "A mocked connection")
		if err != nil {
			Fail("Failed to mock db!")
		}
		SetDb(db)
		BeforeEach(func() {
			mocket.Catcher.Reset()
		})
		Describe("Department Struct", func() {
			var (
				testDept Department
			)
			BeforeEach(func() {
				testDept = Department{
					Name:   "Fictitious Activities",
					Number: "d042",
				}
			})
			Describe("GetEmployeeByTimeRange", func() {
				Context("Simple lookups", func() {
					It("Translates department assignments to employees", func() {
						mocket.Catcher.Attach([]*mocket.FakeResponse{
							{
								Once:    true,
								Pattern: "SELECT * FROM \"dept_emp\"  WHERE",
								Response: []map[string]interface{}{
									{"emp_no": 1},
									{"emp_no": 2},
									{"emp_no": 3},
								},
							},
							{
								Once:    true,
								Pattern: "SELECT * FROM \"employees\"",
								Response: []map[string]interface{}{
									{"emp_no": 1},
									{"emp_no": 2},
									{"emp_no": 3},
								},
							},
						})
						emps, _ := testDept.GetEmployeeByTimeRange(makeTime("2011-01-01"), makeTime("2011-04-01"))

						var ids []uint
						for _, emp := range emps {
							ids = append(ids, emp.Number)
						}
						Expect(ids).To(ConsistOf([]uint{1, 2, 3}))
					})
				})
			})
			Describe("GetSalaryByTimeRange", func() {
				Context("Simple employment/salary", func() {
					It("Calculates expected departmental salary for time range", func() {
						mocket.Catcher.Attach([]*mocket.FakeResponse{
							{
								Once:    false,
								Pattern: "SELECT * FROM \"dept_emp\"  WHERE",
								Response: []map[string]interface{}{
									{"emp_no": 1},
									{"emp_no": 2},
									{"emp_no": 3},
									{"emp_no": 4},
								},
							},
							{
								Once:    false,
								Pattern: "SELECT * FROM \"employees\"",
								Response: []map[string]interface{}{
									{"emp_no": 1},
									{"emp_no": 2},
									{"emp_no": 3},
									{"emp_no": 4},
								},
							},
							{
								Once:    false,
								Pattern: "SELECT * FROM \"salaries\"  WHERE",
								Response: []map[string]interface{}{
									{
										"salary":    52000.00,
										"from_date": makeTime("2011-01-01"),
										"to_date":   makeTime("2012-01-01"),
									},
								},
							},
						})
						Q1, _ := testDept.GetSalaryByTimeRange(makeTime("2011-01-01"), makeTime("2011-04-01"))
						Q2, _ := testDept.GetSalaryByTimeRange(makeTime("2011-04-01"), makeTime("2011-07-01"))
						Q3, _ := testDept.GetSalaryByTimeRange(makeTime("2011-07-01"), makeTime("2011-10-01"))
						Q4, _ := testDept.GetSalaryByTimeRange(makeTime("2011-10-01"), makeTime("2012-01-01"))
						Expect(Q1 + Q2 + Q3 + Q4).To(Equal(208000.00))
					})
				})
			})
		})
		Describe("Employee Struct", func() {
			var (
				easyEmployee Employee
			)
			BeforeEach(func() {
				mocket.Catcher.Reset()
				easyEmployee = Employee{
					Number:    5,
					FirstName: "Sam",
					LastName:  "Baker",
					Gender:    "f",
				}
			})
			Describe("GetSalaryByTimeRange", func() {
				Context("Simple time ranges", func() {
					It("Returns expected sum of quarterly salaries", func() {
						mocket.Catcher.Attach([]*mocket.FakeResponse{
							{
								Once:    false,
								Pattern: "SELECT * FROM \"salaries\"  WHERE",
								Response: []map[string]interface{}{
									{
										"emp_no":    5,
										"salary":    52000.00,
										"from_date": makeTime("2009-01-01"),
										"to_date":   makeTime("2010-01-01"),
									},
								},
							},
						})
						Q1, _ := easyEmployee.GetSalaryByTimeRange(makeTime("2009-01-01"), makeTime("2009-04-01"))
						Q2, _ := easyEmployee.GetSalaryByTimeRange(makeTime("2009-04-01"), makeTime("2009-07-01"))
						Q3, _ := easyEmployee.GetSalaryByTimeRange(makeTime("2009-07-01"), makeTime("2009-10-01"))
						Q4, _ := easyEmployee.GetSalaryByTimeRange(makeTime("2009-10-01"), makeTime("2010-01-01"))
						Expect(Q1 + Q2 + Q3 + Q4).To(Equal(52000.00))
					})
				})
			})
		})
		Describe("Salary Struct", func() {
			var (
				testSalary Salary
			)
			BeforeEach(func() {
				testSalary = Salary{
					EmployeeNumber: 0,
					Salary:         75000.00,
					StartDate:      makeTime("2015-01-01"),
					EndDate:        makeTime("2016-01-01"),
				}
			})
			Describe("Basic Helpers", func() {
				It("Should Calculate salary duration", func() {
					Expect(testSalary.Duration()).To(Equal(uint(365)))
				})
				It("Should Calculate correct daily rate", func() {
					Expect(testSalary.DailyRate()).To(Equal(75000.00 / 365))
				})
			})
			Describe("Edge cases", func(){
				It("Should recognize leap years", func(){
					testSalary.StartDate = makeTime("2016-01-01")
					testSalary.EndDate = makeTime("2017-01-01")
					Expect(testSalary.Duration()).To(Equal(uint(366)))
				})
				It("handles fractional years", func(){
					testSalary.StartDate = makeTime("2015-01-01")
					testSalary.EndDate = makeTime("2015-02-01")
					Expect(testSalary.Duration()).To(Equal(uint(31)))
				})
				It("handles misaligned date boundries", func(){
					testSalary.StartDate = makeTime("2015-06-15")
					testSalary.EndDate = makeTime("2015-11-09")
					Expect(testSalary.Duration()).To(Equal(uint(147)))
				})
			})
		})
	})
})

func makeTime(stringTime string) time.Time {
	timeStruct, parseErr := time.Parse("2006-01-02", stringTime)
	if parseErr != nil {
		Fail("Invalid time specified!")
	}

	return timeStruct
}
