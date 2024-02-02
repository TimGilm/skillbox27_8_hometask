/*
Цель задания - научиться работать с композитными типами данных: структурами и картами
Что нужно сделать - напишите программу, которая считывает ввод с stdin, создаёт структуру
student и записывает указатель на структуру в хранилище map[studentName] *Student.
type Student struct {
name string
age int
grade int
}
Программа должна получать строки в бесконечном цикле, создать структуру *Student через функцию
newStudent, далее сохранить указатель на эту структуру в map, а после получения EOF (ctrl + d)
вывести на экран имена всех студентов из хранилища. Также необходимо реализовать методы put, get.
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// За образец возьмем пример из 7 урока 27 модуля

type Storage interface {
	Get(student *Student) bool
	Put(student *Student)
	Print()
}

// заглушка
type StubStorage struct{}

func (fs *StubStorage) Get(student *Student) bool {
	return true
}
func (fs *StubStorage) Put(student *Student) {
	return
}
func (fs *StubStorage) Print() {
	student := Student{name: "timur", age: 45, grade: 5}
	fmt.Println(student)
}

// создадим "базовую" структуру Student с необходимыми полями:
type Student struct {
	name  string
	age   int
	grade int
}

// функция newStudent создает ссылку на структуру Student, т.е. затем можно оперируя с этой
// ссылкой записать данные в структуру Student, а не записать случайно в ее копию.
func newStudent(name string, age int, grade int) *Student {
	return &Student{
		name:  name,
		age:   age,
		grade: grade,
	}
}

type App struct {
	repository Storage
}

//перепишем метод Run в соответствии с нашей задачей

func (a *App) Run() {
	for {
		if student, ok := a.inputNextStudent(); ok {
			a.storeStudent(student)
		} else {
			break
		}
	}
}

func (a *App) inputNextStudent() (*Student, bool) {
	fmt.Println("Введите имя, возраст и оценку студента через пробел или `end` для выхода (либо ctrl+D в терминале)")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		err := io.EOF
		if len(parts) == 3 {
			name := parts[0]
			age, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Повторите ввод, введен неверный формат данных")
				continue
			}
			grade, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Повторите ввод, введен неверный формат данных")
				continue
			}
			return newStudent(name, age, grade), true
		} else if len(parts) == 1 && parts[0] == "end" {
			fmt.Println("Завершение ввода данных, выход из программы")
			a.repository.Print()
			break
		} else if err != nil {
			fmt.Println("Завершение ввода данных, выход из программы")
			a.repository.Print()
			break
		} else {
			fmt.Println("Повторите ввод, введены недостоверные данные")
			continue
		}
	}
	return nil, false
}

func (a *App) storeStudent(student *Student) {
	//проверяем есть ли в карте студент с таким именем (так как по определению языка Go, в карте не может быть дублей)
	msg := "Студент с именем %v уже присутствует в хранилище\n"
	if a.repository.Get(student) { //для цели указанной выше применяем метод Add
		msg = "Данные студента с именем %v успешно добавлены в хранилище\n"
		a.repository.Put(student)
	}
	fmt.Printf(msg, student.name)
}

// Перепишем структуру MemStorage с хранилищем students под наше условие
type MemStorage struct {
	students map[string]*Student
}

// метод NewMemStore создает указатель на карту
func NewMemStore() *MemStorage {
	return &MemStorage{
		students: make(map[string]*Student),
	}
}

// перепишем метод Add под наше условие, метод добавляет карту с именем студента, которого не было раньше в хранилище,
// используя вспомогательный метод contains
func (ms *MemStorage) Get(student *Student) bool {
	if ms.contains(student) {
		return false //т.е. дословно, если ms.contains(studentName) - true, то не добавляем и возвращаем false
	}
	return true
}

// метод Put помещает структуру в карту
func (ms *MemStorage) Put(student *Student) {
	students := student.name
	ms.students[students] = student //добавляем структуру в карту
}

// метод print выводит на печать данные из карты в случае окончания цикла ввода данных
func (ms *MemStorage) Print() {
	fmt.Println("Вывод на печать введенных данных")
	counter := 1
	for _, v := range ms.students {
		fmt.Printf("%v.student name: %v, age: %v, grade: %v\n", counter, v.name, v.age, v.grade)
		counter += 1
	}
}

// данный метод определяет есть ли уже в хранилище карта с таким именем студента
func (ms *MemStorage) contains(student *Student) bool {
	for st, _ := range ms.students {
		if student.name == st {
			return true
		} //так как в карте может быть только один ключ с уникальным именем!
	}
	return false
}

func main() {
	//repository := &StubStorage{}
	repository := NewMemStore()
	app := &App{repository}
	app.Run()
}
