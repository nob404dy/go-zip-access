package main

import (
	"fmt"
	"io"
	"os"

)

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func search_signature(f *os.File, file_length int64, byte_size int, x []byte) int64 {

	for i := 1; i < int(file_length); i++ {

		// scan every byte_block, starting from the end
		f.Seek(file_length - int64(byte_size + i), 0)

		b := make([]byte, byte_size)

		n2, err := io.ReadAtLeast(f, b, byte_size)
		check(err)

		if string(b[:n2]) == string(x) {
			index := file_length - int64(i)
			return index
			break
		}
	}

	panic("signature" + string(x) + " not found")
}

func main() {
	zip_path := "library_tiny.zip"
	fmt.Println("opening " + zip_path)

	// read file
	f, err := os.Open(zip_path)
	check(err)

	// get length
	file_info, err := os.Stat(zip_path)
	check(err)
	file_length := file_info.Size()
	fmt.Print("length ")
	fmt.Print(int64(file_length))
	fmt.Println(" Bytes")

	// End of central directory signature =  0x06054b50
	EoCD := search_signature(f, file_length, 4, []byte{0x50,0x4b,0x05,0x06})
	fmt.Print(EoCD)
	fmt.Println(" - End of central directory")

	// Central directory file header signature = 0x02014b50
	CDFH := search_signature(f, file_length, 4, []byte{0x50,0x4b,0x01,0x02})
	fmt.Print(CDFH)
	fmt.Println(" - Central directory file header ")


	f.Close()
}
