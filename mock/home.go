package mock

// WorkInterface  work interface
type WorkInterface interface {
}

// Work worker struct
type Work struct {
	workName string
	money    int
}

// one test method
func (work *Work) sayWorkName() string {
	return work.workName
}

// NewWork  work construct
func NewWork(workName string, money int) *Work {
	return &Work{workName: workName, money: money}
}

// Homework  homework struct
type Homework struct {
	workName string
}

// one test method
func (homework *Homework) sayWorkName() string {
	return homework.workName
}

// NewHomework a construct
func NewHomework(workName string) *Homework {
	return &Homework{workName: workName}
}

// Father a father struct
type Father struct {
	fatherName string
	age        int
	work       *Work
}

// NewFatherWithAllParam  father construct
func NewFatherWithAllParam(work *Work, fatherName string, age int) *Father {
	return &Father{fatherName: fatherName, age: age, work: work}
}

// NewFatherWithPre  father construct
func NewFatherWithPre(work *Work) *Father {
	return &Father{fatherName: "张三", age: 33, work: work}
}

// NewFatherWithStruct  father construct
func NewFatherWithStruct(work Work) *Father {
	return &Father{fatherName: "张三", age: 33, work: &work}
}

// NewFatherWithInterface  father construct
func NewFatherWithInterface(work WorkInterface, fatherName string, age int) *Father {
	return &Father{fatherName: "张三", age: 33, work: work.(*Work)}
}

// Mother a mother struct
type Mother struct {
	motherName string
	age        int
	homework   Homework
}

// NewMother construct
func NewMother(motherName string, age int, homework Homework) *Mother {
	return &Mother{motherName: motherName, age: age, homework: homework}
}

// NewMotherWithInterface  construct
func NewMotherWithInterface(motherName string, age int, work WorkInterface) *Mother {
	return &Mother{motherName: motherName, age: age, homework: work.(Homework)}
}

// Child a child
type Child struct {
	childName string
	age       int
}

// NewChild construct
func NewChild(childName string, age int) *Child {
	return &Child{childName: childName, age: age}
}

// Family  struct
type Family struct {
	familyName string
	father     Father
	children   []Child
	mother     Mother
}

// NewFamily construct
func NewFamily(familyName string, father Father, children []Child, mother Mother) *Family {
	return &Family{familyName: familyName, father: father, children: children, mother: mother}
}
