package structures

import (
	"encoding/binary"
	"fmt"
	util "server/util"
	"strings"
)

// FDISK estructura que representa el comando fdisk con sus parámetros
type FDISK struct {
	Size int    // Tamaño de la partición
	Unit string // Unidad de medida del tamaño (B, K o M); por defecto K
	Path string // Ruta del archivo del disco
	TypE  string // Tipo de partición (P, E, L)
	Fit  string // Tipo de ajuste (BF, FF, WF); por defecto WF
	Name string // Nombre de la partición
}

func CommandFdisk(fdisk *FDISK) (string, error) {
	
	// Convertir el tamaño a bytes
	sizeBytes, err := util.ConvertToBytes(fdisk.Size, fdisk.Unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return "Error converting size en Fdisk", err
	}

	var msg string

	// para crear la particion primaria
	if(fdisk.TypE == "P"){
		msg, err = CreatePrimaryPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creating primary partition:", err)
			return msg, err
		}
	}else if(fdisk.TypE == "E"){
		fmt.Println("Particion extendida")
	}else if(fdisk.TypE == "L"){
		fmt.Println("Particion logica")
	}
	return "",nil
}


func CreatePrimaryPartition(fdisk *FDISK, sizeBytes int)(string, error){
	
	var mbr MBR
	
	msg, err := mbr.DeserializeMBR(fdisk.Path)
	if err != nil {
		return msg, fmt.Errorf("error leyendo el MBR del disco: %s", err)
	}

	// se obtiene la primera particion libre
	particionDisponible, inicioParticion, indexParticion, msg:= mbr.GetFirstPartitionAvailable()
	if particionDisponible == nil {
		return msg, fmt.Errorf("no hay particiones disponibles")
	}

	// solo para pruebas

	//print para verificar que la particion está disponible
	// fmt.Println("\nParticion disponible:")
	// particionDisponible.Print()

	// crear la particion con los parámetros proporcionados 
	particionDisponible.CreatePartition(inicioParticion, sizeBytes, fdisk.TypE, fdisk.Fit, fdisk.Name)

	// verificar que la particion se haya creado correctamente
	// fmt.Println("\nParticion creada(modificada):")
	// particionDisponible.Print()

	// montar la particion
	mbr.Mbr_partitions[indexParticion] = *particionDisponible //asignar la particion al MBR

	// imprimir las particiones del MBR
	// fmt.Println("\nParticiones del MBR:")
	// mbr.PrintPartitions()

	// Serialiazar el MBR modificado
	msg, err = mbr.SerializeMBR(fdisk.Path)
	if err != nil {
		return msg, fmt.Errorf("error escribiendo el MBR al disco: %s", err)
	}	
	return "",nil
}


//funcion para crear la particion
func CreatePartition(fdisk *FDISK, sizeBytes int) (string, error){

	// 1. leer el MBR del disco
	mbr := &MBR{}
	msg, err := mbr.DeserializeMBR(fdisk.Path)
	if err != nil {
		return msg, fmt.Errorf("error leyendo el MBR del disco: %s", err)	
	}

	//verficar que el nombre de la particion no exista
	partitionName := fdisk.Name //nombre de la particion
    for _, partition := range mbr.Mbr_partitions { //recorrer las particiones
        if partition.Part_status[0] != 'N' { // Asegurarse de que la partición esté en uso
            partName := strings.Trim(string(partition.Part_name[:]), "\x00") //obtener el nombre de la particion, sin caracteres nulos
            if partName == partitionName { //si el nombre de la particion ya existe
                return "ERROR: ya existe una particion con el nombre",fmt.Errorf("ya existe una partición con el nombre: %s", fdisk.Name)
            }
        }
    }

	//2. Encontrar un espacio disponible para la nueva particion
	startPosition := int32(binary.Size(mbr)) //el inicio despues del mbr
	for i:= 0; i<len(mbr.Mbr_partitions); i++{ //recorrer las particiones
		partition := mbr.Mbr_partitions[i] //obtener la particion
		if partition.Part_status[0]=='N'{ // si la particion esta vacia
			if partition.Part_start	== -1 || partition.Part_size == -1 { //si la particion esta vacia
				break
			}
		}
		startPosition = partition.Part_start + partition.Part_size //actualizar la posicion de inicio
	}

	//verificar que haya espacion suficiente
	if startPosition + int32(sizeBytes) > mbr.Mbr_size{ //si no hay espacio suficiente
		return "no hay espacio suficiente para la particion",fmt.Errorf("no hay espacio suficiente para la particion") 
	}
	
	//3. Crear la particion
	newPartition := PARTITION{
		Part_status: [1]byte{'1'}, 			//1 = activa, 0 = inactiva
		Part_type: [1]byte{fdisk.TypE[0]}, 	//P, E, L
		Part_fit: [1]byte{fdisk.Fit[0]}, 	//BF, FF, WF
		Part_start: startPosition, 			//inicio de la particion
		Part_size: int32(sizeBytes), 		//tamaño de la particion
		Part_name: util.ConvertToFixedSizeArray(fdisk.Name, 16), //nombre de la particion
	}

	// asignar la nueva particion al MBR
	for i:= range mbr.Mbr_partitions{ //recorrer las particiones
		if mbr.Mbr_partitions[i].Part_status[0] == 'N'{ //si la particion esta vacia
			mbr.Mbr_partitions[i] = newPartition //asignar la nueva particion
			break
		}
	}

	// 4. escribir el mbr actualizado de nuevo al disco
	var mensaje string
	mensaje, err = mbr.SerializeMBR(fdisk.Path)
	if err != nil {
		return mensaje ,fmt.Errorf("error escribiendo el MBR al disco: %s", err)
	}

	fmt.Println("particion creada con éxito")
	
	return "",nil

}