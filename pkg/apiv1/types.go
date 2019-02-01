package apiv1

type User struct {
	Name string
	Age  int
	Test *string
}

type TypeB struct {
	Name string
	Age  int
}
type TypeA struct {
	TypeB TypeB
}
