package mock

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

func NewFather(fatherName string, age int, work *Work) *Father {
	return &Father{fatherName: fatherName, age: age, work: work}
}
type Mother struct {
	motherName string
	age        int
	homework   Homework
}

type Child struct {
	childName string
	age       int
}

type Family struct {
	father   Father
	mother   Mother
	children []Child
}
