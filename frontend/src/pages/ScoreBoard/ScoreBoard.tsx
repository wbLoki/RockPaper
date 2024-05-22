import React from 'react'

type Player = {
    score: number;
};

type Board = {
    [playerId: string]: Player;
};


function ScoreBoard({ board }: { board: Board | boolean }) {

    if (typeof board === 'boolean'){
        return <div>0-0</div>
    }

    const player1 = board[Object.keys(board)[0]].score
    const player2 = board[Object.keys(board)[1]].score

    return (
        <>

            <div>ScoreBoard</div>
            <div>{player1}-{player2}</div>
        </>
    )
}

export default ScoreBoard