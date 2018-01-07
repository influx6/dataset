Below are the different configuration of the Geckodataset CLI using TOML:

- Using Javascript Processor with JSON source file


```toml
interval= "60s"
pull_batch = 100
push_batch = 100
api_key = "your_api_key"

[[datasets]]
driver = "json-file"
dataset = "user_sales_freq"

[[datasets.fields]]
name = "user"
type = "string"

[[datasets.fields]]
name = "scores"
type = "number"

[datasets.conf]
source = "./fixtures/sales/user_sales.json"

[datasets.conf.js]
target = "transformDocument"
main = "./fixtures/transforms/js/user_sales.js"
libraries = ["./fixtures/transforms/js/support/types.js"]
```

- Using Javascript Processor with JSON source directory


```toml
interval= "60s"
pull_batch = 100
push_batch = 100
api_key = "your_api_key"

[[datasets]]
driver = "json-dir"
dataset = "user_sales_freq"

[[datasets.fields]]
name = "user"
type = "string"

[[datasets.fields]]
name = "scores"
type = "number"

[datasets.conf]
source_dir = "./fixtures/sales"

[datasets.conf.jd]
 target = transformDocument
 main = "./fixtures/transforms/js/user_sales.js"
 libraries = ["./fixtures/transforms/js/support/types.js"]
```

- Using Javascript Processor with MongoDB source


```toml
interval= "60s"
pull_batch = 100
push_batch = 100
api_key = "your_api_key"

[[datasets]]
driver = "mongodb"
dataset = "user_sales_freq"

[[datasets.fields]]
name = "user"
type = "string"

[[datasets.fields]]
name = "scores"
type = "number"

[datasets.conf]
source = "user_sales_collection"
dest = "user_sales_metrics" # optional, we want to save transformed records here

[datasets.conf.db]
authdb = "admin"
db = "machines_sales"
user = "tobi_mach"
password = "xxxxxxxxxxxx"
host = "db.mongo.com:4500"

[datasets.conf.js]
 target = transformDocument
 main = "./fixtures/transforms/js/user_sales.js"
 libraries = ["./fixtures/transforms/js/support/types.js"]
```

- Binary Processor with JSON source file

```toml
interval= "60s"
pull_batch = 100
push_batch = 100
api_key = "your_api_key"

[[datasets]]
driver = "json-file"
dataset = "user_sales_freq"

[datasets.conf]
source = "./fixtures/sales/user_sales.json"

[datasets.conf.binary]
bin = "echo"

[[datasets.fields]]
name = "user"
type = "string"

[[datasets.fields]]
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
dataset = "user_sales_freq"

[[datasets.fields]]
name = "user"
type = "string"

[[datasets.fields]]
name = "scores"
type = "number"

[datasets.conf]
source = "user_sales_collection"
dest = "user_sales_metrics" # optional, we want to save transformed records here

[datasets.conf.db]
authdb = "admin"
db = "machines_sales"
user = "tobi_mach"
password = "xxxxxxxxxxxx"
host = "db.mongo.com:4500"

[datasets.conf.binary]
bin = "echo"
```
