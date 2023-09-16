package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func getCurrentVolume() uint8 {
	amixer := exec.Command("amixer", "sget", "'Master'")
	outputBytes, err := amixer.Output()

	if err != nil {
		log.Fatal(err)
	}

	output := string(outputBytes)
	volume, err := strconv.ParseUint(strings.Split(strings.Split(output, "[")[1], "%")[0], 10, 8)

	if err != nil {
		log.Fatal(err)
	}

	return uint8(volume)
}

func increaseVolume(amount uint8) {
	currVolume := getCurrentVolume()
	if currVolume+amount > 100 {
		setVolume(100)
	} else {
		setVolume(currVolume + amount)
	}
}

func decreaseVolume(amount uint8) {
	currVolume := getCurrentVolume()
	if amount > currVolume {
		setVolume(0)
	} else {
		setVolume(currVolume - amount)
	}
}

func toggleMute() {
	currVolume := getCurrentVolume()
	filePath := os.Getenv("HOME") + "/.config/prev_vol"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(filePath)

		if err != nil {
			log.Fatal(err)
		}

		file.Close()
	}

	if currVolume > 0 {
		output := fmt.Sprintf("%v", currVolume)
		err := os.WriteFile(filePath, []byte(output), 0666)

		if err != nil {
			log.Fatal(err)
		}
		setVolume(0)
	} else {
		prevVolBytes, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		prevVol, err := strconv.ParseUint(string(prevVolBytes), 10, 8)
		if err != nil {
			log.Fatal(err)
		}

		setVolume(uint8(prevVol))
	}
}

func setVolume(level uint8) {
	if level > 100 {
		return
	}

	amixer := exec.Command("amixer", "sset", "'Master'", fmt.Sprintf("%v%%", level))
	err := amixer.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func cli() {
	helpMsg := `Usage: volumectl-go [OPTIONS]

Options:
  -g, --get          Get current volume
  -i, --inc <VALUE>  Increase volume by this value
  -d, --dec <VALUE>  Decrease volume by this value
  -t, --toggle-mute  Toggle mute
  -h, --help         Print help
`

	args := os.Args
	if len(args) == 1 {
		fmt.Println(helpMsg)
		os.Exit(0)
	}

	for idx, arg := range args[1:] {
		switch arg {
		case "-i", "--inc":
			amount, err := strconv.ParseUint(args[idx+2], 10, 8)
			if err != nil {
				log.Fatalln(err)
			}

			increaseVolume(uint8(amount))
			return
		case "-d", "--dec":
			amount, err := strconv.ParseUint(args[idx+2], 10, 8)
			if err != nil {
				log.Fatalln(err)
			}

			decreaseVolume(uint8(amount))
			return
		case "-t", "--toggle-mute":
			toggleMute()
			return
		case "-g", "--get":
			fmt.Printf("Current volume: %v\n", getCurrentVolume())
			return
		case "-h", "--help":
			fmt.Println(helpMsg)
			return
		}
	}

	log.Fatalln("You provided wrong flags.")
}

func main() {
	cli()
}
