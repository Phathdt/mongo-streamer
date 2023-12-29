# mongo-streamer: Golang Authentication Service with OAuth and Dex IDP Support

## Description

The mongo-streamer project is an authentication service written in Golang, using PostgreSQL as the database, and Redis for storing session information. The project supports OAuth with Dex IDP, providing flexible and secure authentication capabilities.

## Getting Started

Use the following guide to get started with mongo-streamer on your machine.

### Requirements

1. Golang: [Install Golang](https://golang.org/doc/install)
2. PostgreSQL: [Install PostgreSQL](https://www.postgresql.org/download/)
3. Redis: [Install Redis](https://redis.io/download)

### Installation

1. Clone the project from the repository:

```bash
git clone https://github.com/phathdt/mongo-streamer.git
cd mongo-streamer
```

2. Install dependencies:

```bash
go mod tidy
```

3. Set up configurations in the `.env` file. Change the values to reflect your configuration:

```
DB_DSN=postgresql://postgres:123123123@localhost:15432/mongo-streamer?sslmode=disable
REDIS_URI=redis://localhost:16379
DEX_CLIENT_ID=mongo-streamer
DEX_CLIENT_SECRET=c03879da7f12f890a537b3cacef1569a8493c471
```

4. Run the service:

```bash
task mongo-streamer
```

5. The project will run at `http://localhost:4000`.

## Usage

Certainly! Below is an enhanced usage section with descriptions for the mentioned APIs:

## Usage

The mongo-streamer authentication service provides the following APIs for user authentication:

- **Signup**: Create a new user account.

```http
POST /auth/signup
```

Example Request:
```bash
curl -X POST http://localhost:4000/auth/signup \
 -H "Content-Type: application/json" \
 -d '{"email": "exampleuser", "password": "secretpassword"}'
```

- **Login**: Authenticate an existing user.

```http
POST /auth/login
```

Example Request:
```bash
curl -X POST http://localhost:4000/auth/login \
 -H "Content-Type: application/json" \
 -d '{"email": "exampleuser", "password": "secretpassword"}'
```

- **OAuth Connect**: Initiate OAuth authentication.

```http
GET /auth/connect
```

Example Request:
```bash
curl -X GET http://localhost:4000/auth/connect?connector_id=xxx
```

- **OAuth Callback**: Handle OAuth callback after authentication.

```http
GET /auth/callback
```

Example Request:
```bash
curl -X GET http://localhost:4000/auth/callback
```

- **Get User Profile**: Retrieve the user's profile information.

```http
GET /auth/me
```

Example Request:
```bash
curl -X GET http://localhost:4000/auth/me \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

- **Check Token Validity**: Check the validity of an access token.

```http
GET /auth/valid
```

Example Request:
```bash
curl -X GET http://localhost:4000/auth/valid \
 -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

Please ensure proper authentication and authorization headers as required by each API. Refer to the [API documentation](API.md) for detailed information on each endpoint.

## Contribution

We welcome contributions from the community. If you wish to contribute, please read the [Contribution Guide](CONTRIBUTING.md).

## Release Notes

See [Changelog](CHANGELOG.md) for details about versions and updates.

## Author

- Phathdt
- Email: phathdt379@gmail.com
- GitHub: [phathdt](https://github.com/phathdt)

## License

The project is distributed under the MIT License - See the [LICENSE](LICENSE) file for details.

## Contact

- If you have any questions or feedback, please contact us via email: phathdt379@gmail.com

## Additional Resources

- [Dex IDP Documentation](https://dexidp.io/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Redis Documentation](https://redis.io/documentation)
