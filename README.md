
# Twilio Plays Pokemon


## Installation

- install git
- clone the repo
- Run the install script
- Edit the scale of the gvbam to `2`
  - or copy the gvbam config to `~/.config/gvbam/config`
- Copy in your legally-attained gbc ROM to the `twitch_config` directory
  - update the settings file with the name of the file

## Running it

- fill out your `settings.sh`
- go into the `twitch_config` directory
- in window 1:
  - `bash load_xvfb.sh`
  - `bash load_emulator.sh`
- in window 2:
  - `bash twitch_streamer.sh`
- in window 1:
  - `make`
  - `bash all_windows.sh`
  - find the window ID for gvbam
  - `./pokebot WINDOW_ID`



### Help

Used scripts from [stritch's repo](https://github.com/strich/HeadlessCrowdEmulator) to get a headless emluator streamed to Twitch.tv.


