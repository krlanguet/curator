package lib

import(
    "fmt"
    "path"
    //"gopkg.in/yaml.v2"
)

type mode string
type pkgType string
type ioType string

type Condition struct {
    Hosts []string
    Pkgs []string
}

type defaults struct {
    Order uint
    User string
    Group string
    Fmode mode
    Dmode mode
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
    Defaults defaults
    Unlinks []unlink
    Directories []dirToSync
    Files []fileToSync
    Symlinks []linkToCreate
}

func ReadManifest() {
  fmt.Println("Reading manifest")
}

func normalizeManifest() {
    fmt.Println("Normalizing manifest")

    x := "dir"
    fmt.Println(x)
    fmt.Println(path.Base(x))
    return
}
