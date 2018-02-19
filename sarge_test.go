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
		BeforeEach(func() {
			mocket.Catcher.Register()
			db, err := gorm.Open(mocket.DRIVER_NAME, "A mocked connection")
			if err != nil {
				Fail("Failed to mock db!")
			}
			SetDb(db)
		})
		Describe("Department Struct", func() {
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
										"from_date": time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
										"to_date":   time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC),
									},
								},
							},
						})
						Q1, _ := easyEmployee.GetSalaryByTimeRange(
							time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
							time.Date(2009, time.April, 1, 0, 0, 0, 0, time.UTC),
						)
						Q2, _ := easyEmployee.GetSalaryByTimeRange(
							time.Date(2009, time.April, 1, 0, 0, 0, 0, time.UTC),
							time.Date(2009, time.July, 1, 0, 0, 0, 0, time.UTC),
						)
						Q3, _ := easyEmployee.GetSalaryByTimeRange(
							time.Date(2009, time.July, 1, 0, 0, 0, 0, time.UTC),
							time.Date(2009, time.October, 1, 0, 0, 0, 0, time.UTC),
						)
						Q4, _ := easyEmployee.GetSalaryByTimeRange(
							time.Date(2009, time.October, 1, 0, 0, 0, 0, time.UTC),
							time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC),
						)
						Expect(Q1 + Q2 + Q3 + Q4).To(Equal(52000.00))
					})
				})
			})
		})
		Describe("Salary Struct", func() {
			Describe("Basic Helpers", func() {
				var (
					testSalary Salary
				)
				BeforeEach(func() {
					testSalary = Salary{
						EmployeeNumber: 0,
						Salary:         75000.00,
					}
					if start, parseErr := time.Parse("2006-01-02", "2015-01-01"); parseErr == nil {
						testSalary.StartDate = start
					} else {
						Fail("Invalid Start Date")
					}
					if end, parseErr := time.Parse("2006-01-02", "2016-01-01"); parseErr == nil {
						testSalary.EndDate = end
					} else {
						Fail("Invalid End Date")
					}
				})

				It("Should Calculate salary duration", func() {
					Expect(testSalary.Duration()).To(Equal(uint(365)))
				})
				It("Should Calculate correct daily rate", func() {
					Expect(testSalary.DailyRate()).To(Equal(75000.00 / 365))
				})
			})
		})
	})
})
