package cataloge

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Cataloge struct {
	Id   int
	Name string
	Icon string
}

type Item struct {
	Id          int
	Name        string
	Seller      string
	CatalogeId  int
	Description string
	Price       float64
	Icon        string
	Reviews     []int
	Discount    int
	Latitude    float64
	Longitude   float64
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "storage/catalog.db")
	if err != nil {
		return nil, err
	}

	// Создание таблицы Cataloge
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS Cataloge (
            Id INTEGER PRIMARY KEY,
            Name TEXT,
            Icon TEXT
        )
    `)
	if err != nil {
		return nil, err
	}

	// Создание таблицы Item
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS Item (
            Id INTEGER PRIMARY KEY,
            Name TEXT,
            Seller TEXT,
            CatalogeId INTEGER,
            Description TEXT,
            Price REAL,
            Icon TEXT,
            Reviews TEXT,
            Discount INTEGER,
			Latitude    REAL,
			Longitude   REAL,
            FOREIGN KEY (CatalogeId) REFERENCES Cataloge(Id)
        )
    `)
	if err != nil {
		return nil, err
	}

	Cataloges = []Cataloge{
		{Id: 0, Name: "TV y Audio", Icon: "https://imgbly.com/ib/v5moMxV6uV.png"},
		{Id: 1, Name: "Tecnología", Icon: "https://imgbly.com/ib/AtvY66GPBW.png"},
		{Id: 2, Name: "Electrodomésticos", Icon: "https://imgbly.com/ib/GBblkvb9qb.png"},
		{Id: 3, Name: "Hogar", Icon: "https://imgbly.com/ib/42ltgbDET7.png"},
		{Id: 4, Name: "Deportes y Fitness", Icon: "https://imgbly.com/ib/qiRy62mf8P.png"},
		{Id: 5, Name: "Juguetes, Bebés y Niños", Icon: "https://imgbly.com/ib/idUwF9NyU6.png"},
		{Id: 6, Name: "Belleza y Cuidado Personal", Icon: "https://imgbly.com/ib/kJmjW5EzTh.png"},
		{Id: 7, Name: "Herramientas y Construcción", Icon: "https://imgbly.com/ib/jqYgf7msIL.png"},
	}
	Items = []Item{
		{Id: 0, CatalogeId: 0, Name: "Smart TV 32\" HD Samsung UN32T4300A", Seller: "Samsung", Price: 289.999, Description: "the super smart tv", Reviews: []int{5, 5, 4, 3, 5}, Discount: 50, Icon: "https://imgbly.com/ib/40vMjKAnoX.png", Latitude: 37.78825, Longitude: -122.4324},
		{Id: 1, CatalogeId: 0, Name: "Smart TV 32\" HD Samsung UN32T4300A", Seller: "Samsung", Price: 289.999, Description: "the super smart tv", Reviews: []int{5, 5, 4, 3, 5}, Icon: "https://imgbly.com/ib/40vMjKAnoX.png"},
	}

	err = initData(db, Cataloges, Items)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initData(db *sql.DB, catalogs []Cataloge, items []Item) error {
	// Проверка и инициализация таблицы Cataloge
	catalogMap := make(map[int]bool)
	existingCatalogs, err := getCatalogeFromDB(db)
	if err != nil {
		return err
	}
	for _, c := range existingCatalogs {
		catalogMap[c.Id] = true
	}
	for _, c := range catalogs {
		if !catalogMap[c.Id] {
			stmt, err := db.Prepare("INSERT INTO Cataloge(Id, Name, Icon) VALUES(?, ?, ?)")
			if err != nil {
				return err
			}
			_, err = stmt.Exec(c.Id, c.Name, c.Icon)
			if err != nil {
				return err
			}
			stmt.Close()
		}
	}

	// Проверка и инициализация таблицы Item
	itemMap := make(map[int]bool)
	existingItems, err := getItemsFromDB(db)
	if err != nil {
		return err
	}
	for _, i := range existingItems {
		itemMap[i.Id] = true
	}
	for _, i := range items {
		if !itemMap[i.Id] {
			stmt, err := db.Prepare("INSERT INTO Item(Id, Name, Seller, CatalogeId, Description, Price, Icon, Reviews, Discount, Latitude, Longitude) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			if err != nil {
				return err
			}
			reviewsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(i.Reviews)), ","), "[]")
			_, err = stmt.Exec(i.Id, i.Name, i.Seller, i.CatalogeId, i.Description, i.Price, i.Icon, reviewsStr, i.Discount, i.Latitude, i.Longitude)
			if err != nil {
				return err
			}
			stmt.Close()
		}
	}

	return nil
}

