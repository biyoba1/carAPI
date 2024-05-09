package service

import (
	"TZ"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"math"
	"net/http"
	"strconv"
)

type Pagination struct {
	TotalItems   int
	ItemsPerPage int
	Page         int
}

type PageInfo struct {
	Items      []TZ.Car
	Pagination *Pagination
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.ItemsPerPage
}

func (p *Pagination) GetLimit() int {
	return p.ItemsPerPage
}

func (p *Pagination) GetTotalPages() int {
	return int(math.Ceil(float64(p.TotalItems) / float64(p.ItemsPerPage)))
}

func GetCar(db *sql.DB) ([]TZ.Car, error) {
	db, err := sql.Open("postgres", "postgres://car-user:car-password@localhost:33333/car?sslmode=disable")
	if err != nil {
		log.Fatal("Cannot open  database", http.StatusBadRequest)
	}
	defer db.Close()
	page := 1
	pagination := &Pagination{
		TotalItems:   getTotalPages(db),
		ItemsPerPage: 100,
		Page:         page,
	}

	regNums, err := getRegNums(db, pagination.GetOffset(), pagination.GetLimit())
	if err != nil {
		return nil, err
	}
	return regNums, nil
}

func getTotalPages(db *sql.DB) int {
	var totalItems int

	err := db.QueryRow("select count(*) from car").Scan(&totalItems)
	if err != nil {
		log.Fatal("Error getting total items", err)
	}

	return totalItems
}

func getRegNums(db *sql.DB, offset, limit int) ([]TZ.Car, error) {
	rows, err := db.Query("select * from car offset $1 limit $2", offset, limit)
	if err != nil {
		log.Println("Error getting reg numbers:", err)
		return nil, err
	}
	defer rows.Close()

	var cars []TZ.Car
	for rows.Next() {
		var car TZ.Car
		err := rows.Scan(
			&car.ID,
			&car.RegNum,
			&car.Mark,
			&car.Model,
			&car.Year,
			&car.OwnerName,
			&car.OwnerSurname,
			&car.OwnerPatronymic,
		)
		if err != nil {
			log.Println("Error scanning reg number:", err)
			return nil, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

func DeleteCar(c *gin.Context, db *sql.DB) (map[string]string, map[string]string, error) {
	db, err := sql.Open("postgres", "postgres://car-user:car-password@localhost:33333/car?sslmode=disable")
	if err != nil {
		log.Fatal("Cannot open database", http.StatusBadRequest)
	}
	defer db.Close()
	id := c.Query("id")
	var regNum, mark, model, year, ownerName, ownerSurname, ownerPatronymic string
	err = db.QueryRow("SELECT reg_num, mark, model, year, owner_name, owner_surname, owner_patronymic FROM car WHERE id = $1", id).Scan(&regNum, &mark, &model, &year, &ownerName, &ownerSurname, &ownerPatronymic)
	if err != nil {
		log.Println("Error fetching car details before deletion:", err)
		return nil, nil, err
	}

	_, err = db.Exec("DELETE FROM car WHERE id = $1", id)
	if err != nil {
		log.Println("Error deleting a car:", err)
		return nil, nil, err
	}

	deleteCar := map[string]string{}
	deleteCar["Registration_number:"] = regNum
	deleteCar["Mark_of_car:"] = mark
	deleteCar["Model:"] = model
	deleteCar["Year:"] = year
	deleteOwner := map[string]string{}
	deleteOwner["Name_of_owner:"] = ownerName
	deleteOwner["Surname_of_owner:"] = ownerSurname
	deleteOwner["Patronymic_of_owner:"] = ownerPatronymic
	return deleteCar, deleteOwner, nil
}

func PutCar(c *gin.Context, db *sql.DB) (map[string]string, map[string]string, error) {
	db, err := sql.Open("postgres", "postgres://car-user:car-password@localhost:33333/car?sslmode=disable")
	if err != nil {
		log.Fatal("Cannot open  database", http.StatusBadRequest)
	}

	id1 := c.Query("id")
	reg_num1 := c.Query("reg_num")
	mark1 := c.Query("mark")
	model1 := c.Query("model")
	year1 := c.Query("year")
	owner_name := c.Query("owner_name")
	owner_surname := c.Query("owner_surname")
	owner_patronymic := c.Query("owner_patronymic")

	var regNum, mark, model, year, ownerName, ownerSurname, ownerPatronymic string
	err = db.QueryRow("SELECT reg_num, mark, model, year, owner_name, owner_surname, owner_patronymic FROM car WHERE id = $1", id1).Scan(&regNum, &mark, &model, &year, &ownerName, &ownerSurname, &ownerPatronymic)
	if err != nil {
		log.Println("Error fetching internal details before update:", err)
		return make(map[string]string), make(map[string]string), err
	}
	updateQuery := "UPDATE car SET "
	updateValues := []interface{}{}
	if reg_num1 != "" {
		updateQuery += "reg_num=$1, "
		updateValues = append(updateValues, reg_num1)
	}
	if mark1 != "" {
		updateQuery += "mark=$" + strconv.Itoa(len(updateValues)+1)
		updateValues = append(updateValues, mark1)
	}
	if model1 != "" {
		updateQuery += "model=$" + strconv.Itoa(len(updateValues)+1)
		updateValues = append(updateValues, model1)
	}
	if year1 != "" {
		updateQuery += "year=$" + strconv.Itoa(len(updateValues)+1)
		updateValues = append(updateValues, year1)
	}
	if owner_name != "" {
		updateQuery += "owner_name=$" + strconv.Itoa(len(updateValues)+1)
		updateValues = append(updateValues, owner_name)
	}
	if owner_surname != "" {
		updateQuery += "owner_surname=$" + strconv.Itoa(len(updateValues)+1)
		updateValues = append(updateValues, owner_surname)
	}
	if owner_patronymic != "" {
		updateQuery += "owner_patronymic=$" + strconv.Itoa(len(updateValues)+1)
		updateValues = append(updateValues, owner_patronymic)
	}
	updateQuery += "WHERE id=$" + strconv.Itoa(len(updateValues)+1)
	updateValues = append(updateValues, id1)

	_, err = db.Exec(updateQuery, updateValues...)
	if err != nil {
		log.Println("Cannot update car information in the database:", err)
		return make(map[string]string), make(map[string]string), err
	}

	oldInfo := map[string]string{
		"reg_num":          regNum,
		"mark":             mark,
		"model":            model,
		"year":             year,
		"owner_name":       ownerName,
		"owner_surname":    ownerSurname,
		"owner_patronymic": ownerPatronymic,
	}
	newInfo := map[string]string{
		"reg_num":          reg_num1,
		"mark":             mark1,
		"model":            model1,
		"year":             year1,
		"owner_name":       owner_name,
		"owner_surname":    owner_surname,
		"owner_patronymic": owner_patronymic,
	}
	return oldInfo, newInfo, nil
}

func AddNewCars(c *gin.Context) (string, error) {
	var newCarsReq TZ.NewCarsRequest
	if err := c.ShouldBindJSON(&newCarsReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return "", err
	}

	for _, regNum := range newCarsReq.RegNums {
		err := sendRequest(regNum)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to external API"})
			return "", err
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "New cars addeчd successfully"})
	return "Машина успешно добавлена", nil
}

func sendRequest(regNum string) error {
	url := "https://api.example.com/addCar"
	requestBody, err := json.Marshal(map[string]string{
		"regNum": regNum,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
