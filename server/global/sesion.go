package global

import "sync"

// Mapa global para almacenar el estado de las sesiones por partición
var sessionMap = make(map[string]bool)
var mutex sync.Mutex // Para sincronización

// Función para iniciar sesión en una partición (marcar sesión activa en esa partición)
func ActivateSession(partitionID string) {
    mutex.Lock()
    sessionMap[partitionID] = true
    mutex.Unlock()
}

// Función para cerrar todas las sesiones (desactivar todas las particiones)
func DeactivateSession() {
    mutex.Lock()
    for partitionID := range sessionMap {
        sessionMap[partitionID] = false
    }
    mutex.Unlock()
}

// Función para verificar si la sesión está activa en una partición específica
func IsSessionActive(partitionID string) bool {
    mutex.Lock()
    defer mutex.Unlock()
    return sessionMap[partitionID]
}

//función para verificar si hay alguna sesión activa en alguna partición
func IsAnySessionActive() bool {
    mutex.Lock()
    defer mutex.Unlock()
    for _, active := range sessionMap {
        if active {
            return true
        }
    }
    return false
}

func GetIDSession() string {
    mutex.Lock()
    defer mutex.Unlock()
    for partitionID, active := range sessionMap {
        if active {
            return partitionID
        }
    }
    return ""
}