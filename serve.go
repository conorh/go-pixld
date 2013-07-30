package main

import (
  "fmt"
  "net/http"
  "errors"
  "log"
  "os"
  "github.com/nfnt/resize"
  "path/filepath"
  "strconv"
  "flag"
  "io"
  "image/jpeg"
//  "image/png"
  "net/url"
)

  // Functions
  // 1. Request comes in in the format
  // http://images.squaremill.com/s?img=&w=&h=&t=
  // 2. Check the cache for the image. If it exists already then serve it up
  // 3. Otherwise transform the image and serve it up
func getDiskLocation(u string, params map[string]string) (string, string, error) {
  parsedUrl, err := url.Parse(u)
  if err != nil {
    log.Fatal(err)
  }

  originalFilePath := "/tmp/" + parsedUrl.Host + parsedUrl.Path
  paramsSegment := "-"
  // create a params string
  for k, v := range params {
    paramsSegment = paramsSegment + k + "=" + v 
  }

  paramPath := originalFilePath + paramsSegment
  return originalFilePath, paramPath, nil
}

func validateNumberParam(param string, min int64, max int64) (string, error) {
  if(param == "") {
    return "", nil
  }
  value, err := strconv.ParseInt(param, 10, 64)
  if err != nil || value < min || value > max {
    return "", errors.New("not a valid number greater than 0 and less than 20000")
  }
  return param, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
    log.Println("Request:", r.URL)
    params := r.URL.Query()

    // Get the URL of the img from the required img parameter
    imageUrls := params["img"]
    if imageUrls == nil || len(imageUrls) < 1 {
      fmt.Fprintln(w, "No img parameter present")
      return
    }
    imageUrl := imageUrls[0]
    
    // Get the image transformation parameters
    var err error
    imageParams := make(map[string]string, 5)

    imageParams["w"], err = validateNumberParam(params.Get("w"), 1, 20000)
    if err != nil {
      fmt.Fprintln(w, "w param invalid:", err)
      return
    }

    imageParams["h"], err = validateNumberParam(params.Get("h"), 1, 20000)
    if err != nil {
      fmt.Fprintln(w, "h param invalid:", err)
      return
    }

    // Get the disk location for the image
    originalFilePath, filePath, err := getDiskLocation(imageUrl, imageParams)
    log.Println("File Location", originalFilePath, filePath)

    // Check to see if the file xists
    _, err = os.Stat(filePath)
    if err != nil && os.IsNotExist(err) {
      // Does not exist, try and create it
      log.Println("File does not exist, creating")

      // Download the remote file
      dir := filepath.Dir(originalFilePath)
      err = os.MkdirAll(dir, 0700)
      if err != nil {
        log.Println("Could not create directory", dir, err)
      }

      out, err := os.Create(originalFilePath)
      if err != nil {
        log.Println("Could not open", filePath, " for writing: ", err)
        return
      }
      defer out.Close()

      resp, err := http.Get(imageUrl)
      if err != nil {
        log.Println("Could not download", imageUrl, ": ", err)
        return
      }
      defer resp.Body.Close()

      n, err := io.Copy(out, resp.Body)
      if err != nil || n <= 0 {
        log.Println("Could not write out download", err)
        return
      }

      file, err := os.Open(originalFilePath)
      if err != nil {
        log.Println("Could not open downloaded file", err)
      }
      defer file.Close()

      img, err := jpeg.Decode(file)
      if err != nil {
        log.Println("Could not decode jpeg", err)
      }
      file.Close()

      width, _ := strconv.ParseInt(imageParams["w"], 10, 64)
      height, _ := strconv.ParseInt(imageParams["h"], 10, 64)
      m := resize.Resize(uint(width), uint(height), img, resize.Bicubic)

      out, err = os.Create(filePath)
      if err != nil {
        log.Println("Could not create file for writing", err)
      }
      defer out.Close()

      // write new image to file
      jpeg.Encode(out, m, nil)
    } else if err != nil {
      // There was an error accessing the file, note and return a 500 response
      log.Println("File system error accessing file", err)
      fmt.Fprintln(w, "There was an error accessing the file")
      return
    }

    http.ServeFile(w, r, filePath)
}

func main() {
  var port = flag.Int("port", 8080, "port to listen on")
  //dec := json.NewDecoder(os.Stdin)

  http.HandleFunc("/s", handler)
  log.Println("Listening on port", *port)
  http.ListenAndServe(":" + strconv.Itoa(*port), nil)
  log.Println("Terminating")
}