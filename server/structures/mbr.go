package structures

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// MBR estructura del MBR
type MBR struct {
	Mbr_tamano   		int32
	Mbr_fecha_creacion 	[19]byte
	Mbr_disk_signature 	int32
	Dsk_fit 			[1]byte
	Mbr_particion 		[4]Partition
}



// Partition estructura de una partición
type Partition struct {
	Part_status 	[1]byte
	Part_type 		[1]byte
	Part_fit 		[1]byte
	Part_start 		int32
	Part_size 		int32
	Part_name 		[16]byte
	//part_correlative int32
	Part_id		    [4]byte
}

func ObtenerFechaHora () string{

	fechaHora := time.Now()

	fecha_hora := fmt.Sprintf("%02d-%02d-%d %02d:%02d", fechaHora.Day(), fechaHora.Month(), fechaHora.Year(), fechaHora.Hour(), fechaHora.Minute())
	
	return fecha_hora
}


// CrearMBR crea un MBR con los valores inicializados
func CrearMBR(tamano int32, fit string) MBR {
	fechaHora := ObtenerFechaHora()
	var fechaHoraArray [19]byte
	copy(fechaHoraArray[:], fechaHora)

	var fitArray [1]byte
	copy(fitArray[:], fit)

	mbr := MBR{
		Mbr_tamano: tamano,
		Mbr_fecha_creacion: fechaHoraArray,
		Mbr_disk_signature: rand.Int31(),
		Dsk_fit: fitArray,

	}

	return mbr
}

// ReadMBRFromFile lee el MBR desde el archivo binario y lo retorna
func ReadMBRFromFile(path string) (MBR, error) {
	file, err := os.Open(path)
	if err != nil {
		return MBR{}, err // Si no se pudo abrir el archivo, retornar el error
	}
	defer file.Close() // Cerrar el archivo

	var mbr MBR
	err = binary.Read(file, binary.LittleEndian, &mbr)
	if err != nil {
		return MBR{}, err // Si no se pudo leer el MBR, retornar el error
	}

	return mbr, nil
}

// PrintMBR imprime la información del MBR
func PrintMBR(mbr MBR) {
	fmt.Println("MBR Info:")
	fmt.Printf("Tamaño del disco: %d bytes\n", mbr.Mbr_tamano)
	fmt.Printf("Fecha de creación: %s\n", string(mbr.Mbr_fecha_creacion[:]))
	fmt.Printf("Disk Signature: %d\n", mbr.Mbr_disk_signature)
	fmt.Printf("Fit: %s\n", string(mbr.Dsk_fit[:]))
	for i, part := range mbr.Mbr_particion {
		fmt.Printf("Partición %d:\n", i+1)
		fmt.Printf("  Status: %s\n", string(part.Part_status[:]))
		fmt.Printf("  Type: %s\n", string(part.Part_type[:]))
		fmt.Printf("  Fit: %s\n", string(part.Part_fit[:]))
		fmt.Printf("  Start: %d\n", part.Part_start)
		fmt.Printf("  Size: %d\n", part.Part_size)
		fmt.Printf("  Name: %s\n", string(part.Part_name[:]))
	}
}
