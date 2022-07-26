# GOSTORAGE

## Описание
Пакет предоставляет простое "ключ-значение" generic хранилище объектов. Есть возможность указания времени жизни объектов, переодической очистки устпревших данных, секционирования и записи/считывания данных с диска.

## Пример использования

    // Создание экземпляра с указанием типа данных
    // временем жизни и переодической очисткой данных
	stor := NewStorage[string]().DefaultExpiration(5 * time.Second).WithCleaner(5 * time.Minute)

	// Запись значения
	stor.Set("key", "value")

    // Чтение значения
    if value, ok := stor.Get("key"); ok {
		fmt.Printf("value - %v", value)
	} else if value != testValue {
		fmt.Print("there is no value")
	}

    // Сохранение в файл
	err := stor.SaveFile("savefile")
	if err != nil {
		// обработка ошибок
	}

	// Создание из файла
	stor2 := NewStorage[string]()
	err = stor2.LoadFile("savefile")
	if err != nil {
		// обработка ошибок
	}

	// Хеш секционирование
	stor := NewStorageShards[string](5)
	stor.Set(testKey, testValue)
