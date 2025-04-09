# RAG-Courses
A RAG-based lookup for USFCA Courses uses chromadb vector embedding
## Start the database:
    ```sh
    rm -rf chromadb && docker compose up
    ```

## Running with prompt for querying
    ```sh
    go run .
    ```

## Running the test
    ```sh
    go test
    ```
# Credits
All code is made by me for my CS-272 courses (Software Development)
