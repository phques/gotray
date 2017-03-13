package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/sys/windows/registry"
)

func rot13(p []byte) {
	for i := 0; i < len(p); i++ {
		if (p[i] >= 'A' && p[i] < 'N') || (p[i] >= 'a' && p[i] < 'n') {
			p[i] += 13
		} else if (p[i] > 'M' && p[i] <= 'Z') || (p[i] > 'm' && p[i] <= 'z') {
			p[i] -= 13
		}
	}
}

func usage() {
	fmt.Println("usage")
	fmt.Println("without parameters, lists entries, else :")
	fmt.Println(" -entry (0..)")
	fmt.Println(" -value (0,1,2)")
	fmt.Println("   0 = only show notifications")
	fmt.Println("   1 = hide icon and notifications")
	fmt.Println("   2 = show icon and notifications")
	log.Fatal("invalid parameters")
}

func makePath(pathSlice []byte) []byte {
	// make a copy (src is 2 bytes / char), leaving original alone
	var path []byte
	for i := 0; i < len(pathSlice); i++ {
		b := pathSlice[i*2]
		if b == 0 {
			break
		}
		path = append(path, b)
	}

	// decode path
	rot13(path)
	return path
}

func main() {

	// check for cmdline params
	//	if len(os.Args) != 1 && len(os.Args) != 3 {
	//		log.Fatal("parameters: [entryIndex, visibilityValue (0-2)]")
	//	}
	var entryIndex = flag.Int("entry", -1, "index of entry to modify")
	var newVisibilityValue = flag.Int("value", -1, "new visibility value (0-2)")

	flag.Parse()

	hasParams := false
	if len(os.Args) != 1 {
		if *entryIndex < 0 || *newVisibilityValue < 0 || *newVisibilityValue > 2 {
			usage()
		}
		hasParams = true
	}

	// -------------
	keypath := `Software\Classes\Local Settings\Software\Microsoft\Windows\CurrentVersion\TrayNotify`
	valuename := "IconStreams"

	// open the registry key
	//access := registry.QUERY_VALUE | registry.SET_VALUE
	k, err := registry.OpenKey(registry.CURRENT_USER, keypath, registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	// read the value (binary)
	value, _, err := k.GetBinaryValue(valuename)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("value is %v\n", v)

	// skip header (see readme for the data format)
	data := value[20:]

	// go through each entry
	for i := 0; i < len(data)/1640; i++ {

		// get the EXE name
		pathSlice := data[:527]
		path := makePath(pathSlice)

		// get the visibility flag (value is in 1st byte)
		visibility := data[528 : 528+4]

		if !hasParams || i == *entryIndex {
			fmt.Printf("%02d %v %s\n", i, visibility, string(path))
		}

		if hasParams && i == *entryIndex {
			fmt.Println(" changing tray icon visibility => ", *newVisibilityValue)
			visibility[0] = byte(*newVisibilityValue)

			err := k.SetBinaryValue(valuename, value)
			if err != nil {
				log.Fatal(err.Error())
			}

			break
		}

		data = data[1640:]
	}
}
