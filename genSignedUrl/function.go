package function

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"

    "cloud.google.com/go/storage"
)

const (
    bucketName = "kikidoc-stg"
    fileName   = "setup.sh"
)

type Result struct {
    Code int
    URL  string
    Msg  string
}

func toJson(code int, url string, msg string) []byte {
    in := Result{code, url, msg}
    var buffer bytes.Buffer
    e := json.NewEncoder(&buffer)
    e.SetEscapeHTML(false)
    e.Encode(in)
    return buffer.Bytes()
}

func Handler(w http.ResponseWriter, r *http.Request) {
    privKey := os.Getenv("PRIVKEY")
    fmt.Print(privKey)
    expires := time.Now().Add(5*time.Minute)
    opts := &storage.SignedURLOptions{
        GoogleAccessID: "terraform@kikidoc-stg.iam.gserviceaccount.com",
        PrivateKey:     []byte(privKey),
        Method:         http.MethodGet,
        Expires:        expires,
    }
    url, err := storage.SignedURL(bucketName, fileName, opts)
    if err != nil {
        panic(err)
    }
    fmt.Print(url)
    w.Header().Set("Content-Type", "application/json")
    res := toJson(200, url, "success")
    w.Write(res)
}
