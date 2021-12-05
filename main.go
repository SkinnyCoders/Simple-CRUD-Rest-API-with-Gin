package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//connection to database
var db *gorm.DB

func init() {
	var err error
	db, err =
		gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/biodata?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("Koneksi gagal!")
	}
	db.AutoMigrate(&biodata{})
}

//routing app
func main() {
	r := gin.Default()

	geni := r.Group("/api/biodata")
	{
		geni.GET("/", getBio)
		geni.POST("/", createBio)
		geni.GET("/:id", getBioById)
		// geni.PUT("/:id", updateBio)
		geni.DELETE("/:id", deleteBio)
	}
	r.Run()
}

//create a model
type (
	biodata struct {
		gorm.Model
		Nama   string `json:"nama"`
		Alamat string `json:"alamat"`
		Hobi   string `json:"hobi"`
	}
	transformedBiodata struct {
		ID     uint   `json:"id"`
		Nama   string `json:"nama"`
		Alamat string `json:"alamat"`
		Hobi   string `json:"hobi"`
	}
)

//function for create data
func createBio(c *gin.Context) {
	//buat variable untuk map database
	var bio biodata

	//variable ambil nama dari post
	namaBio := c.PostForm("nama")

	//error cek saat ada nama yang sama
	err := db.Model(&bio).Where("nama = ?", namaBio).Find(&bio).Error

	//
	if err == nil {
		c.JSON(500, gin.H{
			"message": "Failed",
			"result":  "Data dengan nama " + namaBio + " sudah ada",
		})
		return
	} else {
		nama := c.PostForm("nama")
		alamat := c.PostForm("alamat")
		hobi := c.PostForm("hobi")

		bio := biodata{Nama: nama, Alamat: alamat, Hobi: hobi}

		// db.Save(&bio)

		err := db.Model(&bio).Save(&bio).Error

		if err != nil {
			c.JSON(500, gin.H{
				"message": "Failed",
				"result":  "Tambah Gagal",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Berhasil Input data",
				"result":  bio,
			})
		}
	}

	// result = gin.H{
	// 	"message": "Berhasil Input data",
	// 	"result": bio,
	// }
	// c.JSON(200, result)
}

//function for get all data
func getBio(c *gin.Context) {
	var bio []biodata
	var _bio []transformedBiodata

	db.Find(&bio)

	if len(bio) <= 0 {
		c.JSON(500, gin.H{
			"message": "Failed",
			"result":  "Tidak Ada Data",
		})
		return
	}

	for _, item := range bio {

		_bio = append(_bio, transformedBiodata{ID: item.ID, Nama: item.Nama, Alamat: item.Alamat, Hobi: item.Hobi})
	}

	c.JSON(200, gin.H{
		"message": "Success",
		"result":  _bio,
	})
}

// //function for get data by id
func getBioById(c *gin.Context) {
	var bio biodata
	bioID := c.Param("id")

	db.First(&bio, bioID)

	if bio.ID == 0 {
		c.JSON(500, gin.H{
			"message": "Failed",
			"result":  "Tidak Ada Data",
		})
		return
	}

	_bio := transformedBiodata{ID: bio.ID, Nama: bio.Nama, Alamat: bio.Alamat, Hobi: bio.Hobi}
	c.JSON(200, gin.H{
		"message": "Success",
		"result":  _bio,
	})
}

// //function for update data
// func updateBio(c *gin.Context) {
// 	nama := c.PostForm("nama")
// 	alamat := c.PostForm("alamat")
// 	hobi := c.PostForm("hobi")

// 	var _bio Bio
// 	var newBio Bio
// 	bioID := c.Param("id")

// 	err := db.First(&_bio, bioID).Error
// 	if err != nil {
// 		c.JSON(500, gin.H{
// 			"message": "Failed",
// 			"result":  "Tidak Ada Data",
// 		})
// 		return
// 	}

// 	newBio.Nama = nama
// 	newBio.Alamat = alamat
// 	newBio.Hobi = hobi
// 	err = db.Model(&_bio).Updates(newBio).Error
// 	if err != nil {
// 		c.JSON(500, gin.H{
// 			"message": "Failed",
// 			"result":  "Update Gagal",
// 		})
// 	} else {
// 		c.JSON(200, gin.H{
// 			"message": "Success",
// 			"result":  "Update Berhasil",
// 		})
// 	}
// }

//function for delete data
func deleteBio(c *gin.Context) {
	var bio biodata
	bioID := c.Param("id")

	db.First(&bio, bioID)

	if bio.ID == 0 {
		c.JSON(500, gin.H{
			"message": "Failed",
			"result":  "Tidak Ada Data",
		})
		return
	}

	db.Delete(&bio)
	c.JSON(200, gin.H{
		"message": "Success",
		"result":  "Data Berhasil Dihapus",
	})
}
