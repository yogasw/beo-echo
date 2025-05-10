package repositories

import (
	"sync"

	"mockoon-control-panel/backend_new/src/types"
)

// MockInstanceRepository manages running mock instances
type MockInstanceRepository struct {
	instances map[int]types.MockInstance
	mutex     sync.RWMutex
}

var MockInstanceRepo = &MockInstanceRepository{
	instances: make(map[int]types.MockInstance),
}

// Add adds a new mock instance to the repository
func (r *MockInstanceRepository) Add(port int, instance types.MockInstance) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.instances[port] = instance
}

// Get retrieves a mock instance by port
func (r *MockInstanceRepository) Get(port int) (types.MockInstance, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	instance, exists := r.instances[port]
	return instance, exists
}

// Remove removes a mock instance from the repository
func (r *MockInstanceRepository) Remove(port int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if instance, exists := r.instances[port]; exists {
		// Close log file if it exists
		if instance.LogFile != nil {
			instance.LogFile.Close()
		}

		delete(r.instances, port)
	}
}

// GetAll returns all running mock instances
func (r *MockInstanceRepository) GetAll() []types.MockInstance {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	instances := make([]types.MockInstance, 0, len(r.instances))
	for _, instance := range r.instances {
		instances = append(instances, instance)
	}

	return instances
}

// IsRunning checks if a mock server is running on a specific port
func (r *MockInstanceRepository) IsRunning(port int) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.instances[port]
	return exists
}

// IsRunningWithConfig checks if a mock server with a specific config is running
func (r *MockInstanceRepository) IsRunningWithConfig(configFile string) (bool, int) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for port, instance := range r.instances {
		if instance.ConfigFile == configFile {
			return true, port
		}
	}

	return false, 0
}
