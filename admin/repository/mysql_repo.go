package repository

import (
	"database/sql"
	"fmt"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

//MysqlAdminRepository struct implements AdminRepository interface
type MysqlAdminRepository struct {
	conn *sql.DB
}

//NewAdminRepository creates new object of MysqlAdminRepository
func NewAdminRepository(conn *sql.DB) *MysqlAdminRepository {
	return &MysqlAdminRepository{conn: conn}
}

//StoreAdmin creates admin
func (mr *MysqlAdminRepository) StoreAdmin(ad *entity.Admin) (*entity.Admin, []error) {
	fmt.Println("inside store admin")
	totalNumOfUsers := mr.CountMember("admin")

	for mr.SearchMember("admin", ad.AID) {
		totalNumOfUsers++
		ad.AID = fmt.Sprintf("AID%d", totalNumOfUsers)
	}

	_, err := mr.conn.Exec("INSERT INTO admin (AID,Firstname,Lastname,Password,Phonenumber,Email) values(?, ?, ?,?,?,?)", ad.AID, ad.Firstname, ad.Lastname, ad.Password, ad.Phonenumber, ad.Email)

	if err != nil {
		panic(err)
	}

	return ad, nil
}

// CountMember is a method that is used for counting the member of a table where our table name is provided as an argument.
func (mr *MysqlAdminRepository) CountMember(tableName string) (totalNumOfMembers int) {

	stmt, err := mr.conn.Prepare("SELECT COUNT(*) FROM " + tableName)
	if err != nil {
		return
	}
	row := stmt.QueryRow()
	row.Scan(&totalNumOfMembers)
	return

}

//Admin gets admin
func (mr *MysqlAdminRepository) Admin(email, password string) (*entity.Admin, []error) {
	row := mr.conn.QueryRow(" SELECT * FROM  admin WHERE Email = ?  AND Password = ?;", email, password)
	a := entity.Admin{}
	err := row.Scan(&a.AID, &a.Firstname, &a.Lastname, &a.Password, &a.Phonenumber, &a.Email)
	if err != nil {
		panic(err)
	}
	return &a, nil
}

// UpdateAdmin updates a given object with a new data
func (mr *MysqlAdminRepository) UpdateAdmin(a *entity.Admin) (*entity.Admin, []error) {

	_, err := mr.conn.Exec("UPDATE admin SET Firstname=?,Lastname=?, Password=?,Phonenumber=?,Email=? WHERE AID=?", a.Firstname, a.Lastname, a.Password, a.Phonenumber, a.Email, a.AID)
	if err != nil {
		panic(err)
	}
	return a, nil
}

// DeleteAdmin removes admin from a database by its id
func (mr *MysqlAdminRepository) DeleteAdmin(id string) (*entity.Admin, []error) {
	row := mr.conn.QueryRow(" SELECT * FROM  admin WHERE AID=?;", id)
	a := entity.Admin{}
	err := row.Scan(&a.AID, &a.Firstname, &a.Lastname, &a.Password, &a.Phonenumber, &a.Email)

	_, err = mr.conn.Exec("DELETE FROM admin WHERE AID=?", id)
	if err != nil {
		panic(err)
	}

	return &a, nil
}

// SearchMember is a method that is used for searching the member of a table where our table name is provided as an argument.
func (mr *MysqlAdminRepository) SearchMember(tableName string, columnValue string) bool {

	if tableName == "admin" {
		stmt, _ := mr.conn.Prepare("SELECT COUNT(*) FROM " + tableName + " WHERE AID=?")
		var totalNumOfMembers int
		row := stmt.QueryRow(columnValue)
		row.Scan(&totalNumOfMembers)

		if totalNumOfMembers > 0 {
			return true
		}
	}
	return false

}
