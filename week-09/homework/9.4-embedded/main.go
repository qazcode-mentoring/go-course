package main

import "fmt"

// Person представляет базовую информацию о человеке
type Person struct {
	Name string
	Age  int
}

// Greet возвращает приветствие
func (p Person) Greet() string {
	// TODO: реализуй метод
	// Формат: "Привет, меня зовут {Name}"
	greed := fmt.Sprintf("Привет, меня зовут %v", p.Name)
	return greed
}

// Employee представляет сотрудника (встраивает Person)
type Employee struct {
	Person   // встраивание — поля и методы Person доступны напрямую
	Position string
	Salary   float64
}

// Work возвращает описание работы сотрудника
func (e Employee) Work() string {
	// TODO: реализуй метод
	// Формат: "{Name} работает как {Position}"
	w := fmt.Sprintf("%v работает как %v", e.Name, e.Position)
	return w
}

// GiveRaise повышает зарплату на указанный процент
func (e *Employee) GiveRaise(percent float64) {
	// TODO: реализуй метод
	// Увеличь Salary на percent процентов
	p := (percent * e.Salary) / 100
	e.Salary += p
}

// Manager представляет менеджера (встраивает Employee)
type Manager struct {
	Employee // встраивание Employee (который встраивает Person)
	Team     []string
}

// AddToTeam добавляет сотрудника в команду
func (m *Manager) AddToTeam(name string) {
	// TODO: реализуй метод
	m.Team = append(m.Team, name)
}

// TeamSize возвращает размер команды
func (m Manager) TeamSize() int {
	// TODO: реализуй метод
	return len(m.Team)
}

func main() {
	fmt.Println("=== Встраивание структур ===")

	// Создаём менеджера
	manager := Manager{
		Employee: Employee{
			Person:   Person{Name: "Алексей", Age: 35},
			Position: "Tech Lead",
			Salary:   150000,
		},
		Team: []string{},
	}

	// Методы Person (доступны напрямую через встраивание)
	fmt.Println("\n--- Методы Person ---")
	fmt.Println(manager.Greet())
	fmt.Println("Имя:", manager.Name) // поле Person доступно напрямую
	fmt.Println("Возраст:", manager.Age)

	// Методы Employee
	fmt.Println("\n--- Методы Employee ---")
	fmt.Println(manager.Work())
	fmt.Printf("Зарплата до повышения: %.2f\n", manager.Salary)
	manager.GiveRaise(10)
	fmt.Printf("Зарплата после повышения на 10%%: %.2f\n", manager.Salary)

	// Методы Manager
	fmt.Println("\n--- Методы Manager ---")
	fmt.Println("Размер команды:", manager.TeamSize())
	manager.AddToTeam("Мария")
	manager.AddToTeam("Иван")
	manager.AddToTeam("Ольга")
	fmt.Println("После добавления сотрудников:")
	fmt.Println("Команда:", manager.Team)
	fmt.Println("Размер команды:", manager.TeamSize())
}
