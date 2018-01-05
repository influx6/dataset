GeckoDataset
----------------
GeckoDataset is a project similar to the [Sql-Dataset](https://github.com/geckoboard/sql-dataset), which allows the exporting of collection data into your Geckobaord account's dataset using a combination of processors, pullers and pushers.

It provides a more involved approach to the transformation of incoming data which then gets saved into a dataset of your choosing on your Geckoboard API account. It takes this approach to both allow flexibility and user control of what the data gets transformed to, in both shape and form. 


## Install

```bash
go get -u github.com/influx6/geckodataset/...
```

## Run

GeckoDataset provides a CLI tooling called `geckoboard-dataset` which is central means of using the project:

```bash
> geckoboard-dataset 
Usage: dataset [flags] [command] 

⡿ COMMANDS:
	⠙ push	Push data from a mongodb to the geckoboa


⡿ HELP:
	Run [command] help

⡿ OTHERS:
	Run 'dataset printflags' to print all flags of all commands.

⡿ WARNING:
	Uses internal flag package so flags must precede command name. 
	e.g 'dataset -cmd.flag=4 run'

```

It which exposes a `push` command which handles the necessary logic to process provided configuration for the retrieval, transformation and update of dataset data.

```bash
> geckoboard-dataset push -config config.toml
```

*GeckoDataset relies on [Toml](https://github.com/toml-lang/toml) for it's configuration.*

## Configuration

Listing below are different configuration for usage of the geckodataset CLI tooling for sourcing data either through a [MongoDB](htts://mongodb.com) database collection or through a json file or directory:

- Using Javascript Processor with JSON source file


```toml
interval= "60s"
pull_batch = 100
push_batch = 100
api_key = "your_api_key"

[[datasets]]
driver = "json-file"

	[datasets.conf]
	dataset = "user_sales_freq"
	source = "./fixtures/sales/user_sales.json"

	[datasets.conf.js]
	target = "transformDocument"
	main = "./fixtures/transforms/js/user_sales.js"
	libraries = ["./fixtures/transforms/js/support/types.js"]

	[[datasets.conf.fields]]
	name = "user"
	type = "string"

	[[datasets.conf.fields]]
	name = "scores"
	type = "number"
```

- Using Javascript Processor with MongoDB source


```toml
interval= "60s"
pull_batch = 100
push_batch = 100
api_key = "your_api_key"

[[datasets]]
driver = "mongodb"

	[datasets.conf]
	dataset = "user_sales_freq"
	source = "user_sales_collection"
	dest = "user_sales_metrics" # optional, we want to save transformed records here

	[datasets.conf.db]
	authdb = "admin"
	db = "machines_sales"
	user = "tobi_mach"
	password = "xxxxxxxxxxxx"
	host = "db.mongo.com:4500"

	[datasets.conf.js]
	target = "transformDocument"
	main = "./fixtures/transforms/js/user_sales.js"
	libraries = ["./fixtures/transforms/js/support/types.js"]

	[[datasets.conf.fields]]
	name = "user"
	type = "string"

	[[datasets.conf.fields]]
	name = "scores"
	type = "number"
```

- Binary Processor with JSON source file

```toml
interval= "60s"
pull_batch = 100
push_batch = 100
api_key = "your_api_key"

[[datasets]]
driver = "json-file"

	[datasets.conf]
	dataset = "user_sales_freq"
	source = "./fixtures/sales/user_sales.json"

	[datasets.conf.js]
	target = "transformDocument"
	main = "./fixtures/transforms/js/user_sales.js"
	libraries = ["./fixtures/transforms/js/support/types.js"]

	[[datasets.conf.fields]]
	name = "user"
	type = "string"

	[[datasets.conf.fields]]
	name = "scores"
	type = "number"
```

- Using Binary Processor with MongoDB source


```toml
interval= "60s"
pull_batch = 100
push_batch = 100
api_key = "your_api_key"

[[datasets]]
driver = "mongodb"

[datasets.conf]
dataset = "user_sales_freq"
source = "user_sales_collection"
dest = "user_sales_metrics" # optional, we want to save transformed records here

	# set fields that dataset must have

	[[datasets.conf.fields]]
	name = "user"
	type = "string"

	[[datasets.conf.fields]]
	name = "scores"
	type = "number"


	[datasets.conf.db]
	authdb = "admin"
	db = "machines_sales"
	user = "tobi_mach"
	password = "xxxxxxxxxxxx"
	host = "db.mongo.com:4500"

	[datasets.conf.binary]
	binary = "echo"

```

## Processors/Procs

GeckoDataset employs the idea of processors termed `procs` which provided functions internally that will take either a single record or a batch of records from the scanned mongodb collection and return as desired appropriate JSON response which will be stored into the user's Geckoboard dataset account.

#### Javascript

The type of processor is based on the usage of javascript file, which exposes a function which would be called to transform the provided json of incoming records into desired format, which then is transformed into json then is returned to the dataset system which umarshals and attempts to save into user's Geckboard dataset account.

GeckoDataset uses [Otto](https://github.com/robertkrimen/otto) which is a golang javascript runtime for executing javascript, it does not support event loops based functions like those of `setInterval` and `setTimeout`, but does provide support for majority of the javascript runtime code. See project page for more details.

#### Binaries

This type of processor is based on the the usage of executable binary, which either is written to read from stdin a json of a record list or has a function which reads from stdin a json of record list, which will process and return appropriate json list containing the formated records which then is pushed up to the Geckboard user's dataset account.

## Disclaimer

We strongly recommend that the user account you use with the dataset project binaries has the lowest level of permission necessary for retrieving records from the database if possible. 

We also strongly recommend that any storage pusher which saves data into your database has only record writing permission and no removal or updating of existing record permission.

Although the `dataset` project contains no code to perform any adverse effect on your database, but it still reads and saves (if asked to) processed records, hence to ensure no adverse effect, we highly recommend this advice is taking. I accept no responsibility for any adverse changes to your database due to accidentally running users with inappropriate permissions.
