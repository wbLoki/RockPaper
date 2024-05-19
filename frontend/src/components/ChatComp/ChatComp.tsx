import React, { useEffect, useState } from 'react'
import { sendMsg } from '../../api'
import { useParams } from 'react-router-dom';

type MessageType = {
    type: number;
    message: string;
};

function Message({ messageJson }: { messageJson: MessageType }) {
    const messageText = messageJson.message
    const messageType = messageJson.type
    if (messageType === 3) {
        return <p><b>{messageText}</b></p>
    }
    return <p>{messageText}</p>
}


function Chat({ History }: { History: Array<MessageType> }) {
    const [ChatInput, setChatInput] = useState("")
    const { gameId } = useParams();

    const handleOnSubmit = (e: any) => {
        e.preventDefault()
        var Message = { "type": 1, "message": ChatInput, "gameId": gameId }
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