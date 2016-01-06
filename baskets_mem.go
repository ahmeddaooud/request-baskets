package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

/// Basket interface ///

type memoryBasket struct {
	sync.RWMutex
	token      string
	config     BasketConfig
	requests   []*RequestData
	totalCount int
}

func (basket *memoryBasket) applyLimit() {
	// Keep requests up to specified capacity
	if len(basket.requests) > basket.config.Capacity {
		basket.requests = basket.requests[:basket.config.Capacity]
	}
}

func (basket *memoryBasket) Config() BasketConfig {
	return basket.config
}

func (basket *memoryBasket) Update(config BasketConfig) {
	basket.Lock()
	defer basket.Unlock()

	basket.config = config
	basket.applyLimit()
}

func (basket *memoryBasket) Authorize(token string) bool {
	return token == basket.token
}

func (basket *memoryBasket) Add(req *http.Request) *RequestData {
	basket.Lock()
	defer basket.Unlock()

	data := ToRequestData(req)
	// insert in front of collection
	basket.requests = append([]*RequestData{data}, basket.requests...)

	// keep total number of all collected requests
	basket.totalCount++
	// apply limits according to basket capacity
	basket.applyLimit()

	return data
}

func (basket *memoryBasket) Clear() {
	basket.Lock()
	defer basket.Unlock()

	// reset collected requests and total counter
	basket.requests = make([]*RequestData, 0, basket.config.Capacity)
	basket.totalCount = 0
}

func (basket *memoryBasket) Size() int {
	return len(basket.requests)
}

func (basket *memoryBasket) GetRequests(max int, skip int) RequestsPage {
	basket.RLock()
	defer basket.RUnlock()

	size := basket.Size()
	last := skip + max

	requestsPage := RequestsPage{
		Count:      size,
		TotalCount: basket.totalCount,
		HasMore:    last < size}

	if skip < size {
		if last > size {
			last = size
		}
		requestsPage.Requests = basket.requests[skip:last]
	}

	return requestsPage
}

/// BasketsDatabase interface ///

type memoryDatabase struct {
	sync.RWMutex
	baskets map[string]*memoryBasket
	names   []string
}

func (db *memoryDatabase) Create(name string, config BasketConfig) (BasketAuth, error) {
	db.Lock()
	defer db.Unlock()

	auth := BasketAuth{}

	_, exists := db.baskets[name]
	if exists {
		return auth, fmt.Errorf("Basket with name '%s' already exists", name)
	}

	token, err := GenerateToken()
	if err != nil {
		return auth, fmt.Errorf("Failed to generate token: %s", err.Error())
	}

	basket := new(memoryBasket)
	basket.token = token
	basket.config = config
	basket.requests = make([]*RequestData, 0, config.Capacity)
	basket.totalCount = 0

	db.baskets[name] = basket
	db.names = append(db.names, name)
	// Uncomment if sorting is expected
	// sort.Strings(db.names)

	auth.Token = token

	return auth, nil
}

func (db *memoryDatabase) Get(name string) Basket {
	basket, exists := db.baskets[name]
	if exists {
		return basket
	} else {
		return nil
	}
}

func (db *memoryDatabase) Delete(name string) {
	db.Lock()
	defer db.Unlock()

	delete(db.baskets, name)
	for i, v := range db.names {
		if v == name {
			db.names = append(db.names[:i], db.names[i+1:]...)
			break
		}
	}
}

func (db *memoryDatabase) Size() int {
	return len(db.names)
}

func (db *memoryDatabase) GetNames(max int, skip int) BasketNamesPage {
	db.RLock()
	defer db.RUnlock()

	size := len(db.names)
	last := skip + max

	namesPage := BasketNamesPage{
		Count:   size,
		HasMore: last < size}

	if skip < size {
		if last > size {
			last = size
		}

		namesPage.Names = db.names[skip:last]
	}

	return namesPage
}

func (db *memoryDatabase) Release() {
	log.Printf("Releasing in-memory database resources")
}

// NewMemoryDatabase creates an instance of in-memory Baskets Database
func NewMemoryDatabase() BasketsDatabase {
	return &memoryDatabase{baskets: make(map[string]*memoryBasket), names: make([]string, 0)}
}