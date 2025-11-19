package main

import (
	"fmt"
	"reflect"
	"strings"
)

// ===== БАЗОВЫЕ СТРУКТУРЫ ДЛЯ ПРИМЕРОВ =====
type User struct {
	ID       int                    `json:"id" db:"user_id" validate:"required,min=1"`
	Name     string                 `json:"name" db:"user_name" validate:"required"`
	Email    string                 `json:"email" db:"email" validate:"email"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Admin struct {
	User
	Level  int    `json:"level" validate:"min=1,max=10"`
	secret string // приватное поле
}

// ===== ОСНОВЫ REFLECT =====
func reflectBasics() {
	fmt.Println("\n=== ОСНОВЫ REFLECT ===")

	user := User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
		Tags:  []string{"go", "backend"},
	}

	// получение Type и Value
	t := reflect.TypeOf(user)
	v := reflect.ValueOf(user)

	fmt.Printf("Type: %v\n", t)
	fmt.Printf("Kind: %v\n", t.Kind())
	fmt.Printf("Value: %v\n", v)
	fmt.Printf("Is struct: %v\n", t.Kind() == reflect.Struct)
}

// ===== ИССЛЕДОВАНИЕ СТРУКТУРЫ =====
func inspectStruct(obj interface{}) {
	fmt.Println("\n=== ИССЛЕДОВАНИЕ СТРУКТУРЫ ===")

	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	// проверка, что это структура
	if t.Kind() != reflect.Struct {
		fmt.Println("Not a struct!")
		return
	}

	fmt.Printf("Struct: %s\n", t.Name())
	fmt.Printf("Fields count: %d\n", t.NumField())

	// итерация по полям
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		fmt.Printf("\nField %d: %s\n", i, field.Name)
		fmt.Printf("  Type: %v\n", field.Type)
		fmt.Printf("  Kind: %v\n", field.Type.Kind())
		fmt.Printf("  CanSet: %v\n", fieldValue.CanSet())
		fmt.Printf("  IsExported: %v\n", field.IsExported())

		// БЕЗОПАСНОЕ получение значения
		if field.IsExported() {
			// для экспортированных полей можно использовать Interface()
			fmt.Printf("  Value: %v\n", fieldValue.Interface())
		} else {
			// для приватных полей показываем только тип
			fmt.Printf("  Value: <unexported field: %s>\n", field.Type)
		}

		// чтение тегов
		if tag := field.Tag; tag != "" {
			fmt.Printf("  Tags:\n")
			if jsonTag := tag.Get("json"); jsonTag != "" {
				fmt.Printf("    json: %s\n", jsonTag)
			}
			if dbTag := tag.Get("db"); dbTag != "" {
				fmt.Printf("    db: %s\n", dbTag)
			}
			if validateTag := tag.Get("validate"); validateTag != "" {
				fmt.Printf("    validate: %s\n", validateTag)
			}
		}
	}
}

// ===== БЕЗОПАСНАЯ РАБОТА С ПРИВАТНЫМИ ПОЛЯМИ =====
func inspectStructSafe(obj interface{}) {
	fmt.Println("\n=== БЕЗОПАСНОЕ ИССЛЕДОВАНИЕ ===")

	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	if t.Kind() != reflect.Struct {
		fmt.Println("Not a struct!")
		return
	}

	fmt.Printf("Struct: %s\n", t.Name())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		fmt.Printf("\nField: %s\n", field.Name)
		fmt.Printf("  Type: %v\n", field.Type)
		fmt.Printf("  Exported: %v\n", field.IsExported())

		if field.IsExported() {
			// безопасный доступ к экспортированным полям
			fmt.Printf("  Value: %v\n", fieldValue.Interface())
		} else {
			// для приватных полей - только информация о типе
			fmt.Printf("  Value: [PRIVATE - type: %v]\n", field.Type)

			// демонстрация: можно проверять конкретные типы без доступа к значению
			switch field.Type.Kind() {
			case reflect.String:
				fmt.Printf("  Hint: private string field\n")
			case reflect.Int:
				fmt.Printf("  Hint: private int field\n")
			case reflect.Bool:
				fmt.Printf("  Hint: private bool field\n")
			default:
				fmt.Printf("  Hint: private %s field\n", field.Type.Kind())
			}
		}
	}
}

// ===== ИЗМЕНЕНИЕ ЗНАЧЕНИЙ =====
func modifyStruct(obj interface{}) {
	fmt.Println("\n=== ИЗМЕНЕНИЕ ЗНАЧЕНИЙ ===")

	v := reflect.ValueOf(obj)

	// должен передавать указатель для изменения!
	if v.Kind() != reflect.Ptr {
		fmt.Println("Must pass a pointer to modify!")
		return
	}

	// разыменовываем указатель
	v = v.Elem()

	// изменяем публичные поля
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		if !field.CanSet() {
			fmt.Printf("Field %s: cannot set (private or invalid)\n", fieldType.Name)
			continue
		}

		// изменяем в зависимости от типа
		switch field.Kind() {
		case reflect.String:
			if fieldType.Name == "Name" {
				field.SetString("Modified " + field.String())
				fmt.Printf("Modified %s to: %s\n", fieldType.Name, field.String())
			}
		case reflect.Int:
			if fieldType.Name == "ID" {
				field.SetInt(field.Int() + 100)
				fmt.Printf("Modified %s to: %d\n", fieldType.Name, field.Int())
			}
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				// добавляем элемент в срез
				newSlice := reflect.Append(field, reflect.ValueOf("reflected"))
				field.Set(newSlice)
				fmt.Printf("Modified %s to: %v\n", fieldType.Name, field.Interface())
			}
		}
	}
}

// ===== ВЫЗОВ МЕТОДОВ ЧЕРЕЗ REFLECT =====
type Calculator struct {
	Offset int
}

// value receiver для демонстрации
func (c Calculator) Add(a, b int) int {
	return a + b + c.Offset
}

// pointer receiver
func (c *Calculator) Multiply(a, b int) int {
	return a * b * c.Offset
}

// приватный метод
func (c *Calculator) privateMethod() {
	fmt.Println("This is private")
}

// публичный метод без параметров
func (c *Calculator) GetOffset() int {
	return c.Offset
}

func callMethods() {
	fmt.Println("\n=== ВЫЗОВ МЕТОДОВ ===")

	calc := Calculator{Offset: 10}
	v := reflect.ValueOf(calc)
	vPtr := reflect.ValueOf(&calc)

	fmt.Printf("Methods count (value): %d\n", v.NumMethod())
	fmt.Printf("Methods count (pointer): %d\n", vPtr.NumMethod())

	// вызов метода Add (value receiver)
	methodAdd := v.MethodByName("Add")
	if methodAdd.IsValid() {
		result := methodAdd.Call([]reflect.Value{
			reflect.ValueOf(5),
			reflect.ValueOf(3),
		})
		fmt.Printf("Add result: %v\n", result[0].Int())
	}

	// вызов метода Multiply (pointer receiver)
	methodMultiply := vPtr.MethodByName("Multiply")
	if methodMultiply.IsValid() {
		result := methodMultiply.Call([]reflect.Value{
			reflect.ValueOf(5),
			reflect.ValueOf(3),
		})
		fmt.Printf("Multiply result: %v\n", result[0].Int())
	}

	// вызов метода без параметров
	methodGetOffset := vPtr.MethodByName("GetOffset")
	if methodGetOffset.IsValid() {
		result := methodGetOffset.Call(nil) // nil для методов без параметров
		fmt.Printf("GetOffset result: %v\n", result[0].Int())
	}

	// попытка вызова приватного метода (не сработает)
	privateMethod := vPtr.MethodByName("privateMethod")
	fmt.Printf("Private method valid: %v\n", privateMethod.IsValid())
}

// ===== РАБОТА С ТЭГАМИ - ВАЛИДАЦИЯ =====
func validateStruct(obj interface{}) []string {
	fmt.Println("\n=== ВАЛИДАЦИЯ ПО ТЭГАМ ===")

	var errors []string
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// пропускаем приватные поля
		if !fieldType.IsExported() {
			continue
		}

		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		rules := strings.Split(validateTag, ",")
		for _, rule := range rules {
			switch {
			case rule == "required":
				if isZero(field) {
					errors = append(errors,
						fmt.Sprintf("%s is required", fieldType.Name))
				}
			case rule == "email":
				if field.Kind() == reflect.String {
					email := field.String()
					if !strings.Contains(email, "@") {
						errors = append(errors,
							fmt.Sprintf("%s must be a valid email", fieldType.Name))
					}
				}
			case strings.HasPrefix(rule, "min="):
				var min int
				fmt.Sscanf(rule, "min=%d", &min)
				if field.Kind() == reflect.Int && field.Int() < int64(min) {
					errors = append(errors,
						fmt.Sprintf("%s must be at least %d", fieldType.Name, min))
				}
			case strings.HasPrefix(rule, "max="):
				var max int
				fmt.Sscanf(rule, "max=%d", &max)
				if field.Kind() == reflect.Int && field.Int() > int64(max) {
					errors = append(errors,
						fmt.Sprintf("%s must be at most %d", fieldType.Name, max))
				}
			}
		}
	}

	if len(errors) > 0 {
		fmt.Println("Validation errors:", errors)
	} else {
		fmt.Println("Validation passed!")
	}

	return errors
}

// проверка на нулевое значение
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	default:
		return false
	}
}

// ===== ДИНАМИЧЕСКОЕ СОЗДАНИЕ СТРУКТУР =====
func createDynamic() {
	fmt.Println("\n=== ДИНАМИЧЕСКОЕ СОЗДАНИЕ ===")

	// создание slice через reflect
	intSliceType := reflect.SliceOf(reflect.TypeOf(0))
	intSlice := reflect.MakeSlice(intSliceType, 0, 10)
	intSlice = reflect.Append(intSlice, reflect.ValueOf(1))
	intSlice = reflect.Append(intSlice, reflect.ValueOf(2))
	intSlice = reflect.Append(intSlice, reflect.ValueOf(3))
	fmt.Printf("Dynamic slice: %v\n", intSlice.Interface())

	// создание map через reflect
	mapType := reflect.MapOf(
		reflect.TypeOf(""),
		reflect.TypeOf(0),
	)
	stringMap := reflect.MakeMap(mapType)

	// добавляем несколько пар ключ-значение
	keys := []string{"answer", "count", "score"}
	values := []int{42, 7, 100}

	for i, key := range keys {
		stringMap.SetMapIndex(
			reflect.ValueOf(key),
			reflect.ValueOf(values[i]),
		)
	}
	fmt.Printf("Dynamic map: %v\n", stringMap.Interface())

	// создание указателя через reflect
	intPtrType := reflect.PtrTo(reflect.TypeOf(0))
	intPtr := reflect.New(intPtrType.Elem())
	intPtr.Elem().SetInt(255)
	fmt.Printf("Dynamic pointer: %v -> %v\n", intPtr.Interface(), intPtr.Elem().Interface())
}

// ===== ПРАКТИЧЕСКИЙ ПРИМЕР: СЕРИАЛИЗАЦИЯ В MAP =====
func structToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !fieldType.IsExported() {
			continue // Пропускаем приватные поля
		}

		// используем json тег как ключ, или имя поля
		key := fieldType.Name
		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
			key = strings.Split(jsonTag, ",")[0]
		}

		result[key] = field.Interface()
	}

	return result
}

func main() {
	// базовые операции
	reflectBasics()

	// исследование структуры User
	user := User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
		Tags:  []string{"golang", "backend"},
		Metadata: map[string]interface{}{
			"role":    "admin",
			"created": "2024",
		},
	}
	inspectStruct(user)

	// безопасное исследование структуры Admin с приватными полями
	admin := Admin{
		User:   user,
		Level:  5,
		secret: "super-secret-key",
	}
	inspectStructSafe(admin)

	// изменение структуры
	fmt.Println("\n--- До модификации ---")
	fmt.Printf("User: %+v\n", user)
	modifyStruct(&user)
	fmt.Printf("--- После модификации ---\n")
	fmt.Printf("User: %+v\n", user)

	// вызов методов через рефлексию
	callMethods()

	// валидация структуры
	fmt.Println("\n--- Тест валидации ---")
	invalidUser := User{
		ID:    0,               // нарушает required,min=1
		Name:  "",              // нарушает required
		Email: "invalid-email", // нарушает email
	}
	validateStruct(&invalidUser)

	// динамическое создание объектов
	createDynamic()

	// сериализация в map
	fmt.Println("\n=== СЕРИАЛИЗАЦИЯ В MAP ===")
	serialized := structToMap(&user)
	fmt.Printf("User as map: %+v\n", serialized)

	// дополнительные примеры
	fmt.Println("\n=== ДОПОЛНИТЕЛЬНЫЕ ВОЗМОЖНОСТИ ===")

	// проверка наличия поля
	adminValue := reflect.ValueOf(admin)
	levelField := adminValue.FieldByName("Level")
	if levelField.IsValid() {
		fmt.Printf("Field 'Level' exists: %v\n", levelField.Interface())
	}

	secretField := adminValue.FieldByName("secret")
	if secretField.IsValid() {
		fmt.Printf("Field 'secret' exists but is private\n")
	}

	// получение имени типа
	fmt.Printf("Type name of user: %v\n", reflect.TypeOf(user).Name())
	fmt.Printf("Type name of admin: %v\n", reflect.TypeOf(admin).Name())
}
