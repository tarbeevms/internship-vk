package data

import (
	"errors"
	"sync"

	"github.com/tarantool/go-tarantool"
)

// SetData записывает данные в спейс data (если ключ уже существует, происходит перезапись)
func (dr *DataRepository) SetData(key string, value interface{}) error {
	_, err := dr.tConn.ReplaceAsync("data", []interface{}{key, value}).Get()
	if err != nil {
		return err
	}
	return nil
}

// GetData получает данные из спейса data по ключу
func (dr *DataRepository) GetData(key string) (interface{}, error) {
	resp, err := dr.tConn.SelectAsync("data", "primary", 0, 1, tarantool.IterEq, []interface{}{key}).Get()
	if err != nil {
		return nil, err
	}
	if len(resp.Tuples()) == 0 {
		return nil, errors.New("data not found")
	}
	return resp.Tuples()[0][1], nil
}

// SetDataBatch записывает данные пачками
func (dr *DataRepository) SetDataBatch(data map[string]interface{}) (batchErr map[string]error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	//хранит, какой ключ вызвал какую ошибку
	batchErr = make(map[string]error)

	for key, value := range data {
		wg.Add(1)
		go func(key string, value interface{}) {
			defer wg.Done()
			err := dr.SetData(key, value)
			if err != nil {
				mu.Lock()
				batchErr[key] = err
				mu.Unlock()
			}
		}(key, value)
	}

	wg.Wait()
	return batchErr
}

// GetDataBatch читает данные пачками по ключам
func (dr *DataRepository) GetDataBatch(keys []string) (result map[string]interface{}, batchErr map[string]error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	result = make(map[string]interface{})
	//хранит, какой ключ вызвал какую ошибку
	batchErr = make(map[string]error)

	for _, key := range keys {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			data, err := dr.GetData(key)
			if err != nil {
				if err.Error() == "data not found" {
					return
				}
				mu.Lock()
				batchErr[key] = err
				mu.Unlock()
				return
			}
			mu.Lock()
			result[key] = data
			mu.Unlock()
		}(key)
	}

	wg.Wait()
	return result, batchErr
}
