package commands

import (
	"errors"
	"fmt"
	"regexp"
	structures "server/structures"
	"strings"
)

func ParseRep(tokens []string)(*structures.REP, string, error){
	
	cmd := &structures.REP{}

	args := strings.Join(tokens, " ")

	re := regexp.MustCompile(`-size=\d+|-unit=[kKmM]|-fit=[bBfF]{2}|-path="[^"]+"|-path=[^\s]+|-type=[pPeElL]|-name="[^"]+"|-name=[^\s]+`)

	matches := re.FindAllString(args, -1)

	
	for _, math := range matches{
		kv := strings.SplitN(math, "=", 2)
		if len(kv) != 2 {
			return nil, "ERROR: formato de parámetro inválido", fmt.Errorf("formato de parámetro inválido: %s", math)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		switch key{
			case "-path":
				if value == "" {
					return nil, "ERROR: el path es obligatorio", errors.New("el path es obligatorio")
				}
				cmd.Path = value
			
			default: 
				return nil, "ERROR: parámetro no reconocido", fmt.Errorf("parámetro no reconocido: %s", key)
		}
	}

	if cmd.Path == "" {
		return nil, "ERROR: el path es obligatorio", errors.New("el path es obligatorio")
	}

	msg, err := structures.CommandRep(cmd)
	if err != nil {
		return nil, msg, err
	}

	
	return cmd, "", nil
}