package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"sort"
	"time"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
)

func (d *chromaDB) getOrCreateCollection(name string) (*chroma.Collection, bool) {
	existingCollection, err := d.client.GetCollection(context.TODO(), name, d.openaiEf)
	if err != nil {
		log.Printf("Error getting collection: %s ; Making new one\n", err)
	} else {
		return existingCollection, true
	}

	newCollection, err := d.client.CreateCollection(
		context.TODO(),
		name,
		nil,
		true,
		d.openaiEf,
		types.L2,
	)
	if err != nil {
		log.Fatalf("Error creating collection: %s \n", err)
	}
	return newCollection, false
}

// hashStrings returns a hash of the sorted strings and ensures the hash is less than 63 characters
func hashStrings(records []string) string {
	sortedRecords := make([]string, len(records))
	copy(sortedRecords, records)

	sort.Strings(sortedRecords)

	hash := sha256.New()

	for _, record := range sortedRecords {
		hash.Write([]byte(record))
	}

	hashString := hex.EncodeToString(hash.Sum(nil))

	// Ensure the hash is less than 63 characters
	if len(hashString) > 62 {
		hashString = hashString[:62]
	}

	return hashString
}

func (d *chromaDB) makeCollectionWithRecords(records []string) (*chroma.Collection, error) {
	startTime := time.Now()
	name := hashStrings(records)

	batchSize := 1000

	rs, err := types.NewRecordSet(
		types.WithEmbeddingFunction(d.openaiEf),
		types.WithIDGenerator(types.NewULIDGenerator()),
	)
	if err != nil {
		log.Fatalf("Error creating record set: %s \n", err)
	}
	// remove duplicates
	uniqueRecords := make(map[string]struct{})
	var deduplicatedRecords []string

	for _, record := range records {
		if _, exists := uniqueRecords[record]; !exists && record != "" {
			uniqueRecords[record] = struct{}{}
			deduplicatedRecords = append(deduplicatedRecords, record)
		}
	}

	records = deduplicatedRecords
	recordLength := len(records)

	collection, existed := d.getOrCreateCollection(name)
	if existed {
		log.Printf("Found \"%s\" with %v records in %.2f seconds", name[:5], len(records), time.Since(startTime).Seconds())
		return collection, nil
	}

	// Insert records in batches of `batchSize` records, and the last batch will have the remaining records
	for start := 0; start < recordLength; start += batchSize {
		end := start + batchSize
		if end > recordLength {
			end = recordLength
		}
		for _, record := range records[start:end] {
			rs.WithRecord(types.WithDocument(record))
		}
		_, err = rs.BuildAndValidate(context.TODO())
		if err != nil {
			log.Fatalf("Error validating record set: %s \n", err)
		}

		_, err = collection.AddRecords(context.Background(), rs)
		if err != nil {
			log.Fatalf("Error adding documents: %s \n", err)
		}
	}
	log.Printf("Created \"%s\" with %v records in %.2fs", name[:5], len(records), time.Since(startTime).Seconds())
	return collection, nil
}

func (d *chromaDB) query(collectionToQuery *chroma.Collection, query string, numResults int) ([]string, error) {
	startTime := time.Now()
	defer func() {
		elapsedTime := time.Since(startTime)
		log.Printf("Queried \"%s\" for \"%s\" in %.2fs\n", collectionToQuery.Name, query, elapsedTime.Seconds())
	}()
	qr, err := collectionToQuery.Query(context.TODO(), []string{query}, int32(numResults), nil, nil, nil)
	if err != nil {
		return nil, err
	}
	var results []string
	for _, doc := range qr.Documents {
		results = append(results, doc...)
	}
	return results, nil
}
