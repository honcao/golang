package exam

var (
	firstName = "Hongbin"
	lastName  = "Cao"
)

const (
	expectScore = 98
	//ActualScore represent actual score
	ActualScore = 80
)

// Tester interface
type Tester interface {
	Test(expect int) (actual int, err error)
}

// Student struct
type Student struct {
	name        string
	expectScore int
}

// New Return student object
func New(name string, expectScore int) (student Student, err error) {
	return Student{
		name:        name,
		expectScore: expectScore,
	}, nil
}

//Test student take the test
func (s Student) Test(expect int) (actual int, err error) {
	return s.expectScore - 20, nil
}
