package mock

// Animal is a test struct
type Animal struct {
	name     string
	age      int
	category string
}

// NewAnimal animal construct
func NewAnimal(name string, age int, category string) *Animal {
	return &Animal{name: name, age: age, category: category}
}

// NewAnimalAndParam animal construct
func NewAnimalAndParam(name string, age int, category string) (*Animal, string, int, string) {
	return &Animal{name: name, age: age, category: category}, name, age, category
}

// Say a test method
func (animal *Animal) Say() string {
	return animal.name
}
