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

		msg, err = CreateExtendedPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creating extended partition:", err)
			return msg, err
		}


	}else if(fdisk.TypE == "L"){

		msg, err = CreateLogicalPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creating logical partition:", err)
			return msg, err
		}

	}
	return "",nil
}


func CreatePrimaryPartition(fdisk *FDISK, sizeBytes int)(string, error){
	
	var mbr MBR
	
	msg, err := mbr.DeserializeMBR(fdisk.Path)
	if err != nil {
		return msg, fmt.Errorf("error leyendo el MBR del disco: %s", err)
	}

	// Contar el número de particiones primarias y extendidas
	primaryCount := 0
	extendedExists := false
	for _, partition := range mbr.Mbr_partitions {
		if partition.Part_status[0] != '2' {
			if partition.Part_type[0] == 'P' {
				primaryCount++
			} else if partition.Part_type[0] == 'E' {
				extendedExists = true
			}
		}
	}

	// Verificar que no se exceda el límite de 4 particiones
	if primaryCount >= 3 && !extendedExists {
		return "ERROR: No se pueden crear más de 3 particiones primarias si no hay una extendida", fmt.Errorf("límite de particiones primarias alcanzado")
	}

	if primaryCount >= 4 {
		return "ERROR: No se pueden crear más particiones primarias", fmt.Errorf("límite de particiones primarias alcanzado")
	}

	// se obtiene la primera particion libre
	particionDisponible, inicioParticion, indexParticion, msg:= mbr.GetFirstPartitionAvailable()
	if particionDisponible == nil {
		return msg, fmt.Errorf("no hay particiones disponibles")
	}

	// crear la particion con los parámetros proporcionados 
	particionDisponible.CreatePartition(inicioParticion, sizeBytes, fdisk.TypE, fdisk.Fit, fdisk.Name)

	// montar la particion
	mbr.Mbr_partitions[indexParticion] = *particionDisponible //asignar la particion al MBR

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


func CreateExtendedPartition(fdisk *FDISK, sizeBytes int)(string, error){
	
	var mbr MBR

	msg, err := mbr.DeserializeMBR(fdisk.Path)
	if err != nil {
		return msg, fmt.Errorf("error leyendo el MBR del disco: %s", err)
	}

	// verificar si ya existe una particion extendida
	for _, partition := range mbr.Mbr_partitions {
		if partition.Part_type[0] != '0' && partition.Part_type[0] == 'E' {
			return "ERROR: ya existe una particion extendida",fmt.Errorf("ya existe una partición extendida")
		}
	}

	// se obtiene la primera particion libre
	particionDisponible, inicioParticion, indexParticion, msg:= mbr.GetFirstPartitionAvailable()
	if particionDisponible == nil {
		return msg, fmt.Errorf("no hay particiones disponibles")
	}

	// crear la particion con los parámetros proporcionados 
	particionDisponible.CreatePartition(inicioParticion, sizeBytes, fdisk.TypE, fdisk.Fit, fdisk.Name)

	// montar la particion
	mbr.Mbr_partitions[indexParticion] = *particionDisponible //asignar la particion al MBR

	// Serialiazar el MBR modificado
	msg, err = mbr.SerializeMBR(fdisk.Path)
	if err != nil {
		return msg, fmt.Errorf("error escribiendo el MBR al disco: %s", err)
	}	
	return "",nil

}


func CreateLogicalPartition(fdisk *FDISK, sizeBytes int)(string, error){

	var mbr MBR

	msg, err := mbr.DeserializeMBR(fdisk.Path)
	if err != nil {
		return msg, fmt.Errorf("error leyendo el MBR del disco: %s", err)
	}

	// verificar si ya existe una particion extendida
	var extendedPartition *PARTITION
	for i := range mbr.Mbr_partitions{
		if mbr.Mbr_partitions[i].Part_type[0] == 'E'{
			extendedPartition = &mbr.Mbr_partitions[i]
			break
		}
	}

	// verificar si no existe una particion extendida
	if extendedPartition == nil{
		return "ERROR: no existe una particion extendida para crear la lógica",fmt.Errorf("no existe una partición extendida para crear la lógica")
	}

	//buscar el primer EBR dentro de la particion extendida
	var ebr EBR
	position := extendedPartition.Part_start //posicion del ebr

	//leer el primer EBR
	_, err = ebr.DeserializeEBR(fdisk.Path, position)

	if err != nil {
		// si no existe, crear el primer EBR
		ebr = EBR{
			Part_fit: [1]byte{fdisk.Fit[0]},
			Part_start: position,
			Part_size: int32(sizeBytes),
			Part_next: -1,
			Part_name: util.ConvertToFixedSizeArray(fdisk.Name, 16),
		}
		return ebr.SerializeEBR(fdisk.Path, position)
	}

	// recorrer la lista de EBRs hasta encontrar un espacio
	for ebr.Part_next != -1 {
		position = ebr.Part_next
		msg, err = ebr.DeserializeEBR(fdisk.Path, position)
		if err != nil {
			return msg, fmt.Errorf("error leyendo el EBR: %s", err)
		}
	}

	// crear el nuevo EBR despues del último encontrado
	newEBR := EBR{
		Part_fit: [1]byte{fdisk.Fit[0]},
		Part_start: position+ebr.Part_size, // La nueva partición empieza después de la última
		Part_size: int32(sizeBytes),
		Part_next: -1,
		Part_name: util.ConvertToFixedSizeArray(fdisk.Name, 16),
	}

	// verificar que haya espacio suficiente en la particion extendida
	if newEBR.Part_start + newEBR.Part_size > extendedPartition.Part_start + extendedPartition.Part_size {
		return "ERROR: No hay espacio suficiente para crear la partición lógica", fmt.Errorf("no hay espacio suficiente para crear la partición lógica")
	}

	// escribir el nuevo EBR en el disco
	msg, err = newEBR.SerializeEBR(fdisk.Path, newEBR.Part_start)
	if err != nil {
		return msg, fmt.Errorf("error escribiendo el EBR al disco: %s", err)
	}

	// actualizar el apuntador del EBR anterior
	ebr.Part_next = newEBR.Part_start
	msg, err = ebr.SerializeEBR(fdisk.Path, ebr.Part_start)
	if err != nil {
		return msg, fmt.Errorf("error actualizando el EBR al disco: %s", err)
	}
		
	return "",nil
}