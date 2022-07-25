package main

/*
// windows via vcpkg
//#cgo CFLAGS: -I C:\vcpkg-master\installed\x64-windows\include\libxml2
//#cgo LDFLAGS: -L C:\vcpkg-master\installed\x64-windows\lib\libxml2

// linux via pkg manager
#cgo CFLAGS: -I/usr/include/libxml2
#cgo LDFLAGS: -lxml2

#include <stdlib.h>
#include <libxml/parser.h>
#include <libxml/tree.h>

const xmlChar* toXmlChar(const char* str){
	return (const xmlChar*)str;
}
const char* toChar(const xmlChar* str){
	return (const char*)str;
}
*/
import "C"
import (
	"fmt"
	"log"
	"unsafe"
)

/*
	Ubuntu libxml installing

	sudo apt install libxml2
	sudo apt install libxml2-dev

	C example source and tutorial for libxml2:
	http://web.mit.edu/outland/share/doc/libxml2-2.4.30/html/tutorial/xmltutorial.pdf
*/

type XMLError struct {
	info string
}

func (err *XMLError) Error() string {
	return err.info
}

func parse(name string) error {
	nameStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameStr))

	doc := C.xmlReadFile(nameStr, nil, 0)
	defer C.xmlFreeDoc(doc)

	if doc == nil {
		return &XMLError{"File can't be openned"}
	}

	// cur := C.xmlDocGetRootElement(doc)
	root := C.xmlDocGetRootElement(doc)
	if root == nil {
		return &XMLError{"Empty document"}
	}
	storyStr := C.CString("story")
	defer C.free(unsafe.Pointer(storyStr))

	if C.xmlStrcmp(root.name, C.toXmlChar(storyStr)) != 0 {
		return &XMLError{"Wrong type"}
	}

	infoStr := C.CString("storyinfo")
	defer C.free(unsafe.Pointer(infoStr))

	for cur := root.children; cur != nil; cur = cur.next {
		if C.xmlStrcmp(cur.name, C.toXmlChar(infoStr)) == 0 {
			parseStory(doc, cur)
		}
	}

	return nil
}

func parseStory(doc C.xmlDocPtr, cur C.xmlNodePtr) {
	keywordStr := C.CString("keyword")
	defer C.free(unsafe.Pointer(keywordStr))

	for cur = cur.children; cur != nil; cur = cur.next {
		if C.xmlStrcmp(cur.name, C.toXmlChar(keywordStr)) == 0 {
			str := C.GoString(C.toChar(C.xmlNodeListGetString(doc, cur.children, 1)))
			fmt.Printf("keyword: %s\n", str)
		}
	}
}

func main() {

	err := parse("data.xml")
	if err != nil {
		log.Fatalf(err.Error())
	}

}
