package main

import ("net/http"
	"github.com/gorilla/mux"
        "fmt"
        iconv "github.com/djimenez/iconv-go"
        "strings"
        "io/ioutil"
	"runtime"
	"log"
	"os"
        )

func main() {
	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)
	log.Println("Number of CPUs: ", nCPU)

	rtr := mux.NewRouter()
	rtr.HandleFunc("/search/{year}",searchbyYear).Methods("GET").Queries("id","{id:[0-9]+}","fn","{fn}","ln","{ln}")
	http.Handle("/", rtr)
	/*******************************
	for Openshift
	*******************************/

	bind := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

	/*******************************
		End for Openshift
	*******************************/

//	bind := ":8080"

	log.Println("Listening:" + bind + "...")
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}

}

func searchbyYear(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	Year := params["year"]
	ID := params["id"]
	FN := params["fn"]
	LN := params["ln"]
	//fmt.Printf("%s %s %s %s",Year,ID,FN,LN)
	Surl := "http://refundedcheque.rd.go.th/itp_x_tw/SearchTaxpayerServlet"
	query := fmt.Sprintf("nid=%s&fName=%s&lName=%s&taxYear=%s&searchType=null&effDate=null",ID,FN,LN,Year)
	log.Println(query)
        req, err := http.NewRequest("POST", Surl, strings.NewReader(query))
    req.Header.Set("Referer", "http://refundedcheque.rd.go.th/itp_x_tw/pages/ITPtaxresult.jsp")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("X-Requested-With","XMLHttpRequest")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    output,_ := iconv.ConvertString(string (body), "tis-620", "utf-8")
    fmt.Println("response Body:", string(output))

	if len(output) == 0 {
		log.Println("Data not found")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
	}
	w.Write([]byte (output))

}


/*
func main() {
	http.HandleFunc("/", hello)
	bind := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	fmt.Printf("listening on %s...", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "hello, world from %s", runtime.Version())
}
*/
