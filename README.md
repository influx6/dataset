MgoDataset
----------------
MgoDataset is a project similar to the [SqlDataset](https://github.com/geckoboard/sql-dataset), which allows the exporting of collection data into your Geckobaord account's dataset.

It provides a more involved approach to the processing of incoming data both to allow flexibility and user control of what the data gets transformed into both in shape and form. 

## Install

```bash
go get -u github.com/influx6/mgo-dataset
```

## Procs

MgoDataset employs the idea of processors termed `procs` which provided functions internally that will take either a single record or a batch of records from the scanned mongodb collection and return as desired appropriate JSON response which will be stored into the user's Geckoboard dataset account.

### Javascript file

### Binary file