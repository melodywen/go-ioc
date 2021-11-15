package mock

type Animal struct {
	name     string
	age      int
	category string
}

func NewAnimal(name string, age int, category string) *Animal {
	return &Animal{name: name, age: age, category: category}
}

func NewAnimalAndParam(name string, age int, category string) (*Animal, string, int, string) {
	return &Animal{name: name, age: age, category: category}, name, age, category
}

func (animal *Animal) Say() string {
	return animal.name
}
