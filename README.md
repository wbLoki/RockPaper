# RockPaper
Making RockPaperScissor game using GO &amp; websockets

## Setup
### Backend
For running in `Local` 
Navigate to `backend` directory and run the following commands
```bash
make run # or just make
```
or 

```bash
export CLIENT_URL=http://localhost:3000/ # client Url here
go run ./cmd/main.go
```
### Frontend
Navigate to `frontend` directory </br>
install dependencies & create `.env` file see `.env.example`
```bash
npm install 
```
Run ReactApp
```bash
npm run start 
```