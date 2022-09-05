# GOSTORAGE

## Описание
Пакет предоставляет простое "ключ-значение" generic хранилище объектов. Есть возможность указания времени жизни объектов, переодической очистки устпревших данных, секционирования и записи/считывания данных с диска.

## Пример использования

	// Создание экземпляра с указанием типа данных
	// временем жизни и переодической очисткой данных
	newStorage := NewStorage[string]().WithExpiration(5 * time.Second).WithCleaner(5 * time.Minute)

	// Запись значения
	newStorage.Set("key", "value")

	// Чтение значения
	if value, ok := newStorage.Get("key"); ok {
		fmt.Printf("value - %v", value)
	} else if value != testValue {
		fmt.Print("there is no value")
	}

	value, ok := stor.GetFetch(testKey, func(s string) (string, error) {
		return testValue, nil
	})

	// Сохранение в файл
	err := newStorage.SaveFile("savefile")
	if err != nil {
		// обработка ошибок
	}

	// Создание из файла
	newStorage2 := NewStorage[string]()
	err = newStorage2.LoadFile("savefile")
	if err != nil {
		// обработка ошибок
	}

	// Хеш секционирование
	newStorageShards := NewStorageShards[string](5)
	stor.Set(testKey, testValue)
