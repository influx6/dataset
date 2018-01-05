Dataset
----------------
Dataset is a project similar to the [SqlDataset](https://github.com/geckoboard/sql-dataset), which allows the exporting of collection data into your Geckobaord account's dataset using a combination of processors, pullers and pushers.

It provides a more involved approach to the processing of incoming data both to allow flexibility and user control of what the data gets transformed in both shape and form by the user. 

It will bundle 

## Install

```bash
go get -u github.com/influx6/dataset/...
```

## Commands

#### `mgo`


#### `mgo`

## Processors/Procs

Dataset employs the idea of processors termed `procs` which provided functions internally that will take either a single record or a batch of records from the scanned mongodb collection and return as desired appropriate JSON response which will be stored into the user's Geckoboard dataset account.

#### Javascript

The type of processor is based on the usage of javascript file, which exposes a function which would be called to transform the provided json of incoming records into desired format, which then is transformed into json then is returned to the dataset system which umarshals and attempts to save into user's Geckboard dataset account.

#### Binaries

This type of processor is based on the the usage of executable binary, which either is written to read from stdin a json of a record list or has a function which reads from stdin a json of record list, which will process and return appropriate json list containing the formated records which then is pushed up to the Geckboard user's dataset account.

## Disclaimer

We strongly recommend that the user account you use with the dataset project binaries has the lowest level of permission necessary for retrieving records from the database if possible. 

We also strongly recommend that any storage pusher which saves data into your database has only record writing permission and no removal or updating of existing record permission.

Although the `dataset` project contains no code to perform any adverse effect on your database, but it still reads and saves (if asked to) processed records, hence to ensure no adverse effect, we highly recommend this advice is taking. I accept no responsibility for any adverse changes to your database due to accidentally running users with inappropriate permissions.
