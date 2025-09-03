package exercises

import (
	"fmt"
	"sync"
	"time"
)

type MockConnection struct {
	id int
}

func NewMockConnection(id int) *MockConnection {
	fmt.Printf("Connection %d created\n", id)
	return &MockConnection{id: id}
}

func (c *MockConnection) Execute(query string) {
	fmt.Printf("Connection %d executing query %s\n ", c.id, query)
	time.Sleep(100 * time.Millisecond)
}

type ConnectionPool struct {
	connections chan *MockConnection
	size        int
	mu          sync.Mutex
}

func NewConnectionPool(size int) *ConnectionPool {
	pool := &ConnectionPool{
		connections: make(chan *MockConnection, size),
		size:        size,
	}
	for i := 0; i < size; i++ {
		pool.connections <- NewMockConnection(i)
	}
	return pool
}

func (p *ConnectionPool) GetConnection() *MockConnection {
	conn := <-p.connections
	fmt.Printf("Connection %d acquired. Connections in pool: %d\n", conn.id, len(p.connections))
	return conn
}

func (p *ConnectionPool) ReleaseConnection(conn *MockConnection) {
	p.connections <- conn
	fmt.Printf("Connection %d released. Connections in pool: %d\n", conn.id, len(p.connections))
}

func (p *ConnectionPool) CloseAll() {
	close(p.connections)
	for conn := range p.connections {
		fmt.Printf("Connection %d closed\n", conn.id)
	}
	fmt.Println("Connection pool closed")
}

func worker(id int, pool *ConnectionPool, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d trying to get a connection...\n", id)
	conn := pool.GetConnection()
	if conn != nil {
		defer pool.ReleaseConnection(conn)
		conn.Execute(fmt.Sprintf("SELECT * from users WHERE id=%d\n", id))
	} else {
		fmt.Printf("Worker %d could not get a connection\n", id)
	}
}

func ConnectionPoolExercise() {
	poolSize := 3
	numTasks := 10
	var wg sync.WaitGroup

	connPool := NewConnectionPool(poolSize)
	fmt.Printf("\nStarting %d tasks with a pool of size %d...\n\n", numTasks, poolSize)

	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		go worker(i, connPool, &wg)
		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println("All tasks completed. Cleaning up pool")
}
