package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "encoding/json"
    "io"
    "io/ioutil"
    "os"
    "encoding/csv"
    "bufio"
)

type RequestMessage struct {
    Carmark string
    Carmodel string
    Numberdays string
    Numberunits string
}


type ResponseMessage struct {
    Carmark string
    Carmodel string
    Numberdays string
    Numberunits string
}

func main() {

router := mux.NewRouter().StrictSlash(true)
router.HandleFunc("/", Index)
router.HandleFunc("/endpoint/{param}", endpointFunc)
router.HandleFunc("/endpoint2/{param}", endpointFunc2JSONInput)

log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Println(w, "Service OK")
}

func endpointFunc(w http.ResponseWriter, r *http.Request) {
    //vars := mux.Vars(r)
    //param := vars["param"]
    
    file, err := os.Open("rentals.csv")
    if err!=nil {
       json.NewEncoder(w).Encode(err)
       return
    }
    reader := csv.NewReader(bufio.NewReader(file))
    for {
        record, err := reader.Read()
        if err == io.EOF {
                break
        }
    fmt.Fprintf(w, "La marca del cote es %q \nEl model del cotxe es %q \nEl nombre de dies es %q \nEl nombre d'unitats es %q \n" , record[0], record[1], record[2], record[3])
    //res := ResponseMessage{Carmark: "Text1", Field2: param}
    //json.NewEncoder(w).Encode(res)
    }

}

func endpointFunc2JSONInput(w http.ResponseWriter, r *http.Request) {
    var requestMessage RequestMessage
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &requestMessage); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    } else {
        fmt.Println("Marca cotxe: " + requestMessage.Carmark)
        fmt.Println("Model: " + requestMessage.Carmodel)
        fmt.Println("Nombre de dies: " + requestMessage.Numberdays)
        fmt.Println("Nombre d'unitats " + requestMessage.Numberunits)
	file, err := os.OpenFile("rentals.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
        if err!=nil {
	   json.NewEncoder(w).Encode(err)
	   return
        }
        writer := csv.NewWriter(file)
        var data1 = []string{requestMessage.Carmark, requestMessage.Carmodel, requestMessage.Numberdays, requestMessage.Numberunits }
        writer.Write(data1)
        writer.Flush()
        file.Close()
     }
}



	
