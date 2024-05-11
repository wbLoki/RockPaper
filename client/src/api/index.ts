var gameId = window.location.pathname.replace("/", "")
var apiDomain = process.env.REACT_APP_API_DOMAIN
var socket = new WebSocket(`ws://${apiDomain}/${gameId}`)

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
        window.location.href = "/"
    }

    socket.onerror = (error) => {
        console.log("socket error: ", error)
        window.location.href = "/"
    }
}

let disconnect = () => {
    socket.close()
}

let sendMsg = (msg: string) => {
    socket.send(msg)
}

export { connect, sendMsg, disconnect }