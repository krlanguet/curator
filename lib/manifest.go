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
    Hosts []string          `yaml:"Hosts"`
    Pkgs []string           `yaml:"Pkgs"`
}

type unlink struct {
    PathToRemove string     `yaml:"PathToRemove"`
    Order uint              `yaml:"Order"`
    If condition            `yaml:"If"`
}

type fileToSync struct {
    SourcePath string       `yaml:"SourcePath"`
    DestinationPath string  `yaml:"DestinationPath"`
    Order uint              `yaml:"Order"`
    User string             `yaml:"User"`
    Group string            `yaml:"Group"`
    Fmode mode              `yaml:"Fmode"`
    If condition            `yaml:"If"`
}

type dirToSync struct {
    fileToSync              `yaml:",inline"`
    Dmode mode              `yaml:"Dmode"`
}

type linkToCreate struct {
    LinkPath string         `yaml:"LinkPath"`
    TargetPath string       `yaml:"TargetPath"`
    If condition            `yaml:"If"`
}

type Manifest struct {
    OriginRoot string       `yaml:"OriginRoot"`
    TargetRoot string       `yaml:"TargetRoot"`
    PkgType pkgType         `yaml:"PkgType"`
    IoType ioType           `yaml:"IoType"`
    Defaults struct {
        Order uint          `yaml:"Order"`
        User string         `yaml:"User"`
        Group string        `yaml:"Group"`
        Fmode mode          `yaml:"Fmode"`
        Dmode mode          `yaml:"Dmode"`
    }                       `yaml:"Defaults"`
    Unlinks []unlink        `yaml:"Unlinks"`
    Directories []dirToSync `yaml:"Directories"`
    Files []fileToSync      `yaml:"Files"`
    Symlinks []linkToCreate `yaml:"Symlinks"`
}

func ReadManifest(pathToManifest string) (Manifest) {
    fmt.Println("Reading manifest")

    fmt.Println("Opening manifest file")
    dat, err := ioutil.ReadFile(pathToManifest)
    Handle(err)

    fmt.Println("Parsing manifest as YAML")
    var mfst Manifest
    err = yaml.UnmarshalStrict(dat, &mfst)
    Handle(err)

    return mfst
}

func (mfst *Manifest) Normalize() (error) {
    fmt.Println("Normalizing manifest")

    //fmt.Printf("%+v\n", mfst)
    spew.Dump(mfst)

    return nil
}

