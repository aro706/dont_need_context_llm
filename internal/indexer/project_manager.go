package indexer

import "sync"

var activeProjects = make(map[string]string)
var projectMutex sync.Mutex

func AddProject(path, name string) {
	projectMutex.Lock()
	defer projectMutex.Unlock()
	activeProjects[path] = name
}

func GetAllProjects() map[string]string {
	projectMutex.Lock()
	defer projectMutex.Unlock()

	// Return a copy to avoid concurrent map reads/writes
	copyMap := make(map[string]string)
	for k, v := range activeProjects {
		copyMap[k] = v
	}
	return copyMap
}