func getCatalogeFromDB(db *sql.DB) ([]Cataloge, error) {
	rows, err := db.Query("SELECT Id, Name, Icon FROM Cataloge")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var catalogs []Cataloge
	for rows.Next() {
		var c Cataloge
		if err := rows.Scan(&c.Id, &c.Name, &c.Icon); err != nil {
			return nil, err
		}
		catalogs = append(catalogs, c)
	}

	return catalogs, nil
}

func getItemsFromDB(db *sql.DB) ([]Item, error) {
	rows, err := db.Query("SELECT Id, Name, Seller, CatalogeId, Description, Price, Icon, Reviews, Discount, latitude, longitude FROM Item")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var i Item
		var reviewsStr string
		if err := rows.Scan(&i.Id, &i.Name, &i.Seller, &i.CatalogeId, &i.Description, &i.Price, &i.Icon, &reviewsStr, &i.Discount, &i.Latitude, &i.Longitude); err != nil {
			if strings.Contains(err.Error(), "converting NULL to float64 is unsupported") {
				// Если ошибка связана с преобразованием NULL в float64, заменяем значения на 0.0
				i.Latitude = 0.0
				i.Longitude = 0.0
			} else {
				return nil, err
			}
		}
		i.Reviews = parseReviews(reviewsStr)
		items = append(items, i)
	}

	return items, nil
}

func parseReviews(reviewsStr string) []int {
	if reviewsStr == "" {
		return []int{}
	}

	reviewsSlice := strings.Split(reviewsStr, ",")
	reviews := make([]int, len(reviewsSlice))

	for i, review := range reviewsSlice {
		parsedReview, err := strconv.Atoi(strings.TrimSpace(review))
		if err != nil {
			// обработка ошибки при неверном формате числа
			return []int{}
		}
		reviews[i] = parsedReview
	}

	return reviews
}

var Cataloges []Cataloge
var Items []Item

func Setup(rg *gin.RouterGroup) {

	db, err := initDB()
	if err != nil {
		// Обработка ошибки
		log.Fatal(err)
	}
	defer db.Close()

	Cataloges, err = getCatalogeFromDB(db)
	if err != nil {
		// Обработка ошибки
		log.Fatal(err)
	}

	Items, err = getItemsFromDB(db)
	if err != nil {
		// Обработка ошибки
		log.Fatal(err)
	}

	api := rg.Group("cataloge")
	api.GET("", getCataloge)
	api.GET("items", getItems)
	api.POST("", addCataloge)  // Новый маршрут для добавления Cataloge
	api.POST("items", addItem) // Новый маршрут для добавления Item

}

func getCataloge(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	catalogs, err := getCatalogeFromDB(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"catalogs": catalogs})
}

func getItems(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	items, err := getItemsFromDB(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func addCataloge(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var cataloge Cataloge
	if err := c.BindJSON(&cataloge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("INSERT INTO Cataloge(Name, Icon) VALUES(?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(cataloge.Name, cataloge.Icon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cataloge added successfully"})
}

func addItem(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	var item Item
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("INSERT INTO Item(Name, Seller, CatalogeId, Description, Price, Icon, Reviews, Discount) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	reviewsStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(item.Reviews)), ","), "[]")

	_, err = stmt.Exec(item.Name, item.Seller, item.CatalogeId, item.Description, item.Price, item.Icon, reviewsStr, item.Discount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added successfully"})
}

// func getCataloge(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"Cataloge": Cataloges})
// }

// func getItems(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{"Items": Items})
// }
