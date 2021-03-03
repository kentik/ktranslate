package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

/**
Generate input file with

for file in `ls /home/pye/src/observium/observium/mibs/cisco/`; do snmptranslate -M ./mibs:/usr/share/snmp/mibs:~/tmp/mibs:~/src/observium/observium/mibs:/home/pye/src/observium/observium/mibs/cisco -m /home/pye/src/observium/observium/mibs/cisco/$file -IR -On -Totd; done > numbers

for file in `ls /home/pye/src/observium/observium/mibs/cisco/`; do snmptranslate -M ./mibs:/usr/share/snmp/mibs:~/tmp/mibs:~/src/observium/observium/mibs:/home/pye/src/observium/observium/mibs/cisco -m /home/pye/src/observium/observium/mibs/cisco/$file -IR -On -Tp; done > tree

*/

func main() {
	path := os.Args[1]
	oids := os.Args[2]

	db, err := leveldb.OpenFile(path, &opt.Options{})
	if err != nil {
		log.Fatalf("Cannot open db: %s -> %v", path, err)
	}
	defer db.Close()

	log.Printf("Starting to process %s", oids)
	st := time.Now()
	file, err := os.Open(oids)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var next []byte
	found := 0
	batch := new(leveldb.Batch)
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if bytes.HasPrefix(line, []byte(".")) {
			next = line
		} else {
			batch.Put(next, line)
			found++
		}
		if found%1000 == 0 {
			err = db.Write(batch, &opt.WriteOptions{})
			if err != nil {
				log.Fatalf("Cannot write batch db: %v", err)
			}
			batch.Reset()
		}
	}
	err = db.Write(batch, &opt.WriteOptions{})
	if err != nil {
		log.Fatalf("Cannot finish batch db: %v", err)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Done with %d oids in %v", found, time.Now().Sub(st))

	iter := db.NewIterator(util.BytesPrefix([]byte(".1.3.6.1.2.1.27")), nil)
	for iter.Next() {
		log.Printf("%s -> %s", string(iter.Key()), string(iter.Value()))
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		log.Fatalf("Cannot itter db: %v", err)
	}
}
