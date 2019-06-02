package lib

import(
    "fmt"
    "github.com/davecgh/go-spew/spew"
    "os"
    "syscall"
    s "strings"
    "regexp"
    "path/filepath"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

func Handle(err error) {
    if err != nil {
        fmt.Println("Error:", err)
    }
}

type Manifest struct {
    // Exported
    DryRun bool             `yaml:"Dry_Run"`
    SrcRoot filePath           `yaml:"Src_Relative_To"`
    DestRoot filePath `yaml:"Dest_Relative_To"`
    PkgSys pkgSys           `yaml:"Pkg_Sys"`
    SysType sysType        `yaml:"Sys_Type"`
    Defaults struct {
        Order uint          `yaml:"Order"`
        User string         `yaml:"User"`
        Group string        `yaml:"Group"`
        FileMode os.FileMode   `yaml:"File_Mode"`
        DirMode os.FileMode   `yaml:"Dir_Mode"`
    }                       `yaml:"Defaults"`
    Removals []removal      `yaml:"Removals"`
    Directories []dirToSync `yaml:"Directories"`
    Files []fileToSync      `yaml:"Files"`
    Symlinks []linkToCreate `yaml:"Symlinks"`

    // Unexported
    pathToManifest string
    mfstUid uint32
    mfstGid uint32
}

// Represents a condition which must be met for an action in the manifest to occur
type condition struct {
    Hosts []string          `yaml:"Hosts"`
    Pkgs []string           `yaml:"Pkgs"`
}

// Represents a file path which may contain environment variables
type filePath string

type envVarError struct {
    varName  string
}
func (e *envVarError) Error() string {
    return fmt.Sprintf("Environment variable %s not found.", e.varName)
}

var envVarRegEx = regexp.MustCompile("(?:[$])([a-zA-Z_]+[a-zA-Z0-9_]*)")
func (path *filePath) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var filePathStr string
    err := unmarshal(&filePathStr)
    if err != nil {
        return err
    }

    if s.Contains(filePathStr, "$") {
        fmt.Println("Found environment variable in", filePathStr)

        // Returns list of [ <whole match>, <submatch1>, <submatch2>, ... ]
        envVars := envVarRegEx.FindAllStringSubmatch(filePathStr, -1)
        
	      for _, result := range envVars {
    	      varName := result[1] // First (only) submatch, i.e. the variable name without $

    	      varValue, found := os.LookupEnv(varName)
    	      if !found {
        	      return &envVarError{varName}
    	      }
    	      fmt.Println("Replacing", result[0], "with", varValue)
    	      filePathStr = s.Replace(filePathStr, result[0], varValue, 1)
    	      fmt.Println(filePathStr)
 	      }
    }

    *path = filePath(filePathStr)
    return nil
}

type pkgSys string
var pkgSysEnum = [...]string{"pacman", "dpkg", "homebrew", "pkgng", "ignore"}


type sysType string
var sysTypeEnum = [...]string{"linux", "macos"}



type removal struct {
    PathToRemove filePath `yaml:"Path_To_Remove"`
    Order uint              `yaml:"Order"`
    If condition            `yaml:"If"`
}

type fileToSync struct {
    SrcPath filePath        `yaml:"Src_Path"`
    DestPath filePath   `yaml:"Dest_Path"`
    Order uint              `yaml:"Order"`
    User string             `yaml:"User"`
    Group string            `yaml:"Group"`
    FileMode os.FileMode       `yaml:"File_Mode"`
    If condition            `yaml:"If"`
}

type dirToSync struct {
    fileToSync              `yaml:",inline"`
    DirMode os.FileMode       `yaml:"Dir_Mode"`
}

type linkToCreate struct {
    LinkPath filePath          `yaml:"Link_Path"`
    TargetPath filePath        `yaml:"Target_Path"`
    If condition            `yaml:"If"`
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

    return mfst.normalize()
}

/*
    ** Normalizing manifest **

    Ensures conformity to configuration rules, including:
        * Existant shell environment variables
        * Only absolute paths (set from relative paths in file)
        * Valid operation order numbers
        * Valid permission modes
        * Valid pkgSys and sysType choices
        * No empty fields (set with defaults)
*/
func (mfst *Manifest) normalize() (error) {
    fmt.Println("Normalizing manifest")

    // Existant shell environment variables
    //SrcRoot
    //DestRoot
    


//
// Only absolute paths (set from relative paths in file)
// Valid operation order numbers
// Valid permission modes
// Valid pkgSys and sysType choices
// No empty fields (set with defaults)


    //fmt.Printf("%+v\n", mfst)
    spew.Dump(mfst)

    return nil
}

