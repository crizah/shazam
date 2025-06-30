import os
import spotipy
from spotipy.oauth2 import SpotifyClientCredentials
import yt_dlp
import ffmpeg
import numpy as np
from pathlib import Path
import hashlib
import struct
import json


# SPOTIPY_CLIENT_ID = '4cfb127f9a3549a598aad3e5bda188f2'
# SPOTIPY_CLIENT_SECRET = 'a6eeffd19d4a471dabe79fbbea15ab0f'


MP3_DIR = Path("downloaded_mp3")
WAV_DIR = Path("downloaded_wav")
MP3_DIR.mkdir(exist_ok=True)
WAV_DIR.mkdir(exist_ok=True)

# Spotify Authentication
sp = spotipy.Spotify(auth_manager=SpotifyClientCredentials(
    client_id=SPOTIPY_CLIENT_ID,
    client_secret=SPOTIPY_CLIENT_SECRET
))



SONG_ID_FILE = "song_ids.json"
YT_ID_FILE = "yt_ids.json"

song_id_list = {} # key is song id and value is song title
yt_id_list = {} # key is the song ID and value is the yt video id


def load_ids():
    global song_id_list, yt_id_list
    if Path(SONG_ID_FILE).exists():
        with open(SONG_ID_FILE, "r") as f:
            song_id_list = json.load(f)
    if Path(YT_ID_FILE).exists():
        with open(YT_ID_FILE, "r") as f:
            yt_id_list = json.load(f)

def save_ids():
    with open(SONG_ID_FILE, "w") as f:
        json.dump(song_id_list, f)
    with open(YT_ID_FILE, "w") as f:
        json.dump(yt_id_list, f)


def get_playlist_tracks(playlist_url):
    results = sp.playlist_tracks(playlist_url)
    tracks = []
    names = [] # names of the songs 

    while results:
        for item in results['items']:
            track = item['track']
            if track is None:
                continue
            name = track['name']
            names.append(name)
            artists = ', '.join([a['name'] for a in track['artists']])
            full_title = f"{name} - {artists}"
            print(full_title, ", ")
            tracks.append(full_title)

        if results['next']:
            results = sp.next(results)
        else:
            results = None
    return tracks

def IDalreadyexists(id, list):
    for key in list.keys():
        if key == id:
            return True
        
    return False

def createSongID(song_title):
    #unique for every song_title

    # uint32 type for songID
    sha1 = hashlib.sha1(song_title.encode('utf-8')).digest()  # 20 bytes

    # Take the first 4 bytes and convert to uint32_t
    return struct.unpack('>I', sha1[:4])[0]  




def download_mp3(search_query, output_path):

    # if songtitle already exists in the folder, dont download again

    ydl_opts = {
        'format': 'bestaudio/best',
        'outtmpl': str(output_path),
        'noplaylist': True,
        'quiet': True,
        'postprocessors': [{
            'key': 'FFmpegExtractAudio',
            'preferredcodec': 'mp3',
            'preferredquality': '192',
        }],
    }

   

    with yt_dlp.YoutubeDL(ydl_opts) as ydl:

        try:
        
            info = ydl.extract_info(f"ytsearch1:{search_query}", download=True)

            if 'entries' in info:
               video_info = info['entries'][0]
            else:
               video_info = info  # fallback if not a search

            video_id = video_info.get('id')
            
            # print(type(video_id)) # string
            song_id= createSongID(search_query) # unique

            # should probably seperate the song name and artist name insteadof one big title but oh well

            if( IDalreadyexists(song_id, song_id_list) == False):
                #dd to list only if not already present in there                
                song_id_list[song_id] = search_query
                yt_id_list[song_id] = video_id

        except Exception as e:
            print(f"Error downloading {search_query}: {e}")


def convert_mp3_to_wav(mp3_file: Path, wav_file: Path):
    try:
        ffmpeg.input(str(mp3_file)).output(str(wav_file)).run(overwrite_output=True, quiet=True)
    except ffmpeg.Error as e:
        print(f"FFmpeg error: {e}")

def main(playlist_url):
    print("Fetching Spotify playlist metadata...")
    load_ids()  # Load previous song IDs
    tracks = get_playlist_tracks(playlist_url)
    print(f"Found {len(tracks)} tracks.")

    for track in tracks:
        print(f"Processing: {track}")
        safe_filename = track.replace("/", "_").replace("?", "").replace("\\", "_")
        mp3_path = MP3_DIR / f"{safe_filename}.mp3"
        wav_path = WAV_DIR / f"{safe_filename}.wav"

        if not mp3_path.exists():
            download_mp3(track, mp3_path)

        if mp3_path.exists() and not wav_path.exists():
            convert_mp3_to_wav(mp3_path, wav_path)

    save_ids()  # Save after processing
    print("All tracks processed.")


if __name__ == "__main__":
    playlist_link = "54urz9eVTb5kDaAhAh2vHY"
    main(playlist_link)

    # assign a unique ID to the song along with the youtube video id for it 
    # didnt convert all to WAV file

    # get_playlist_tracks(playlist_link)
    # download_mp3("Sheep - 2018 Remix - Pink Floyd", "downloaded_mp3" )

