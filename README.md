# mongo-streamer: Golang streaming change from Mongodb to Nats jetstream


## Getting Started

Use the following guide to get started with mongo-streamer on your machine.

### Requirements

1. Golang
2. Nats
3. Mongodb

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

```bash
APP_ENV="dev"
COLLECTION_NAME=demo-collection
DB_NAME=demo-database
FIBER_PORT=4000
LOG_LEVEL="trace"
SUBJECT=demo
FullDocumentBeforeChange="true" # only available for mongodb 6.0
MONGO_URI="mongodb://localhost:27017"
NATS_SUB_URI="nats://localhost:4222"

```

4. Create a jetstream in nats

```bash
nats str add demo --config examples/demo.json
```

5. Run the service:

```bash
task mongo-streamer
```

6. The project will run at `http://localhost:4000`.

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
