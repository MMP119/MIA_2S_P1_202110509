package global

import "sync"

// Variable global para almacenar el estado de la sesión
var sessionActive bool
var mutex sync.Mutex // Para sincronización en caso de múltiples goroutines

// Función para iniciar sesión (marcar la sesión como activa)
func ActivateSession() {
    mutex.Lock()
    sessionActive = true
    mutex.Unlock()
}

// Función para cerrar sesión (marcar la sesión como inactiva)
func DeactivateSession() {
    mutex.Lock()
    sessionActive = false
    mutex.Unlock()
}

// Función para verificar si la sesión está activa
func IsSessionActive() bool {
    mutex.Lock()
    defer mutex.Unlock()
    return sessionActive
}