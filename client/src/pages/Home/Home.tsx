import { useState } from "react"

export default function Home() {
    const [gameId, setgameId] = useState("")
    return <div>
        <div>
            <h1>
                Home
            </h1>
        </div>

        <div>
            {/* 
                Send Post Request to /game 
                and get redirected to (302) /game/:gameId 
                Let the backend make the Game and store it in the hub
                if the player tried to join a game that isn't been made yet
                return 404
            */}
            <input type="submit" value="New Game" />
        </div>

        <div>
            <form action={`./game/${gameId}`} method="GET">
                {/* 
                    Redirect to /game/:gameId
                */}
                <input onChange={(e)=> setgameId(e.target.value)} type="text" value={gameId} placeholder="Game Id" required />
                <input type="submit" value="Join" />
            </form>
        </div>
    </div>
}