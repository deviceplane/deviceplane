package utils

import "github.com/deviceplane/deviceplane/pkg/models"

// map[DeviceID]map[ApplicationID]map[Service]*State
func DeviceServiceStatesListToMap(deviceServiceStates []models.DeviceServiceState) (map[string]map[string]map[string]*models.DeviceServiceState, error) {
	stateMap := make(map[string]map[string]map[string]*models.DeviceServiceState)

	for i, deviceServiceState := range deviceServiceStates {
		if _, exists := stateMap[deviceServiceState.DeviceID]; !exists {
			stateMap[deviceServiceState.DeviceID] = make(map[string]map[string]*models.DeviceServiceState)
		}
		if _, exists := stateMap[deviceServiceState.DeviceID][deviceServiceState.ApplicationID]; !exists {
			stateMap[deviceServiceState.DeviceID][deviceServiceState.ApplicationID] = make(map[string]*models.DeviceServiceState)
		}
		stateMap[deviceServiceState.DeviceID][deviceServiceState.ApplicationID][deviceServiceState.Service] = &deviceServiceStates[i]
	}
	return stateMap, nil
}

// map[DeviceID]map[ApplicationID]*Status
func DeviceApplicationStatusesListToMap(deviceApplicationStatuses []models.DeviceApplicationStatus) (map[string]map[string]*models.DeviceApplicationStatus, error) {
	statusMap := make(map[string]map[string]*models.DeviceApplicationStatus)

	for i, deviceApplicationStatus := range deviceApplicationStatuses {
		if _, exists := statusMap[deviceApplicationStatus.DeviceID]; !exists {
			statusMap[deviceApplicationStatus.DeviceID] = make(map[string]*models.DeviceApplicationStatus)
		}
		statusMap[deviceApplicationStatus.DeviceID][deviceApplicationStatus.ApplicationID] = &deviceApplicationStatuses[i]
	}
	return statusMap, nil
}

// map[ApplicationID]*Release
func ReleasesListToMap(releases []models.Release) (map[string]*models.Release, error) {
	releaseMap := make(map[string]*models.Release)

	for i, release := range releases {
		releaseMap[release.ApplicationID] = &releases[i]
	}
	return releaseMap, nil
}
