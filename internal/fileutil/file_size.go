package fileutil

import "fmt"

var (
	KiloByte int64 = 1024
	MegaByte int64 = 1024 * KiloByte
	GigaByte int64 = 1024 * MegaByte
)

func ByteToAppropriateUnit(byte int64) string {
	if byte >= GigaByte {
		return fmt.Sprintf("%.2f GB", float64(byte)/float64(GigaByte))
	} else if byte >= MegaByte {
		return fmt.Sprintf("%.2f MB", float64(byte)/float64(MegaByte))
	} else if byte >= KiloByte {
		return fmt.Sprintf("%.2f KB", float64(byte)/float64(KiloByte))
	}
	return fmt.Sprintf("%d B", byte)
}
