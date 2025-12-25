# ProdBooster Tools

## Seed Database

To populate the database with sample data (5 todos, 15 events, 20 notes):

```bash
go run tools/cmd_seed.go
```

This will add realistic sample data to help you test and see the full potential of ProdBooster!

## Reset Database

To start fresh:

```bash
rm ~/.prodbooster/data.db
./prodbooster  # Creates new empty database
```

Then optionally seed it:

```bash
go run tools/cmd_seed.go
```
