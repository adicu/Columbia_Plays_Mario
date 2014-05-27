


# Installation

- install git
- clone the repo
- Run the install script
- Edit the scale of the gvbam to `2`
  - or copy the gvbam config to `~/.config/gvbam/config`
- Copy in your legally-attained gbc ROM
  - update the settings file with the name of the file

# Running it

- fill out your `settings.sh`
- in window 1:
  - `bash load_xvfb.sh`
- in window 2:
  - `bash load_emulator.sh`
- in window 3:
  - `bash twitch_streamer.sh`
- in window 4:
  - `bash all_windows.sh`
  - find the window ID for gvbam
  - `go run *.go WINDOW_ID`


## Help

Used scripts from [stritch's repo](https://github.com/strich/HeadlessCrowdEmulator) to get a headless emluator streamed to Twitch.tv.


