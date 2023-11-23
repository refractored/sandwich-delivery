## 🥪 Sandwich Delivery 🥪
Sandwich Delivery is a Discord bot in which people can order any sandwich they desire (even non-existent ones), have it made in our Discord by our top-of-the-line sandwich artists, and then have it delivered to your Discord server.

[Join our Discord](https://discord.gg/uHVzYBt8jD) for quick responses, support, and to see the bot in action!



## Self Hosting
### Bare Metal
#### Requirements

* MariaDB 10 or later
* A Discord server
* [Discord bot token](https://discord.com/developers/applications)

#### Setting Up 

To get started, you first get an executable, you can either compile yourself (see [compilation](#compilation)) or download the latest release: [stable](../../releases/latest) or [nightly/beta](https://nightly.link/refractored/sandwich-delivery/workflows/go/master/sandwich-delivery-linux-amd64.zip).

Once you have the executable, you need to create a config file in the working directory, you can use the [example config](config.json.example) as a template.
You can find a configuration reference [here](#configuration).

Once you have the config file, you can run the executable and you're off!

### Docker
#### Requirements

* Docker
* Docker Compose (optional)
* MariaDB 10 or later (must be accessible from the container)
* A Discord server
* [Discord bot token](https://discord.com/developers/applications)

#### Setting Up

For Docker, you can use our example [docker-compose.yml](docker-compose.yml.example) file as a reference.
You can use our image, or build one yourself, see [compilation](#compilation) (long story short, you can just change the line containing `image` to `build: .`).

Once you have a container, you need to create a config file in the volume attached to `/app/work`, you can use the [example config](config.json.example) as a template.
You can find a configuration reference [here](#configuration).

Once you have the config file, you can run the container with `docker-compose up -d` (or `docker-compose up` if you want to see the logs).



## Configuration

### Config File

The config file is a JSON file that contains all the configuration for the bot.
The config file is located in the working directory and is called `config.json`.

### Config Reference

| Key                     | Type                          | Description                                                                                                        | Optional?                                                 |
|-------------------------|-------------------------------|--------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------|
| `token`                 | `string`                      | The Discord bot token                                                                                              | No                                                        |
| `owners`                | `array` of `string`           | The Discord IDs of the bot owners                                                                                  | One is required                                           |
| `guildID`               | `string`                      | The Discord ID of the guild that certain commands will be registered in (e.g. `/blacklist`)                        | No                                                        |
| `database`              | `object`                      | The database configuration                                                                                         | No                                                        |
| `database.host`         | `string`                      | The database host                                                                                                  | No, unless `database.URL` is set, then it must be empty.  |
| `database.port`         | `int`                         | The database port                                                                                                  | No, unless `database.URL` is set, then it must be empty.  |
| `database.user`         | `string`                      | The database user                                                                                                  | No, unless `database.URL` is set, then it must be empty.  |
| `database.password`     | `string`                      | The database password                                                                                              | No, unless `database.URL` is set, then it must be empty.  |
| `database.database`     | `string`                      | The database name                                                                                                  | No, unless `database.URL` is set, then it must be empty.  |
| `database.extraOptions` | `map` of `string` -> `string` | Extra options to pass to the database connection                                                                   | Yes, unless `database.URL` is set, then it must be empty. |
| `database.URL`          | `string`                      | The database URL/DSN. A reference can be found [here](https://github.com/go-sql-driver/mysql#dsn-data-source-name) | Yes, unless other database entries are not set.           |
| `tokensPerOrder`        | `int`                         | The amount of tokens a user gets per order. Must be > 0 and defaults to 1.                                         | Yes                                                       |
| `dailyTokens`           | `int`                         | The amount of tokens a user gets per day. Must be > 0 and defaults to 1.                                           | Yes                                                       |

### !! Warning !!
If you are using the `database.URL` option and do not add `parseTime=true` as an option, you will run into issues.
This will also occur if you add an entry to `database.extraOptions` with the key `parseTime` and the value `false`.



## Compilation
### Getting the Source

To get the source code, you can either download the source code from the [latest release](../../releases/latest) or clone the repository with git:

```bash
git clone https://github.com/refractored/sandwich-delivery
```

Then you can `cd` into the directory.

### Bare Metal Executable
#### Requirements

* Go 1.21.4 or later

#### Building

After you have the source, building is as easy as one command:

```bash
go build -v -o sandwich-delivery src/main
```

A new executable should be created in the working directory called `sandwich-delivery`.

### Docker Image
#### Requirements

* Docker
* Docker Compose (optional)

#### Building

After you have the source, building is as easy as one command:

```bash
docker build -t sandwich-delivery .
```

**OR**

You can use the [docker-compose.yml](docker-compose.yml.example) file as a reference, but replace `image: ` with `build: .`

Then you can run the container with `docker-compose up -d` (or `docker-compose up` if you want to see the logs).



## Bugs

If you find any bugs or have a suggestion, please [create an issue](../../issues/new)!

For quick responses and community support, you can also [join our Discord](https://discord.gg/uHVzYBt8jD).



## Contributing

If you want to contribute, you can create a pull request, currently, however, there are no guidelines for contributing apart from the [code of conduct](CODE_OF_CONDUCT.md).



## License

This project is licensed under the [GPLv3 License](LICENSE).



## Credits

* [Refractored](https://refractored.net) - Project Owner & Lead Developer
* [Bacon](https://baconing.tech) - Developer
* [All our contributors](../../graphs/contributors)
* [All our dependencies](../../network/dependencies)
* And you for using it ❤️