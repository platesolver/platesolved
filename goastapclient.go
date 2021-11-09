// saves wcs and ini files, sunday night, monday morning, need sky to test now. 
package main

import (
        "bytes"
        "fmt"
        "io"
        "io/ioutil"
        "mime/multipart"
        "log"
        "net/http"
        "os"
        "flag"
        "strings"
)


const submitUri string =  "https://dev.platesolver.com/api/submit"
const solverUri string =  "https://dev.platesolver.com/api/_server.php"




func BaseName(s string) string {
   n := strings.LastIndexByte(s, '.')
   if n == -1 { return s }
   return s[:n]
}

func usage() {
        fmt.Fprintf(os.Stderr, "\n\tPLATESOLVER.COM (c)2021\n\tAkadata Limited\n")
        fmt.Fprintf(os.Stderr, "\t\tonline astrometry plate solveer\n\n")
    fmt.Fprintf(os.Stderr, "\t\tusage: platesolver.exe\n")
    fmt.Fprintf(os.Stderr, "\t\tcmdline: -f -r -z -fov -ra -spd -s\n");
    fmt.Fprintf(os.Stderr, "\t\twith: N.I.N.A \n")
    fmt.Fprintf(os.Stderr, "\t\tversion: v.0.1.alpha \n")
    fmt.Fprintf(os.Stderr, "\t\tby: Andrew Smalley\n")
    fmt.Fprintf(os.Stderr, "\t\taudiance: RESTRICTED\n")
    fmt.Fprintf(os.Stderr, "--------------------------------------\n")
    fmt.Fprintf(os.Stderr, "shush, dont tell everyone but there are\n")
    fmt.Fprintf(os.Stderr, "options available on request....")
        flag.PrintDefaults()
    os.Exit(1)
}

