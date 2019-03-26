package lib

import(
    "fmt"
    //"path"
    //"os"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

func handle(err error) {
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

type dirToSync struct {
    SourcePath string       `yaml:"src"`
    DestinationPath string  `yaml:"dst"`
    Order uint
    User string
    Group string
    Fmode mode
    Dmode mode
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

type linkToCreate struct {
    LinkPath string         `yaml:"dst"`
    TargetPath string       `yaml:"src"`
    OnlyIf condition
}

type manifest struct {
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

func ReadManifest(pathToManifest string) (manifest) {
    fmt.Println("Reading manifest")

    fmt.Println("Opening manifest file")
    dat, err := ioutil.ReadFile(pathToManifest)
    handle(err)

    fmt.Println("Parsing manifest as YAML")
    var mfst manifest
    err = yaml.Unmarshal(dat, &mfst)
    handle(err)

    fmt.Println("Normalizing manifest")
    err = mfst.normalize()
    handle(err)

    return mfst
}

func (mfst *manifest) normalize() (error) {
    fmt.Println("Normalizing manifest")

    return nil
}

