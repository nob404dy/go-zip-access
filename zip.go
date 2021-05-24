package zip

import (
	"archive/zip"
	//"fmt"
  //"os"
	"io/ioutil"
  //"strings"
)

type myCloser interface {
	Close() error
}

// closeFile is a helper function which streamlines closing
// with error checking on different file types.
func closeFile(f myCloser) {
	err := f.Close()
	check(err)
}

// Read_item is a wrapper function for ioutil.ReadAll. It accepts a zip.File as
// its parameter, opens it, reads its content and returns it as a byte slice.
func Read_item(file *zip.File) []byte {
	fc, err := file.Open()
	check(err)
	defer closeFile(fc)

	content, err := ioutil.ReadAll(fc)
	check(err)

	return content
}

// check is a helper function which streamlines error checking
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Get_Directory(zipFile string) ([]string){
	//zipFile := "../libgen.scimag29999000-29999999.zip"
  //zipFile := "../library_tiny.zip"

	zf, err := zip.OpenReader(zipFile)
	check(err)
	defer closeFile(zf)

  var directory []string
	for _, file := range zf.File {
    directory = append(directory,string(file.Name))
		//fmt.Printf("%s\n", file.Name)

    // FILE CONTENT
    //if strings.Contains(file.Name, "/"){
  	//	os.MkdirAll(strings.Split(file.Name, "/")[0], os.ModePerm)
  	//}
    //err := ioutil.WriteFile(file.Name, Read_item(file), 0644)
    //fmt.Println(err)
		//fmt.Printf("%s\n\n", Read_item(file)) // file content
	}

  return directory
}
