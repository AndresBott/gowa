package gowa

import (
	log "github.com/sirupsen/logrus"
	"fmt"
	"os/user"
	"os"
	"encoding/base64"
	"github.com/asticode/go-astilectron"
)

type settings struct {
	installLocation string
	electronName string
}

var config = settings{
	installLocation: "",
	electronName: "electron",
}

func preStartChecks(){

	path := os.Getenv("GOWAPATH")
	if path == ""{

		usr, err := user.Current()
		if err != nil {
			log.Fatal( err )
		}
		fmt.Println( usr.HomeDir )
		config.installLocation = usr.HomeDir+"/.gowa"

		log.Info("GOWAPATH: "+config.installLocation)
	}else{
		config.installLocation = path
		log.Info("GOWAPATH: "+config.installLocation)
	}


	if _, err := os.Stat(config.installLocation); os.IsNotExist(err) {
		log.Info("Creating dir GOWAPATH: "+config.installLocation)
		os.MkdirAll(config.installLocation, 0777)
		os.MkdirAll(config.installLocation+"/"+config.electronName, 0777)
	}

	if _, err := os.Stat(config.installLocation+"/icon.png"); os.IsNotExist(err) {

		log.Info("Creating icon GOWAPATH: "+config.installLocation+"/icon.png")
		dec, err := base64.StdEncoding.DecodeString(gowaIcon)
		if err != nil {
			panic(err)
		}

		f, err := os.Create(config.installLocation+"/icon.png")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			panic(err)
		}
		if err := f.Sync(); err != nil {
			panic(err)
		}
	}
}

func Start()  {

	preStartChecks()

	 //Initialize astilectron
	var a, err = astilectron.New(astilectron.Options{
		AppName: "Gowa",
		AppIconDefaultPath: config.installLocation+"/icon.png", // If path is relative, it must be relative to the data directory
		BaseDirectoryPath: config.installLocation+"/"+config.electronName,
	})
	if err != nil {
		fmt.Println(err)
	}

	defer a.Close()

	// Start astilectron (download electron)
	log.Info("Starting GOWA, this can take some time")
	a.Start()

	// Blocking pattern
	// Create a new window
	var w, werr = a.NewWindow("https://outlook.office365.com/owa/", &astilectron.WindowOptions{
		Center: astilectron.PtrBool(true),
		Height: astilectron.PtrInt(900),
		Width:  astilectron.PtrInt(1200),
		//inject electron into window

		//Custom: &astilectron.WindowCustomOptions{
		//
		//	Script: "alert(\"G\")",
		//
		//},
	})

	if err != nil{
		log.Error(werr)
	}
	w.Create()
	//w.OpenDevTools()
	a.Wait()
}


