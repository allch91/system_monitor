package handlers

import (
  "html/template"
  "strconv"
  "io/ioutil"
  "strings"
  "fmt"
  "net/http"
  "encoding/json"
  "time"
  "math/rand"
  "github.com/thedevsaddam/renderer"
  "github.com/shirou/gopsutil/host"
  "github.com/shirou/gopsutil/load"
  linuxproc "github.com/c9s/goprocinfo/linux"
  ps "github.com/mitchellh/go-ps"
  helpers "../helpers"
  repos "../repos"
)

var rnd *renderer.Render

func init() {
	opts := renderer.Options{
		ParseGlobPattern: "./static/*.html",
	}

	rnd = renderer.New(opts)
}


func CpuProcessHandler(response http.ResponseWriter, request *http.Request){
  rnd.HTML(response, http.StatusOK, "cpu_graph", nil)
}

func RamProcessHandler(response http.ResponseWriter, request *http.Request){
  rnd.HTML(response, http.StatusOK, "ram_graph", nil)
}


func AdminPageHandler(response http.ResponseWriter, request *http.Request){
  rnd.HTML(response, http.StatusOK, "admin", nil)
  //var body, _= helpers.LoadFile("static/chartjs_example3.html")
  //fmt.Fprintf(response, body)
}

//LoginPage GET
func LoginPageHandler(response http.ResponseWriter, request *http.Request){
  var body, _= helpers.LoadFile("static/index.html")
  fmt.Fprintf(response, body)
}

type Values struct{
  X string`json:"x"`
  Y int `json:"y"`
}

//======= HOME DATA

type Process struct{
  Pid int
  Usuario string
  Estado string
  Pram int
  Name string
}

type HomeData struct{
  Procesos uint64
  Ejecucion int
  Suspendido int
  Detenido int
  Zombie int
  Processes []Process
}


func HomePageHandler(response http.ResponseWriter, request *http.Request){
  processList, _ := ps.Processes()
  tmpl, _ := template.ParseFiles("static/admin_home.html")
  Proce := []Process{}
  for x:= range processList{
    var process ps.Process
    process = processList[x]
    Proce = append(Proce,Process{Pid: process.Pid(),Usuario:"dennis-pc",Estado:"activo",Pram:rand.Intn(100), Name: process.Executable()})
  }
  infoStat, _ := host.Info()
  miscStat, _ := load.Misc()
  data:=HomeData{
    Procesos:infoStat.Procs,
    Ejecucion:miscStat.ProcsRunning,
    Suspendido:rand.Intn(int(infoStat.Procs)),
    Detenido:rand.Intn(int(infoStat.Procs)),
    Zombie:rand.Intn(int(infoStat.Procs)),
    Processes: Proce}
  tmpl.Execute(response,data)
  //var body, _= helpers.LoadFile("static/index.html")
  //fmt.Fprintf(response, body)
}

//======= RAM DATA

type RamStruct struct{
  X int `json:"x"`
  Y int `json:"y"`
  Z int `json:"z"`
}

func RamData(response http.ResponseWriter, request *http.Request){
  memory, _ := linuxproc.ReadMemInfo("/proc/meminfo")
  total:=memory.MemTotal/1000
  consumida:= (memory.MemTotal - memory.MemFree)/1000
  per_consumo := (float64(memory.MemTotal - memory.MemFree)/float64(memory.MemTotal))*100
  ram_data := RamStruct{X:int(total),Y:int(consumida),Z:int(per_consumo)}
  byteArray, _ := json.Marshal(ram_data)
  fmt.Fprintf(response,string(byteArray))
}

//======= CPU DATA

func getCPUSample() (idle, total uint64) {
    contents, err := ioutil.ReadFile("/proc/stat")
    if err != nil {
        return
    }
    lines := strings.Split(string(contents), "\n")
    for _, line := range(lines) {
        fields := strings.Fields(line)
        if fields[0] == "cpu" {
            numFields := len(fields)
            for i := 1; i < numFields; i++ {
                val, err := strconv.ParseUint(fields[i], 10, 64)
                if err != nil {
                    fmt.Println("Error: ", i, fields[i], err)
                }
                total += val // tally up all the numbers to get total ticks
                if i == 4 {  // idle is the 5th field in the cpu line
                    idle = val
                }
            }
            return
        }
    }
    return
}

type CpuStruct struct{
  X int `json:"x"`
  Y int `json:"y"`
}

func CpuData(response http.ResponseWriter, request *http.Request){
  idle0, total0 := getCPUSample()
  time.Sleep(3 * time.Second)
  idle1, total1 := getCPUSample()

  idleTicks := float64(idle1 - idle0)
  totalTicks := float64(total1 - total0)
  cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

  cpu_data := RamStruct{X:int(totalTicks-idleTicks),Y:int(cpuUsage)}
  byteArray, _ := json.Marshal(cpu_data)
  fmt.Fprintf(response,string(byteArray))
}



func AdminHandler(response http.ResponseWriter, request *http.Request){

  values := Values{X:time.Now().Format("2006-01-02T15:04:05Z"),Y:rand.Intn(100)}

  byteArray, _ := json.Marshal(values)

  fmt.Fprintf(response,string(byteArray))
}

//LoginPage POST
func LoginHandler(response http.ResponseWriter, request *http.Request){
  username:= request.FormValue("username")
  pass:= request.FormValue("pass")
  redirectTarget:= "/adminPage"

  if !helpers.IsEmpty(username) && !helpers.IsEmpty(pass){
    _userIsValid := repos.UserIsValid(username, pass)
    if _userIsValid {
        redirectTarget = "/adminPage"
    } else {
        redirectTarget = "/"
    }
  }
  http.Redirect(response, request, redirectTarget, 302)
}
