package structures

import (
	"encoding/binary"
	"fmt"
	"os"
)

type PARTITION struct {
	Part_status      [1]byte  // Estado de la partición
	Part_type        [1]byte  // Tipo de partición
	Part_fit         [1]byte  // Ajuste de la partición
	Part_start       int32    // Byte de inicio de la partición
	Part_size        int32    // Tamaño de la partición
	Part_name        [16]byte // Nombre de la partición
	Part_correlative int32    // Correlativo de la partición
	Part_id          [4]byte    // ID de la partición
}

// 0 = creado, 1 = activo, 2 = Disponible

//crear particion con parámetros proporcionados
func (p *PARTITION) CreatePartition(partStart, partSize int, partType, partFit, partName string) {
	// Asignar status de la partición
	p.Part_status[0] = '0' // El valor '0' indica que la partición ha sido creada

	// Asignar el byte de inicio de la partición
	p.Part_start = int32(partStart)

	// Asignar el tamaño de la partición
	p.Part_size = int32(partSize)

	// Asignar el tipo de partición
	if len(partType) > 0 {
		p.Part_type[0] = partType[0]
	}

	// Asignar el ajuste de la partición
	if len(partFit) > 0 {
		p.Part_fit[0] = partFit[0]
	}

	// Asignar el nombre de la partición
	copy(p.Part_name[:], partName)
}

func (p *PARTITION) MountPartition(correlative int, id string) error {
	// Asignar correlativo a la partición
	p.Part_correlative = int32(correlative) + 1

	//cambiar el status de la particion, para que esté montada
	p.Part_status[0] = '1' //indica que la particion se ha montado

	// Asignar ID a la partición
	copy(p.Part_id[:], id)

	return nil
}

// Serializa la partición lógica y la escribe en el disco en la posición dada
func (p *PARTITION) SerializePartition(path string, position int32) (string, error) {
    file, err := os.OpenFile(path, os.O_RDWR, 0644)
    if err != nil {
        return "error abriendo el archivo para serializar la partición lógica", fmt.Errorf("error abriendo el archivo: %s", err)
    }
    defer file.Close()

    file.Seek(int64(position), 0)
    err = binary.Write(file, binary.BigEndian, p)
    if err != nil {
        return "error escribiendo la partición lógica", fmt.Errorf("error escribiendo la partición lógica: %s", err)
    }
    return "", nil
}

func (p *PARTITION) Print() {
	fmt.Printf("Part_status: %c\n", p.Part_status[0])
	fmt.Printf("Part_type: %c\n", p.Part_type[0])
	fmt.Printf("Part_fit: %c\n", p.Part_fit[0])
	fmt.Printf("Part_start: %d\n", p.Part_start)
	fmt.Printf("Part_size: %d\n", p.Part_size)
	fmt.Printf("Part_name: %s\n", string(p.Part_name[:]))
	fmt.Printf("Part_correlative: %d\n", p.Part_correlative)
	fmt.Printf("Part_id: %s\n", string(p.Part_id[:]))
}