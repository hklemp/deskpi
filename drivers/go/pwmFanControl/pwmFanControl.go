package main

import "fmt"
import "github.com/jacobsa/go-serial/serial"
import "log"
import "time"
import "bufio"
import "os"
//import "os/signal"
import "bytes"
import "io"
import "strconv"
import "strings"
func main() {
  // setup signal catching
  // sigs := make(chan os.Signal, 1)

  // // catch all signals since not explicitly listing
  // //signal.Notify(sigs)
  // signal.Notify(sigs,syscall.SIGQUIT)

  // // method invoked upon seeing signal
  // go func() {
  //   s := <-sigs
  //   log.Printf("RECEIVED SIGNAL: %s",s)
  //   //AppCleanup()
  //   os.Exit(1)
  // }()


    // Set up options.
  options := serial.OpenOptions{
    PortName: "/dev/ttyUSB0",
    BaudRate: 9600,
    DataBits: 8,
    StopBits: 1,
    MinimumReadSize: 4,
  }
  
  // Open the port.
  port, err := serial.Open(options)
  if err != nil {
    log.Fatalf("serial.Open: %v", err)
  }
  
  // Make sure to close it later.
  defer port.Close()

  var conf_info [8]int
  conf_info[0]=40
	conf_info[1]=25

	conf_info[2]=50
	conf_info[3]=50

	conf_info[4]=65
	conf_info[5]=75

	conf_info[6]=75
	conf_info[7]=100
  oldPwm := "pwm_000"
  for { 

    conf, err := readConfig("/etc/deskpi.conf")
    
    if err != nil {
      log.Fatalf("Error in config %v", err)
    } else {
      for i := 0; i < 8; i++ {
        conf_info[i],_ = strconv.Atoi(conf[i])
      }
    }

    cpu_temp, err := read_cpu_tmp()
    pwm := "pwm_000"

    if(cpu_temp < conf_info[0]) {
			pwm = "pwm_000"
		} else if(cpu_temp >= conf_info[0] && cpu_temp < conf_info[2]) {
      pwm = fmt.Sprintf("%s%03d", "pwm_" , conf_info[1]) 
		}	else if(cpu_temp >= conf_info[2] && cpu_temp < conf_info[4]) {
			pwm = fmt.Sprintf("%s%03d", "pwm_" , conf_info[3]) 
		}	else if(cpu_temp >= conf_info[4] && cpu_temp < conf_info[6]) {
			pwm = fmt.Sprintf("%s%03d", "pwm_" , conf_info[5]) 
		}	else if(cpu_temp >= conf_info[6]) {
			pwm = fmt.Sprintf("%s%03d", "pwm_" , conf_info[7]) 
		}

    if (pwm != oldPwm) {
      b := []byte(pwm)
      _, err = port.Write(b)
      if err != nil {
        log.Fatalf("port.Write: %v", err)
      }
      fmt.Println(pwm)
      oldPwm = pwm
    }
    time.Sleep(time.Second)
  }
}

func readConfig(path string) (lines []string, err error) {
    var (
        file *os.File
        part []byte
        prefix bool
    )
    if file, err = os.Open(path); err != nil {
        return
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    buffer := bytes.NewBuffer(make([]byte, 0))
    for {
        if part, prefix, err = reader.ReadLine(); err != nil {
            break
        }
        buffer.Write(part)
        if !prefix {
            lines = append(lines, buffer.String())
            buffer.Reset()
        }
    }
    if err == io.EOF {
        err = nil
    }
    return
}

func read_cpu_tmp() (temp int, err error) {
  dat, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
  if err != nil {
    log.Fatalf("Error in config %v", err)
  }
  s := strings.TrimSpace(string(dat))
  i, _ := strconv.Atoi(s)
  temp = i /1000
  return
}
