import React, { useEffect, useState } from 'react'
import { sendMsg } from '../../api'
import { useParams } from 'react-router-dom';

type MessageType = {
    type: number;
    message: string;
    player: PlayerInfo
};

type PlayerInfo = {
    name: string,
    score: number
}

function Message({ messageJson }: { messageJson: MessageType }) {
    const messageText = messageJson.message
    const messageType = messageJson.type
    if (messageType !== 1) {
        return <p><b>{messageText}</b></p>
    }
    return <p>{messageJson.player.name}: {messageText}</p>
}


function Chat({ History }: { History: Array<MessageType> }) {
    const [ChatInput, setChatInput] = useState("")
    const { gameId } = useParams();
    const playerInfo = localStorage.getItem("player")

    const handleOnSubmit = (e: any) => {
        e.preventDefault()
        if (ChatInput == "") {
            return
        }
        let player = {}
        if (playerInfo) {
            player = JSON.parse(playerInfo)
        }
        var Message = {
            "type": 1, "message": ChatInput, "gameId": gameId, "player": player
        }
        sendMsg(JSON.stringify(Message))
        setChatInput("")
    }

    useEffect(() => {
        console.log("New Record")

    }, [History])
    return (
        <div>
            <div className='Chat' id="chat-box">
                {History.map((msg, key) => <Message key={key} messageJson={msg} />)}
            </div>
            <form onSubmit={handleOnSubmit}>
                <input type="text" onChange={(e) => setChatInput(e.target.value)} value={ChatInput} />
                <input type="submit" value="Send" />
            </form>
        </div>
    )
}

export default Chat