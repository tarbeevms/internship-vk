package logic

import "fmt"

// WriteDataBatch добавляет пачку данных в Tarantool.
func (ll *LogicLayer) WriteDataBatch(data map[string]interface{}) (err error) {

	// Валидация данных перед вставкой
	err = ValidateData(data)
	if err != nil {
		return err
	}

	batchErr := ll.DataRepo.SetDataBatch(data)

	if len(batchErr) != 0 {
		var errorMessages []string
		for key, batchError := range batchErr {
			errorMessages = append(errorMessages, fmt.Sprintf("(key: '%s', error: '%v') ", key, batchError))
		}
		return fmt.Errorf("Failed to write batch. Errors: %v", errorMessages)
	}

	return nil
}

// ReadDataBatch читает пачку данных по ключам из Tarantool.
func (ll *LogicLayer) ReadDataBatch(keys []string) (result map[string]interface{}, err error) {

	result, batchErr := ll.DataRepo.GetDataBatch(keys)

	if len(batchErr) != 0 {
		var errorMessages []string
		for key, readError := range batchErr {
			errorMessages = append(errorMessages, fmt.Sprintf("key: %s, error: %v", key, readError))
		}
		return nil, fmt.Errorf("Failed to read batch. Errors: %v", errorMessages)
	}

	return result, nil
}

// ValidateData проверяет, что значения являются только скалярными типами.
func ValidateData(data map[string]interface{}) error {
	for key, value := range data {
		switch value.(type) {
		case string, int, float64, bool:
			// Валидные скалярные типы, ничего не делаем
		default:
			// Если значение не скалярного типа, возвращаем ошибку
			return fmt.Errorf("Bad Request: invalid value type for key %s (only scalar values are allowed)", key)
		}
	}
	return nil
}
