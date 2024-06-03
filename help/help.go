package help

import "fmt"

func PrintHelp() {
	fmt.Println("----------------------------------------------------")
	fmt.Println("		   LEA encryption					 		 ")
	fmt.Println("----------------------------------------------------")
	fmt.Println("")
	fmt.Println("*file* 		source file path")
	fmt.Println("-E 		encrypt source file")
	fmt.Println("-D 		decrypt source file")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Usage")
	fmt.Println("./lea test.txt -E")
	fmt.Println("./lea test.txt -D")
	fmt.Println("")
}
