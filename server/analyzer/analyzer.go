package analyzer

import (
	commands "server/commands"                  
	"fmt"                       
	"os"                        
	"os/exec"                   
	"strings"                  
)

// Analyzer analiza el comando de entrada y ejecuta la acción correspondiente
// Analyzer analiza múltiples comandos de entrada y ejecuta las acciones correspondientes
func Analyzer(inputs []string) (map[string]string, map[string]string) {
    results := make(map[string]string)
    errors := make(map[string]string)

    for i, input := range inputs {
        tokens := strings.Fields(input)
        if len(tokens) == 0 {
            errors[fmt.Sprintf("Comando %d", i)] = "No se proporcionó ningún comando"
            continue
        }
        tokens[0] = strings.ToLower(tokens[0])
        var msg string
        var err error

        switch tokens[0] {
        case "mkdisk":
            _, msg, err = commands.ParserMkdisk(tokens[1:])
        case "rmdisk":
            _, msg, err = commands.ParserRmdisk(tokens[1:])
        case "fdisk":
            _, msg, err = commands.ParserFdisk(tokens[1:])
        case "mount":
            _, msg, err = commands.ParserMount(tokens[1:])
        case "clear":
            cmd := exec.Command("clear")
            cmd.Stdout = os.Stdout
            err = cmd.Run()
            if err == nil {
                msg = "Terminal limpia"
            }
        default:
            err = fmt.Errorf("comando desconocido: %s", tokens[0])
        }

        if err != nil {
            errors[fmt.Sprintf("Comando %d", i)] = err.Error()
        } else {
            results[fmt.Sprintf("Comando %d", i)] = msg
        }
    }

    return results, errors
}
