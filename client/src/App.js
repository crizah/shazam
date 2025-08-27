import React, { useState } from "react";





import sha1 from "js-sha1"; 

const access_token ="BQDUz0dE1GVnZXHyfMa_MPn3FC_oLWNfWxlBFA2k0QQnBSakSX39dr_gIGAnhMbSV0of0ATDn_8d3tgbTIhEVj8gG61l9AAqi28z1SOEjj7PMIeGRB_Ag_dHuM2JB94V-JRWQ6V-xQ8"
// curl -X POST "https://accounts.spotify.com/api/token" \
//      -H "Content-Type: application/x-www-form-urlencoded" \
//      -d "grant_type=client_credentials&client_id=4cfb127f9a3549a598aad3e5bda188f2&client_secret=a6eeffd19d4a471dabe79fbbea15ab0f"

function createSongID(name, artist, album) {
  const input = name + artist + album;
  const hash = sha1.arrayBuffer(input); // gives you raw bytes
  const view = new DataView(hash);
  return view.getUint32(0); // get first 4 bytes as uint32
}

function App() {
  const [playlistURL, setPlaylistURL] = useState(""); // const [state, setState] = useState(initialState)
  const [tracks, setTracks] = useState([]);
  const [Results, setResults] = useState({});  // object, send this to Go server 
  const [status, setStatus] = useState("");

  const extractPlaylistId = (url) => {
    // https://open.spotify.com/playlist/54urz9eVTb5kDaAhAh2vHY

    const excess = "https://open.spotify.com/playlist/"
    return url.slice(excess.length); 
  };




const getPlaylistTracks = async (playlistId) => {
  const headers = {
    Authorization: `Bearer ${access_token}`,
  };

  // let result = {};

  const result = [];

  let nextUrl = `https://api.spotify.com/v1/playlists/${playlistId}/tracks`;

  while (nextUrl) {
    const res = await fetch(nextUrl, { headers });
    const data = await res.json();

    if (data.items) {
      for (let i = 0; i < data.items.length; i++) {
        const currentItem = data.items[i];

        if (currentItem.track) {
          const track = currentItem.track;
          const name = track.name;
          

          const album = track.album.name;
          

          let artists = "";
          

          for (let j = 0; j < track.artists.length; j++) {
            artists += track.artists[j].name;
            if (j < track.artists.length - 1) {
              artists += ", ";
            }
          }

          const fullTitle = `${name} - ${artists}`;
          const songId = createSongID(name, artists, album)



          // result[songId] = {
          //   name: name,
          //   artist: artists,
          //   album: album,
          // };


          result.push(fullTitle);
        }
      }
    }

    nextUrl = data.next;
  }

  return result;
};


const sendSongsToGoServer = async (results) => {
  try{
    
    const ws = new WebSocket("ws://localhost:1058");

    ws.onopen = () => {
      console.log("Connected to server ");
      ws.send(results);
    } 
  }catch (err){

    console.error("Error sending to server:", err);
    setStatus("Failed to send data to server");
  }
};




    





//   try {
//     const response = await fetch("http://localhost:1058/get_songs", {
//       method: "POST",
//       headers: {
//         "Content-Type": "application/json",
//       },
//       body: JSON.stringify(results),
//     });

//     if (!response.ok) {
//       throw new Error("Failed to send data");
//     }

//     const resultText = await response.text(); // or response.json() if server returns JSON
//     console.log("Server response:", resultText);
//     setStatus("Data sent to server successfully!");
//   } catch (err) {
//     console.error("Error sending to server:", err);
//     setStatus("Failed to send data to server");
//   }
// };




  



  const handleSubmit = async (e) => {
    e.preventDefault();
    setStatus("Fetching...");
    const playlistId = extractPlaylistId(playlistURL);

    if (!playlistId) {
      setStatus("Invalid playlist URL");
      return;
    }

    try {
      const result = await getPlaylistTracks(playlistId);     
      // setResults(result)  // send to Go server  
      // setTracks(Object.values(result)); // all the names, artists, and albums      
    
      // setStatus(`Fetched ${result.size} tracks`);

      const r = JSON.stringify(result);

      // await sendSongsToGoServer(result);
      await sendSongsToGoServer(r);

      // i dont want to send an object. i want to convert it into a string

    } catch (err) {
      console.error(err);
      setStatus("Failed to fetch tracks");
    }
  };


  return (
    <div>
      <h2>Spotify Playlist Fetcher</h2>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={playlistURL}
          placeholder="Paste Spotify playlist URL"
          onChange={(e) => setPlaylistURL(e.target.value)}
          style={{ width: "400px" }}
        />
        <button type="submit">Get Tracks</button>
      </form>
      <p>{status}</p>
      <ul>
         {tracks.map((t, i) => (
            <li key={i}>
             {t.name} - {t.artist} - {t.album}
            </li>
          ))}
      </ul>
    </div>
  );
}



export default App;



// import { useEffect } from "react";

// function App() {
//   useEffect(() => {
//     const ws = new WebSocket("ws://localhost:1058");

//     ws.onopen = () => {
//       console.log("Connected to server ");
//       ws.send("Hello from client");
//     };

//     // ws.onmessage = (event) => {
//     //   console.log("Received:", event.data);
//     // };

//     // ws.onclose = () => {
//     //   console.log("Disconnected");
//     // };

//     ws.onerror = (err) => {
//       console.error("Error:", err);
//     };
//   }, []);

//   return <h1>WebSocket Test</h1>;
// }

// export default App;
