package lib

import(
    "fmt"
    "github.com/davecgh/go-spew/spew"
    //"path"
    //"os"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

func Handle(err error) {
    if err != nil {
        fmt.Println("Error:", err)
    }
}

type mode string
type pkgType string
type ioType string

type condition struct {
    Hosts []string
    Pkgs []string
}

type unlink struct {
    PathToRemove string     `yaml:"src"`
    Order uint
    OnlyIf condition        `yaml:"if"`
}

type fileToSync struct {
    SourcePath string       `yaml:"src"`
    DestinationPath string  `yaml:"dst"`
    Order uint
    User string
    Group string
    Fmode mode
    OnlyIf condition        `yaml:"if"`
}

type dirToSync struct {
    fileToSync              `yaml:",inline"` // Inherits fileToSync fields
    Dmode mode
}

type linkToCreate struct {
    LinkPath string         `yaml:"dst"`
    TargetPath string       `yaml:"src"`
    OnlyIf condition
}

type Manifest struct {
    OriginRoot string
    TargetRoot string
    PkgType pkgType
    IoType ioType
    Defaults struct {
        Order uint
        User string
        Group string
        Fmode mode
        Dmode mode
    }
    Unlinks []unlink
    Directories []dirToSync
    Files []fileToSync
    Symlinks []linkToCreate
}

func ReadManifest(pathToManifest string) (Manifest) {
    fmt.Println("Reading manifest")

    fmt.Println("Opening manifest file")
    dat, err := ioutil.ReadFile(pathToManifest)
    Handle(err)

    fmt.Println("Parsing manifest as YAML")
    var mfst Manifest
    err = yaml.Unmarshal(dat, &mfst)
    Handle(err)

    return mfst
}

func (mfst *Manifest) Normalize() (error) {
    fmt.Println("Normalizing manifest")

    //fmt.Printf("%+v\n", mfst)
    spew.Dump(mfst)

    return nil
}

