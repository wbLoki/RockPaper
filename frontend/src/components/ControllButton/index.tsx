import React from 'react'
import { sendMsg } from '../../api'

function ControllButton({val}:{val:string}) {
    const HandleOnClick = ()=> {
        var Message = {"type":2,"message":val}
        sendMsg(JSON.stringify(Message))
    }
  return (
    <button onClick={HandleOnClick}>{val}</button>
  )
}

export default ControllButton