func UploadTheImageToSolve(ps_filename string, solverUri string, apikey string, ps_radius string, ps_z string, ps_fov string, ps_ra string, ps_spd string, ps_s string) error {
//func UploadTheImageToSolve(ps_filename string, solverUri string, apikey string, ps_radius string, ps_z string, ps_fov string, ps_ra string, ps_spd string, ps_s string, ps_t string, ps_m string, ps_speed, ps_o string, ps_analyse string, ps_extract string, ps_log string, ps_progress string, ps_update string, ps_wcs string) {
        client := &http.Client{}
        bodyBuf := &bytes.Buffer{}
        bodyWriter := multipart.NewWriter(bodyBuf)

        fileWriter, err := bodyWriter.CreateFormFile("platesolve", ps_filename)
        if err != nil {
                fmt.Println("error writing to buffer")
                return err
        }
        f, err := os.Open(ps_filename)
        if err != nil {
                fmt.Println("error open file")
                return err
        } 
        _, err = io.Copy(fileWriter, f)
        if err != nil {
                return err
        }
        contentType := bodyWriter.FormDataContentType()
        bodyWriter.Close()


        req, err := http.NewRequest("POST", solverUri, bodyBuf)
        req.Header.Set("Content-Type", contentType)
        req.Header.Set("User-Agent", "go-platesolver/0.1.1/dev")
        req.Header.Set("X-CLIENT-PLATESOLVER-VERSION", "com.plateserver.client.c000") // set the header of the useragenet for the platesolver client
    req.Header.Set("X-CLIENT-PLATESOLVER-ABOUT", "your favorite platesolver platesolved you")
    req.Header.Set("X-CLIENT-PLATESOLVER-F", "\""+ps_filename+"\"") //" filename  {fits, tiff, png, jpg files}",
    req.Header.Set("X-CLIENT-PLATESOLVER-FOV", ps_fov) // "diameter_field[degrees]",
    req.Header.Set("X-CLIENT-PLATESOLVER-Z", ps_z) // "downsample_factor[0,1,2,3,4] {Downsample prior to solving. 0 is auto}",
    req.Header.Set("X-CLIENT-PLATESOLVER-S", ps_s) //"max_number_of_stars  {default 500}",
    req.Header.Set("X-CLIENT-PLATESOLVER-R", ps_radius) //"radius_area_to_search[degrees]",
    req.Header.Set("X-CLIENT-PLATESOLVER-RA", ps_ra) //"center_right ascension[hours]",
    req.Header.Set("X-CLIENT-PLATESOLVER-SPD", ps_spd) //"center_south_pole_distance[degrees]",

        resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
    }
    defer resp.Body.Close()


   // fmt.Println("response Status:", resp.Status)
   // fmt.Println("response Headers:", resp.Header)

    if resp.StatusCode == http.StatusOK {
        bodyBytes, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                log.Fatal(err)
        }
        //bodyString := string(bodyBytes)
        bodyStrings := strings.Split(string(bodyBytes), "<=====>") // string(bodyBytes)
   
    //  fmt.Println(bodyStrings[0]) 
                ini_output_file:=BaseName(ps_filename)+".ini"
        //      fmt.Println(ini_output_file)
                fw, err := os.Create(ini_output_file)
                if err != nil {
                        fmt.Println(err)
                }
                l, err := fw.WriteString(bodyStrings[0])
                if err != nil {
                        fmt.Println(err)
                        fw.Close()
                }
                //fmt.Println(l, "INFO: "+ini_output_file+" written successfully")
                err = fw.Close()
                if err != nil {
                        fmt.Println(err)
                }
                wcs_output_file:=BaseName(ps_filename)+".wcs"
                //fmt.Println(wcs_output_file)
                fi, err := os.Create(wcs_output_file)
                if err != nil {
                        fmt.Println(err)
                }
                fl, err := fi.WriteString(bodyStrings[1])
                if err != nil {
                        fmt.Println(err)
                        fi.Close()
                }
        //      fmt.Println(fl, "INFO: "+wcs_output_file+" written successfully")
                err = fi.Close()
                if err != nil {
                        fmt.Println(err)
                }
        }
        return err
} 

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
        // variables declaration  
    // platesolver.com system settings
   // var apiSolverUri string
   // apiSolverUri = "https://dev.platesolver.com/api/login.php"
    // uri used for login preauth/ get session key / get solve id presolve. 
    //var authUri string
    //authUri = "api/login.php"
    // uri used to solve the image, will have solve id, presolve id added.
       //s var formData string
    // platesolver.com user settings
    var uname string    
    var pass string
    var apikey string
    var config string

    // flags user declaration using flag package
    flag.StringVar(&uname, "u", "username", "Specify username. Default is yours")    
    flag.StringVar(&pass, "p", "password", "Specify password. Default is yours")
    flag.StringVar(&apikey, "apikey", "apikey", "Specify apikey. Default is yours")
    flag.StringVar(&config, "config", "config", "define config. Default is yours")


    // platesolver.com auth vars
    //  mode platesolver_simple, platesolver_advanced
    var mode string
    mode = "platesolver_simple"  

    // latesolver interface to ASTAP
    // basic N.I.N.A options uses, these are free. 
    // platesolver_simple
        var ps_filename string
        var ps_radius string
        var ps_z string
        var ps_fov string
        var ps_ra string
        var ps_spd string
        var ps_s string
        // platesolver interface to ASTAP
        // options available to the free version 
        // -f -r -z -fov -ra -spd -s 
        //
        flag.StringVar(&ps_filename, "f", "", "-f  filename  fits, tiff, png, jpg files ")
        flag.StringVar(&ps_radius, "r", "30", "-r  stdin  radius area to search [degrees]")
        flag.StringVar(&ps_z, "z", "0", "-z  downsample_factor[0,1,2,3,4] Downsample prior to solving. 0 is auto")
        flag.StringVar(&ps_fov, "fov", "", "-fov diameter_field[degrees")
        flag.StringVar(&ps_ra, "ra", "", "-ra  center_right ascension[hours]")
        flag.StringVar(&ps_spd, "spd", "", "-spd center_south_pole_distance[degrees]")
        flag.StringVar(&ps_s, "s", "500", "-s  max_number_of_stars  (default 500)")
        // other ASTAP options, these could cost.
        // platesolver_advanced
        var ps_t string
        var ps_m string
        var ps_speed string
        var ps_o string
        var ps_analyse string
        var ps_extract string
        var ps_log string
        var ps_progress string
        var ps_update string
        var ps_wcs string  
        if(mode == "platesolver_advanced") {
                flag.StringVar(&ps_t, "t", "0.007", "-t  tolerance tolerance default 0.007")
                flag.StringVar(&ps_m, "m", "1.5", "-m  minimum_star_size[] default 1.5")
                flag.StringVar(&ps_speed, "speed", "auto", "-speed mode[auto/slow] (Slow is forcing small search steps to improve detection.)")
                flag.StringVar(&ps_o, "o", "", "-o  file Name the output files with this base path & file name")
                flag.StringVar(&ps_analyse, "analyse", "", "-analyse snr_min Analyse only and report median HFD and number of stars used")
                flag.StringVar(&ps_extract, "extract", "", "-extract snr_min As -analyse but additionally write a .csv file with the detected stars info")
                flag.StringVar(&ps_log, "log", "", "-log   Write the solver log to file")
                flag.StringVar(&ps_progress, "progress", "", "-progress  Log all progress steps and messages")
                flag.StringVar(&ps_update, "update", "", "-update  update the FITS header with the found solution. Jpg, png, tiff will be written as fits")
                flag.StringVar(&ps_wcs, "wcs", "", "-wcs Write a .wcs file  in similar format as Astrometry.net. Else text style.")   
        }
    flag.Parse()  
    UploadTheImageToSolve(ps_filename, solverUri, apikey, ps_radius, ps_z, ps_fov, ps_ra, ps_spd, ps_s)
    os.Exit(0)
}
