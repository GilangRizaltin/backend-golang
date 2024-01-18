![Logo](https://res.cloudinary.com/doncmmfaa/image/upload/v1705476586/samples/Frame_13_ksk8wi.png)

# Backend Coffee Shop with Golang

This project is about to show you on my performance in developing backend architecture using Golang. It has couple of features and API also several security authorization. It is a website for purchasing cinema tickets with main features including a list of films and their details, ordering cinema tickets based on the desired time and place. There are 2 roles, namely Normal user and Admin. Its has authentication and authorization for several accessible pages based on role.

## Features

- Gin Gonic \
  Gin Gonic is a lightweight and fast web framework for Golang. It simplifies the process of building web applications by providing essential routing features and middleware support. In your code, import Gin and utilize its powerful features to effortlessly handle HTTP requests and responses.

- JSON Web Token \
  JSON Web Tokens provide a secure and compact way to transmit information between parties. In your project, JWTs can be employed for user authentication and authorization. Generate a token when a user logs in and include it in subsequent requests to ensure secure communication between the client and server.

- Cloudinary \
  Cloudinary is a cloud-based service for managing and optimizing images and videos. Integrate Cloudinary into your project to effortlessly upload, store, and manipulate media assets. Leverage its API to dynamically transform images, ensuring optimal performance and user experience.

- Midtrans \
  Midtrans is a payment gateway service that simplifies online transactions. Integrate Midtrans into your application to facilitate secure and seamless payment processing. Utilize its APIs to handle payment requests, confirmations, and other transactions, providing users with a reliable and efficient payment experience.

- Govalidator \
  Govalidator is a versatile validation library for Golang. Integrate it into your project to easily validate user input and ensure data integrity. Employ Govalidator's functions to validate fields such as email addresses, URLs, and other form inputs, enhancing the robustness of your application.

## Installation

Install my-project with dependencies based on go mod

```bash
  go get .
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

| Method   | Endpoint | Description                     |
| :------- | :------- | :------------------------------ |
| `get`    | `"/"`    | Get all users (admin only)      |
| `post`   | `"/"`    | Add user (admin only)           |
| `patch`  | `"/"`    | Update users detail and profile |
| `delete` | `"/:id"` | Deleting user                   |

#### Products

```http
  /product
```

| Method   | Endpoint         | Description                                         |
| :------- | :--------------- | :-------------------------------------------------- |
| `get`    | `"/"`            | Get all product                                     |
| `get`    | `"/productstat"` | Get statistic for selled product (admin only)       |
| `get`    | `"/popular"`     | Get all popular and favourite product               |
| `get`    | `"/:id"`         | Get all product detail **Required** order_id        |
| `post`   | `"/"`            | Add a product (admin only)                          |
| `patch`  | `"/:id"`         | Edit a product **Required** order_id (admin only)   |
| `delete` | `"/:id"`         | Deleting product **Required** order_id (admin only) |

#### Promos

```http
  /promo
```

| Method   | Endpoint | Description      |
| :------- | :------- | :--------------- |
| `get`    | `"/"`    | Get all promos   |
| `post`   | `"/"`    | Post some promos |
| `patch`  | `"/:id"` | Edit some promo  |
| `delete` | `"/:id"` | Deleting promo   |

#### Orders

```http
  /order
```

| Method   | Endpoint | Description                  |
| :------- | :------- | :--------------------------- |
| `get`    | `"/"`    | Get all orders               |
| `post`   | `"/"`    | Create transaction           |
| `patch`  | `"/:id"` | Updating orders statuses     |
| `delete` | `"/:id"` | Deleting order (soft delete) |

## Documentation

[Postman Documentation](https://documenter.getpostman.com/view/29696636/2s9YsRbURh)

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`DB_HOST`,
`DB_NAME`,
`DB_USER`,
`DB_PASSWORD`,
`JWT_KEY`,
`ISSUER`,
`CLOUDINARY_NAME`,
`CLOUDINARY_KEY`,
`CLOUDINARY_SECRET`,
`MIDTRANS_ID_MERCHANT`,
`MIDTRANS_CLIENT_KEY`,
`MIDTRANS_SERVER_KEY`

## Run Locally

Clone the project

```bash
  git clone https://github.com/GilangRizaltin/backend-golang
```

Go to the project directory

```bash
  cd my-project
```

Install dependencies

```bash
  go get .
```

Start the server

```bash
  go run ./cmd/main.go
```

## Running Tests

To run tests, run the following command

```bash
  go test
```

## Front End Project

https://github.com/GilangRizaltin/Coffee-Shop-React

## Back End Project with Express Js

https://github.com/GilangRizaltin/CoffeeShop

## Support

For support, email gilangzaltin@gmail.com or join our Slack channel.

## Authors

Authors By Me _as known as_ Gilang Muhamad Rizaltin \
**Github link**

- [@Gilang Muhamad Rizaltin](https://github.com/GilangRizaltin)
