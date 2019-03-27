package main

import(
    "fmt"
    //"errors"
    "os"
    "syscall"
    "path/filepath"
    //"github.com/davecgh/go-spew/spew"
    "github.com/krlanguet/curator/lib"
)

func main() {
    fmt.Println("Starting krlanguet/curator")

    /*
        ** Get the current working directory **
    */
    
    /*ex, err := os.Executable() // Absolute path to executable of current process - for go 1.8+
    if err != nil {
        panic(err)
    }
    workingDirectory := filepath.Dir(ex)
    */
    workingDirectory, err := os.Getwd()
    if err != nil {
        panic(err)
    }

    /*
        ** Use manifest specified or look for default manifest **
    */
    var mfstPath string
    if len(os.Args) > 1 {
        // Convert any specified path to absolute
        mfstPath = os.Args[1]
        if filepath.IsAbs(mfstPath) {
            fmt.Println("Absolute path to manifest file specified:", mfstPath)
        } else {
            mfstPath = filepath.Join(workingDirectory, mfstPath)
            fmt.Println("Absolute path to specified manifest file:", mfstPath)
        }
    } else {
        mfstPath = filepath.Join(workingDirectory, "manifest.yaml")
        fmt.Println("Default manifest file used:", mfstPath)
    }

    // Attempt to read filesystem info on file at specified path
    fileInfo, err := os.Stat(mfstPath)
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Println("Manifest", mfstPath, "does not exist")
            os.Exit(1)
        } else {
            // Panic if an error unrelated to file existence occured
            panic(err)
        }
    } else {
        stat, ok := fileInfo.Sys().(*syscall.Stat_t)
        if !ok {
            fmt.Println("Not a syscall.Stat_t")
        } else {
            fmt.Println("Manifest Uid:", stat.Uid)
            fmt.Println("Manifest Gid:", stat.Gid)
        }
    }

    /*
        ** Reading manifest **
        
        Returns pointer to lib.Manifest type, marshalled from YAML.
    */
    mfst := lib.ReadManifest(mfstPath)

    /*
        ** Normalizing manifest **

        Ensures conformity to configuration rules, including:
            * Valid operation order numbers
            * Valid permission modes
            * Valid ioType and pkgTypes
    */
    err = mfst.Normalize()
    lib.Handle(err)

    /*
        ** Performing manifest install DRY-RUN **

        Checks whether manifest can be installed COMPLETELY.
    */
    err = mfst.CheckInstall()
    lib.Handle(err)
    
    /*
        ** Conditionally installing manifest **
    */
    if mfst.IoType != "noop" {
        err = mfst.Install()
        lib.Handle(err)
    }
}

