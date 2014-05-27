#!/bin/sh

# Settings that will work for Twitch streaming:
VI_RES="320x288" # Captured area of the desktop
VI_OFFSET="0,0"
VI_FREQ="30"
 
VO_RES="320x288"
VO_TUNE="ultrafast"
VO_BITS="1000k"
 
AO_FREQ="44100"
AO_BITS="128k"
 
BUFFER="512k"
THREADS="2"
 
KEY=$TWITCH_KEY
SERVER="rtmp://live-ams.justin.tv/app/"

avconv \
-f x11grab -s "$VI_RES" -r "$VI_FREQ" -i :0.0+$VI_OFFSET \
-codec:v libx264 -s "$VO_RES" -preset $VO_TUNE -b:v "$VO_BITS" -pix_fmt yuv420p -g 2 \
-codec:a libmp3lame -ar "$AO_FREQ" -b:a "$AO_BITS" \
-f flv "$SERVER$KEY" -threads $THREADS -bufsize "$BUFFER"
#-f alsa -ac 2 -i pulse \
