package main

import (
	"fmt"
	"io"
	"os"
	"encoding/binary"
)

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)
	}
}


func search_signature(f *os.File, file_length int64, start_offset int64, byte_size int, target []byte) int64 {

	for i := 1; i < int(file_length); i++ {

		// scan every byte_block, starting from the end with start_offset
		f.Seek(file_length - start_offset - int64(byte_size + i ), 0)

		b := make([]byte, byte_size)

		n, err := io.ReadAtLeast(f, b, byte_size)
		check(err)

		if string(b[:n]) == string(target) {
			index := file_length - int64(i)
			return index
			break
		}
	}

	panic("signature not found")
}

func read_bytes(f *os.File, amount int, location int64) []byte{
	f.Seek(location, 0)
	b := make([]byte, amount)
	n, err := io.ReadAtLeast(f, b, amount)
	check(err)
	return b[:n]
}

func list_files(f *os.File, file_length int64, start_offset int64, target []byte, interrupt []byte) []int64{

	var list []int64

	for i := 1; i < int(file_length); i++ {

		// scan every byte_block, starting from the end with start_offset
		test_location := file_length - start_offset - int64(4 + i )
		test_slice := read_bytes(f, 4, test_location)

		if string(test_slice) == string(target) {
			index := file_length - int64(i) - 4
			Central_FH_slice := read_bytes(f, 4, index)
			fmt.Print("found Central directory file header")
			fmt.Println(Central_FH_slice)

			version_slice := read_bytes(f, 2, index+4)
			fmt.Print("Version ")
			fmt.Println(version_slice)

			name_len_slice := read_bytes(f, 2, index+28)
			name_len := int(binary.LittleEndian.Uint16(name_len_slice))
			fmt.Print("Name Length: ")
			fmt.Print(name_len_slice)
			fmt.Print(" -> ")
			fmt.Println(name_len)

			name := read_bytes(f, name_len, index+46)
			fmt.Print("Name: ")
			fmt.Println(string(name))

			list = append(list,index)
		}
		if string(test_slice) == string(interrupt){
			fmt.Println("found interrupt")
			return list
			break
		}
	}

	panic("signature not found")
}

func main() {
	zip_path := "library.zip"
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
	EoCD := search_signature(f, file_length, 0, 4, []byte{0x50,0x4b,0x05,0x06})
	fmt.Print(EoCD)
	fmt.Println(" - End of central directory")



	// Local file header signature = 0x04034b50
	//Local_FH := search_signature(f, file_length, 0, 4, []byte{0x50,0x4b,0x03,0x04})
	//fmt.Print(Local_FH)
	//fmt.Println(" - Local file header ")


	// interrupts at the first local file header 0x04034b50
	fmt.Print(list_files(f, file_length, 0, []byte{0x50,0x4b,0x01,0x02}, []byte{0x50,0x4b,0x03,0x04}))

	f.Close()
}
