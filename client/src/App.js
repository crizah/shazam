import React, { useState } from "react";

const access_token = "BQBAejPAh-rj0JsLkmFVOY6JIGeABgMnngyrDEu_tjtInlqrWoUW3saI_QhPgzchkQ7-2s7FDJAPmMntJY0bZwkSnAeL8XUiT7XCE4RSAuED6olfUrhktzfRh2SKMHsWI9KlRc7hbLY"; 


// curl -X POST "https://accounts.spotify.com/api/token" \
//      -H "Content-Type: application/x-www-form-urlencoded" \
//      -d "grant_type=client_credentials&client_id=4cfb127f9a3549a598aad3e5bda188f2&client_secret=a6eeffd19d4a471dabe79fbbea15ab0f"



function App() {
  const [playlistURL, setPlaylistURL] = useState(""); // const [state, setState] = useState(initialState)
  const [tracks, setTracks] = useState([]);
  const [status, setStatus] = useState("");

  const extractPlaylistId = (url) => {
    // https://open.spotify.com/playlist/54urz9eVTb5kDaAhAh2vHY
    const excess = "https://open.spotify.com/playlist/"
    return url.slice(excess.length); 
  };

  const getPlaylistTracks = async(playlistId) => {
    // strta here 

    const headers = {
      Authorization: `Bearer ${access_token}`,
    };

    let allTracks = [];
    let nextUrl = `https://api.spotify.com/v1/playlists/${playlistId}/tracks`;

    while (nextUrl) {

      const res = await fetch(nextUrl, { headers }); // stop and wait for fetch to go to nextUrl
      const data = await res.json();

      if (data.items) {
        const currentTracks = data.items
          .filter((item) => item.track)
          .map((item) => {
            const track = item.track;
            const name = track.name;
            const artists = track.artists.map((a) => a.name).join(", ");
            return `${name} - ${artists}`;
          });
        allTracks = allTracks.concat(currentTracks);
      }

      nextUrl = data.next;
    }

    return allTracks;
  };

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
      setTracks(result);
      setStatus(`Fetched ${result.length} tracks`);
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
          <li key={i}>{t}</li>
        ))}
      </ul>
    </div>
  );
}

export default App;
