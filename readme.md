# GOSSIP-API
### A PART OF GOSSIP INTERNAL ARCHITECTURE...

An API server that is written in go and provide APIs for signup, signin, etc.

### API Documentation

#### Base Path:
* `/api`

#### Paths
* `/v1/sendotp` :

    ```
        Method: POST
        In: Body
        Request-Type: json
            Properties:
                "number": string
        Response-Type: json
            Properties:
                "id": string
    ```

* `/v2/sendotp/{number}` :

    ```
        Method: GET
        In: Path
            Parameter:
                "number" : string
        Request-Type: string
        Response-Type: protobuf
            Properties:
                Status: type -> bool
                Disc: type -> string
                Data: type -> string
    ```





