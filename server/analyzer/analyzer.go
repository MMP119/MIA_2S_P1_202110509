package analyzer

import (
	commands "server/commands"                  
	"fmt"                       
	"os"                        
	"os/exec"                   
	"strings"                  
)


func Analyzer(inputs []string) ([]string, []string) {
    var results []string
    var errors []string

    for i, input := range inputs {

        //ignorar líneas en blanco y comentarios
        inputs := strings.TrimSpace(input)
        if inputs == "" || strings.HasPrefix(inputs, "#") {
            continue //ignorar comentarios y líneas en blanco
        }

        tokens := strings.Fields(input)
        if len(tokens) == 0 {
            errors = append(errors, fmt.Sprintf("Comando %d: No se proporcionó ningún comando", i))
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
            errors = append(errors, err.Error())
        } else {
            results = append(results, msg)
        }
    }

    return results, errors
}
