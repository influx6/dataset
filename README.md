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

#### `mgo-dataset`


## Procs

MgoDataset employs the idea of processors termed `procs` which provided functions internally that will take either a single record or a batch of records from the scanned mongodb collection and return as desired appropriate JSON response which will be stored into the user's Geckoboard dataset account.

### Javascript

### Binaries