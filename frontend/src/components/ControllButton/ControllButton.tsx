import React from 'react'
import { sendMsg } from '../../api'
import { useParams } from 'react-router-dom';

function ControllButton({ val }: { val: string }) {
  const { gameId } = useParams();
  const HandleOnClick = () => {
    var Message = { "type": 2, "message": val, "gameId": gameId }
    sendMsg(JSON.stringify(Message))
  }
  return (
    <button onClick={HandleOnClick}>{val}</button>
  )
}

export default ControllButton