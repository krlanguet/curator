package main

import(
    "fmt"
    "github.com/krlanguet/curator/lib"
)

func main() {
    fmt.Println("Starting krlanguet/curator")

    lib.ReadManifest("/mnt/d/Projects/src/github.com/krlanguet/curator/manifest.example.yaml")
}
