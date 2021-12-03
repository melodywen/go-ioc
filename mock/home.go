package mock

type WorkInterface interface {
	// 请告知名称
	sayWorkName()
}


type Work struct {
	workName string
	money    int
}

func NewWork(workName string, money int) *Work {
	return &Work{workName: workName, money: money}
}

type Homework struct {
	workName string
}

func NewHomework(workName string) *Homework {
	return &Homework{workName: workName}
}

type Father struct {
	fatherName string
	age        int
	work       *Work
}

func NewFatherWithAllParam(work *Work, fatherName string, age int) *Father {
	return &Father{fatherName: fatherName, age: age, work: work}
}

func NewFatherWithPre(work *Work) *Father {
	return &Father{fatherName: "张三", age: 33, work: work}
}
func NewFatherWithStruct(work Work) *Father {
	return &Father{fatherName: "张三", age: 33, work: &work}
}

type Mother struct {
	motherName string
	age        int
	homework   Homework
}

func NewMother(motherName string, age int, homework Homework) *Mother {
	return &Mother{motherName: motherName, age: age, homework: homework}
}


type Child struct {
	childName string
	age       int
}

func NewChild(childName string, age int) *Child {
	return &Child{childName: childName, age: age}
}

type Family struct {
	familyName string
	father   Father
	children []Child
	mother   Mother
}

func NewFamily(familyName string, father Father, children []Child, mother Mother) *Family {
	return &Family{familyName: familyName, father: father, children: children, mother: mother}
}

