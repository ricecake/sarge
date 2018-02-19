package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/ricecake/sarge/models"
	"time"
)

var _ = Describe("Sarge", func() {
	Describe("DB Model", func() {
		Describe("Department Struct", func() {
		})
		Describe("Employee Struct", func() {

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
				It("Should Calculate correct daily rate", func(){
					Expect(testSalary.DailyRate()).To(Equal(75000.00/365))
				})
			})
		})
	})
})
