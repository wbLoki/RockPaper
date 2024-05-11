export async function isGameExist(gameId: string): Promise<boolean> {
    const apiURL = process.env.REACT_APP_API_URL as string
    const response = await fetch(`${apiURL}game/${gameId}/valid`, { method: "GET" })
    if (response.status != 200) {
        return false
    }
    return true
}