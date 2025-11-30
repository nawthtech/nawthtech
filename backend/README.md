#Noot Tech - E-Services Platform

## 🚀 Overview

This is the back-server of the Noth Tech e-services platform, built in the Go language using the Gin framework.

## 📁 Project Structure

Backend/

— api/

APIs #

(API)

— internal/

The internal code of the application #

— cmd/

Application entry points #

Email Templates # /templates

Deployments/ #files

## ⚡ Quick Start

### Basic requirements

- Go 1.21+

- MongoDB 5.0+

- Cloudinary Account

- Cloudflare account

- Office 365 account

### Installation

1. **Repository reproduction**

```bash

Git clone https://github.com/nawthtech/backend.git

Cd backend

#Fistaling credits

Make deps

#Preparing the environment

Cp .env.example .env

# Modify .env with your settings

#Run the application

Make run

#Available orders

# Building the application

Make build

# Development operation

Make run

# Running tests

Make a test

# migrate the database

Make migrate

#Docker building

Make docker

Supported Services

• Main database - MongoDB

Cloudinary - File Upload Service •

Cloudflare - CDN and Protection

• Email service - Office 365

5 Documentation

• Documentation API

• Publication directory

• Development guide

C contribution

Project Fork.1

Git checkout -b) create a branch of the feature .2

(Feature/ AmazingFeature)

Git push origin.

Feature/ AmazingFeature)

5. Opening a merger request

| Licensing

This project is licensed under the MIT license - see file

LICENSED

Communication

Email: support@nawthtech.com

• Website: https://nawthtech.com

### **`backend/go.mod`**

```go

Module github.com/nawthtech/nawthtech/backend

Go 1.21

Require (

Github.com/gin-gonic/gin v1.9.1

Github.com/nawthtech/nawthtech/backend v0.0.0

Go.mongodb.org/mongo-driver v1.12.1

Github.com/cloudinary/cloudinary-go/v2 v2.5.1

Github.com/urfave/cli/v2 v2.25.7

Gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df

)

Require (

// Other appropriations...

)