# SEATUDY

Welcome to the Seatudy Backend repository! This project is the backend part of the Seatudy application, and it is responsible for handling the platform's core logic, data management, and API services.

## Table of Contents

- [About](#about)
- [Developers](#developers)
- [Documentation](#documentation)
- [Tech Stack](#tech-stack)
- [Features](#features)
- [Installation](#installation)

## About

Seatudy is an educational platform that allows users to manage their courses, track their progress, and provide feedback through reviews. This backend project offers a robust API to support these functionalities, handling all interactions with the database and ensuring secure and efficient data flow.

## Developers

### Back-end

- **Benardo** - Institut Teknologi Bandung
- **I Putu Natha Kusuma** - Universitas Brawijaya

### Front-end
- **Elgin Brian Wahyu Brahmandhika** - Universitas Brawijaya
- **Dindin Imanudin** - Institut Teknologi Nasional Bandung

## Documentation

### Entity Relationship Diagram (ERD)

![image](https://github.com/user-attachments/assets/775dc28c-9942-456d-bed0-ab323b33d7d8)

### Use Case Diagram

![image](https://github.com/user-attachments/assets/59f7d18a-fc82-45a3-9daa-75b9ec36985f)

### System Design Diagram

![image](https://github.com/user-attachments/assets/3c2418dd-ff42-48f9-b349-1c85f606b41d)

## Tech Stack

- **Golang**: The programming language used for the core backend logic.
- **Gin**: A high-performance HTTP web framework for building the RESTful API.
- **GORM**: An ORM library for Golang, used to interact with the PostgreSQL database.
- **PostgreSQL**: The primary database for storing all application data.
- **Redis**: Used for OTP management.
- **AWS S3**: For storing and serving static assets like course materials and profile pictures.
- **Docker**: Containerization for easy deployment and management.
- **JWT**: For handling user authentication and authorization.
- **Github Actions**: For test and deploy automation.

## Features

- **User Authentication**: Secure authentication using JWT.
- **Course Management**: Create, update, and manage courses, assignments, and materials.
- **Progress Tracking**: Track user progress and store completed assignments.
- **Review System**: Submit and manage course reviews.
- **Community Forum**: Discuss courses in a community forum.
- **Role-Based Access Control**: Different access levels for students and instructors.
- **File Uploads**: Securely manage file uploads to AWS S3.

## Installation

Before you begin, ensure you have the these installed on your machine:
- Docker

### Steps
1. **Clone the Repository:**
   ```bash
   git clone https://github.com/highfive-compfest/seatudy-backend.git
   ```
2. **Navigate to the Project Directory:**
   ```bash
   cd seatudy-backend
   ```
3. **Set Up Environment Variables:**  
   Create a `.env` file in the root directory and provide the necessary environment variables. See `.env.example` file for reference.

5. **Start the Server:**
   ```bash
   docker compose up
   ```
