package main

import(
    "fmt"
    "os"
    //"github.com/davecgh/go-spew/spew"
    "github.com/krlanguet/curator/lib"
)

func main() {
    fmt.Println("Starting krlanguet/curator")

    /*
        ** Lookup manifest **
        
        Uses manifest specified or defaults to "manifest.yaml" in working directory.
        Returns pointer to mostly empty lib.Manifest type.
    */
    var mfst *lib.Manifest
    if len(os.Args) > 1 {
        mfst = lib.LocateManifest(os.Args[1])
    } else {
        mfst = lib.LocateManifest("")
    }

    /*
        ** Load manifest **
        
        Populates lib.Manifest type by marshalling from YAML.
    */
    err := mfst.Load()

    /*
        ** Normalize manifest **

        Ensures conformity to configuration rules, including:
            * Valid operation order numbers
            * Valid permission modes
            * Valid ioType and pkgTypes
    */
    err = mfst.Normalize()
    lib.Handle(err)

    /*
        ** Check manifest installabilityy **

        Checks whether manifest can be installed COMPLETELY.
    */
    err = mfst.CheckInstall()
    lib.Handle(err)
    
    /*
        ** Conditionally instal manifest **
    */
    if mfst.IoType != "noop" {
        err = mfst.Install()
        lib.Handle(err)
    }
}

