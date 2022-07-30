# GOSSIP-API
### A PART OF GOSSIP INTERNAL ARCHITECTURE...

An API server that is written in go and provide APIs for signup, signin, etc.

<video src='utils/video/myvideo.mp4' width=1000 height=300/>


### API Documentation

#### Base Path:
* `/api`

#### Paths
* `/v1/sendotp` :

    Discription: Send OTP to the number the is send by user.

    ```yaml
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

    Discription: Send OTP to the number which is provided by user in path.

    ```yaml
        Method: GET
        In: Path
            Parameter:
                "number" : string
        Request-Type: string
        Response-Type: protobuf
            Properties:
                "Status": type -> bool
                "Disc": type -> string
                "Data": type -> string
    ```

* `/v2/createnewuser` :

    Discription: Create new user and save the informations in the database.

    ```yaml
        Method: POST
        In: Body
        Request-Type: json
            properties:
                "id": string
                "name": string
                "dob": string
                "phoneno": string
                "email": string
                "profilepic": string
                "mainkey": pem-format-string
                "gender": string
                "password": base64 string
        Response-Type: protobuf
            Properties:
                "Status": type -> bool
                "Disc": type -> string
                "Data": type -> string


    ```





