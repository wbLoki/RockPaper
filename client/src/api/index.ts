var gameId = window.location.pathname.replace("/","")
var socket = new WebSocket(`ws://localhost:8080/game/${gameId}`)
console.log(gameId)

let connect = (cb: any) => {
    console.log("connecting")
    socket.onopen = () => {
        console.log("Successfully connected")

    }

    socket.onmessage = (msg) => {
        console.log("Message from socket: ", msg)   
        cb(msg.data)
    }

    socket.onclose = (event) => {
        console.log("Socket closed connected: ", event)
    }

    socket.onerror = (error) => {
        console.log("socket error: ", error)
    }
}

let disconnect = ()=> {
    socket.close()
}

let sendMsg = (msg: string) => {
    socket.send(msg)
}

export {connect, sendMsg, disconnect}