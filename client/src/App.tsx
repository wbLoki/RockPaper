import React, { useEffect, useState } from 'react';
import './App.css';
import { connect, disconnect, sendMsg } from './api';
import ControllButton from './components/ControllButton';
import Chat from './components/ChatComp';

function App() {
  const [IncommingMsg, setIncommingMsg] = useState("")
  const [MessagesList, setMessagesList] = useState([] as any)
  const [Refresher, setRefresher] = useState("")

  const setMoreMessages = (msg: any) => {
    setRefresher(msg)
    const messageJson = JSON.parse(msg)
    const message = messageJson.message
    MessagesList.push(messageJson)
    if (messageJson.type === 3){
      setIncommingMsg(message)
    }
  }
  useEffect(()=> {
    connect(setMoreMessages)
  },[IncommingMsg, Refresher])

  return (
    <div className="App">
      <h1>ROckPAperSCissor Game</h1>
      <h4>{IncommingMsg}</h4>
      <ControllButton val="rock" />
      <ControllButton val="paper" />
      <ControllButton val="scissors" />
      <br />
      {/* <button onClick={HandleOnClick}>Send Message</button> */}
      <br />
      <br />
      <br />

      <Chat History={MessagesList} />

      <br />
      <button onClick={disconnect}>Disconnect</button>
    </div>
  );
}

export default App;
