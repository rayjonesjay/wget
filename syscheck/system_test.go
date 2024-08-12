package syscheck

import (
	"testing"

)



func TestCheckOperatingSystem(t *testing.T){
	operatingSystems := []string{"linux","windows","darwin"}

	for _, os := range operatingSystems {
		t.Run(os, func(t *testing.T) {

			err := CheckOperatingSystem()
			
			// if the operaring system is windows , error will be returned
			if os == "windows" {

				// check if an error was detected for windows
				if err != nil {
					t.Errorf("expected an error for %s got nil",os)
				}
			}else{
				// if the os was not windows
				
			}
		})
	}
}
