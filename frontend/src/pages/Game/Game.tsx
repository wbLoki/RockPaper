import React, { useEffect, useState } from 'react';
import './App.css';
import { connect, disconnect } from '../../api';
import Chat from '../../components/ChatComp/ChatComp';

function App() {
  const [IncommingMsg, setIncommingMsg] = useState("")
  const [MessagesList, _setMessagesList] = useState([] as any)
  const [Refresher, setRefresher] = useState("")
  const [username, setUsername] = useState("")
  const [score, setScore] = useState(0)


  const setMoreMessages = (msg: any) => {
    setRefresher(msg)
    const messageJson = JSON.parse(msg)
    const message = messageJson.message
    if (messageJson.type === 3 || messageJson.type === 2) {
      setIncommingMsg(message)
    } else if (messageJson.type === 4) {
      setScore(messageJson.score)
      setIncommingMsg(message)
    } else if (messageJson.type === 5) {
      setUsername(messageJson.name)
      localStorage.setItem('player', JSON.stringify({ "name": messageJson.name, "score": messageJson.score }));
      return
    }

    MessagesList.push(messageJson)
    return
  }
  useEffect(() => {
    connect(setMoreMessages)
    const chatBox = document.getElementById("chat-box")
    if (chatBox) {
      chatBox.scrollTop = chatBox?.scrollHeight
    }
  }, [IncommingMsg, Refresher])

  return (
    <div className="App">
      <h1>ROckPAperSCissor Game</h1>
      <h2>Welcome {username}</h2>
      <h5>Score: {score}</h5>
      <ControllButton val="rock" />
      <ControllButton val="paper" />
      <ControllButton val="scissors" />
      <br />
      <Chat History={MessagesList} />
      <br />
      <button onClick={disconnect}>Disconnect</button>
    </div>
  );
}

export default App;
