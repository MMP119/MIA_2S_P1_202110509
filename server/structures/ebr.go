package structures


import(
	"encoding/binary"
	"fmt"
	"os"
)


//structura del EBR

type EBR struct {
	Part_mount	[1]byte
	Part_fit	[1]byte
	Part_start	int32
	Part_size	int32
	Part_next	int32
	Part_name	[16]byte
}


// Serializa el EBR y lo escribe en el disco en la posición dada
func (ebr *EBR) SerializeEBR(path string, position int32) (string, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return "error abriendo el archivo para serializacion el ebr",fmt.Errorf("error abriendo el archivo: %s", err)
	}
	defer file.Close()

	file.Seek(int64(position), 0)
	err = binary.Write(file, binary.BigEndian, ebr)
	if err != nil {
		return "error escribiendo el EBR",fmt.Errorf("error escribiendo el EBR: %s", err)
	}
	return "", nil
}

// Deserializa el EBR desde el disco en la posición dada
func (ebr *EBR) DeserializeEBR(path string, position int32) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "Error abriendo el archivo para deserialización del EBR", fmt.Errorf("error abriendo el archivo: %s", err)
	}
	defer file.Close()

	file.Seek(int64(position), 0)
	err = binary.Read(file, binary.BigEndian, ebr)
	if err != nil {
		return "Error leyendo el EBR", fmt.Errorf("error leyendo el EBR: %s", err)
	}

	// Verificar si el EBR es válido (por ejemplo, si Part_size es mayor que 0)
	if ebr.Part_size <= 0 {
		return "EBR inválido o no encontrado", fmt.Errorf("EBR inválido o no encontrado")
	}

	return "", nil
}

// imprime las particiones lógicas dentro de la extendida
func ImprimirParticionesLogicas (path string, startExtendida int32){
	
	var ebr EBR
	posicionActual := startExtendida
	fmt.Println("\nParticiones lógicas dentro de la extendida:")

	for {
		
		//deseralizar el ebr en la posicion actual
		msg, err := ebr.DeserializeEBR(path, posicionActual)
		if err != nil {
			fmt.Println(msg)
			return 
		}	
			

		// Verificar si el EBR es válido (es decir, si tiene una partición lógica)
		if ebr.Part_size > 0 {
			fmt.Printf("Nombre: %s, Inicio: %d, Tamaño: %d, Siguiente EBR: %d\n\n\n",
				string(ebr.Part_name[:]), ebr.Part_start, ebr.Part_size, ebr.Part_next)
			
		}


		// Si Part_next es -1, no hay más particiones lógicas
		if ebr.Part_next == -1 {
			break
		}

		// Actualizar la posición actual para leer el siguiente EBR
		posicionActual = ebr.Part_next	
	}

}
