package lib

import(
    "fmt"
    "github.com/davecgh/go-spew/spew"
    "os"
    "syscall"
    "path/filepath"
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
    // Exported
    OriginRoot string       `yaml:"OriginRoot"`
    DestRoot string         `yaml:"DestRoot"`
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

    // Unexported
    pathToManifest string
    mfstUid uint32
    mfstGid uint32
}

func LocateManifest(mfstPath string) (*Manifest) {
    fmt.Println("Resolving manifest location")

    /*
        ** Get the current working directory **
    */
    // See os.Executable() for go 1.8+

    workingDirectory, err := os.Getwd()
    if err != nil {
        panic(err)
    }

    // If not specified by user
    if mfstPath == "" {
        mfstPath = "manifest.yaml"
    }

    // Convert path to absolute
    if !filepath.IsAbs(mfstPath) {
        mfstPath = filepath.Join(workingDirectory, mfstPath)
    }

    fmt.Println("Manifest location resolved:", mfstPath)

    fmt.Println("Checking file existence")

    // Attempt to read filesystem info for specified path
    fileInfo, err := os.Stat(mfstPath)
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Println("Manifest", mfstPath, "does not exist")
            os.Exit(1)
        } else {
            // Panic if an error unrelated to file existence occured
            panic(err)
        }
    }

    // Attempt to load with current ownership and permissions
    stat, ok := fileInfo.Sys().(*syscall.Stat_t)
    if !ok { // Don't know why this returns a bool instead of an error!?
        fmt.Println("Not a syscall.Stat_t")
        os.Exit(1)
    }

    fmt.Println("Allocating for manifest")
    // Given existant file, create empty Manifest to load it into
    mfst := Manifest{
        pathToManifest: mfstPath,
        mfstUid: stat.Uid,
        mfstGid: stat.Gid,
    }

    return &mfst
}

func (mfst *Manifest) Load() (error) {
    fmt.Println("Loading manifest")

    fmt.Println("Opening manifest file")
    dat, err := ioutil.ReadFile(mfst.pathToManifest)
    if err != nil {
        return err
    }
    
    fmt.Println("Parsing manifest as YAML")
    err = yaml.UnmarshalStrict(dat, mfst)
    if err != nil {
        return err
    }

    return nil
}

/*
    ** Normalizing manifest **

    Ensures conformity to configuration rules, including:
        * Existant shell environment variables
        * Only absolute paths (set from relative paths in file)
        * Valid operation order numbers
        * Valid permission modes
        * Valid ioType and pkgTypes
        * No empty fields (set with defaults)
*/
func (mfst *Manifest) Normalize() (error) {
    fmt.Println("Normalizing manifest")

    //os.GetEnv(variable)

    //fmt.Printf("%+v\n", mfst)
    spew.Dump(mfst)

    return nil
}

