var s;
var gameId;
var state = new GameState()
function connectToWs() {
    const socket = new WebSocket(`ws://localhost:8080/ws?token=${gameId}`);
    s = socket
    // Connection opened
    socket.addEventListener("open", (event) => {
        console.log("Connected to ws://localhost:8080")
    });

    // Listen for messages
    socket.addEventListener("message", (event) => {
        console.log(event.data)
        // if (event.data == "gameon") {
        //     document.getElementsByTagName("body")[0].className = "gameon"
        // }
        state.set(event.data)

    });
}

function send_message() {
    event.preventDefault()
    var message = document.getElementById("msg")
    var roiMessage = message.value
    s.send(roiMessage)
    message.value = ""
    return 1
}
function generateMixedId(length) {
    const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    let mixedId = "";

    for (let i = 0; i < length; i++) {
        const randomIndex = Math.floor(Math.random() * charset.length);
        mixedId += charset[randomIndex];
    }

    return mixedId;
}

function hideElements() {
    document.getElementById("ch1").className = "hide"
    document.getElementById("ch2").className = "hide"
    document.getElementById("gameid").className = "hide"
}

function createGame() {
    gameId = generateMixedId(10)
    var phaseElem = document.createElement("div")
    var gameContainer = document.getElementsByClassName("game-container")[0]
    var headerElement = document.createElement("h2")
    var underLine = document.createElement("u")
    var instrElem = document.createElement("p")
    phaseElem.className = "ph1"
    instrElem.textContent = "waiting for the 2nd player ..."
    underLine.textContent = gameId
    headerElement.textContent = `Game Id `
    headerElement.appendChild(underLine)
    phaseElem.appendChild(instrElem)
    phaseElem.appendChild(headerElement)
    gameContainer.appendChild(phaseElem)
    gameId += 1 // playerOne
    connectToWs()
    hideElements()
    setCookie('gameToken', gameId, 1);
}

function joinGame() {
    gameIdValue = document.getElementById("gameid").value
    if (gameIdValue.length != 10) {
        alert("Enter the game id")
        return false
    }
    gameId = gameIdValue
    gameId += 2 // playerTwo
    connectToWs()
    hideElements() 
    setCookie('gameToken', gameId, 1);
}

