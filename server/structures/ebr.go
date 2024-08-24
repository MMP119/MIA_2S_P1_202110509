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