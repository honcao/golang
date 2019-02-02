package apiv1

type User struct {
	Name string
	Age  int
	Test *string
}

type TypeB struct {
	Name         string
	Age          int
	Uint         uint
	Float32      float32
	StringPtr    *string
	IntPtr       *int
	StructPtr    *TypeD
	IntSlice     []int
	StringSlice  []string
	TypeCSlice   []TypeC
	StringArray  [5]string
	StringMap    map[string]string
	StringMapMap map[string]map[string]string
}
type TypeA struct {
	TypeB TypeB
}

type TypeC struct {
	CName string
}

type TypeD struct {
	DName string
}
