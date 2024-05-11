import React, { useState } from "react"
import { isGameExist } from "./functions"




export default function Home() {
    const [gameId, setgameId] = useState("")
    const [status, setStatus] = useState("")

    const url: string = process.env.REACT_APP_API_URL as string


    const handleSubmitJoinGame = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setStatus("Loading ...")
        const gameExist = await isGameExist(gameId)
        if (!gameExist) {
            setStatus("Game not Found !")
            return
        }
        window.location.href = `./game/${gameId}`


    }
    console.log(url)
    return <div>
        <div>
            <h1>
                Home
            </h1>
        </div>

        <div>
            <form action={`${url}game`} method="POST">
                <input type="submit" value="New Game" />
            </form>
        </div>

        <div>
            <form onSubmit={handleSubmitJoinGame}>
                <input onChange={(e) => setgameId(e.target.value)} type="text" value={gameId} placeholder="Game Id" required />
                <input type="submit" value="Join" />
                <div>
                    <span>{status}</span>
                </div>
            </form>
        </div>
    </div>
}