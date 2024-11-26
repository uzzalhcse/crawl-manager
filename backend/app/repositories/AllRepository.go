package repositories

import (
	"context"
	"crawl-manager-backend/app/models"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const DBName = "crawl_manager"

type Repository struct {
	DB *mongo.Client
}

func NewRepository(db *mongo.Client) *Repository {
	return &Repository{DB: db}
}

// SiteCollection CRUD
var siteCollection models.SiteCollection
var crawlingHistoryCollection models.CrawlingHistory

func (r *Repository) GetAllSiteCollections() ([]models.SiteCollection, error) {
	collection := r.DB.Database(DBName).Collection(siteCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := []models.SiteCollection{}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
func (r *Repository) CreateSiteCollection(siteCollection *models.SiteCollection) error {
	collection := r.DB.Database(DBName).Collection(siteCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, siteCollection)
	return err
}
func (r *Repository) CreateCrawlingHistory(crawlingHistory *models.CrawlingHistory) error {
	collection := r.DB.Database(DBName).Collection(crawlingHistoryCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, crawlingHistory)
	return err
}

func (r *Repository) GetSiteCollectionByID(siteID string) (*models.SiteCollection, error) {
	collection := r.DB.Database(DBName).Collection(siteCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var siteCollection models.SiteCollection
	err := collection.FindOne(ctx, bson.M{"site_id": siteID}).Decode(&siteCollection)
	if err != nil {
		return nil, err
	}
	return &siteCollection, nil
}

func (r *Repository) GetCrawlerFromHistory(instanceName string) (*models.CrawlingHistory, error) {
	collection := r.DB.Database(DBName).Collection(crawlingHistoryCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var crwCollection models.CrawlingHistory
	err := collection.FindOne(ctx, bson.M{"instance_name": instanceName, "status": "running"}).Decode(&crwCollection)
	if err != nil {
		return nil, err
	}
	return &crwCollection, nil
}

func (r *Repository) GetCrawlingHistoryByID(siteID string, runningOnly bool) ([]models.CrawlingHistory, error) {
	collection := r.DB.Database(DBName).Collection(crawlingHistoryCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var crwCollection []models.CrawlingHistory
	filter := bson.M{"site_id": siteID}
	if runningOnly {
		filter["status"] = "running"
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &crwCollection); err != nil {
		return nil, err
	}

	return crwCollection, nil
}

func (r *Repository) GetCrawlingHistory() ([]models.CrawlingHistory, error) {
	collection := r.DB.Database(DBName).Collection(crawlingHistoryCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	crwCollection := []models.CrawlingHistory{}

	// Use MongoDB sort to get the latest data based on StartDate
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"start_date", -1}}) // Sort by StartDate in descending order

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &crwCollection); err != nil {
		return nil, err
	}

	return crwCollection, nil
}

func (r *Repository) UpdateSiteCollection(siteID string, siteCollection *models.SiteCollection) error {
	collection := r.DB.Database(DBName).Collection(siteCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert siteCollection to a bson.M for partial updates
	updateData, err := bson.Marshal(siteCollection)
	if err != nil {
		return err
	}

	var updateFields bson.M
	if err := bson.Unmarshal(updateData, &updateFields); err != nil {
		return err
	}

	_, err = collection.UpdateOne(ctx, bson.M{"site_id": siteID}, bson.M{"$set": updateFields})
	return err
}
func (r *Repository) UpdateCrawlingHistory(instanceName string, update bson.M) error {
	collection := r.DB.Database(DBName).Collection(crawlingHistoryCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(ctx, bson.M{"instance_name": instanceName}, bson.M{"$set": update})
	return err
}

func (r *Repository) DeleteSiteCollection(siteID string) error {
	collection := r.DB.Database(DBName).Collection(siteCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"site_id": siteID})
	return err
}

// Collection CRUD
var collectionModel models.Collection

func (r *Repository) GetAllCollections() ([]models.Collection, error) {
	collection := r.DB.Database(DBName).Collection(collectionModel.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.Collection
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
func (r *Repository) CreateCollection(collection *models.Collection) error {
	collectionColl := r.DB.Database(DBName).Collection(collection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collectionColl.InsertOne(ctx, collection)
	return err
}

func (r *Repository) GetCollectionByID(collectionID string) (*models.Collection, error) {
	collectionColl := r.DB.Database(DBName).Collection(collectionModel.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var collection models.Collection
	err := collectionColl.FindOne(ctx, bson.M{"collection_id": collectionID}).Decode(&collection)
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

func (r *Repository) UpdateCollection(collectionID string, update bson.M) error {
	collectionColl := r.DB.Database(DBName).Collection(collectionModel.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collectionColl.UpdateOne(ctx, bson.M{"collection_id": collectionID}, bson.M{"$set": update})
	return err
}

func (r *Repository) DeleteCollection(collectionID string) error {
	collectionColl := r.DB.Database(DBName).Collection(collectionModel.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collectionColl.DeleteOne(ctx, bson.M{"collection_id": collectionID})
	return err
}

// UrlCollection CRUD
var urlCollection models.UrlCollection

func (r *Repository) CreateUrlCollection(urlCollection *models.UrlCollection) error {
	urlCollectionColl := r.DB.Database(DBName).Collection(urlCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := urlCollectionColl.InsertOne(ctx, urlCollection)
	return err
}

func (r *Repository) GetUrlCollectionByID(collectionID string) (*models.UrlCollection, error) {
	urlCollectionColl := r.DB.Database(DBName).Collection(urlCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var urlCollection models.UrlCollection
	err := urlCollectionColl.FindOne(ctx, bson.M{"collection_id": collectionID}).Decode(&urlCollection)
	if err != nil {
		return nil, err
	}
	return &urlCollection, nil
}

func (r *Repository) UpdateUrlCollection(collectionID string, update bson.M) error {
	urlCollectionColl := r.DB.Database(DBName).Collection(urlCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := urlCollectionColl.UpdateOne(ctx, bson.M{"collection_id": collectionID}, bson.M{"$set": update})
	return err
}

func (r *Repository) DeleteUrlCollection(collectionID string) error {
	urlCollectionColl := r.DB.Database(DBName).Collection(urlCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := urlCollectionColl.DeleteOne(ctx, bson.M{"collection_id": collectionID})
	return err
}

var siteSecretCollection models.SiteSecret

func (r *Repository) CreateSecretCollection(siteSecret *models.SiteSecret) error {
	collection := r.DB.Database(DBName).Collection(siteSecret.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"site_id": siteSecret.SiteID}
	opts := options.Replace().SetUpsert(true)

	_, err := collection.ReplaceOne(ctx, filter, siteSecret, opts)
	return err
}
func (r *Repository) GetAllSiteSecretCollections() ([]models.SiteSecret, error) {
	collection := r.DB.Database(DBName).Collection(siteSecretCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.SiteSecret
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
func (r *Repository) GetSiteSecretCollectionByID(siteID string) (*models.SiteSecret, error) {
	collection := r.DB.Database(DBName).Collection(siteSecretCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var secretCollection models.SiteSecret
	err := collection.FindOne(ctx, bson.M{"site_id": siteID}).Decode(&secretCollection)
	if err != nil {
		return nil, err
	}
	return &secretCollection, nil
}

var globalSecretCollection models.GlobalSecret

func (r *Repository) GetAllGlobalSecretCollections() ([]models.GlobalSecret, error) {
	collection := r.DB.Database(DBName).Collection(globalSecretCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.GlobalSecret
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

/*
Proxy Service
*/

var proxyCollection models.Proxy
var siteProxyCollection models.SiteProxy

//	func (r *Repository) CreateProxy(proxy *models.Proxy) error {
//		collection := r.DB.Database(DBName).Collection(proxy.GetTableName())
//		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//		defer cancel()
//
//		filter := bson.M{"server": proxy.Server}
//		opts := options.Replace().SetUpsert(true)
//
//		_, err := collection.ReplaceOne(ctx, filter, proxy, opts)
//		return err
//	}
func (r *Repository) GetAllProxy() ([]models.Proxy, error) {
	collection := r.DB.Database(DBName).Collection(proxyCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.Proxy
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// Fetch SiteProxies for each Proxy
	for i := range results {
		siteProxies, err := r.GetSiteProxiesByProxyID(results[i].ID) // Assuming you have a method to get SiteProxies
		if err != nil {
			return nil, err
		}
		results[i].SiteProxies = siteProxies
	}

	return results, nil
}

// GetSiteProxiesByProxyID fetches SiteProxies related to the provided ProxyID
func (r *Repository) GetSiteProxiesByProxyID(proxyID primitive.ObjectID) ([]models.SiteProxy, error) {
	collection := r.DB.Database(DBName).Collection(siteProxyCollection.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{"proxy_id": proxyID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var siteProxies []models.SiteProxy
	if err := cursor.All(ctx, &siteProxies); err != nil {
		return nil, err
	}

	return siteProxies, nil
}

// GetSiteProxiesBySiteID fetches SiteProxies related to the provided SiteID
func (r *Repository) GetSiteProxiesBySiteID(SiteID string) ([]models.Proxy, error) {
	// Initialize collections
	siteProxyCollections := r.DB.Database(DBName).Collection(siteProxyCollection.GetTableName())
	proxyCollections := r.DB.Database(DBName).Collection(proxyCollection.GetTableName())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Step 1: Fetch SiteProxy entries that match the given proxyID
	cursor, err := siteProxyCollections.Find(ctx, bson.M{"site_id": SiteID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var siteProxies []models.SiteProxy
	if err := cursor.All(ctx, &siteProxies); err != nil {
		return nil, err
	}

	// Step 2: Extract ProxyIDs and fetch corresponding Proxy documents
	var proxies []models.Proxy
	for _, siteProxy := range siteProxies {
		var proxy models.Proxy
		err = proxyCollections.FindOne(ctx, bson.M{"_id": siteProxy.ProxyID}).Decode(&proxy)
		if err != nil {
			return nil, err
		}
		if proxy.Valid {
			proxies = append(proxies, proxy)
		}
	}

	return proxies, nil
}

func (r *Repository) FindProxy(id string) (*models.Proxy, error) {
	// Initialize collections
	proxyCollections := r.DB.Database(DBName).Collection(proxyCollection.GetTableName())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert id to ObjectID if necessary
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	var proxy models.Proxy
	err = proxyCollections.FindOne(ctx, bson.M{"_id": objectID}).Decode(&proxy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("proxy with ID %s not found", id)
		}
		return nil, err
	}
	return &proxy, nil
}

func (r *Repository) DeleteProxy(id string) error {
	collectionName := proxyCollection.GetTableName()
	if collectionName == "" {
		return errors.New("collection name for Proxy is empty")
	}

	collection := r.DB.Database(DBName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert the ID string to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return err
	}

	// Create the filter
	filter := bson.M{"_id": objectID}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting proxy with ID %s: %v", id, err)
		return err
	}

	// Check if a document was actually deleted
	if result.DeletedCount == 0 {
		log.Printf("No proxy found with ID %s to delete", id)
		return mongo.ErrNoDocuments
	}

	log.Printf("Successfully deleted proxy with ID %s", id)
	return nil
}
func (r *Repository) UpdateProxies(proxies []models.Proxy) error {
	m := models.Proxy{}
	collection := r.DB.Database(DBName).Collection(m.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var bulkOps []mongo.WriteModel
	for _, proxy := range proxies {
		// Create a filter based on unique fields
		filter := bson.M{
			"proxy_address": proxy.ProxyAddress, // Replace with the actual field name
			"port":          proxy.Port,         // Replace with the actual field name
			"proxy_id":      proxy.ProxyID,      // Replace with the actual field name
		}

		// Create update operation
		update := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(bson.M{"$set": proxy}).SetUpsert(true)
		bulkOps = append(bulkOps, update)
	}

	// Execute bulk write operation
	if len(bulkOps) > 0 {
		_, err := collection.BulkWrite(ctx, bulkOps)
		if err != nil {
			log.Printf("Failed to bulk update proxies: %v", err)
			return err
		}
	}

	return nil
}

func (r *Repository) UpdateProxy(proxyID string, proxy *models.Proxy) error {
	m := models.Proxy{}
	collection := r.DB.Database(DBName).Collection(m.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert the ID string to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(proxyID)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return err
	}

	// Convert the proxy model to bson for updating fields
	update := bson.M{"$set": proxy}

	// Update the proxy with the provided fields
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		log.Printf("Failed to update proxy: %v", err)
		return err
	}

	return nil
}
func (r *Repository) AssignProxiesToSite(siteID string, proxyCount int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// If proxyCount is 0, remove all existing proxies for the site
	if proxyCount == 0 {
		_, err := r.DB.Database(DBName).Collection(siteProxyCollection.GetTableName()).DeleteMany(ctx, bson.M{"site_id": siteID})
		if err != nil {
			return fmt.Errorf("failed to remove proxies for site %s: %w", siteID, err)
		}
		return nil // No proxies to assign, so return immediately
	}

	// Check if the site already has enough proxies assigned
	existingProxyCount, err := r.DB.Database(DBName).Collection(siteProxyCollection.GetTableName()).CountDocuments(ctx, bson.M{"site_id": siteID})
	if err != nil {
		return err
	}
	if int(existingProxyCount) >= proxyCount {
		return nil // Proxies are already assigned, so exit the function
	}

	// Fetch all available proxies
	cursor, err := r.DB.Database(DBName).Collection(proxyCollection.GetTableName()).Find(ctx, bson.M{"valid": true})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var availableProxies []models.Proxy
	if err := cursor.All(ctx, &availableProxies); err != nil {
		return err
	}

	// Ensure there are available proxies
	if len(availableProxies) == 0 {
		return errors.New("no available proxies found")
	}

	// Count the usage of each proxy
	proxyUsageCounts := make(map[primitive.ObjectID]int)
	for _, proxy := range availableProxies {
		countCursor, err := r.DB.Database(DBName).Collection(siteProxyCollection.GetTableName()).CountDocuments(ctx, bson.M{"proxy_id": proxy.ID})
		if err != nil {
			return err
		}
		proxyUsageCounts[proxy.ID] = int(countCursor)
	}

	// Create a slice to sort proxies by usage count
	type proxyUsage struct {
		Proxy models.Proxy
		Count int
	}

	var sortedProxies []proxyUsage
	for _, proxy := range availableProxies {
		sortedProxies = append(sortedProxies, proxyUsage{Proxy: proxy, Count: proxyUsageCounts[proxy.ID]})
	}

	// Sort proxies by usage count in ascending order (least used first)
	sort.Slice(sortedProxies, func(i, j int) bool {
		return sortedProxies[i].Count < sortedProxies[j].Count
	})

	// Distribute proxies to sites
	for i := 0; i < proxyCount-int(existingProxyCount); i++ { // Adjust based on existing count
		// Find the least used proxy
		for _, proxyUsage := range sortedProxies {
			currentProxy := proxyUsage.Proxy

			// Check if this site proxy already exists
			existing := r.DB.Database(DBName).Collection(siteProxyCollection.GetTableName()).FindOne(ctx, bson.M{"site_id": siteID, "proxy_id": currentProxy.ID})
			if existing.Err() == mongo.ErrNoDocuments {
				// Create the site proxy assignment
				siteProxy := models.SiteProxy{
					SiteID:  siteID,
					ProxyID: currentProxy.ID,
				}

				// Insert the assignment into the site_proxies collection
				_, err := r.DB.Database(DBName).Collection(siteProxyCollection.GetTableName()).InsertOne(ctx, siteProxy)
				if err != nil {
					log.Println("Error inserting site proxy:", err)
					return err
				}

				break // Break out to assign the next site proxy
			}
		}
	}

	return nil
}

func (r *Repository) SaveCrawlingPerformance(crawlingPerformance *models.CrawlingPerformance) error {
	collection := r.DB.Database(DBName).Collection(crawlingPerformance.GetTableName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, crawlingPerformance)
	return err
}
