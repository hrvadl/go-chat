# Simple TCP client/server chat written in Go.

## To run this application
1. **Clone** the repo.
2. Run the server with the following command from the root of repo:
  ```bash
    go run server/cmd/main.go
  ```

## To connect to the client 
- Run client using similar comand
    ```bash
      go run client/cmd/main.go
    ```
- Or connect to the server directly using [netcat](https://www.geeksforgeeks.org/introduction-to-netcat/)(or other similar tool)
    ```bash
      nc localhost 5000
    ```
