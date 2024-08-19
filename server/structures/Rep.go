package structures

import(
	"fmt"
)


type REP struct {
	//name
	Path string 
	//id
	//path_file_ls
}

func CommandRep(rep *REP) (string, error) {
	
	mbr := MBR{}

	msg, err := mbr.DeserializeMBR(rep.Path)
	if err != nil {
		return msg, fmt.Errorf("error leyendo el MBR del disco: %s", err)
	}

	// se imprime la info del mbr
	fmt.Println("\nMBR:")
	mbr.Print()

	// info de las particiones
	fmt.Println("\nParticiones:")
	mbr.PrintPartitions()

	return "",nil
}