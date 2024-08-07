package structures


//structura del EBR

type EBR struct {
	Part_mount	[1]byte
	Part_fit	[1]byte
	Part_start	int32
	Part_size	int32
	Part_next	int32
	Part_name	[16]byte
}