package commands

import (
	"fmt"
	global "server/global"
)

func ParseLogout(tokens []string) (string, string, error) {
	if len(tokens) != 0 {
		return "", "", fmt.Errorf("logout: número de parámetros incorrecto")
	}

	if !global.IsSessionActive() {
		return "", "Comando LOGOUT: ERROR: No hay ninguna sesión activa", nil
	}

	global.DeactivateSession()
	return "", "Comando LOGOUT realizado correctamente: Sesión cerrada exitosamente", nil
}