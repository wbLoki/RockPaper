import React, { useState } from 'react'
import { sendMsg } from '../../api'
import { useParams } from 'react-router-dom';

function ControllButton({ val }: { val: string }) {
  const [clicked, setclicked] = useState("")
  const { gameId } = useParams();
  const HandleOnClick = () => {
    setclicked("HandImageClicked")
    var Message = { "type": 2, "message": val, "gameId": gameId }
    sendMsg(JSON.stringify(Message))
    setTimeout(() => setclicked(""), 200)
  }
  let Hand = null

  switch (val) {
    case "rock": {
      Hand = <RockComp />
      break
    }
    case "paper": {
      Hand = <PaperComp />
      break
    }
    case "scissors": {
      Hand = <ScissorComp />
      break
    }
    default: {
      console.error("ControllButton Error")
      break
    }
  }
  return (
    <div className={`Hand ${clicked}`} onClick={HandleOnClick}>
      {Hand}
    </div>
  )
}

function RockComp() {
  return (<div><img alt="Rock" src="/assets/Rock.jpg" /><p>Rock</p></div>)
}
function PaperComp() {
  return (<div><img alt="Paper" src="/assets/Paper.jpg" /><p>Paper</p></div>)
}
function ScissorComp() {
  return (<div>
    <img alt="Scissor" src="/assets/Scissor.jpg" />
    <p>Scissor</p>
  </div>)
}

export default ControllButton