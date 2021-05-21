package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"encoding/binary"
	"strings"
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
			index := file_length - int64(i) - int64(byte_size)
			return index
			break
		}
	}

	panic("signature not found")
}

// note: filesize limit of int range ~ 4GB
func read_bytes(f *os.File, amount int, location int64) []byte{
	f.Seek(location, 0)
	b := make([]byte, amount)
	n, err := io.ReadAtLeast(f, b, amount)
	check(err)
	return b[:n]
}

func Find(array []string, x string) int {
    for i, n := range array {
        if x == n {
            return i
        }
    }
    panic("item not found")
}

func get_item(f *os.File, index int64, name string) []byte{
		item_name_len_slice := read_bytes(f,2, index+26)
		item_name_len := int(binary.LittleEndian.Uint16(item_name_len_slice))

		item_name_slice := read_bytes(f,item_name_len, index+30)
		item_name := string(item_name_slice)

		if item_name == name{
			// confirmed
			fmt.Print("")
		} else {
			panic("item name doesnt match")
		}

		extra_len_slice := read_bytes(f,2, index+28)
		extra_len := int64(binary.LittleEndian.Uint16(extra_len_slice))
		data_index := index + 30 + extra_len

		data_len_slice := read_bytes(f,4, index+18)
		data_len := int(binary.LittleEndian.Uint32(data_len_slice))

		data_slice := read_bytes(f, data_len, data_index)
		return data_slice

}

func list_files(f *os.File, file_length int64, start_offset int64, num_items int, target []byte, interrupt []byte) ([]int64,[]string){

	var list_index []int64
	var list_name []string

	for i := 1; i < int(file_length); i++ {

		// scan every byte_block, starting from the end with start_offset
		test_location := file_length - start_offset - int64(4 + i )
		test_slice := read_bytes(f, 4, test_location)

		if string(test_slice) == string(target) {
			index := file_length - int64(i) - 4

			//Central_FH_slice := read_bytes(f, 4, index)

			//version_slice := read_bytes(f, 2, index+4)

			Local_FH_index_slice := read_bytes(f, 4, index+42)
			Local_FH_index := int64(binary.LittleEndian.Uint32(Local_FH_index_slice))

			check_bytes_slice := read_bytes(f, 4, Local_FH_index)
			if string(check_bytes_slice) == string([]byte{0x50,0x4b,0x03,0x04}){
				// confirmed
				fmt.Print("")
			}	else {
				panic("Local file header signature not found")
			}

			name_len_slice := read_bytes(f, 2, index+28)
			name_len := int(binary.LittleEndian.Uint16(name_len_slice))

			name := read_bytes(f, name_len, index+46)
			fmt.Print("Name: ")
			fmt.Println(string(name))

			list_index = append(list_index,Local_FH_index)
			list_name = append(list_name,string(name))

			if len(list_index) == num_items{
				fmt.Println("All Entries found")
				return list_index, list_name
				break
			}
		}
		if string(test_slice) == string(interrupt){
			fmt.Println("found interrupt")
			return list_index, list_name
			break
		}
	}

	panic("signature not found")
}



func main() {
	//zip_path := "library.zip"
	zip_path := "library_tiny.zip"
	fmt.Println("opening " + zip_path)
	//requested_file := "10.1002/0471264180.or083.01.pdf"
	requested_file := "10.4269/00000002.pdf"
	fmt.Println("searching for " + requested_file)

	// read file
	f, err := os.Open(zip_path)
	check(err)

	// get length
	file_info, err := os.Stat(zip_path)
	check(err)
	file_length := file_info.Size()

	// End of central directory signature =  0x06054b50
	EoCD := search_signature(f, file_length, 0, 4, []byte{0x50,0x4b,0x05,0x06})

	// Total number of central directory records
	num_CD_slice := read_bytes(f, 2, EoCD+10)
	num_CD := int(binary.LittleEndian.Uint16(num_CD_slice))
	fmt.Print("Number of Central Directory records: ")
	fmt.Println(num_CD)

	fmt.Println("Listing Files: ")
	// interrupts at the first local file header 0x04034b50
	list_index, list_name := list_files(f, file_length, 0, num_CD, []byte{0x50,0x4b,0x01,0x02}, []byte{0x50,0x4b,0x03,0x04})
	pos_item := Find(list_name, requested_file)
	fmt.Print("requested item at ")
	fmt.Println(list_index[pos_item])

	//item
	data := get_item(f,list_index[pos_item], requested_file)
	f.Close()
	// redirect datastream as virtual filesystem for IPFS at this point

	// save data for now
	if strings.Contains(requested_file, "/"){
		os.MkdirAll(strings.Split(requested_file, "/")[0], os.ModePerm)
	}
	err = ioutil.WriteFile(requested_file, data, 0644)
	check(err)
	fmt.Print("file extracted.")

}
