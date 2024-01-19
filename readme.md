# Backend Coffee Shop with Golang

<div align="center">
  <img src="https://res.cloudinary.com/doncmmfaa/image/upload/v1705476586/samples/Frame_13_ksk8wi.png" alt="Logo" />
</div>

This project is about to show you on my performance in developing backend architecture using Golang. It has couple of features and API also several security authorization. It is a website for purchasing cinema tickets with main features including a list of films and their details, ordering cinema tickets based on the desired time and place. There are 2 roles, namely Normal user and Admin. Its has authentication and authorization for several accessible pages based on role.

## Technologies used in this project

- Gin Gonic \
  Gin Gonic is a lightweight and fast web framework for Golang. \
  [Gin Gonic Documentation](https://pkg.go.dev/github.com/gin-gonic/gin#section-readme)

- JSON Web Token \
  JSON Web Tokens provide a secure and compact way to transmit information between parties. \
  [JSON Web Token](https://jwt.io/introduction)

- Cloudinary \
  Cloudinary is a cloud-based service for managing and optimizing images and videos.
  [Cloudinary Documentation](https://cloudinary.com/documentation)

- Midtrans \
  Midtrans is a payment gateway service that simplifies online transactions. \
  [Midtrans Documentation](https://docs.midtrans.com/)

- Govalidator \
  Govalidator is a versatile validation library for Golang. \
  [Govalidator Documentation](https://github.com/asaskevich/govalidator)

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

```bash
  DB_HOST = "YOUR DB_HOST"
  DB_NAME = "YOUR DB_NAME"
  DB_USER = "YOUR DB_USER"
  DB_PASSWORD = "YOUR DB_PASSWORD"
  JWT_KEY = "YOUR JWT_KEY"
  ISSUER = "YOUR ISSUER"
  CLOUDINARY_NAME = "YOUR CLOUDINARY_NAME"
  CLOUDINARY_KEY = "YOUR CLOUDINARY_KEY"
  CLOUDINARY_SECRET = "YOUR CLOUDINARY_SECRET"
  MIDTRANS_ID_MERCHANT = "YOUR MIDTRANS_ID_MERCHANT"
  MIDTRANS_CLIENT_KEY = "YOUR MIDTRANS_CLIENT_KEY"
  MIDTRANS_SERVER_KEY = "YOUR MIDTRANS_SERVER_KEY"
```

## Run Locally

Clone the project

```bash
  $ git clone https://github.com/GilangRizaltin/backend-golang
```

Go to the project directory

```bash
  $ cd backend-golang
```

Install dependencies

```bash
  $ go get .
```

Start the server

```bash
  $ go run ./cmd/main.go
```

## Running Tests

To run tests, run the following command

```bash
  $ go test
```

## API Reference

#### Authentication & Authorization

```http
  /auth
```

| Method | Endpoint      | Description                        |
| :----- | :------------ | :--------------------------------- |
| `post` | `"/register"` | register user                      |
| `post` | `"/login"`    | get access and identity of user    |
| `post` | `"/logout"`   | delete access and identity of user |

#### Users

```http
  /user
```

| Method | Endpoint | Description                |
| :----- | :------- | :------------------------- |
| `get`  | `"/"`    | Get all users (admin only) |
| `post` | `"/"`    | Add user (admin only)      |

#### Products

```http
  /product
```

| Method | Endpoint     | Description                                    |
| :----- | :----------- | :--------------------------------------------- |
| `get`  | `"/"`        | Get all product                                |
| `get`  | `"/popular"` | Get all popular and favourite product          |
| `get`  | `"/:id"`     | Get all product detail **Required** product_id |

#### Promos

```http
  /promo
```

| Method | Endpoint | Description    |
| :----- | :------- | :------------- |
| `get`  | `"/"`    | Get all promos |

#### Orders

```http
  /order
```

| Method | Endpoint | Description        |
| :----- | :------- | :----------------- |
| `get`  | `"/"`    | Get all orders     |
| `post` | `"/"`    | Create transaction |

## Documentation

[Postman Documentation](https://documenter.getpostman.com/view/29696636/2s9YsRbURh)

## Related Project

[Front End (React JS)](https://github.com/GilangRizaltin/Coffee-Shop-React)

[Backend (Javascript)](https://github.com/GilangRizaltin/CoffeeShop)

## Support

For support, email gilangzaltin@gmail.com.
