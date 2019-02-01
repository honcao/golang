package user_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/honcao/golang/pkg/user"
)

var _ = Describe("User", func() {
	var (
		u1 User
		u2 User
	)

	BeforeEach(func() {
		u1 = User{
			FirstName: "Hongbin",
			LastName:  "Cao",
		}

		u2 = User{
			LastName: "Fox In Socks",
		}
	})

	Describe("GetFullName()", func() {
		Context("Both first and last name", func() {
			It("should first name ", func() {
				Expect(u1.GetFullName()).To(Equal(u1.LastName + " " + u1.FirstName))
			})
		})

		Context("With fewer than 300 pages", func() {
			It("should be a short story", func() {
				Expect(u2.GetFullName()).To(Equal("Fox In Socks"))
			})
		})
	})

})

var _ = Describe("User", func() {
	var (
		u1 User
		u2 User
	)

	BeforeEach(func() {
		u1 = User{
			FirstName: "Hongbin",
			LastName:  "Cao",
		}

		u2 = User{
			LastName: "Fox In Socks",
		}
	})

	Describe("GetFullName()", func() {
		Context("Both first and last name", func() {
			It("should first name ", func() {
				Expect(u1.GetFullName()).To(Equal(u1.LastName + " " + u1.FirstName))
			})
		})

		Context("With fewer than 300 pages", func() {
			It("should be a short story", func() {
				Expect(u2.GetFullName()).To(Equal("Fox In Socks"))
			})
		})
	})

})
